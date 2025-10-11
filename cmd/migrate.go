package cmd

import (
	"fmt"

	"github.com/dingdayu/go-project-template/model/dao"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  "Run database migrations for the current environment.",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		dao.Setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
		// It will change existing column's type if it's size, precision, nullable changed.
		// AutoMigrate will not delete unused columns to protect your data.
		// 参考：https://gorm.io/docs/migration.html#Auto-Migration

		err := dao.GetContextDB(cmd.Context()).AutoMigrate(dao.User{})
		if err != nil {
			fmt.Printf("❌ Failed to apply migrations: %v\n", err)
			return
		}

		fmt.Println("✅ All migrations applied successfully!")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
