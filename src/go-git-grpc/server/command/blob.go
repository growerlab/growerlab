package command

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strings"
	"time"

	"bitbucket.org/creachadair/shell"
	"github.com/asaskevich/govalidator"
	"github.com/bitfield/script"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/growerlab/growerlab/src/common/errors"
)

/*
该blob主要是为了解决在网页上创建文件的需求，
因为很多操作是在内存中，所以这里应该要限制内容的大小
*/

const maxBlobSize = 1024 * 1024 * 2 // 2MB

type Blob struct {
	ctx      *Context
	repoDir  string
	filePath string
	content  []byte
}

func NewBlob(root, filepath string, content []byte, ctx *Context) *Blob {
	repoDir := path.Join(root, ctx.RepoPath)
	return &Blob{repoDir: repoDir, filePath: filepath, content: content, ctx: ctx}
}

func (b *Blob) Commit(author object.Signature, message string, toRef string) (plumbing.Hash, error) {
	if err := b.ctx.Verify(); err != nil {
		return plumbing.ZeroHash, errors.Trace(err)
	}
	if govalidator.IsNull(message) {
		return plumbing.ZeroHash, errors.New("commit message is required.")
	}
	if govalidator.IsNull(toRef) {
		return plumbing.ZeroHash, errors.New("commit to ref/branch name is required.")
	}
	if len(b.content) == 0 || len(b.content) > maxBlobSize {
		return plumbing.ZeroHash, errors.New("invalid blob size, max 2MB")
	}
	if !strings.HasPrefix(toRef, "refs/") {
		toRef = plumbing.NewBranchReferenceName(toRef).String()
	}

	blobHash, err := b.buildBlobObject()
	if err != nil {
		return plumbing.ZeroHash, errors.Trace(err)
	}

	if err = b.buildIndex(blobHash); err != nil {
		return plumbing.ZeroHash, errors.Trace(err)
	}

	parentHash, err := b.getCurrentRefHash(toRef)
	if err != nil {
		return plumbing.ZeroHash, errors.Trace(err)
	}

	treeHash, err := b.doWriteTree()
	if err != nil {
		return plumbing.ZeroHash, errors.Trace(err)
	}

	commitHash, err := b.doCommitTree(author, message, treeHash, parentHash)
	if err != nil {
		return plumbing.ZeroHash, errors.Trace(err)
	}

	err = b.doUpdateRef(toRef, commitHash)

	return commitHash, errors.Trace(err)
}
func (b *Blob) doCommitTree(author object.Signature, message string, treeHash, parentHash plumbing.Hash) (plumbing.Hash, error) {
	cmd := fmt.Sprintf(`%s commit-tree %s -m "%s" -p %s`, b.ctx.Bin, treeHash, message, parentHash)
	envs := []string{
		fmt.Sprintf("GIT_AUTHOR_NAME=%s", author.Name),
		fmt.Sprintf("GIT_AUTHOR_EMAIL=%s", author.Email),
		fmt.Sprintf("GIT_AUTHOR_DATE=%d", time.Now().Unix()),
		fmt.Sprintf("GIT_COMMITTER_NAME=%s", author.Name),
		fmt.Sprintf("GIT_COMMITTER_EMAIL=%s", author.Email),
		fmt.Sprintf("GIT_COMMITTER_DATE=%d", time.Now().Unix()),
		fmt.Sprintf("EMAIL=%s", author.Email),
	}
	commitHash, err := b.runGitWithParams(nil, envs, cmd).String()
	if err != nil {
		return plumbing.ZeroHash, errors.Wrap(err, commitHash)
	}
	return plumbing.NewHash(commitHash), nil
}

func (b *Blob) doUpdateRef(ref string, commitHash plumbing.Hash) error {
	if commitHash.IsZero() {
		return errors.New("zero commit is illegal for update ref.")
	}

	cmd := fmt.Sprintf("%s update-ref %s %s", b.ctx.Bin, ref, commitHash)
	s, err := b.runGit(cmd).String()
	if err != nil {
		return errors.Wrap(err, s)
	}
	return nil
}

func (b *Blob) doWriteTree() (plumbing.Hash, error) {
	cmd := fmt.Sprintf("%s write-tree", b.ctx.Bin)
	treeHash, err := b.runGit(cmd).String()
	if err != nil {
		return plumbing.ZeroHash, errors.Wrap(err, treeHash)
	}
	return plumbing.NewHash(treeHash), nil
}

func (b *Blob) getCurrentRefHash(ref string) (plumbing.Hash, error) {
	cmd := fmt.Sprintf("%s show-ref --heads --hash %s", b.ctx.Bin, ref)
	hash, err := b.runGit(cmd).String()
	if err != nil {
		return plumbing.ZeroHash, errors.Wrap(err, hash)
	}
	if len(strings.TrimSpace(hash)) == 0 {
		return plumbing.ZeroHash, errors.Errorf("can't get ref hash for '%s'", ref)
	}
	return plumbing.NewHash(hash), nil
}

// 添加index文件
func (b *Blob) buildIndex(hash plumbing.Hash) error {
	cmd := fmt.Sprintf(`%s update-index --add --cacheinfo 100644 %s "%s"`, b.ctx.Bin, hash, b.filePath)
	s, err := b.runGit(cmd).String()
	if err != nil {
		return errors.Wrap(err, s)
	}
	return nil
}

func (b *Blob) buildBlobObject() (plumbing.Hash, error) {
	// write to objects/
	cmd := fmt.Sprintf("%s hash-object -w --stdin", b.ctx.Bin)
	blobHash, err := b.runGitWithParams(bytes.NewReader(b.content), nil, cmd).String()
	if err != nil {
		return plumbing.ZeroHash, errors.Wrap(err, blobHash)
	}
	return plumbing.NewHash(blobHash), nil
}

func (b *Blob) runGit(cmdline string) *script.Pipe {
	return b.runGitWithParams(nil, nil, cmdline)
}

func (b *Blob) runGitWithParams(stdin io.Reader, envs []string, cmdline string) *script.Pipe {
	pipe := script.NewPipe()
	if stdin != nil {
		pipe = pipe.WithReader(stdin)
	}
	pipe = pipe.Filter(func(r io.Reader, w io.Writer) error {
		args, ok := shell.Split(cmdline)
		if !ok {
			return errors.Errorf("unbalanced quotes or backslashes in [%s]", cmdline)
		}
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Env = append(cmd.Env, envs...)
		cmd.Stdin = r
		cmd.Stdout = w
		cmd.Stderr = w
		cmd.Dir = b.repoDir
		err := cmd.Start()
		if err != nil {
			fmt.Fprintln(w, err)
			return err
		}
		return cmd.Wait()
	})
	return pipe
}
