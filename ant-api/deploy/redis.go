package deploy


import (
	"os/exec"
  // "fmt"
  "strings"
	"path/filepath"
	"os"
	"github.com/wonderivan/logger"
)

// 检查服务器是否安装过redis
func RedisInstallCheck(ip string) bool {
  res := true
  ssh := "root@" + ip
  // 查看服务是否存在这些文件或目录
  files := []string{"/usr/bin/redis-server", "/usr/bin/redis-sentinel"}
  for _, f := range files{
    cmd := exec.Command("ssh", ssh, "ls " + f)
    err := cmd.Run()
    if err == nil {
			logger.Error(err)
      res = false
      break
    }
  }
  // 查看是否有redis进程
  cmd := exec.Command("ssh", ssh, "ps aux |grep redis |grep -v grep | wc -l")
  bytes,err := cmd.Output()
  if err != nil {
		logger.Error(err)
  }
  // 去换行
  resp := strings.Replace(string(bytes), "\n", "", -1)
  if resp != "0" {
    res = false
  }
  return res
}


// 上传配置文件和包到服务器
func RedisPackageUpload(ip string, version string, redis bool, sentinel bool) bool {
	res := true
	ssh := "root@" + ip
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Error(err)
	}
	packagedir := dir + "/packages/redis/"
	tarname := version + ".tar.gz"
	destinationpath := "/opt/install/"
	// 查看是否有/opt/install目录，没有则创建一个
	cmd := exec.Command("ssh", ssh, "ls " + destinationpath)
	err = cmd.Run()
	if err != nil {
		cmd := exec.Command("ssh", ssh, "mkdir " + destinationpath)
		err = cmd.Run()
		// 创建目录失败
		if err != nil {
			logger.Error(err)
			res = false
		}
	}

	// 上传包到服务器
	command := "scp " + packagedir + tarname + " " + ssh + ":" + destinationpath
	cmd = exec.Command("/bin/bash", "-c", command)
	_,err = cmd.Output()
	if err != nil {
	     logger.Error(err)
			 res = false
	}

	// 上传配置文件
	if redis {
		command = "scp " + packagedir + "redis-6380.conf " + ssh + ":" + destinationpath
		cmd = exec.Command("/bin/bash", "-c", command)
		_,err = cmd.Output()
		if err != nil {
				 logger.Error(err)
				 res = false
		}
	}
	if sentinel {
		command = "scp " + packagedir + "sentinel.conf " + ssh + ":" + destinationpath
		cmd = exec.Command("/bin/bash", "-c", command)
		_,err = cmd.Output()
		if err != nil {
				 logger.Error(err)
				 res = false
		}
	}


  return res
}
