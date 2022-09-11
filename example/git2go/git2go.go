package main

import (
	"C"
	"fmt"
	"io/ioutil"

	git2go "github.com/libgit2/git2go/v28"
	"github.com/spf13/viper"
)
import "os/exec"

func commitExist(repoDir, commitId string) bool {
	cmd := exec.Command("git", "rev-list", fmt.Sprintf("HEAD..%s", viper.GetString("branch")))
	cmd.Dir = repoDir

	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Git exit code 128
			return false
		}
	}
	return true
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	dir, _ := ioutil.TempDir("", "temp_dir")
	fmt.Println(dir)
	options := &git2go.CloneOptions{
		CheckoutBranch: "main",
	}

	_, err = git2go.Clone(viper.GetString("url"), dir, options)
	if err != nil {
		panic(err)
	}
	// Fields that start with lower case characters are package internal and not exposed :(
	// rep.CheckoutIndex(&git2go.Index{
	// 	ptr:  viper.GetString("commit"),
	// 	repo: rep}, &git2go.CheckoutOptions{})
	// repo, err := git2go.OpenRepository(dir)
	// if err != nil {
	// 	panic(err)
	// }

	// commitId := git2go.Oid{}
	// copy(commitId[:], viper.GetString("commit"))

	// commit, err := repo.LookupCommit(&commitId)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(commit)
	// repo.CheckoutIndex()

	// TODO: differenciate between branch and commit
	if commitExist(dir, viper.GetString("branch")) {
		// Poor man solution to checkout to specific commit id :(
		cmd := exec.Command("git", "checkout", viper.GetString("branch"))
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

// Evaluation: Work, but checking out commit is not straight forward + require additional package as dependency (libgit2-dev)
// and gcc, pkg-config
