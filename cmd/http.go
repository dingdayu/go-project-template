package cmd

import (
	"github.com/dingdayu/go-project-template/api"
	"github.com/dingdayu/go-project-template/model/dao"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var httpAsync bool

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start Singbox Adapter http server",
	Long:  "Start Singbox Adapter http server",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// redis.Init()
		if viper.GetString("db") != "" {
			dao.Setup()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if httpAsync {
			go api.AsyncRun(cmd.Context())
		}

		api.Run(cmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	// add --async flag to control whether to start async processing
	httpCmd.Flags().BoolVar(&httpAsync, "async", true, "Start async processing at server start")
	_ = viper.BindPFlag("http.async", httpCmd.Flags().Lookup("async"))
}
