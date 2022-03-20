package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

var (
	calldepth   = 2
	projectName = "hulk"
)

func NewLogger(prefix string) *Logger {
	// TODO 目前日志会放在仓库目录下，应转到其他目录统一起来
	f, err := os.OpenFile(fmt.Sprintf("%s.log", projectName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(errors.WithStack(err))
	}
	return &Logger{out: f, prefix: prefix}
}

type Logger struct {
	out    io.Writer
	prefix string
}

func (l *Logger) Flush() {
	if o, ok := l.out.(io.Closer); ok {
		_ = o.Close()
	}
}

func (l *Logger) Write(b []byte) (n int, err error) {
	var sb bytes.Buffer

	_, file, line, ok := runtime.Caller(calldepth)

	sb.WriteString(fmt.Sprintf("[%s] ", projectName))
	if len(l.prefix) > 0 {
		sb.WriteString(fmt.Sprintf("[%s]", l.prefix))
	}
	sb.WriteString(time.Now().Format(time.RFC3339))
	sb.WriteString(" ")
	if ok {
		sb.WriteString(fmt.Sprintf("%s:%d ", file, line))
	}
	sb.Write(b)

	nb, err := sb.WriteTo(l.out)
	return int(nb), err
}
