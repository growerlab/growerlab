package main

import (
	"context"
	"log"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	gggrpc "github.com/growerlab/growerlab/src/go-git-grpc"
)

func main() {
	clientCtx := context.Background()
	repoPath := "testrepo_bare"
	store, closeFn, err := gggrpc.NewStoreClient(clientCtx, "localhost:9001", repoPath)
	if err != nil {
		panic(err)
	}
	defer closeFn.Close()

	repo, err := git.Open(store, nil)
	if err != nil {
		panic(err)
	}

	// testReferences(repo)
	// testTags(repo)
	// testCommits(repo)
	// // tag下的文件列表
	// testFileTreesInTag(repo, "v1.0")

	// testAddTag(repo, "v10.0")
	testAddBranch(repo, "master2")

	time.Sleep(500 * time.Millisecond)
}

func testAddTag(repo *git.Repository, tagName string) {
	ref, err := repo.CreateTag(tagName, plumbing.NewHash("93450f0c98eeb96155cf8252ca9e07ca6c5ecbc2"),
		&git.CreateTagOptions{
			Tagger: &object.Signature{
				Name:  "moli2",
				Email: "mox2@out.com",
				When:  time.Now(),
			},
			Message: "hello",
			SignKey: nil,
		},
	)
	if err != nil {
		log.Fatalf("create tag was err: %+v", err)
	}
	log.Println("create tag: ", ref.String())
}

func testAddBranch(repo *git.Repository, branchName string) {
	err := repo.CreateBranch(&config.Branch{
		Name:   branchName,
		Remote: "origin",
		Merge:  "", // plumbing.NewBranchReferenceName("master"),
		Rebase: "",
	})
	if err != nil {
		log.Fatalf("create branch was err: %+v", err)
	}
}

func testFileTreesInTag(repo *git.Repository, tagName string) {
	iter, err := repo.TreeObjects()
	if err != nil {
		log.Fatalf("get trees was err:%+v", err)
	}
	var lastTree *object.Tree
	iter.ForEach(func(t *object.Tree) error {
		log.Printf("for each tag: %s \n", t.Hash.String())
		lastTree = t
		return nil
	})
	log.Printf("last tree: %s \n", lastTree.Hash.String())

	tag, err := repo.Tag(tagName)
	if err != nil {
		log.Fatalf("get tag '%s' was err: %+v", tagName, err)
	}
	log.Println("tag: ", tag.Strings())

	realTag, err := repo.TagObject(tag.Hash())
	if err != nil {
		log.Fatalf("get tag err: %+v", err)
	}

	tree, err := realTag.Tree()
	if err != nil {
		log.Fatalf("get tree '%s' was err: %+v", tagName, err)
	}

	tree.Files().ForEach(func(file *object.File) error {
		log.Printf("file %s in tag '%s'", file.Name, realTag.Name)
		return nil
	})
}

func testCommits(repo *git.Repository) {
	commitIter, err := repo.CommitObjects()
	if err != nil {
		log.Fatalf("get commits was err: %+v", err)
	}
	n := 0
	_ = commitIter.ForEach(func(c *object.Commit) error {
		n++
		log.Printf("commit committer: %s message: %s hash: %s\n",
			c.Committer.String(),
			c.Message,
			c.Hash.String(),
		)
		fileIter, err := c.Files()
		if err != nil {
			log.Fatalf("get files in %s was err: %+v", c.Hash.String(), err)
		}
		fileIter.ForEach(func(file *object.File) error {
			cts, err := file.Contents()
			if err != nil {
				log.Fatalf("get file %s in %s was err: %+v", file.Name, c.Hash.String(), err)
			}
			log.Printf("file '%s' content: %+v\n", file.Name, cts)
			return nil
		})
		return nil
	})
}

func testTags(repo *git.Repository) {
	refs, err := repo.Tags()
	if err != nil {
		log.Fatalf("get tags was err: %+v", err)
	}
	n := 0
	_ = refs.ForEach(func(r *plumbing.Reference) error {
		n++
		log.Printf("tag name: %s type: %d hash: %s string: %s target: %s\n",
			r.Name(),
			r.Type(),
			r.Hash().String(),
			r.String(),
			r.Target())
		return nil
	})
	if n == 0 {
		log.Fatalf("Not found tags")
	}
}

func testReferences(repo *git.Repository) {
	refs, err := repo.References()
	if err != nil {
		log.Fatalf("get branchs was err: %+v", err)
	}
	n := 0
	_ = refs.ForEach(func(r *plumbing.Reference) error {
		n++
		log.Printf("branch name: %s type: %d hash: %s string: %s target: %s\n",
			r.Name(),
			r.Type(),
			r.Hash().String(),
			r.String(),
			r.Target())
		return nil
	})
	if n == 0 {
		log.Fatalf("Not found branchs")
	}
}
