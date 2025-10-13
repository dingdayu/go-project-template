package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dingdayu/go-project-template/api/cron"
	"github.com/dingdayu/go-project-template/pkg/logger"
	pkgCron "github.com/robfig/cron/v3"
)

var c *pkgCron.Cron

func CronRun(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	c = pkgCron.New(pkgCron.WithLogger(logger.NewCronLogger(logger.Logger())))
	// 添加任务到调度器
	c.AddFunc("@every 5m", cron.CronTimer)

	// 启动 Cron 调度器
	c.Start()
	fmt.Println("\u001B[1;30;42m[info]\u001B[0m Cron started.")

	// 等待退出信号
	<-ctx.Done()

	// 关闭 Cron 调度器
	cronCtx := c.Stop()
	// 这里可以等待所有任务完成，或者设置一个超时时间
	<-cronCtx.Done()
	fmt.Println("\u001B[1;30;42m[info]\u001B[0m CRON exited")
}
