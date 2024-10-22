package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ginLambdaALB *ginadapter.GinLambdaALB

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	r := configureRoutes(configureLogger())

	ginLambdaALB = ginadapter.NewALB(r)
	zap.L().Info("Gin router initialized")
}

func configureRoutes(logger *zap.Logger) *gin.Engine {
	zap.L().Info("Configuring routes")
	// Suppress GIN built-in init logs
	gin.SetMode(gin.ReleaseMode)

	// Gin router
	router := gin.New()

	// Middlewares
	router.Use(gin.Recovery())

	// Routes
	// Healthcheck: GET path = /healthcheck
	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.Writer.WriteHeader(http.StatusOK)
	})

	tests := router.Group("/tests")
	tests.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	tests.POST("", CreateItem())
	tests.GET("/:test_id", GetItem())
	tests.DELETE("/:test_id", DeleteItem())
	tests.PUT("/:test_id", UpdateItem())
	return router
}

func configureLogger() *zap.Logger {
	fmt.Println("Configuring logger")
	logLevel := viper.GetString("logs.level")
	conf := zap.NewProductionConfig()
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		log.Printf("fatal failed to set the log level, defaulting to info level: %s", logLevel)
		conf.Level.SetLevel(zap.InfoLevel)
	} else {
		conf.Level = level
	}
	conf.OutputPaths = []string{"stdout"}

	logger, err := conf.Build()
	if err != nil {
		log.Panicf("fatal failed to build the logger: %s", err)
	}

	zap.ReplaceGlobals(logger)
	return logger
}

func Handler(ctx context.Context, req events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	zap.L().Info("Handler called")
	// If no name is provided in the HTTP request body, throw an error
	return ginLambdaALB.ProxyWithContext(ctx, req)
}

func main() {
	zap.L().Info("Starting lambda application")
	lambda.Start(Handler)
}

func GetItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("GetItem called")
		reqCtx := ctx.Request.Context()
		defer reqCtx.Done()

		ctx.JSON(200, gin.H{"GetItem": "abc"})
	}
}

func UpdateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("UpdateItem called")
		reqCtx := ctx.Request.Context()
		defer reqCtx.Done()

		ctx.JSON(200, gin.H{"UpdateItem": "abc"})
	}
}

func DeleteItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("DeleteItem called")
		reqCtx := ctx.Request.Context()
		defer reqCtx.Done()

		ctx.JSON(200, gin.H{"DeleteItem": "abc"})
	}
}

func CreateItem() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zap.L().Info("CreateItem called")
		reqCtx := ctx.Request.Context()
		defer reqCtx.Done()

		ctx.JSON(200, gin.H{"CreateItem": "abc"})
	}
}
