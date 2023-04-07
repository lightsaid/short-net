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