# 修复前故障复现（Docker）

## 项目与标准命令

Maintenance Work-Order Service 为设施团队提供设备资产查询、维修工单创建、维修人员指派、完工和日报查询的本地 HTTP 服务。在修复前源码状态的仓库根目录可使用以下标准命令：

```sh
go build ./...
go run ./cmd/maintenance-server
go test ./...
```

## 环境构建与编译

已实际执行以下 linux/amd64 命令：

```sh
docker buildx build --platform linux/amd64 --load -t mwos-bug001-base:amd64 -f benzhi.Dockerfile .
docker run --rm --platform linux/amd64 mwos-bug001-base:amd64 go build ./...
```

已实际执行以下 linux/arm64 命令：

```sh
docker buildx build --platform linux/arm64 --load -t mwos-bug001-base:arm64 -f benzhi.Dockerfile .
docker run --rm --platform linux/arm64 mwos-bug001-base:arm64 go build ./...
```

两个平台的镜像构建和容器内编译均成功，目标故障在下节命令中触发。

## 故障触发步骤

在修复前源码状态的仓库根目录，执行：

```sh
docker run --rm --platform linux/arm64 mwos-bug001-base:arm64 go test ./...
```

## 实际错误输出

```text
?   	github.com/1260124186-cc/maintenance-work-order-service/cmd/maintenance-server	[no test files]
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/domain	0.003s
--- FAIL: TestCreateRejectsNilAssetFromRepository (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x1264c0]

goroutine 21 [running]:
testing.tRunner.func1.2({0x15c2e0, 0x2c8dd0})
	/usr/local/go/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x318
panic({0x15c2e0?, 0x2c8dd0?})
	/usr/local/go/src/runtime/panic.go:860 +0x12c
github.com/1260124186-cc/maintenance-work-order-service/internal/service.(*WorkOrderService).Create(0x22621b2a3f28, {0x198ac0, 0x2fa400}, {{0x187a3b, 0x6}, {0x18929f, 0xc}, {0x187a53, 0x6}, {0x0, ...}})
	/workspace/internal/service/work_orders.go:41 +0xd0
github.com/1260124186-cc/maintenance-work-order-service/internal/service_test.TestCreateRejectsNilAssetFromRepository(0x22621b2e66c8)
	/workspace/internal/service/work_orders_test.go:46 +0xa4
testing.tRunner(0x22621b2e66c8, 0x195af0)
	/usr/local/go/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x3a8
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/service	0.006s
--- FAIL: TestMissingAssetIsNotReportedAsFound (0.00s)
    memory_test.go:53: GetAsset() = (<nil>, true), want (nil, false)
FAIL
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/store	0.004s
--- FAIL: TestUnknownAssetReturnsNotFound (0.00s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x1cbd50]

goroutine 22 [running]:
testing.tRunner.func1.2({0x23f080, 0x4b9d30})
	/usr/local/go/src/testing/testing.go:1974 +0x1a0
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1977 +0x318
panic({0x23f080?, 0x4b9d30?})
	/usr/local/go/src/runtime/panic.go:860 +0x12c
github.com/1260124186-cc/maintenance-work-order-service/internal/service.(*WorkOrderService).Create(0x68c1ba62a460, {0x2a78e0, 0x4f0ea0}, {{0x68c1ba614668, 0x4}, {0x68c1ba614670, 0x7}, {0x68c1ba614678, 0x6}, {0x0, ...}})
	/workspace/internal/service/work_orders.go:41 +0xd0
github.com/1260124186-cc/maintenance-work-order-service/internal/transport.(*Server).createWorkOrder(0x68c1ba61c690, {0x2a7608, 0x68c1ba628d40}, 0x68c1ba679180)
	/workspace/internal/transport/server.go:60 +0x130
net/http.HandlerFunc.ServeHTTP(0x68c1ba646180?, {0x2a7608?, 0x68c1ba628d40?}, 0x1d0f0c?)
	/usr/local/go/src/net/http/server.go:2286 +0x38
net/http.(*ServeMux).ServeHTTP(0x2a78e0?, {0x2a7608, 0x68c1ba628d40?}, 0x68c1ba679180)
	/usr/local/go/src/net/http/server.go:2828 +0x190
github.com/1260124186-cc/maintenance-work-order-service/internal/transport.(*Server).ServeHTTP(...)
	/workspace/internal/transport/server.go:29
github.com/1260124186-cc/maintenance-work-order-service/internal/transport.TestUnknownAssetReturnsNotFound(0x68c1ba6586c8)
	/workspace/internal/transport/server_test.go:52 +0x3ac
testing.tRunner(0x68c1ba6586c8, 0x2a1c40)
	/usr/local/go/src/testing/testing.go:2036 +0xc4
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:2101 +0x3a8
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/transport	0.006s
FAIL
命令退出结果：1
```

## 期望行为

对不存在或已退役的设备提交建单请求时，服务应返回可处理的未找到结果，不应中断；不得创建任何维修工单。正常设备的建单流程和停用设备的拦截行为应保持可用。
