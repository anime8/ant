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
type RedisController struct {
	beego.Controller
}

// 安装redis
func (u *RedisController) Install() {
	var redis models.Redis
	var res ResResult

	// 将body中的数据存入redis中
	json.Unmarshal(u.Ctx.Input.RequestBody, &redis)


	ips := []string{redis.ClusterNode01, redis.ClusterNode02, redis.ClusterNode03}
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

		// 判断服务是否安装了redis
		RedisInstallResult := deploy.RedisInstallCheck(ip)
		if RedisInstallResult == false {
			logger.Error(ip + " 服务器上已经安装了redis")
			res.Status = "Failure"
			res.Data = "目标服务器上已经安装了redis"
			u.Ctx.WriteString(res.ConvertToString())
			return
		}
  }

	// 判断3个ip是否相同
	if redis.ClusterNode01==redis.ClusterNode02 || redis.ClusterNode01==redis.ClusterNode03 || redis.ClusterNode02==redis.ClusterNode03 {
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
		RedisUploadResult := true
		switch i {
		case 1:
			RedisUploadResult = deploy.RedisPackageUpload(ip, redis.RedisVersion, redis.ClusterNodeRedisChecked01, redis.ClusterNodeSentinelChecked01)
		case 2:
			RedisUploadResult = deploy.RedisPackageUpload(ip, redis.RedisVersion, redis.ClusterNodeRedisChecked02, redis.ClusterNodeSentinelChecked02)
		case 3:
			RedisUploadResult = deploy.RedisPackageUpload(ip, redis.RedisVersion, redis.ClusterNodeRedisChecked03, redis.ClusterNodeSentinelChecked03)
		}
    if RedisUploadResult == false {
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
	redisjson, _ := json.Marshal(redis)
	redisstr := string(redisjson)
	signature := &tasks.Signature{
	  Name: "DelopyRedis",
	  Args: []tasks.Arg{
	    {
	      Type:  "string",
	      Value: redisstr,
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
	redis.TaskId = taskState.TaskUUID
	redis.DeployStatus = taskState.State
	// 将数据插入数据库
	models.RedisInsert(redis)

	logger.Info(taskState.TaskUUID + "提交任务成功")
	res.Status = "Success"
	res.Data = "提交任务成功,请点击详情查看部署进度。"
	u.Ctx.WriteString(res.ConvertToString())
	return
}


// 获取所有redis信息
func (u *RedisController) GetAll() {
	var res ResResult

	// 查询所有redis信息
	redisList := models.RedisGetAll()
	redisJson,_ := json.Marshal(redisList)
	//设置返回值
	res.Status = "Success"
	res.Data = string(redisJson)
	u.Ctx.WriteString(res.ConvertToString())
	return
}


// 获取任务执行状态，如果发生改变，则更新表
func (u *RedisController) GetOne() {
	var redis models.Redis
	var res ResResult
	var taskState TaskState

	// 获取body数据
	json.Unmarshal(u.Ctx.Input.RequestBody, &redis)
	// 通过id查询完整的记录
	redis = models.RedisGetOne(redis)


	// 取出redis中的安装结果,并将结果存入taskState
	taskResult, haskey := cache.Get(redis.TaskId)
	// 如果取到了结果
	if haskey {
		json.Unmarshal([]byte(taskResult), &taskState)

		// 如果任务的执行结果为SUCCESS，则将部署结果设置为安装结果
		if taskState.State == "SUCCESS" {
			installRes := taskState.GetDeployResult()
			redis.DeployResult = installRes.Data
			// 将安装结果的返回值变的和任务执行结果一致
			switch installRes.Status {
			case "Success":
				redis.DeployStatus = "SUCCESS"
			case "Failure":
				redis.DeployStatus = "FAILURE"
			}
		}else{   // 否则将部署结果设置为任务执行结果
			redis.DeployStatus = taskState.State
		}
		// 更新部署进度
		redis = models.RedisUpdate(redis)
		redisJson,_ := json.Marshal(redis)

		//设置返回值
		res.Status = "Success"
		res.Data = string(redisJson)
		u.Ctx.WriteString(res.ConvertToString())
		return
	}else{ // 没有取到结果
		res.Status = "Failure"
		res.Data = "获取任务执行结果异常"
		u.Ctx.WriteString(res.ConvertToString())
		return
	}
}
