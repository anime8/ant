package cache


import (
    "github.com/astaxie/beego"
    "github.com/garyburd/redigo/redis"
    "github.com/wonderivan/logger"
)



func Set(key string, value string) string {
  // 创建连接
  c, err := redis.Dial("tcp", beego.AppConfig.String("redisHost") + ":" + beego.AppConfig.String("redisPort"))
  if err != nil {
     logger.Error(err)
     return "Failure"
  }
  // 设置key
  _, err = c.Do("SET", key, value)
  if err != nil {
      logger.Error(err)
      return "Failure"
  }
  // 设置过期时间
  _, err = c.Do("EXPIRE", key, beego.AppConfig.String("expire"))
  if err != nil {
      logger.Error(err)
      return "Failure"
  }
  // 关闭连接
  defer c.Close()

  return "Success"
}

func Get(key string) (string, bool) {
  haskey := true
  vs := ""
  // 创建连接
  c, err := redis.Dial("tcp", beego.AppConfig.String("redisHost") + ":" + beego.AppConfig.String("redisPort"))
  if err != nil {
     logger.Error(err)
     return vs, false
  }
  // 获取key
  v, err := c.Do("GET", key)
  if err != nil {
      logger.Error(err)
      return vs, false
  }
  // 判断任务结果是否过期，过期则返回nil
  if v != nil {
    vs = string(v.([]byte))
  }else{
    logger.Error(key + " 这个key已经过期了")
    haskey = false
  }

  // 关闭连接
  defer c.Close()
  return vs, haskey
}
