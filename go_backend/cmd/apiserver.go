/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	router "go_backend/api/rest"
	"go_backend/database/mysql"
	"go_backend/database/redis"
	"go_backend/global"
	logger "go_backend/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var apiserverCmd = &cobra.Command{
	Use:   "apiserver",
	Short: "start apiserver",
	Long:  `start apiserver`,
	Run:   RunApiserver,
}

func RunApiserver(cmd *cobra.Command, _ []string) {
	// init logger
	var err error
	if global.Logger, err = logger.InitLogger(); err != nil {
		panic(fmt.Sprintf("Init logger error: %s\n", err))
	}

	// init gin mode
	switch viper.GetString("server.runMode") {
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case gin.DebugMode:
		gin.SetMode(gin.DebugMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	default:
		errMsg := fmt.Sprintf("Unknown server runMode: %s\n", viper.GetString("server.runMode"))
		global.Logger.Error(errMsg)
		panic(errMsg)
	}

	// init mysql
	db, err := mysql.InitMySQL(cmd.Context())
	if err != nil {
		errMsg := fmt.Sprintf("Init MySQL error: %s\n", err)
		global.Logger.Error(errMsg)
		panic(errMsg)
	}

	// init redis
	rdClient, err := redis.InitRedisClient(cmd.Context())
	if err != nil {
		errMsg := fmt.Sprintf("Init Redis error: %s\n", err)
		global.Logger.Error(errMsg)
		panic(errMsg)
	}

	// init router
	r := router.InitRouter(gin.Default(), db, rdClient)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("server.httpPort")),
		Handler: r,
	}

	// start server in goroutine
	go func() {
		global.Logger.Infof("Start server %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Fatalf("Listen error: %s\n", err)
		}
	}()

	// graceful shutdown server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Logger.Info("Shutdown server ...")

	// waiting max 5 seconds, then force shutdown
	ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(viper.GetInt("server.shutdownTimeout"))*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		global.Logger.Fatalf("Server shutdown error: %s\n", err)
	}

	// catching ctx.Done()
	<-ctx.Done()
	global.Logger.Infof("timeout of %d seconds.\n", viper.GetInt("server.shutdownTimeout"))
	global.Logger.Info("Server exiting")
}

func init() {
	// Add apiserverCmd to rootCmd, start on terminal: go run main.go apiserver
	rootCmd.AddCommand(apiserverCmd)
}
