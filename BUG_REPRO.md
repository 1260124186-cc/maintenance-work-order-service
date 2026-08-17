# 修复前故障复现（Docker）

## 项目与标准命令

本项目为设施维护团队提供设备目录、维修工单创建、技师指派、完工状态流转和日报汇总服务。在仓库根目录执行以下标准命令：

```sh
go build ./...
go run ./cmd/maintenance-server
go test ./...
```

## 环境构建与编译

本次基线状态已在两个平台完成镜像构建和容器内编译。以下命令均在仓库根目录执行：

```sh
docker buildx build --platform linux/amd64 --load -t maintenance-work-order-service-bug-002-base-amd64 -f benzhi.Dockerfile .
docker run --rm --platform linux/amd64 maintenance-work-order-service-bug-002-base-amd64 go build ./...
docker buildx build --platform linux/arm64 --load -t maintenance-work-order-service-bug-002-base-arm64 -f benzhi.Dockerfile .
docker run --rm --platform linux/arm64 maintenance-work-order-service-bug-002-base-arm64 go build ./...
```

`linux/amd64` 和 `linux/arm64` 的镜像构建及容器内 `go build ./...` 均成功。目标故障通过下节命令触发。

## 故障触发步骤

在完成 `linux/arm64` 基线镜像构建后，于仓库根目录执行：

```sh
docker run --rm --platform linux/arm64 maintenance-work-order-service-bug-002-base-arm64 go test ./...
```

## 实际错误输出

命令退出状态：1。

```text
?   	github.com/1260124186-cc/maintenance-work-order-service/cmd/maintenance-server	[no test files]
--- FAIL: TestNewWorkOrderCopiesLabels (0.00s)
    rules_test.go:18: order labels were changed through caller input: "changed"
FAIL
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/domain	0.005s
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/service	0.015s
--- FAIL: TestMemoryRepositoryCopiesWorkOrderLabels (0.00s)
    memory_test.go:28: stored labels were changed through caller input: "changed"
FAIL
FAIL	github.com/1260124186-cc/maintenance-work-order-service/internal/store	0.004s
ok  	github.com/1260124186-cc/maintenance-work-order-service/internal/transport	0.005s
FAIL
```

## 期望行为

维护人员在创建或查询工单后修改自己用于筛选的标签列表时，已创建记录、后续查询结果和日报中的标签应保持提交时的内容不变；同一输入运行测试命令应以成功状态结束。
