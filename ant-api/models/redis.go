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


type Redis struct {
	Id                              int     `orm:"unique"`
	ClusterName                     string
	Remark                          string
	ClusterNode01                   string
	ClusterNodeRedisChecked01       bool
	ClusterNodeSentinelChecked01    bool
	ClusterNode02                   string
	ClusterNodeRedisChecked02       bool
	ClusterNodeSentinelChecked02    bool
	ClusterNode03                   string
	ClusterNodeRedisChecked03       bool
	ClusterNodeSentinelChecked03    bool
	RedisVersion                    string
	SentinelName                    string
	RedisData                       string
	SentinelData                    string
	RedisLog                        string
	RedisConf                       string
	RedisAuthentication             bool
	RedisPassword                   string
	TaskId                          string
	DeployStatus                    string
	DeployResult                    string
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

// 插入数据
func RedisInsert(r Redis) int64 {

	o := orm.NewOrm()
	id, err := o.Insert(&r)
	idString := strconv.FormatInt(id,10)
	if err == nil {
		logger.Info(idString + " redis信息插入数据成功")
	}else{
		logger.Error(idString + " redis信息插入数据失败")
	}
	return id
}

// 更新数据
func RedisUpdate(r Redis) Redis {
	o := orm.NewOrm()
	// 将id转换成字符串
	idString := strconv.Itoa(r.Id)
	_, err := o.QueryTable(new(Redis)).Filter("Id", r.Id).Update(orm.Params{
	    "DeployStatus": r.DeployStatus,
			"DeployResult": r.DeployResult,
	})
	if err == nil {
		logger.Info(idString + " redis安装信息更新数据成功")
	}else{
		logger.Error(idString + " redis安装信息更新数据失败")
	}

	// 查询最新的数据
	qs := o.QueryTable(new(Redis))
	qs.Filter("Id", r.Id).One(&r)
	return r
}

// 获取所有redis信息
func RedisGetAll() [] Redis {
	var rediss []Redis
	o := orm.NewOrm()
	qs := o.QueryTable(new(Redis))
	qs.All(&rediss)
	return rediss
}

// 获取所有redis信息
func RedisGetOne(r Redis) Redis {
	var redis Redis
	o := orm.NewOrm()
	qs := o.QueryTable(new(Redis))
	qs.Filter("Id", r.Id).One(&redis)
	return redis
}
