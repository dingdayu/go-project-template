package cmd

import (
	"github.com/dingdayu/singbox-adapter/api"
	"github.com/dingdayu/singbox-adapter/model/dao"
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
