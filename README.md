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
```
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
```
### 安装node
  软件包需要自己去官网下载。
```
tar xvf node-v12.12.0-linux-x64.tar.xz  -C /usr/local/
```
### 安装go
  软件包需要自己去官网下载。
```
tar zxvf go1.14.1.linux-amd64.tar.gz -C /usr/local/
```
### 安装ant
下载并安装ant
```
mkdir -p /opt/app/go/src
git clone https://github.com/anime8/ant.git
cd ant && mv ant* /opt/app/go/src
```
### 安装beego和machinery
```
cd /opt/app/go/src/ant-api
go get github.com/RichardKnop/machinery/v1
go get -u github.com/astaxie/beego
go get -u github.com/beego/bee
```
### ant-api配置
ant-api的配置文件有2个，分别是app.conf和log.json。它们都在/opt/app/go/src/ant-api/conf下面。app.conf主要需要更改连接的reids和mysql。
初次安装需要在user表中插入登录的用户名和密码。
log.json需要更改日志路径。配置好后运行，命令如下：
```
cd /opt/app/go/src/ant-api
bee run

```
*注意：redis连接需要和ant-work配置同一个redis，mysql需要创建一个ant库，表会自动创建。*
### ant-work配置
ant-work的配置文件有2个，分别是app.conf和log.json。它们都在/opt/app/go/src/ant-work/conf下面。app.conf主要需要更改连接的reids。
log.json需要更改日志路径。配置好后运行，命令如下：
```
cd /opt/app/go/src/ant-work
bee run
```
*注意：redis连接需要和ant-api配置同一个redis。*
### ant-web配置
ant-web有一个配置/opt/app/go/src/ant-web/src/conf.js，需要将ip地址改成ant-api的地址。
配置好后运行，命令如下：
```
cd /opt/app/go/src/ant-web
npm install
npm start
```
*注意：npm install为安装依赖，只需要在首次安装时使用。*

### nginx配置
因为只有同源请求才会带上cookies，因此我们需要使用nginx做一下代理。配置如下：
```
    upstream AntWebServer {
                server 127.0.0.1:3000;
    }
    upstream AntApiServer {
                server 127.0.0.1:8080;
    }
    server {
        listen       80;
        server_name  localhost;
        
        charset utf-8;
        
        location / {
                        proxy_http_version 1.1;
                        proxy_set_header Connection "";
                        proxy_set_header Upgrade $http_upgrade;
                        proxy_set_header Connection "Upgrade";
                        proxy_set_header        Host            $host:$server_port;
                        #proxy_set_header        X-Real-IP       $remote_addr;
                        #proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
                        proxy_pass http://AntWebServer/;
        }               
        location /api/ {
                        proxy_http_version 1.1;
                        proxy_set_header Connection "";
                        proxy_set_header X-Real-IP $remote_addr;
                        proxy_set_header        Host            $host:$server_port;
                        proxy_set_header        X-Real-IP       $remote_addr;
                        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
                        proxy_pass http://AntApiServer/;
        }               
        
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }   
        
    }
```
