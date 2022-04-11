package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func operation1(ctx context.Context) error {
	time.Sleep(100 * time.Millisecond)
	return errors.New("failed")
}

func operation2(ctx context.Context) {
	select {
	case <-time.After(500 * time.Millisecond):
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("halted operation2")
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		err := operation1(ctx)
		if err != nil {
			cancel()
		}
	}()

	operation2(ctx)
}

//package context

//type CancelCauseFunc func(cause error)

//func Cause(c Context) error

//func WithDeadlineCause(parent Context, d time.Time, cause error) (Context, CancelFunc)

//func WithTimeoutCause(parent Context, timeout time.Duration, cause error) (Context, CancelFunc)

//ctx, cancel := context.WithCancelCause(parent)
//cancel(myError)

//ctx.Err() // returns context.Canceled
//context.Cause(ctx) // returns myError

//在调用 WithCancelCause 或 WithTimeoutCause 方法后，会返回一个 CancelCauseFunc，而不是 CancelFunc。

//其差异之处在于：可以通过传入对应的 Error 等类型的信息，然后在调用Cause 方法来获取其被取消的根因错误。

//也就是既能得到被取消时的状态（context.Canceled），也能获取到对应的错误信息（myError），以此来解决前文中所提到的场景。
