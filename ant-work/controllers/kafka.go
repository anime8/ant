package controllers

import (
	"github.com/wonderivan/logger"
	"encoding/json"
	"os/exec"
	"strconv"
)

// kafka部署interface
type KafkaInstall interface {
	MakeInstall() bool
	UpdateConf() bool
	StartKafka() bool
}

// kafka编译安装 interface函数
func (kafka Kafka) MakeInstall() bool {
		// 编译安装
		ips := []string{kafka.ClusterNode01, kafka.ClusterNode02, kafka.ClusterNode03}
		for _, ip := range ips{
			ssh := "root@" + ip
			// 解压kafka包
			command := "cd /opt/install/ && tar xvf " + kafka.KafkaVersion + ".tgz"
			cmd := exec.Command("ssh", ssh, command)
			err := cmd.Run()
			if err != nil {
				logger.Error(ip + " tar包解压失败")
				return false
			}
			logger.Info(ip + " tar包解压成功")

			// 创建目录
			dirs := []string{kafka.KafkaPath, kafka.KafkaData}
			for _, dir := range dirs{
				cmd := exec.Command("ssh", ssh, "mkdir -p " + dir)
				err := cmd.Run()
				if err != nil {
					logger.Error(ip + " " + dir + " 创建目录失败")
					return false
				}
				logger.Info(ip + " " + dir + " 创建目录成功")
			}

			// 将软件mv到安装目录下
			command = "mv /opt/install/" + kafka.KafkaVersion + " " + kafka.KafkaPath
			cmd = exec.Command("ssh", ssh, command)
			err = cmd.Run()
			if err != nil {
				logger.Error(ip + " mv文件失败")
				return false
			}
			logger.Info(ip + " mv文件成功")
		}

		return true
}


// 更新配置
func (kafka Kafka) UpdateConf() bool {
	ips := []string{kafka.ClusterNode01, kafka.ClusterNode02, kafka.ClusterNode03}
	i := 1
	for _, ip := range ips{
		ssh := "root@" + ip

		// 备份老配置
		command := "mv " + kafka.KafkaPath + "/" + kafka.KafkaVersion + "/config/server.properties  " + kafka.KafkaPath + "/" + kafka.KafkaVersion + "/config/server.properties.`/usr/bin/date +%Y%m%d`"
		cmd := exec.Command("ssh", ssh, command)
		err := cmd.Run()
		if err != nil {
			logger.Error(ip + " 备份配置文件server.properties失败")
			return false
		}
		logger.Info(ip + " 备份配置文件server.properties成功")

		// 移动配置文件到conf目录下面
		command = "mv /opt/install/server.properties " + kafka.KafkaPath + "/" + kafka.KafkaVersion + "/config/"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(ip + " mv配置文件server.properties失败")
			return false
		}
		logger.Info(ip + " mv配置文件server.properties成功")

		// 更新kafka配置
		command = "sed -i 's#/opt/data/kafka#" + kafka.KafkaData + "#g' " + kafka.KafkaPath + "/" + kafka.KafkaVersion + "/config/server.properties"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(ip + " 更新数据目录配置失败")
			return false
		}
		logger.Info(ip + " 更新数据目录配置成功")

		// 更新kafka连接zk配置
		command = "sed -i 's#zk1:2181,zk2:2181,zk3:2181#" + kafka.KafkaZookeeper + "#g' " + kafka.KafkaPath + "/" + kafka.KafkaVersion + "/config/server.properties"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(ip + " 更新日志目录配置失败")
			return false
		}
		logger.Info(ip + " 更新日志目录配置成功")

		// 更新brokerid
		command = "sed -i 's#broker.id=0#" + "broker.id=" + strconv.Itoa(i) + "#g' " + kafka.KafkaPath + "/" + kafka.KafkaVersion + "/config/server.properties"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(ip + " 更新brokerid失败")
			return false
		}
		logger.Info(ip + " 更新brokerid成功")

		// 计数
		i++
	}

	return true
}


// 启动kafka interface函数
func (kafka Kafka) StartKafka() bool {
	ips := []string{kafka.ClusterNode01, kafka.ClusterNode02, kafka.ClusterNode03}
	for _, ip := range ips{
		ssh := "root@" + ip
		command := kafka.KafkaPath + "/" + kafka.KafkaVersion + "/bin/kafka-server-start.sh -daemon " + kafka.KafkaPath + "/" + kafka.KafkaVersion + "/config/server.properties"
		cmd := exec.Command("ssh", ssh, command)
		err := cmd.Run()
		if err != nil {
			logger.Error(ip + " kafka节点启动失败。")
			return false
		}
		logger.Info(ip + " kafka节点启动成功。")
	}

	return true
}


//  部署kafka的任务
func DelopyKafka(k string) (string, error) {
	var kafka Kafka
	var res ResResult
	res.Status = "Success"
	res.Data = "部署成功"
	// 将字符串转换成struct
	json.Unmarshal([]byte(k), &kafka)

	// 编译安装redis
	logger.Info("开始解压")
	mr := kafka.MakeInstall()
	if mr == false {
		res.Status = "Failure"
		res.Data = "解压失败"
	}

	// 更新配置
	logger.Info("开始更新配置")
	ur := kafka.UpdateConf()
	if ur == false {
		res.Status = "Failure"
		res.Data = "更新配置文件失败"
		data,_ := json.Marshal(res)
		result := string(data)
		return result, nil
	}

	// 启动redis,先启动主，再启动其他的
	logger.Info("开始启动节点")
	ss := kafka.StartKafka()
	if ss == false {
		res.Status = "Failure"
		res.Data = "启动失败"
	}

	// 返回结果到redis
	data,_ := json.Marshal(res)
	result := string(data)
  return result, nil
}
