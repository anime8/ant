package controllers

import (
	"github.com/wonderivan/logger"
	"encoding/json"
	"os/exec"
	"strconv"
)

// redis部署interface
type ZookeeperInstall interface {
	MakeInstall() bool
	UpdateConf() bool
	StartZookeeper() bool
}

// zookeeper编译安装 interface函数
func (zookeeper Zookeeper) MakeInstall() bool {
		// 编译安装
		ips := []string{zookeeper.ClusterNode01, zookeeper.ClusterNode02, zookeeper.ClusterNode03}
		for _, ip := range ips{
			ssh := "root@" + ip
			// 解压zookeeper包
			command := "cd /opt/install/ && tar zxvf " + zookeeper.ZookeeperVersion + ".tar.gz"
			cmd := exec.Command("ssh", ssh, command)
			err := cmd.Run()
			if err != nil {
				logger.Error(ip + " tar包解压失败")
				return false
			}
			logger.Info(ip + " tar包解压成功")

			// 创建目录
			dirs := []string{zookeeper.DeployPath, zookeeper.ZookeeperData, zookeeper.ZookeeperLog}
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
			command = "mv /opt/install/" + zookeeper.ZookeeperVersion + " " + zookeeper.DeployPath
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
func (zookeeper Zookeeper) UpdateConf() bool {
	ips := []string{zookeeper.ClusterNode01, zookeeper.ClusterNode02, zookeeper.ClusterNode03}
	i := 1
	for _, ip := range ips{
		ssh := "root@" + ip
		// 移动配置文件到conf目录下面
		command := "mv /opt/install/zoo.cfg " + zookeeper.DeployPath + "/" + zookeeper.ZookeeperVersion + "/conf/"
		cmd := exec.Command("ssh", ssh, command)
		err := cmd.Run()
		if err != nil {
			logger.Error(ip + " mv配置文件zoo.cfg失败")
			return false
		}
		logger.Info(ip + " mv配置文件zoo.cfg成功")

		// 更新zookeeper配置
		command = "sed -i 's#/opt/data/zookeeper/data#" + zookeeper.ZookeeperData + "#g' " + zookeeper.DeployPath + "/" + zookeeper.ZookeeperVersion + "/conf/zoo.cfg"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(ip + " 更新数据目录配置失败")
			return false
		}
		logger.Info(ip + " 更新数据目录配置成功")

		command = "sed -i 's#/opt/data/zookeeper/log#" + zookeeper.ZookeeperLog + "#g' " + zookeeper.DeployPath + "/" + zookeeper.ZookeeperVersion + "/conf/zoo.cfg"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(ip + " 更新日志目录配置失败")
			return false
		}
		logger.Info(ip + " 更新日志目录配置成功")

		// 更改配置里面的ip
		n := 1
		for _, cfg_ip := range ips{
			command = "sed -i 's#127.0.0." + strconv.Itoa(n) + "#" + cfg_ip + "#g' " + zookeeper.DeployPath + "/" + zookeeper.ZookeeperVersion + "/conf/zoo.cfg"
			cmd = exec.Command("ssh", ssh, command)
			err = cmd.Run()
			if err != nil {
				logger.Error(cfg_ip + " 更改配置里面的ip失败")
				return false
			}
			logger.Info(cfg_ip + " 更改配置里面的ip成功")

			// 计数
			n++
		}

		// 创建myid文件
		command = "echo " + strconv.Itoa(i) + " > " + zookeeper.ZookeeperData + "/myid"
		cmd = exec.Command("ssh", ssh, command)
		err = cmd.Run()
		if err != nil {
			logger.Error(ip + " 创建myid文件失败")
			return false
		}
		logger.Info(ip + " 创建myid文件成功")

		// 计数
		i++
	}

	return true
}


// 启动zookeeper interface函数
func (zookeeper Zookeeper) StartZookeeper() bool {
	ips := []string{zookeeper.ClusterNode01, zookeeper.ClusterNode02, zookeeper.ClusterNode03}
	for _, ip := range ips{
		ssh := "root@" + ip
		command := zookeeper.DeployPath + "/" + zookeeper.ZookeeperVersion + "/bin/zkServer.sh start"
		cmd := exec.Command("ssh", ssh, command)
		err := cmd.Run()
		if err != nil {
			logger.Error(ip + " zookeeper节点启动失败。")
			return false
		}
		logger.Info(ip + " zookeeper节点启动成功。")
	}

	return true
}


//  部署zookeeper的任务
func DelopyZookeeper(z string) (string, error) {
	var zookeeper Zookeeper
	var res ResResult
	res.Status = "Success"
	res.Data = "部署成功"
	// 将字符串转换成struct
	json.Unmarshal([]byte(z), &zookeeper)

	// 编译安装redis
	logger.Info("开始解压")
	mr := zookeeper.MakeInstall()
	if mr == false {
		res.Status = "Failure"
		res.Data = "解压失败"
	}

	// 更新配置
	logger.Info("开始更新配置")
	ur := zookeeper.UpdateConf()
	if ur == false {
		res.Status = "Failure"
		res.Data = "更新配置文件失败"
		data,_ := json.Marshal(res)
		result := string(data)
		return result, nil
	}

	// 启动redis,先启动主，再启动其他的
	logger.Info("开始启动节点")
	ss := zookeeper.StartZookeeper()
	if ss == false {
		res.Status = "Failure"
		res.Data = "启动失败"
	}

	// 返回结果到redis
	data,_ := json.Marshal(res)
	result := string(data)
  return result, nil
}
