package workers

import (
  "github.com/wonderivan/logger"
  "github.com/RichardKnop/machinery/v1/config"
  "github.com/RichardKnop/machinery/v1"
  "github.com/astaxie/beego"
  "ant-work/controllers"
)

// 设置machinery配置
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

func init() {
    // 日志配置
    logger.SetLogger(beego.AppConfig.String("LogJson"))

    // machinery实例化
    server, err := machinery.NewServer(cnf)
    if err != nil {
      logger.Error("machinery实例化失败")
    }
    // 注册task到server
    server.RegisterTask("DelopyRedis", controllers.DelopyRedis)
    server.RegisterTask("DelopyZookeeper", controllers.DelopyZookeeper)
    server.RegisterTask("DelopyKafka", controllers.DelopyKafka)

    // 新建worker并启动
    worker := server.NewWorker("test", 10)
    worker.Launch()

}
