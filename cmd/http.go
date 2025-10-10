package cmd

import (
	"github.com/dingdayu/go-project-template/api"
	"github.com/dingdayu/go-project-template/model/dao"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start Singbox Adapter http server",
	Long:  "Start Singbox Adapter http server",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// redis.Init()
		dao.Init()
	},
	Run: func(cmd *cobra.Command, args []string) {
		api.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
