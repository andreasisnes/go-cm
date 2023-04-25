package cmd

import (
	"encoding/json"
	"os"

	"github.com/andreasisnes/go-configuration-manager/tools/usersecrets/util"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list usage",
	Short: "list",
	Long:  `list long`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			cmd.PrintErrln("Invalid number of arguments")
			os.Exit(1)
		}

		ParseModuleName(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := util.InitializeTree(util.GetModuleDir(), util.GetModuleSecretspath(), make(map[string]interface{})); err != nil {
			exit(cmd, err.Error(), 1)
		}

		secrets, err := util.ReadSecrets()
		if err != nil {
			exit(cmd, err.Error(), 1)
		}

		bytes, err := json.MarshalIndent(secrets, "", "\t")
		if err != nil {
			exit(cmd, err.Error(), 1)
		}

		cmd.Println(string(bytes))
	},
}
