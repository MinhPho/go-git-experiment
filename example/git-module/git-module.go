package main

import (
	"C"
	"fmt"
	"io/ioutil"

	gitmod "github.com/gogs/git-module"
	"github.com/spf13/viper"
)
import "os/exec"

func inputIsBranch(repoDir, commitId string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", fmt.Sprintf("refs/remotes/origin/%s", viper.GetString("branch")))
	cmd.Dir = repoDir

	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Git exit code 128
			return false
		}
	}
	return true
}

func inputIsCommitId(repoDir, commitId string) bool {
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

	options := gitmod.CloneOptions{
		Branch: "main",
	}

	err = gitmod.Clone(viper.GetString("url"), dir, options)
	if err != nil {
		panic(err)
	}

	if inputIsCommitId(dir, viper.GetString("branch")) {
		// Poor man solution to checkout to specific commit id :(
		cmd := exec.Command("git", "checkout", viper.GetString("branch"))
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	} else if inputIsBranch(dir, viper.GetString("branch")) {
		cmd := exec.Command("git", "checkout", viper.GetString("branch"))
		cmd.Dir = dir
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	} else {
		panic("input is neither a valid branch nor commit id")
	}
}

// Evaluation: Work like a charm :)
