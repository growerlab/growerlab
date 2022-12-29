package logger

import (
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql/gqlerrors"
)

var LogWriter = gin.DefaultWriter

func Info(format string, a ...interface{}) {
	output("[Info]", format, a...)
}

func Warn(format string, a ...interface{}) {
	output("[Warn]", format, a...)
	// TODO warning
}

func Error(format string, a ...interface{}) {
	output("[Error]", format, a...)
	// TODO error
}

func GraphQLErrors(errs []gqlerrors.FormattedError) {
	sbErr := strings.Builder{}
	sbErr.WriteString("[GraphQL] ")
	for _, e := range errs {
		sbErr.WriteString(e.Error())
		sbErr.WriteByte('\n')
	}
	Info(sbErr.String())
	return
}

func output(prefix string, format string, a ...interface{}) {
	var sb = strings.Builder{}
	_, _ = sb.WriteString(prefix + " ")
	_, _ = sb.WriteString(fmt.Sprintf(format, a...))
	_ = sb.WriteByte('\n')

	_, _ = io.WriteString(LogWriter, sb.String())
}
