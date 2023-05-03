package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web_app/controller"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/routes"
	"web_app/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	//1. loads the config
	if err := settings.Init(); err != nil {
		fmt.Printf("Init settings failed, err:%v\n", err)
		return
	}

	//2. initialize the logger
	if err := logger.Init(); err != nil {
		fmt.Printf("Init loggers failed, err:%v\n", err)
	}
	defer zap.L().Sync()
	zap.L().Debug("logger inits successfully...")

	//3. initialize the mysql
	if err := mysql.InitDB(); err != nil {
		fmt.Printf("Init mysql failed, err:%v\n", err)
	}
	defer mysql.Close()

	//4. initialize the redis
	if err := redis.Init(); err != nil {
		fmt.Printf("Init redis failed, err:%v\n", err)
	}
	defer redis.Close()

	//5. register the router (Gin framework)
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init transalator failed, err: %v\n", err)
		return
	}
	r := routes.Setup()

	//6. elegantly quits
	srv := &http.Server{
		Addr:                         fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler:                      r,
		DisableGeneralOptionsHandler: false,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutdown Server...")
	// Create a context with timeout of 5s
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Info("Server Shutdown", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
