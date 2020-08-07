package internal

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"payment-system/internal/controller"
	"payment-system/internal/db/clickhouse"
	"payment-system/internal/db/mysql"
	"payment-system/internal/restapi"
)

func StartSystem() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	mysqlSettings := viper.GetStringMapString("mysql")
	mysqlDB, err := mysql.NewDB(mysqlSettings["host"], mysqlSettings["port"], mysqlSettings["user"], mysqlSettings["pass"], mysqlSettings["db"], logger)
	if err != nil {
		logger.Panic("Failed to init Mysql DB connection", zap.Error(err))
	}

	clSettings := viper.GetStringMapString("clickhouse")
	clickHouseDB, err := clickhouse.NewDB(clSettings["host"], clSettings["port"], clSettings["db"], logger)
	if err != nil {
		logger.Panic("Failed to init Clickhouse DB connection", zap.Error(err))
	}

	cntrlr := controller.NewController(logger, mysqlDB, clickHouseDB)

	logger.Info("Starting web server")
	router := restapi.NewRouter(logger, cntrlr)

	quit := cntrlr.StartTimerTask()

	err = http.ListenAndServe(":8080", router)
	logger.With(zap.Error(err)).Error("web server stopped")
	close(quit)
}
