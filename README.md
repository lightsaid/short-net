# short-net

这是一个短网址生成器项目，最近使用Go语言写了一个 [Gotk](https://github.com/lightsaid/gotk) 小工具，本着学习的心态，就有为这个小工具写一个项目测试一下；顺便记录开发过程，如果喜欢，请不要吝啬你的 Star ^_^

## Let's Go

本项目是使用 macOS 系统编写构建，因此基于 Linux 环境，如果你使用是 Windows 系统建议安装 WSL 2 搭建 Docker 环境。

检查你的环境是否有 Docker 和 Docker Compose, 如若没有，先安装好再继续下看
``` bash
# 验证环境
docker --help
docker-compose --help
```

### 基于 Docker 搭建环境
创建项目，初始化 go module 和 git 仓库
``` bash
cd $GOPATH/src/github.com/lightsaid

mkdir short-net && cd short-net

go mod init github.com/lightsaid/short-net

mkdir -p cmd/web

mkdir models dbrepo

touch cmd/web/main.go

touch .gitignore README.md

git init .

```
在 cmd/web/main.go 中编写

``` go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	var serverPort = os.Getenv("HTTP_SERVER_PORT")

	fmt.Println(">>>>>> ", serverPort)
	fmt.Println(">>>>>> ", os.Getenv("VIEW_PATH"))
	fmt.Println(">>>>>> ", os.Getenv("PUBLIC_PATH"))
	fmt.Println(">>>>>> ", os.Getenv("DB_SOURCE"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	log.Println("starting server on :", serverPort)
	err := http.ListenAndServe(fmt.Sprintf(":%s", serverPort), mux)
	if err != nil {
		log.Println("start error: ", err)
	}
}

```
---
创建 Dockerfile、docker-compose.yml 文件

``` bash
touch Dockerfile docker-compose.yml
```
---
编写 Dockerfile 镜像构建文件
``` dockerfile
# 参考：https://hub.docker.com/_/golang

FROM golang:1.20.3

# 设置 GOPROXY 环境变量
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct

# 去到工作目录
WORKDIR /usr/src/app

# 安装热重载插件，方便开发
RUN go install github.com/cosmtrek/air@latest

# 将当前目录所有内容复制到容器中
COPY . .

# 检查依赖并安装
RUN go mod tidy
```

上面使用 github.com/cosmtrek/air 这个工具，以此需要在根目录下添加配置 `.air.toml`, 
去到这个工具githu地址，将 `air_example.toml` 拷贝过来，仅仅需要修改一个小地方即可。
其实就是指定构建go程序main函数包目录。
``` toml
[build]
# Just plain old shell command. You could use `make` as well.
cmd = "go build -o ./tmp/main ./cmd/web"
```

---
编写 docker-compose.yml 文件
``` yml
version: "3.8"

services:
  web:
    build: .          # 设置上下文，会查找当前目录下的 Dockerfile 文件，基于文件描述构建  
    env_file:
      - .env          # 加载配置文件
    ports:
      - ${HTTP_SERVER_PORT}:4000     # web 服务端口映射
    volumes:          
      - .:/usr/src/app   # 数据卷，当前目录:容器目录
    networks:      
      - default  # 使用默认网络，bridge 模式，容器之间可以访问
    depends_on:   # 容器启动先后次序，这里先启动 mysql，但并一定是mysq启动完成后再启动web
      - mysqldb

    # 环境变量，这是设置web服务的环境变量，当执行 go run main.go 时，
    # 这些环境变量是加载到程序的，可以通过 os.Getenv(key) 获取。
    # 建议一些固定的配置可以直接写这里，经常要修改的还是放在配置文件里，如这里的 .env 文件
    # 当然，统一在配置文件来设置也很灵活，查找变量时，就不同看两个地方了。这就仁者见仁...
    environment:  
      - VIEW_PATH=/views
      - PUBLIC_PATH=/public
      - DB_SOURCE=root:${DB_PASSWORD}@tcp(mysqldb:3306)/shortnet?charset=utf8mb4&parseTime=True&loc=Local

    # command: go run cmd/main.go -b 0.0.0.0
    command: air cmd/main.go -b 0.0.0.0     # 使用 air 启动程序，监听文件变化，热重载

  mysqldb:
    image: "mysql:5.7.22"
    ports:
      # mysql端口映射，属主机使用这个mysql服务，用DB_PORT端口；而容器之间使用则用 3306
      - "${DB_PORT}:3306"  
    restart: always  # 总是启动，当docker启动时，就会启动这个服务
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    volumes:
      - ${DB_STORAGE}:/var/lib/mysql
    networks:
      - default

networks:
  default:
  
```
上面设置 env_file 配置，因此需要创建 .env 配置文件，同时添加 mysql 的环境变量
```
# MySQL Config
DB_PORT=3307
DB_NAME=shortnet
DB_PASSWORD=abc123
DB_STORAGE=./storage/mysql

# Web Server Config
HTTP_SERVER_PORT=4000

```
---

当上面一起准备好时，就可以执行 `docker-compose up` 命令构建环境了
看下日志输出，如果有以下信息输出，则说明成功一半了。
``` log
...
web_1      | 
web_1      |   __    _   ___  
web_1      |  / /\  | | | |_) 
web_1      | /_/--\ |_| |_| \_ , built with Go 
web_1      | 
web_1      | mkdir /usr/src/app/tmp
web_1      | watching .
...

web_1      | running...
web_1      | >>>>>>  4000
web_1      | >>>>>>  /views
web_1      | >>>>>>  /public
web_1      | >>>>>>  root:abc123@tcp(mysqldb:3306)/shortnet?charset=utf8mb4&parseTime=True&loc=Local
web_1      | 2023/04/07 12:09:01 starting server on : 4000
```

在浏览器地址栏输入: http://localhost:4000/ 访问，如果看到 Hello! 输出，走在正轨上，对的，你没出轨（:joy:）

最后修改一下 Hello 输出，刷新浏览器也跟着发生变化，很上道哈（:smiley:）。

--- 

下面验证mysql是否正常安装, 如果下面的输出 `STATUS` 是 Up 则运行成功了。

``` 
➜ short-net git:(master) ✗ docker ps
CONTAINER ID   IMAGE    ...          STATUS          PORTS                               NAMES
216b432d22ef   short-net_web         Up 59 minutes   0.0.0.0:4000->4000/tcp              short-net_web_1
07e386fdf7f6   mysql:5.7.22          Up 59 minutes   0.0.0.0:3307->3306/tcp              short-net_mysqldb_1
```

--- 
进入 Web 服务容器
``` bash
# 查看web服务容器名，例如我的是short-net_web_1
docker ps 
# 进入web服务容器
docker exec -it short-net_web_1 bash
# 看一下文件列表
ls  
# 查看 go 版本
go version
# 查看环境
go env
```
---
这里简单解析下当执行 `docker-compose up` 命令时，做了什么。

首先对docker有一点点基础知识。

解析一下以下两个文件：
 - Dockerfile 镜像构建描述文件
 - docker-compose 容器编排, 将一组相关联容器编排组合成一个完整的项目
像现在这个项目，由两个容器(服务)组成，分别是 web 服务和 mysqldb 服务。

执行 `docker-compose up` 的时候，首先在执行命令的当前路径找到 `docker-compose.yml` 配置文件；
然后按照配置文件描述服务的先后顺序构建。

其中 `build` 字段是上下文的意思，也就是基于哪个目录的内容构建服务，比如这里使用 `.` 就是当目录，然后 docker-compose
会根据这个上下文，查找加载 Dockerfile 描述文件的内容参与构建服务。

当然也是可以指定其他目录作为构建上下文，这个时候`build`字段稍稍调整一下，如：
``` yml
build:
    context: ./../my-service
    dockerfile: ./../my-service/my-service.dockerfile
```

其他字段想了解，查看 Dockerfile 和 docker-compose.yml 注释。

--- 

另外简单说一下 Docker 网络

当执行 `docker-compose up` 构建完后，在整个short-net应用里包含web、mysqldb 两个容器服务，他两是相互独立的。
在这里就有几种关系: 宿主机 和 web、宿主机 和 mysqldb、mysqldb 和 web。
- 宿主机访问容器服务，就可以直接使用 127.0.0.1(localhost) + 映射的端口。
- 如访问web服务器：http://localhost:4000/
- 如访问mysql(前提宿主机安装了mysql，要使用mysql指令)：mysql -h127.0.0.1 -P3307 -uroot -pabc123

那么宿主机之间该如何访问呢？那么在 web 服务要连接到mysqldb服务，进行操作数据库如何做呢？
要了解这个问题，要先知道 Docker 中网络类型：分别是 `bridge、host、none` 这三种。
- host 仅允许宿主机和容器之间进行网络连接，容器之间无法通过网络访问
- none 不会分配网络，自能自己玩
- bridge 是默认网络桥接模式，它们是可以访问容器的，前提是他们必须在同个虚拟网络里。

如何确认他们之间能通过网络访问呢？当他们在同一个虚拟网络里就能够相互访问，
如下：web 服务和 mysqldb 服务都在 short-net_default 虚拟网络里，它们是可以相互访问的。
``` bash
➜  ~ docker network ls
# NETWORK ID     NAME                DRIVER    SCOPE
# 0d28d3743d25   bridge              bridge    local
# 222c30f600e0   docker_default      bridge    local
# 0d5dedb7bb5b   host                host      local
# c6b4c5fdd7b2   none                null      local
# dbebfb23a13d   short-net_default   bridge    local
➜  ~ 

# 上面找到 short-net_default 网络，接着查看详情
➜  ~ docker network inspect short-net_default
# ....
# "Containers": {
#   "07e386fdf7f6487c455c3468fa839245b9a7f489a085d1e8796150d2e209a4fc": {
#       "Name": "short-net_mysqldb_1",
#       "EndpointID": "a76e7fcb635d76b5731ff77cd65a3561bba9b20daa9f9400d482dc7bc6faeec0",
#       "MacAddress": "02:42:ac:19:00:02",
#       "IPv4Address": "172.25.0.2/16",
#       "IPv6Address": ""
#   },
#   "216b432d22effd4ab6c1a0f45974677b4c291f377b575699df52a08350abcb2c": {
#       "Name": "short-net_web_1",
#       "EndpointID": "c3df013c0404a313555b58ff6cd385047d7d774964ef9e16e673642a30adea5d",
#       "MacAddress": "02:42:ac:19:00:03",
#       "IPv4Address": "172.25.0.3/16",
#       "IPv6Address": ""
#   }
# },
# ....
```

一般访问服务，我门是通过，IP + Port 来确定一个服务，但是Docker容器之间访问，不建议这么访问，而是通过 服务名 + Port 访问。
如果我在使用 web 服务中使用 GORM 连接 mysqldb 服务，
那么是这样子的 DSN := "root:abc123@tcp(mysqldb:3306)/shortnet?charset=utf8mb4&parseTime=True&loc=Local"
- Host = mysqldb
- Port 是mysqldb服务本地的端口 3306 而不是映射的 3307 

为什么不建议使用 Host + Port 访问，是应为重启服务时 IP 会变变化的。（除非，配置固定IP）


## 数据库设计和功能

为了简单，仅设计一些核心的表, 不会设计商业分析的表

### 数据库设计

数据库分为两个表

users table
- id
- name
- email
- password
- avatar
- role (USER、ADMIN)
- created_at
- updated_at
- deleted_at

links table
- id 
- user_id
- long_url
- short_hash  
- click      (别点击次数)
- created_at
- updated_at
- expired_at (过期时间)


### 基于上面两个表要实现什么功能呢？
- 注册
- 登录
- 更新用户信息
- 创建短网址
- 修改短网址
- 重定向短网址
- 短网址列表
- 短网址删除
- 定时检查短网址，过期超过多长时间即可删除

看着好像没什么功能，如果将每一项功能细化划分，其实代码可不少。

分享一个可以免费生成 Logo 网站：
  https://www.namecheap.com/logo-maker/app/

