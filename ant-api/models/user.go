package models

import (
	// "errors"
	// "strconv"
	"time"
	// "fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego"
)


func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("MysqlConnection"))
	//创建表
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Redis))
	orm.RegisterModel(new(Zookeeper))
	orm.RegisterModel(new(Kafka))
	//自动创建表 参数二为是否开启创建表   参数三是否更新表
	orm.RunSyncdb("default",false,true)
	// 设置为 UTC 时间
	orm.DefaultTimeLoc = time.UTC
}

type User struct {
	Id       int     `orm:"unique"`
	Username string
	Password string
	Email    string
	Phone    string
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
}

// 验证登录信息
func Login(u User) string {
	var user User
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	err := qs.Filter("Username", u.Username).Filter("Password", u.Password).One(&user)
	if err == orm.ErrNoRows {
	    return "Failure"
	} else if err == orm.ErrMissPK {
	    return "Failure"
	} else {
	    return "Success"
	}

}
