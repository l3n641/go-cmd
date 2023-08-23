package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "system_helper",
	Short: "system_helper",
	Long:  `system_helper`,
	Args:  args,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func args(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("请选择要操作的命令!")
	}
	return nil
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(quickTopCmd)
}
