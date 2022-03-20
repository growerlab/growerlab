package controller

//
// import (
// 	"context"
// 	"fmt"
// 	"runtime/debug"
//
// 	gql "github.com/99designs/gqlgen/graphql"
// 	"github.com/99designs/gqlgen/handler"
// 	"github.com/gin-gonic/gin"
// 	"github.com/growerlab/growerlab/src/backend/app/common/errors"
// 	"github.com/growerlab/growerlab/src/backend/app/service/graphql"
// 	"github.com/growerlab/growerlab/src/backend/app/service/graphql/resolver"
// 	"github.com/growerlab/growerlab/src/backend/app/utils/logger"
// 	"github.com/vektah/gqlparser/gqlerror"
// )
//
// type errCauser interface {
// 	Cause() error
// }
// type errFormat interface {
// 	Error() string
// 	Format(s fmt.State, verb rune)
// }
//
// func GraphQL(ctx *gin.Context) {
// 	var userToken = GetUserToken(ctx)
// 	var session = graphql.NewSession(userToken, ctx)
// 	var graphqlOpts = make([]handler.Option, 0)
//
// 	sessionCtx := graphql.BuildContextWithSession(ctx, session)
// 	ctx.Request = ctx.Request.WithContext(sessionCtx)
//
// 	// options
// 	// reqOpt := handler.RequestMiddleware(func(ctx context.Context, next func(ctx context.Context) []byte) []byte {
// 	// 	reqCtx := gql.GetRequestContext(ctx)
// 	// 	opName := reqCtx.OperationName
// 	// 	return next(ctx)
// 	// })
// 	// graphqlOpts = append(graphqlOpts, reqOpt)
//
// 	// resOpt := handler.ResolverMiddleware(func(ctx context.Context, next gql.Resolver) (res interface{}, err error) {
// 	// 	logger.Info("----ResolverMiddleware")
// 	// 	return next(ctx)
// 	// })
// 	// graphqlOpts = append(graphqlOpts, resOpt)
//
// 	errorOpt := handler.ErrorPresenter(func(gqlCtx context.Context, err error) *gqlerror.Error {
// 		logger.Error("graphql err presenter: %+v", err)
//
// 		retErr := gql.DefaultErrorPresenter(gqlCtx, err)
// 		retErr.Message = err.Error()
// 		if retErr.Extensions == nil {
// 			retErr.Extensions = map[string]interface{}{}
// 		}
//
// 		switch e := err.(type) {
// 		case errCauser:
// 			retErr.Extensions["code"] = e.Cause().Error()
// 		case errFormat:
// 			retErr.Extensions["code"] = e.Error()
// 		}
// 		return retErr
// 	})
//
// 	recoverOpt := handler.RecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
// 		logger.Error("graphql recover err: %v\n%+v", err, string(debug.Stack()))
// 		return errors.New(errors.GraphQLError())
// 	})
//
// 	graphqlOpts = append(graphqlOpts, errorOpt, recoverOpt)
//
// 	fn := handler.GraphQL(resolver.NewExecutableSchema(resolver.Config{Resolvers: &resolver.Resolver{}}), graphqlOpts...)
// 	fn.ServeHTTP(ctx.Writer, ctx.Request)
// }
//
// func GraphQLPlayground() gin.HandlerFunc {
// 	fn := handler.Playground("GraphQL playground", "/api/graphql")
// 	return func(ctx *gin.Context) {
// 		fn.ServeHTTP(ctx.Writer, ctx.Request)
// 	}
// }
