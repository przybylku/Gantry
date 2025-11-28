package gitops

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)


func CloneRepo(gitaddr string) (path string, hash string, err error) {
	tempPath, err := os.MkdirTemp("", "gantry-repos-*")
	if err != nil {
		fmt.Println(err) // to remove
		return "", "", err
	}
	var repo *git.Repository
	repo, err = git.PlainClone(tempPath, false, &git.CloneOptions{
		URL:      gitaddr,
		Progress: os.Stdout,
		Depth: 1,
	})
	if err != nil {
		fmt.Println(err) // to remove
		os.RemoveAll(tempPath)
		return "", "", err
	}


	ref, err := repo.Head()

	// hash to 7 chars lenght
	fullHash := ref.Hash().String()
	shortHash := fullHash[:7]

	fmt.Println("Repository cloned to:", tempPath)
	if err != nil {
		fmt.Println(err) // to remove
		os.RemoveAll(tempPath)
		return "", "", err
	}
	return tempPath, shortHash, nil
}