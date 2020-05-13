package controllers

import (
	"ant-api/models"
	"ant-api/deploy"
	"ant-api/cache"
	"ant-api/workers"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/RichardKnop/machinery/v1/tasks"
	"net"
	"github.com/wonderivan/logger"
)


// Operations about Users
type KafkaController struct {
	beego.Controller
}

// 安装kafka
func (u *KafkaController) Install() {
	var kafka models.Kafka
	var res ResResult

	// 将body中的数据存入kafka中
	json.Unmarshal(u.Ctx.Input.RequestBody, &kafka)

	ips := []string{kafka.ClusterNode01, kafka.ClusterNode02, kafka.ClusterNode03}
	for _, ip := range ips{
		// 判断ip地址格式是否正确
    address := net.ParseIP(ip)
		if address == nil {
			logger.Error(ip + " ip地址格式不正确")
			res.Status = "Failure"
			res.Data = "ip地址格式不正确"
			u.Ctx.WriteString(res.ConvertToString())
			return
    }

		// 判断服务是否安装了kafka
		KafkaInstallResult := deploy.KafkaInstallCheck(ip)
		if KafkaInstallResult == false {
			logger.Error(ip + " 服务器上已经安装了kafka")
			res.Status = "Failure"
			res.Data = "目标服务器上已经安装了kafka"
			u.Ctx.WriteString(res.ConvertToString())
			return
		}
  }

	// 判断3个ip是否相同
	if kafka.ClusterNode01==kafka.ClusterNode02 || kafka.ClusterNode01==kafka.ClusterNode03 || kafka.ClusterNode02==kafka.ClusterNode03 {
		logger.Error("ip地址不能相同")
		res.Status = "Failure"
		res.Data = "ip地址不能相同"
		u.Ctx.WriteString(res.ConvertToString())
		return
	}

	// 上传包和配置
	i := 1
	for _, ip := range ips{
		// 声明初始上传结果
		KafkaUploadResult := true
		switch i {
		case 1:
			KafkaUploadResult = deploy.KafkaPackageUpload(ip, kafka.KafkaVersion)
		case 2:
			KafkaUploadResult = deploy.KafkaPackageUpload(ip, kafka.KafkaVersion)
		case 3:
			KafkaUploadResult = deploy.KafkaPackageUpload(ip, kafka.KafkaVersion)
		}
    if KafkaUploadResult == false {
			logger.Error("上传文件失败")
			res.Status = "Failure"
			res.Data = "上传文件失败"
			u.Ctx.WriteString(res.ConvertToString())
			return
    }
		// 计数+1
		i++
	}

	// 提交安装任务
	kafkajson, _ := json.Marshal(kafka)
	kafkastr := string(kafkajson)
	signature := &tasks.Signature{
	  Name: "DelopyKafka",
	  Args: []tasks.Arg{
	    {
	      Type:  "string",
	      Value: kafkastr,
	    },
	  },
	}
	asyncResult, err := workers.Worker.SendTask(signature)
	if err != nil {
		logger.Error(err)
		res.Status = "Failure"
		res.Data = "提交安装任务失败"
		u.Ctx.WriteString(res.ConvertToString())
		return
	}

	// 获取任务状态
	taskState := asyncResult.GetState()
	kafka.TaskId = taskState.TaskUUID
	kafka.DeployStatus = taskState.State
	// 将数据插入数据库
	models.KafkaInsert(kafka)

	logger.Info(taskState.TaskUUID + "提交任务成功")
	res.Status = "Success"
	res.Data = "提交任务成功,请点击详情查看部署进度。"
	u.Ctx.WriteString(res.ConvertToString())
	return
}


// 获取所有kafka信息
func (u *KafkaController) GetAll() {
	var res ResResult

	// 查询所有kafka信息
	kafkaList := models.KafkaGetAll()
	kafkaJson,_ := json.Marshal(kafkaList)
	//设置返回值
	res.Status = "Success"
	res.Data = string(kafkaJson)
	u.Ctx.WriteString(res.ConvertToString())
	return
}


// 获取任务执行状态，如果发生改变，则更新表
func (u *KafkaController) GetOne() {
	var kafka models.Kafka
	var res ResResult
	var taskState TaskState

	// 获取body数据
	json.Unmarshal(u.Ctx.Input.RequestBody, &kafka)
	// 通过id查询完整的记录
	kafka = models.KafkaGetOne(kafka)


	// 取出kafka中的安装结果,并将结果存入taskState
	taskResult, haskey := cache.Get(kafka.TaskId)
	// 如果取到了结果
	if haskey {
		json.Unmarshal([]byte(taskResult), &taskState)

		// 如果任务的执行结果为SUCCESS，则将部署结果设置为安装结果
		if taskState.State == "SUCCESS" {
			installRes := taskState.GetDeployResult()
			kafka.DeployResult = installRes.Data
			// 将安装结果的返回值变的和任务执行结果一致
			switch installRes.Status {
			case "Success":
				kafka.DeployStatus = "SUCCESS"
			case "Failure":
				kafka.DeployStatus = "FAILURE"
			}
		}else{   // 否则将部署结果设置为任务执行结果
			kafka.DeployStatus = taskState.State
		}
		// 更新部署进度
		kafka = models.KafkaUpdate(kafka)
		kafkaJson,_ := json.Marshal(kafka)

		//设置返回值
		res.Status = "Success"
		res.Data = string(kafkaJson)
		u.Ctx.WriteString(res.ConvertToString())
		return
	}else{ // 没有取到结果
		res.Status = "Failure"
		res.Data = "获取任务执行结果异常"
		u.Ctx.WriteString(res.ConvertToString())
		return
	}
}
