# ant文档

## 简介
  因为工作中需要经常部署一些redis哨兵、zookeeper和kafka集群，所以就写了这个项目。ant是有3个组价组成，ant-web为前端页面，ant-api接收前端请求，ant-work进行部署工作。
  
## 功能介绍
  1. redis哨兵部署
  2. zookeeper集群部署
  3. kafka集群部署

软件|版本|前提
--|:--:|--:
redis|redis-2.8.10, redis-3.2.0, redis-5.0.5|无
zookeeper|zookeeper-3.4.14|需要先部署jdk
kafka|kafka_2.11-1.1.0|需要先部署jdk和zk

## 部署
### 环境变量设置
vim /etc/profile
---
export GOROOT=/usr/local/go
export PATH="$PATH:$GOROOT/bin"
export GOPATH=/opt/app/go
export PATH="$PATH:$GOPATH/bin"

export NODE_HOME=/usr/local/node-v12.12.0-linux-x64
export PATH="$PATH:$NODE_HOME/bin"
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
---
source /etc/profile
