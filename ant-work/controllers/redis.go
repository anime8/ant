package controllers

import (
	"github.com/wonderivan/logger"
	"encoding/json"
	"os/exec"
)


// 获取redis列表 函数
func GetRediss(redis Redis) []string {
	var rediss []string
	if redis.ClusterNodeRedisChecked01 {
		rediss = append(rediss,redis.ClusterNode01)
	}
	if redis.ClusterNodeRedisChecked02 {
		rediss = append(rediss,redis.ClusterNode02)
	}
	if redis.ClusterNodeRedisChecked03 {
		rediss = append(rediss,redis.ClusterNode03)
	}
	return rediss
}

// 获取sentinel列表 函数
func GetSentinels(redis Redis) []string {
  var sentinels []string
	if redis.ClusterNodeSentinelChecked01 {
		sentinels = append(sentinels,redis.ClusterNode01)
	}
	if redis.ClusterNodeSentinelChecked02 {
		sentinels = append(sentinels,redis.ClusterNode02)
	}
	if redis.ClusterNodeSentinelChecked03 {
		sentinels = append(sentinels,redis.ClusterNode03)
	}
	return sentinels
}


// redis部署interface
type RedisInstall interface {
	MakeInstall() bool
	UpdateConf() bool
	StartRedis() bool
}

// redis编译安装 interface函数
func (redis Redis) MakeInstall() bool {
		// 编译安装
		ips := []string{redis.ClusterNode01, redis.ClusterNode02, redis.ClusterNode03}
		for _, ip := range ips{
			ssh := "root@" + ip
			// 解压redis包
			command := "cd /opt/install/ && tar zxvf " + redis.RedisVersion + ".tar.gz"
			cmd := exec.Command("ssh", ssh, command)
			err := cmd.Run()
			if err != nil {
				logger.Error(ip + " tar包解压失败")
				return false
			}
			logger.Info(ip + " tar包解压成功")

			// 编译
			command = "cd /opt/install/" + redis.RedisVersion + " && make"
			cmd = exec.Command("ssh", ssh, command)
			err = cmd.Run()
			if err != nil {
				logger.Error(ip + " 编译失败")
				return false
			}
			logger.Info(ip + " 编译成功")

			// 创建目录
			dirs := []string{redis.RedisData, redis.SentinelData, redis.RedisLog, redis.RedisConf}
			for _, dir := range dirs{
				cmd := exec.Command("ssh", ssh, "mkdir -p " + dir)
				err := cmd.Run()
				if err != nil {
					logger.Error(ip + " " + dir + " 创建目录失败")
					return false
				}
				logger.Info(ip + " " + dir + " 创建目录成功")
			}

			// cp文件到/usr/bin下面
			bins := []string{}
			switch redis.RedisVersion {
			case "redis-2.8.10":
				bins = []string{"redis-benchmark", "redis-server", "redis-cli", "redis-sentinel", "redis-check-aof"}
			case "redis-3.2.0":
				bins = []string{"redis-benchmark", "redis-check-rdb", "redis-server", "redis-cli", "redis-sentinel", "redis-trib.rb", "redis-check-aof"}
			case "redis-5.0.5":
				bins = []string{"redis-benchmark", "redis-check-rdb", "redis-server", "redis-cli", "redis-sentinel", "redis-trib.rb", "redis-check-aof"}
			}
			for _, bin := range bins{
				command := "cp /opt/install/" + redis.RedisVersion + "/src/" + bin + " /usr/bin/"
				cmd := exec.Command("ssh", ssh, command)
				err := cmd.Run()
				if err != nil {
					logger.Error(ip + " " + bin + " cp文件失败")
					return false
				}
				logger.Info(ip + " " + bin + " cp文件成功")
			}
		}

		return true
}


// 更新节点redis和sentinel配置 interface函数
func (redis Redis) UpdateConf() bool {
	rediss := GetRediss(redis)
	sentinels := GetSentinels(redis)
	// 处理redis
	i := 1
	for _, redisip := range rediss{
		ssh := "root@" + redisip
		// 移动配置文件到/etc/redis下面
		command := "mv /opt/install/redis-6380.conf " + redis.RedisConf
		cmd := exec.Command("ssh", ssh, command)
		err := cmd.Run()
		if err != nil {
			logger.Error(redisip + " mv配置文件redis-6380.conf失败")
			return false
		}
		logger.Info(redisip + " mv配置文件redis-6380.conf成功")

		// 更新redis配置文件目录
		command = "sed -i 's#/logs/redis#" + redis.RedisLog + "#g' " + redis.RedisConf + "/redis-6380.conf"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(redisip + " 更新日志目录配置失败")
			return false
		}
		logger.Info(redisip + " 更新日志目录配置成功")
		command = "sed -i 's#/opt/data/redis#" + redis.RedisData + "#g' " + redis.RedisConf + "/redis-6380.conf"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(redisip + " 更新数据目录配置失败")
			return false
		}
		logger.Info(redisip + " 更新数据目录配置成功")


		// 如果开启了验证，则添加密码
		if redis.RedisAuthentication {
			command = "sed -i 's#123456#" + redis.RedisPassword + "#g' " + redis.RedisConf + "/redis-6380.conf"
		}else{
			command = "sed -i 's#requirepass 123456##g' " + redis.RedisConf + "/redis-6380.conf && sed -i 's#masterauth 123456##g' " + redis.RedisConf + "/redis-6380.conf"
		}
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(redisip + " 配置redis验证失败，如果没有开启验证，则删除相关配置。")
			return false
		}
		logger.Info(redisip + " 配置redis验证成功，如果没有开启验证，则删除相关配置。")

		// 判断是否有protected-mode
		switch redis.RedisVersion {
		case "redis-2.8.10":
			command = "sed -i 's#protected-mode no##g' " + redis.RedisConf + "/redis-6380.conf"
			cmd = exec.Command("ssh", ssh, command)
			err = cmd.Run()
			if err != nil {
				logger.Error(redisip + " 当前版本为：" + redis.RedisVersion + ",删除protected-mode配置失败。")
				return false
			}
			logger.Info(redisip + " 当前版本为：" + redis.RedisVersion + ",删除protected-mode配置成功。")
		}


		// 判断节点是否为主节点
		if i == 1 {
			command = "sed -i 's#slaveof 127.0.0.1 6380##g' " + redis.RedisConf + "/redis-6380.conf"
		}else{
			command = "sed -i 's#127.0.0.1#" + rediss[0] + "#g' " + redis.RedisConf + "/redis-6380.conf"
		}
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			if i == 1 {
				logger.Error(redisip + " 设置主节点失败。")
			}else{
				logger.Error(redisip + " 设置从节点失败。")
			}
			return false
		}
		if i == 1 {
			logger.Info(redisip + " 设置主节点成功。")
		}else{
			logger.Info(redisip + " 设置从节点成功。")
		}

		i++
	}

	for _, sentinel := range sentinels{
		ssh := "root@" + sentinel
		command := "mv /opt/install/sentinel.conf " + redis.RedisConf
		cmd := exec.Command("ssh", ssh, command)
		err := cmd.Run()
		if err != nil {
			logger.Error(sentinel + " copy配置sentinel.conf失败")
			return false
		}
		logger.Info(sentinel + " copy配置sentinel.conf成功")
		// 如果开启了验证，则添加密码
		if redis.RedisAuthentication {
			command = "sed -i 's#123456#" + redis.RedisPassword + "#g' " + redis.RedisConf + "/sentinel.conf"
		}else{
			command = "sed -i 's#sentinel auth-pass mymaster 123456##g' " + redis.RedisConf + "/sentinel.conf"
		}
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(sentinel + " 配置验证失败，如果配置开启，则删除配置验证失败。")
			return false
		}
		logger.Info(sentinel + " 配置验证成功，如果配置开启，则删除配置验证成功。")

		// 更新sentinel配置文件目录
		command = "sed -i 's#/logs/redis#" + redis.RedisLog + "#g' " + redis.RedisConf + "/sentinel.conf"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(sentinel + " 更新日志目录配置失败")
			return false
		}
		logger.Info(sentinel + " 更新日志目录配置成功")
		command = "sed -i 's#/opt/data/redis/sentinel#" + redis.SentinelData + "#g' " + redis.RedisConf + "/sentinel.conf"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(sentinel + " 更新数据目录配置失败")
			return false
		}
		logger.Info(sentinel + " 更新数据目录配置成功")

		// 判断是否有protected-mode
		switch redis.RedisVersion {
		case "redis-2.8.10":
			command = "sed -i 's#protected-mode no##g' " + redis.RedisConf + "/sentinel.conf"
			cmd = exec.Command("ssh", ssh, command)
			err = cmd.Run()
			if err != nil {
				logger.Error(sentinel + " 当前版本为：" + redis.RedisVersion + ",删除protected-mode配置失败。")
				return false
			}
			logger.Info(sentinel + " 当前版本为：" + redis.RedisVersion + ",删除protected-mode配置成功。")
		}


		// 更改主节点地址
		command = "sed -i 's#127.0.0.1#" + rediss[0] + "#g' " + redis.RedisConf + "/sentinel.conf"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(sentinel + " 更改主节点地址失败。")
			return false
		}
		logger.Info(sentinel + " 更改主节点地址成功。")

		// 更改集群名字
		command = "sed -i 's#mymaster#" + redis.SentinelName + "#g' " + redis.RedisConf + "/sentinel.conf"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(sentinel + " 更改集群名失败。")
			return false
		}
		logger.Info(sentinel + " 更改集群名成功。")
	}

	return true
}


// 启动redis interface函数
func (redis Redis) StartRedis() bool {
	rediss := GetRediss(redis)
	sentinels := GetSentinels(redis)
	command := "redis-server " + redis.RedisConf + "/redis-6380.conf"
	for _, redisip := range rediss{
		ssh := "root@" + redisip
		cmd := exec.Command("ssh", ssh, command)
		err := cmd.Run()
		if err != nil {
			logger.Error(redisip + " redis节点启动失败。")
			return false
		}
		logger.Info(redisip + " redis节点启动成功。")
	}

	command = "redis-sentinel " + redis.RedisConf + "/sentinel.conf"
	for _, sentinel := range sentinels{
		ssh := "root@" + sentinel
		cmd := exec.Command("ssh", ssh, command)
		err := cmd.Run()
		if err != nil {
			logger.Error(sentinel + " sentinel节点启动失败。")
			return false
		}
		logger.Info(sentinel + " sentinel节点启动成功。")
	}

	return true
}


//  部署redis的任务
func DelopyRedis(r string) (string, error) {
	var redis Redis
	var res ResResult
	res.Status = "Success"
	res.Data = "部署成功"
	// 将字符串转换成struct
	json.Unmarshal([]byte(r), &redis)

	// 编译安装redis
	logger.Info("开始编译安装")
	mr := redis.MakeInstall()
	if mr == false {
		res.Status = "Failure"
		res.Data = "编译安装失败"
	}

	// 更新配置
	logger.Info("开始更新配置")
	ur := redis.UpdateConf()
	if ur == false {
		res.Status = "Failure"
		res.Data = "更新配置文件失败"
		data,_ := json.Marshal(res)
		result := string(data)
		return result, nil
	}

	// 启动redis,先启动主，再启动其他的
	logger.Info("开始启动节点")
	ss := redis.StartRedis()
	if ss == false {
		res.Status = "Failure"
		res.Data = "启动失败"
	}

	// 返回结果到redis
	data,_ := json.Marshal(res)
	result := string(data)
  return result, nil
}
