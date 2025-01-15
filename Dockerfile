# 第一阶段：构建 Go 应用
FROM golang:alpine AS builder

# 安装必要的工具
RUN apk update && apk add --no-cache git gcc musl-dev libc-dev

# 设置工作目录
WORKDIR /app

# 设置 Go 代理（可选）
ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=1
ENV CC=gcc

# 复制代码并构建
COPY . /app
RUN go mod tidy && go build -o app

# 第二阶段：精简运行环境
FROM alpine

# 设置工作目录
WORKDIR /app

# 复制编译结果和配置文件
COPY --from=builder /app/app .
COPY --from=builder /app/config /app/config
COPY --from=builder /app/db /app/db
# 暴露端口（根据应用需要）
EXPOSE 8080

# 启动服务
CMD ["./app"]
