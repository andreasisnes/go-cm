package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/andreasisnes/go-configuration-manager/tools/usersecrets/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/mod/modfile"
)

var App = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  ``,
	PersistentPreRun: func(ccmd *cobra.Command, args []string) {
	},
	Run: func(ccmd *cobra.Command, args []string) {
		ccmd.HelpFunc()(ccmd, args)
	},
}

func init() {
	App.PersistentFlags().StringP("module", "m", "", "Go module file")
	App.AddCommand(setCmd)
	App.AddCommand(clearCmd)
	App.AddCommand(listCmd)
	App.AddCommand(removeCmd)
}

func exit(cmd *cobra.Command, message string, statusCode int) {
	cmd.PrintErrln(message)
	os.Exit(statusCode)
}

func ParseModuleName(cmd *cobra.Command, args []string) {
	if cmd.Flag(util.ModuleFlag).Value.String() == "" {
		if pwd, err := os.Getwd(); err == nil {
			cmd.Flag(util.ModuleFlag).Value.Set(path.Join(pwd, "go.mod"))
		} else {
			exit(cmd, "Failed with: "+err.Error(), 1)
		}
	}

	if stat, err := os.Stat(cmd.Flag(util.ModuleFlag).Value.String()); err != nil {
		exit(cmd, "Failed with: "+err.Error(), 1)
	} else {
		if stat.IsDir() {
			cmd.Flag(util.ModuleFlag).Value.Set(path.Join(cmd.Flag(util.ModuleFlag).Value.String(), "go.mod"))
		}
	}

	if _, err := os.Stat(cmd.Flag(util.ModuleFlag).Value.String()); err == nil {
		if content, err := ioutil.ReadFile(cmd.Flag(util.ModuleFlag).Value.String()); err == nil {
			viper.Set(util.ModuleNameKey, modfile.ModulePath(content))
		} else {
			exit(cmd, "Failed with: "+err.Error(), 1)
		}
	} else {
		exit(cmd, "Failed with: "+err.Error(), 1)
	}
}
