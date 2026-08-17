# 修复前故障复现（Docker）

## 项目与标准命令

该项目提供设备资产查询和维护工单创建、指派、完工及日报服务。在仓库根目录可执行以下标准命令：

```sh
go build ./...
go run ./cmd/maintenance-server
go test ./...
```

## 环境构建与编译

已实际执行以下命令，`linux/amd64` 与 `linux/arm64` 镜像均构建成功，且容器内 `go build ./...` 均成功：

```sh
docker buildx build --platform linux/amd64 --load -t maintenance-work-order-service:amd64 -f benzhi.Dockerfile .
docker run --rm --platform linux/amd64 maintenance-work-order-service:amd64 go build ./...
docker buildx build --platform linux/arm64 --load -t maintenance-work-order-service:arm64 -f benzhi.Dockerfile .
docker run --rm --platform linux/arm64 maintenance-work-order-service:arm64 go build ./...
```

## 故障触发步骤

在仓库根目录构建 `linux/arm64` 镜像后，执行：

```sh
docker run --rm --platform linux/arm64 maintenance-work-order-service:arm64 go test ./...
```

## 实际错误输出

```text
?   	github.com/1260124186-cc/maintenance-work-order-service/cmd/maintenance-server	[no test files]
--- FAIL: TestValidateCreateInputRejectsUnsupportedPriority (0.00s)
    rules_test.go:46: ValidateCreateInput() error = <nil>, want unsupported priority
FAIL
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/domain	0.001s
--- FAIL: TestCreateRejectsInvalidPriorityBeforeLookingUpAsset (0.00s)
    work_orders_test.go:85: Create() error = asset not found, want priority validation error
FAIL
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/service	0.001s
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/store	0.002s
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/transport	0.002s
FAIL
```

## 期望行为

当协调员提交不支持的优先级时，系统应立即返回明确的输入错误，不访问设备数据，也不创建或保存工单；`low`、`normal` 和 `urgent` 工单仍应正常完成后续流程。
