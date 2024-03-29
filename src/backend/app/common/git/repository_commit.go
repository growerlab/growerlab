package git

import (
	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/object/commitgraph"
	"github.com/growerlab/growerlab/src/common/errors"
)

type fileHash struct {
	name string
	hash plumbing.Hash
	mode filemode.FileMode
}

type commitAndPaths struct {
	commit commitgraph.CommitNode
	// Paths that are still on the branch represented by commit
	paths []string
	// Set of hashes for the paths
	hashes map[string]fileHash
}

// 这里性能比较差，后期可以考虑给commit加索引
func (r *Repository) getCommitForPaths(repo *git.Repository, hash plumbing.Hash, treePath string, paths []string) (map[fileHash]*object.Commit, error) {
	var result map[fileHash]*object.Commit
	nodeIndex := commitgraph.NewObjectCommitNodeIndex(repo.Storer)
	commitNode, err := nodeIndex.Get(hash)
	if err != nil {
		return nil, errors.Trace(err)
	}
	result, err = r.getLastCommitForPaths(commitNode, treePath, paths)
	return result, errors.Trace(err)
}

func (r *Repository) getLastCommitForPaths(c commitgraph.CommitNode, treePath string, paths []string) (map[fileHash]*object.Commit, error) {
	// We do a tree traversal with nodes sorted by commit time
	heap := binaryheap.NewWith(func(a, b interface{}) int {
		if a.(*commitAndPaths).commit.CommitTime().Before(b.(*commitAndPaths).commit.CommitTime()) {
			return 1
		}
		return -1
	})

	globalFileHashSet := make(map[string]fileHash)
	resultNodes := make(map[string]commitgraph.CommitNode)
	initialHashes, err := r.getFileHashes(c, treePath, paths)
	if err != nil {
		return nil, err
	}
	for fileName, fh := range initialHashes {
		globalFileHashSet[fileName] = fh
	}

	// Start search from the root commit and with full set of paths
	heap.Push(&commitAndPaths{
		commit: c,
		paths:  paths,
		hashes: initialHashes,
	})

	for {
		cIn, ok := heap.Pop()
		if !ok {
			break
		}
		current := cIn.(*commitAndPaths)

		// Load the parent commits for the one we are currently examining
		numParents := current.commit.NumParents()
		var parents []commitgraph.CommitNode
		for i := 0; i < numParents; i++ {
			parent, err := current.commit.ParentNode(i)
			if err != nil {
				break
			}
			parents = append(parents, parent)
		}

		// Examine the current commit and set of interesting paths
		pathUnchanged := make([]bool, len(current.paths))
		parentHashes := make([]map[string]fileHash, len(parents))
		for j, parent := range parents {
			parentHashes[j], err = r.getFileHashes(parent, treePath, current.paths)
			if err != nil {
				break
			}
			for fileName, fh := range parentHashes[j] {
				globalFileHashSet[fileName] = fh
			}

			for i, path := range current.paths {
				if parentHashes[j][path] == current.hashes[path] {
					pathUnchanged[i] = true
				}
			}
		}

		var remainingPaths []string
		for i, path := range current.paths {
			// The results could already contain some newer change for the same path,
			// so don't override that and bail out on the file early.
			if resultNodes[path] == nil {
				if pathUnchanged[i] {
					// The path existed with the same hash in at least one parent so it could
					// not have been changed in this commit directly.
					remainingPaths = append(remainingPaths, path)
				} else {
					// There are few possible cases how can we get here:
					// - The path didn't exist in any parent, so it must have been created by
					//   this commit.
					// - The path did exist in the parent commit, but the hash of the file has
					//   changed.
					// - We are looking at a merge commit and the hash of the file doesn't
					//   match any of the hashes being merged. This is more common for directories,
					//   but it can also happen if a file is changed through conflict resolution.
					resultNodes[path] = current.commit
				}
			}
		}

		if len(remainingPaths) > 0 {
			// Add the parent nodes along with remaining paths to the heap for further
			// processing.
			for j, parent := range parents {
				// Combine remainingPath with paths available on the parent branch
				// and make union of them
				remainingPathsForParent := make([]string, 0, len(remainingPaths))
				newRemainingPaths := make([]string, 0, len(remainingPaths))
				for _, path := range remainingPaths {
					if parentHashes[j][path] == current.hashes[path] {
						remainingPathsForParent = append(remainingPathsForParent, path)
					} else {
						newRemainingPaths = append(newRemainingPaths, path)
					}
				}

				if remainingPathsForParent != nil {
					heap.Push(&commitAndPaths{parent, remainingPathsForParent, parentHashes[j]})
				}

				if len(newRemainingPaths) == 0 {
					break
				} else {
					remainingPaths = newRemainingPaths
				}
			}
		}
	}

	// Post-processing
	result := make(map[fileHash]*object.Commit)
	for path, commitNode := range resultNodes {
		fh, found := globalFileHashSet[path]
		if found {
			result[fh], err = commitNode.Commit()
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func (r *Repository) getFileHashes(c commitgraph.CommitNode, treePath string, paths []string) (map[string]fileHash, error) {
	tree, err := r.getCommitTree(c, treePath)
	if err == object.ErrDirectoryNotFound {
		// The whole tree didn't exist, so return empty map
		return make(map[string]fileHash), nil
	}
	if err != nil {
		return nil, err
	}

	hashes := make(map[string]fileHash)
	for _, path := range paths {
		if path != "" {
			entry, err := tree.FindEntry(path)
			if err == nil {
				hashes[path] = fileHash{
					name: entry.Name,
					hash: entry.Hash,
					mode: entry.Mode,
				}
			}
		} else {
			hashes[path] = fileHash{
				name: "",
				hash: tree.Hash,
				mode: filemode.Dir,
			}
		}
	}

	return hashes, nil
}

func (r *Repository) getCommitTree(c commitgraph.CommitNode, treePath string) (*object.Tree, error) {
	tree, err := c.Tree()
	if err != nil {
		return nil, err
	}

	// Optimize deep traversals by focusing only on the specific tree
	if treePath != "" {
		tree, err = tree.Tree(treePath)
		if err != nil {
			return nil, err
		}
	}

	return tree, nil
}
