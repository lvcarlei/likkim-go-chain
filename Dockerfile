# 使用官方提供的Go语言基础镜像
FROM golang:1.22-alpine as build

# 设置工作目录
WORKDIR /app

# 复制并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制整个项目到工作目录
COPY . .

# 构建应用程序
RUN go build -o main .

# 运行时镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建镜像中复制二进制文件
COPY --from=build /app/main .

# 暴露端口（如果需要）
EXPOSE 8082

# 运行应用程序
CMD ["./main"]

