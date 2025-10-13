package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/dingdayu/async/v4"
)

func init() {
	// Register your async tasks here.
	async.Register(&ExampleTimerAsync{})
}

type ExampleTimerAsync struct {
	ticker *time.Ticker
}

// OnPreRun Before run, panic panic causes registration failure
func (a *ExampleTimerAsync) OnPreRun() {
	a.ticker = time.NewTicker(1 * time.Minute)
	fmt.Printf("\u001B[1;30;42m[info]\u001B[0m ExampleTimerAsync 注册成功，开始运行！\n")
}

// Name async name
func (a ExampleTimerAsync) Name() string {
	return "ExampleTimerAsync"
}

// Handle async logical
func (a *ExampleTimerAsync) Handle(ctx async.Context) {
	defer ctx.Exit()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ExampleTimerAsync received done signal, exiting...")
			return
		case t := <-a.ticker.C:
			fmt.Printf("ExampleTimerAsync timer fired at %v\n", t)
			// Reset timer
			a.ticker.Reset(1 * time.Minute)
		}
	}
}

// OnShutdown on async shutdown
func (a *ExampleTimerAsync) OnShutdown(ctx context.Context) {
	a.ticker.Stop()
	fmt.Printf("\u001B[1;30;42m[info]\u001B[0m ExampleTimerAsync 准备退出！\n")
}
