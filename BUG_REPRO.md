# 修复前故障复现（Docker）

## 项目与标准命令

该项目提供设备资产、维修工单创建和日报汇总的 HTTP 服务。在仓库根目录可执行以下标准命令：

```sh
go build ./...
go run ./cmd/maintenance-server
go test ./...
```

## 环境构建与编译

基线状态已实际完成以下两个平台的镜像构建和容器内编译：

```sh
docker buildx build --platform linux/arm64 --load -t maintenance-work-order-service:delivery004-base-arm64 -f benzhi.Dockerfile .
docker run --rm --platform linux/arm64 maintenance-work-order-service:delivery004-base-arm64 go build ./...
docker buildx build --platform linux/amd64 --load -t maintenance-work-order-service:delivery004-base-amd64 -f benzhi.Dockerfile .
docker run --rm --platform linux/amd64 maintenance-work-order-service:delivery004-base-amd64 go build ./...
```

上述两个平台的镜像构建和容器内 `go build ./...` 均成功。

## 故障触发步骤

在仓库根目录执行以下命令：

```sh
docker run --rm --platform linux/arm64 maintenance-work-order-service:delivery004-base-arm64 go test ./...
```

该命令会稳定触发取消后的工单创建验证失败。

## 实际错误输出

```text
?   	github.com/1260124186-cc/maintenance-work-order-service/cmd/maintenance-server	[no test files]
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/domain	0.003s
--- FAIL: TestCreateHonorsCanceledContext (0.00s)
    work_orders_test.go:65: Create() error = <nil>, want canceled context
FAIL
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/service	0.019s
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/store	0.005s
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/transport	0.003s
FAIL
```

## 期望行为

同一取消时序下，创建操作应返回取消结果且不留下新的待处理工单；日报中的待处理数量不应包含该已撤销请求。
