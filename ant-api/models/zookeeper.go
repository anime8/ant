package models

import (
	// "errors"
	"strconv"
	"time"
	// "fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wonderivan/logger"
)


type Zookeeper struct {
	Id                              int     `orm:"unique"`
	ClusterName                     string
	Remark                          string
	ClusterNode01                   string
	ClusterNode02                   string
	ClusterNode03                   string
	ZookeeperVersion                string
	DeployPath                      string
	ZookeeperData                   string
	ZookeeperLog                    string
	TaskId                          string
	DeployStatus                    string
	DeployResult                    string
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

// 插入数据
func ZookeeperInsert(z Zookeeper) int64 {

	o := orm.NewOrm()
	id, err := o.Insert(&z)
	idString := strconv.FormatInt(id,10)
	if err == nil {
		logger.Info(idString + " zookeeper信息插入数据成功")
	}else{
		logger.Error(idString + " zookeeper信息插入数据失败")
	}
	return id
}

// 更新数据
func ZookeeperUpdate(z Zookeeper) Zookeeper {
	o := orm.NewOrm()
	// 将id转换成字符串
	idString := strconv.Itoa(z.Id)
	_, err := o.QueryTable(new(Zookeeper)).Filter("Id", z.Id).Update(orm.Params{
	    "DeployStatus": z.DeployStatus,
			"DeployResult": z.DeployResult,
	})
	if err == nil {
		logger.Info(idString + " zookeeper安装信息更新数据成功")
	}else{
		logger.Error(idString + " zookeeper安装信息更新数据失败")
	}

	// 查询最新的数据
	qs := o.QueryTable(new(Zookeeper))
	qs.Filter("Id", z.Id).One(&z)
	return z
}

// 获取所有zookeeper信息
func ZookeeperGetAll() [] Zookeeper {
	var zookeepers []Zookeeper
	o := orm.NewOrm()
	qs := o.QueryTable(new(Zookeeper))
	qs.All(&zookeepers)
	return zookeepers
}

// 获取一条zookeeper信息
func ZookeeperGetOne(z Zookeeper) Zookeeper {
	var zookeeper Zookeeper
	o := orm.NewOrm()
	qs := o.QueryTable(new(Zookeeper))
	qs.Filter("Id", z.Id).One(&zookeeper)
	return zookeeper
}
