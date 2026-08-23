# Telemetry Rule Service (Project 06)

纯Go遥测接入与规则分析服务，提供JSON/文本协议接入、窗口计算、阈值/变化率/连续次数/缺失规则、异常事件生命周期和Webhook订阅。默认使用无外部依赖的内存适配器运行，应用层依赖接口，可替换为PostgreSQL/TimescaleDB、Redis与消息队列实现。

## 运行

```bash
go test ./...
HTTP_ADDR=:8086 go run ./cmd/telemetryd
```

服务支持`SIGINT`/`SIGTERM`优雅关闭。配置通过环境变量覆盖，样例见`configs/config.yaml`。

## 快速流程

```bash
curl -X POST http://127.0.0.1:8086/api/v1/sources -H 'Content-Type: application/json' -d '{"id":"plant-a","name":"Plant A","protocol":"json","enabled":true}'
curl -X POST http://127.0.0.1:8086/api/v1/metrics -H 'Content-Type: application/json' -d '{"id":"temperature","name":"Temperature","unit":"celsius","enabled":true}'
curl -X POST http://127.0.0.1:8086/api/v1/rules -H 'Content-Type: application/json' -d '{"id":"hot","name":"High temperature","metric_id":"temperature","type":"threshold","operator":">","threshold":80,"window_seconds":60}'
curl -X POST http://127.0.0.1:8086/api/v1/telemetry -H 'Content-Type: application/json' -d '{"source_id":"plant-a","metric_id":"temperature","value":91}'
curl http://127.0.0.1:8086/api/v1/events
```

文本格式为`source_id,metric_id,value,RFC3339时间`，Content-Type使用`text/plain`。异常可通过`PATCH /api/v1/events/{id}`传递`{"status":"acknowledged"}`确认。规则可通过`PATCH /api/v1/rules/{id}`启停。

## 架构

- `cmd/telemetryd`: 组合根、HTTP服务和优雅关闭
- `internal/domain`: source、metric、rule、event、subscription领域模型
- `internal/application`: 接入管线、规则引擎、事件和订阅用例
- `internal/adapter`: HTTP、协议和内存仓储适配器
- `internal/infrastructure`: 配置、日志、队列和生产存储边界
- `api/openapi.yaml`: API契约
- `migrations`: PostgreSQL/TimescaleDB迁移
- `deploy`: 容器部署资产

## 运维端点

- `GET /healthz`: 存活检查
- `GET /readyz`: 就绪检查
- `GET /metrics`: Prometheus文本指标

生产环境应提供PostgreSQL、Redis及消息队列适配器，并对Webhook目标应用网络访问策略。当前实现的内存模式用于本地运行、开发和无依赖验证，进程退出后数据不会保留。
