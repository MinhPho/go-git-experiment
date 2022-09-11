package main

import (
	"fmt"
	"io/ioutil"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/viper"
)

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

	options := &git.CloneOptions{
		URL:           viper.GetString("url"),
		Depth:         1,
		ReferenceName: "refs/heads/main",
		SingleBranch:  true,
		Tags:          git.NoTags,
	}

	_, err = git.PlainClone(dir, false, options)
	fmt.Println(err)

}

// Evaluation: Not an option for Azure DevOps repository since known issue "unexpected client error: unexpected requesting ... status code: 400"
