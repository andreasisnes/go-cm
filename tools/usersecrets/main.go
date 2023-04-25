package main

import (
	"fmt"
	"os"

	"github.com/andreasisnes/go-configuration-manager/tools/usersecrets/cmd"
)

func main() {
	if err := cmd.App.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
