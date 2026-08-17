# Maintenance Work-Order Service

Maintenance Work-Order Service is a local HTTP service for facility teams that track equipment assets, create repair work orders, assign technicians, complete repairs, and inspect a daily status summary.

## Business Scope

Facility coordinators use the service to:

- Browse active and retired equipment assets before reporting maintenance work.
- Create a work order with an asset, title, priority, and optional labels.
- Assign an open work order to a technician and complete it after the repair.
- Check a daily summary of open, assigned, and completed work orders.

The service stores only bundled in-memory data. A work order can be opened only for an active asset. An order must be assigned before it can be completed, and an empty technician name is rejected.

## Run Locally

```sh
go build ./...
go run ./cmd/maintenance-server
```

In another terminal:

```sh
curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8080/assets
curl -X POST http://127.0.0.1:8080/work-orders \
  -H 'Content-Type: application/json' \
  -d '{"asset_id":"FAN-01","title":"Inspect belt tension","priority":"urgent","labels":["safety"]}'
```

Run the public tests with:

```sh
go test ./...
```
