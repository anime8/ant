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


type Kafka struct {
	Id                              int     `orm:"unique"`
	ClusterName                     string
	Remark                          string
	ClusterNode01                   string
	ClusterNode02                   string
	ClusterNode03                   string
	KafkaVersion                    string
	KafkaPath                       string
	KafkaData                       string
	KafkaZookeeper                  string
	TaskId                          string
	DeployStatus                    string
	DeployResult                    string
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

// 插入数据
func KafkaInsert(k Kafka) int64 {

	o := orm.NewOrm()
	id, err := o.Insert(&k)
	idString := strconv.FormatInt(id,10)
	if err == nil {
		logger.Info(idString + " kafka信息插入数据成功")
	}else{
		logger.Error(idString + " kafka信息插入数据失败")
	}
	return id
}

// 更新数据
func KafkaUpdate(k Kafka) Kafka {
	o := orm.NewOrm()
	// 将id转换成字符串
	idString := strconv.Itoa(k.Id)
	_, err := o.QueryTable(new(Kafka)).Filter("Id", k.Id).Update(orm.Params{
	    "DeployStatus": k.DeployStatus,
			"DeployResult": k.DeployResult,
	})
	if err == nil {
		logger.Info(idString + " kafka安装信息更新数据成功")
	}else{
		logger.Error(idString + " kafka安装信息更新数据失败")
	}

	// 查询最新的数据
	qs := o.QueryTable(new(Kafka))
	qs.Filter("Id", k.Id).One(&k)
	return k
}

// 获取所有kafka信息
func KafkaGetAll() [] Kafka {
	var kafkas []Kafka
	o := orm.NewOrm()
	qs := o.QueryTable(new(Kafka))
	qs.All(&kafkas)
	return kafkas
}

// 获取一条kafka信息
func KafkaGetOne(k Kafka) Kafka {
	var kafka Kafka
	o := orm.NewOrm()
	qs := o.QueryTable(new(Kafka))
	qs.Filter("Id", k.Id).One(&kafka)
	return kafka
}
