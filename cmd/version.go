package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version info",
	Long:  `prints version info based on CI/CD pipeline info, used within 'build.sh'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(AppVersion)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
