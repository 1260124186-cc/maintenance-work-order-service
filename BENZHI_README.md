# 设备维护工单服务 Docker 说明

本项目为设施维护团队提供本地设备目录、维修工单创建、技师指派、完工状态流转和日报汇总能力。服务不依赖外部网络、数据库或第三方服务，容器中始终从仓库源码编译运行。

## 本地标准命令

在仓库根目录执行：

```sh
go build ./...
go run ./cmd/maintenance-server
go test ./...
```

服务默认监听 `8080` 端口。启动后可验收健康检查和设备目录：

```sh
curl --fail http://127.0.0.1:8080/health
curl --fail http://127.0.0.1:8080/assets
```

## Docker 环境

实际 Docker 文件为 `benzhi.Dockerfile`。`go.mod` 固定 Go 语言版本为 `1.26.2`，Dockerfile 使用 `golang:1.26.2-bookworm` 并设置 `GOTOOLCHAIN=local`，避免构建时自动切换工具链。镜像在复制源码后运行 `go build ./...`，不会复制宿主机二进制。

标准双架构验收命令：

```sh
./build_benzhi_docker.sh linux/amd64
./build_benzhi_docker.sh linux/arm64
```

脚本会依次构建指定平台镜像、启动服务容器、在容器内执行 `go build ./...`，然后请求健康检查和设备目录 API。所有命令退出码为 `0` 且两个 HTTP 请求成功时，表示该平台通过。

也可以逐步手工执行：

```sh
docker buildx build --platform linux/amd64 --load \
  -t maintenance-work-order-service:amd64 -f benzhi.Dockerfile .
docker run -d --rm --name maintenance-work-order-service-check \
  -p 8080:8080 maintenance-work-order-service:amd64
docker exec maintenance-work-order-service-check go build ./...
curl --fail http://127.0.0.1:8080/health
curl --fail http://127.0.0.1:8080/assets
docker stop maintenance-work-order-service-check
```

将上述 `linux/amd64` 和镜像标签替换为 `linux/arm64` 与 `maintenance-work-order-service:arm64`，即可进行 arm64 验收。容器内编译、服务成功启动、`/health` 返回成功状态并且 `/assets` 返回 JSON 设备目录，均为通过条件。
