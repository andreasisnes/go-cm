package cmd

import (
	"os"

	"github.com/andreasisnes/go-configuration-manager/tools/usersecrets/util"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	Long:  `set value`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.PrintErrf("Invalid number of arguments. given %d must be 2", len(args))
			os.Exit(1)
		}

		ParseModuleName(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := util.InitializeTree(util.GetModuleDir(), util.GetModuleSecretspath(), make(map[string]interface{})); err != nil {
			exit(cmd, err.Error(), 1)
		}

		secrets, err := util.ReadSecrets()
		if err != nil {
			exit(cmd, err.Error(), 1)
		}

		secrets[key] = value
		if _, err := util.DumpSecrets(secrets); err != nil {
			exit(cmd, err.Error(), 1)
		}
	},
}
