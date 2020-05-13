package workers

import (
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1"
	"github.com/astaxie/beego"
	"github.com/wonderivan/logger"
)

var Worker *machinery.Server

func init() {
	// 配置设置
	var cnf = &config.Config{
		Broker:             beego.AppConfig.String("Broker"),
	  DefaultQueue:       beego.AppConfig.String("DefaultQueue"),
	  ResultBackend:      beego.AppConfig.String("ResultBackend"),
		AMQP:               &config.AMQPConfig{
			Exchange:     "machinery_exchange",
			ExchangeType: "direct",
			BindingKey:   "machinery_task",
		},
	}

	// 初始化server
	server, err := machinery.NewServer(cnf)
	if err != nil {
		logger.Error("连接worker异常")
	}

	// 赋值给Worker，提供给其他package调用
	Worker = server

}
