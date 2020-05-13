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
type ZookeeperController struct {
	beego.Controller
}

// 安装zookeeper
func (u *ZookeeperController) Install() {
	var zookeeper models.Zookeeper
	var res ResResult

	// 将body中的数据存入zookeeper中
	json.Unmarshal(u.Ctx.Input.RequestBody, &zookeeper)

	ips := []string{zookeeper.ClusterNode01, zookeeper.ClusterNode02, zookeeper.ClusterNode03}
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

		// 判断服务是否安装了zookeeper
		ZookeeperInstallResult := deploy.ZookeeperInstallCheck(ip)
		if ZookeeperInstallResult == false {
			logger.Error(ip + " 服务器上已经安装了zookeeper")
			res.Status = "Failure"
			res.Data = "目标服务器上已经安装了zookeeper"
			u.Ctx.WriteString(res.ConvertToString())
			return
		}
  }

	// 判断3个ip是否相同
	if zookeeper.ClusterNode01==zookeeper.ClusterNode02 || zookeeper.ClusterNode01==zookeeper.ClusterNode03 || zookeeper.ClusterNode02==zookeeper.ClusterNode03 {
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
		ZookeeperUploadResult := true
		switch i {
		case 1:
			ZookeeperUploadResult = deploy.ZookeeperPackageUpload(ip, zookeeper.ZookeeperVersion)
		case 2:
			ZookeeperUploadResult = deploy.ZookeeperPackageUpload(ip, zookeeper.ZookeeperVersion)
		case 3:
			ZookeeperUploadResult = deploy.ZookeeperPackageUpload(ip, zookeeper.ZookeeperVersion)
		}
    if ZookeeperUploadResult == false {
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
	zookeeperjson, _ := json.Marshal(zookeeper)
	zookeeperstr := string(zookeeperjson)
	signature := &tasks.Signature{
	  Name: "DelopyZookeeper",
	  Args: []tasks.Arg{
	    {
	      Type:  "string",
	      Value: zookeeperstr,
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
	zookeeper.TaskId = taskState.TaskUUID
	zookeeper.DeployStatus = taskState.State
	// 将数据插入数据库
	models.ZookeeperInsert(zookeeper)

	logger.Info(taskState.TaskUUID + "提交任务成功")
	res.Status = "Success"
	res.Data = "提交任务成功,请点击详情查看部署进度。"
	u.Ctx.WriteString(res.ConvertToString())
	return
}


// 获取所有zookeeper信息
func (u *ZookeeperController) GetAll() {
	var res ResResult

	// 查询所有zookeeper信息
	zookeeperList := models.ZookeeperGetAll()
	zookeeperJson,_ := json.Marshal(zookeeperList)
	//设置返回值
	res.Status = "Success"
	res.Data = string(zookeeperJson)
	u.Ctx.WriteString(res.ConvertToString())
	return
}


// 获取任务执行状态，如果发生改变，则更新表
func (u *ZookeeperController) GetOne() {
	var zookeeper models.Zookeeper
	var res ResResult
	var taskState TaskState

	// 获取body数据
	json.Unmarshal(u.Ctx.Input.RequestBody, &zookeeper)
	// 通过id查询完整的记录
	zookeeper = models.ZookeeperGetOne(zookeeper)


	// 取出zookeeper中的安装结果,并将结果存入taskState
	taskResult, haskey := cache.Get(zookeeper.TaskId)
	// 如果取到了结果
	if haskey {
		json.Unmarshal([]byte(taskResult), &taskState)

		// 如果任务的执行结果为SUCCESS，则将部署结果设置为安装结果
		if taskState.State == "SUCCESS" {
			installRes := taskState.GetDeployResult()
			zookeeper.DeployResult = installRes.Data
			// 将安装结果的返回值变的和任务执行结果一致
			switch installRes.Status {
			case "Success":
				zookeeper.DeployStatus = "SUCCESS"
			case "Failure":
				zookeeper.DeployStatus = "FAILURE"
			}
		}else{   // 否则将部署结果设置为任务执行结果
			zookeeper.DeployStatus = taskState.State
		}
		// 更新部署进度
		zookeeper = models.ZookeeperUpdate(zookeeper)
		zookeeperJson,_ := json.Marshal(zookeeper)

		//设置返回值
		res.Status = "Success"
		res.Data = string(zookeeperJson)
		u.Ctx.WriteString(res.ConvertToString())
		return
	}else{ // 没有取到结果
		res.Status = "Failure"
		res.Data = "获取任务执行结果异常"
		u.Ctx.WriteString(res.ConvertToString())
		return
	}
}
