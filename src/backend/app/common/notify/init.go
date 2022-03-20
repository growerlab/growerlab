// 全局接受系统消息（kill -INT 等消息）
//
package notify

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

var allOfDone = make(chan int, 0)
var notifySubscribes = make([]func(), 0)

func InitNotify() error {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			os.Interrupt,
			os.Kill,
			syscall.SIGQUIT,
			syscall.SIGSTOP,
			syscall.SIGUSR1,
			syscall.SIGUSR2,
		)
		<-c

		for _, sub := range notifySubscribes {
			sub()
		}

		time.Sleep(time.Second) // 等待结果的输出，避免过早结束进程，导致无法看到订阅函数的输出
		close(allOfDone)
	}()
	return nil
}

func Subscribe(fn func()) {
	notifySubscribes = append(notifySubscribes, fn)
}

func Done() {
	<-allOfDone
}
