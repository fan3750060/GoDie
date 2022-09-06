package main

import (
	"goframe/app/db"
	"goframe/app/logger"
	"goframe/config"
	"goframe/route"
	"net/http"
	"time"
)

func main() {
	logger.Logger.Info("服务正在启动...")

	//初始化数据库连接
	var mysql db.Mysql = new(db.MysqlClass)
	mysql.DbStar()

	//加载路由
	var routes route.Routes = new(route.RoutesClass)
	router := routes.LoadRoute()

	//启用http服务器
	http.Handle("/", router)
	svr := http.Server{
		Addr:         "127.0.0.1:" + config.HttpConfig.Port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
	}

	err := svr.ListenAndServe()

	go func() {
		if err != nil {
			if err == http.ErrServerClosed {
				logger.Logger.Print("服务器已根据要求关闭！！")
			} else {
				logger.Logger.Fatal("服务器意外关闭！！")
			}
		}
	}()

	logger.Logger.Println("服务已关闭！")
}
