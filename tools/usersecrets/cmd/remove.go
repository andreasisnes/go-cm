package cmd

import (
	"os"

	"github.com/andreasisnes/go-configuration-manager/tools/usersecrets/util"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove usage",
	Short: "remove",
	Long:  `remove long`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.PrintErrln("Invalid number of arguments")
			os.Exit(1)
		}

		ParseModuleName(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		secrets, err := util.ReadSecrets()
		if err != nil {
			exit(cmd, err.Error(), 1)
		}

		delete(secrets, key)
		if _, err := util.DumpSecrets(secrets); err != nil {
			exit(cmd, err.Error(), 1)
		}
	},
}
