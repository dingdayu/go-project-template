package api

import (
	"context"
	"fmt"

	"github.com/dingdayu/async/v4"

	_ "github.com/dingdayu/go-project-template/api/async"
)

func AsyncRun(ctx context.Context) {
	async.Wait()
	fmt.Println("\u001B[1;30;42m[info]\u001B[0m Task exited")
}
