# 修复前故障复现（Docker）

## 项目与标准命令

本项目是用于管理设施设备、维修工单、技师指派和日报汇总的本地 HTTP 服务。在仓库根目录可使用以下标准命令：

```sh
go build ./...
go run ./cmd/maintenance-server
go test ./...
```

## 环境构建与编译

已实际执行以下 linux/amd64 命令，镜像构建和容器内编译均成功：

```sh
docker buildx build --platform linux/amd64 --load -t maintenance-work-order-service-bug005-base:amd64 -f benzhi.Dockerfile .
docker run --rm --platform linux/amd64 maintenance-work-order-service-bug005-base:amd64 go build ./...
```

已实际执行以下 linux/arm64 命令，镜像构建和容器内编译均成功：

```sh
docker buildx build --platform linux/arm64 --load -t maintenance-work-order-service-bug005-base:arm64 -f benzhi.Dockerfile .
docker run --rm --platform linux/arm64 maintenance-work-order-service-bug005-base:arm64 go build ./...
```

## 故障触发步骤

在仓库根目录基于修复前代码构建 linux/arm64 镜像后，执行：

```sh
docker run --rm --platform linux/arm64 maintenance-work-order-service-bug005-base:arm64 go test ./...
```

## 实际错误输出

```text
?   	github.com/1260124186-cc/maintenance-work-order-service/cmd/maintenance-server	[no test files]
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/domain	0.002s
--- FAIL: TestDailySummaryReturnsCloseError (0.00s)
    work_orders_test.go:50: DailySummary() error = <nil>, want audit close failed
--- FAIL: TestDailySummaryCompletesWithoutWorkOrders (0.00s)
    work_orders_test.go:77: second DailySummary() error = daily audit is already open
FAIL
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/service	0.004s
--- FAIL: TestAuditCloseReturnsConfiguredError (0.00s)
    memory_test.go:55: Close() error = <nil>, want audit close failed
FAIL
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/store	0.001s
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/transport	0.002s
FAIL
```

## 期望行为

查询没有工单的日期日报后，下一次日报查询应能正常继续；当资源关闭出现异常时，请求应返回该异常。工单创建、指派、完成及日报状态统计应保持可用且准确。
