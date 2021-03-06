# KVS (key-value store)

A key value webservice that is cached with Redis and backed by Postgres.

![diagram](kvs.png)

### Installation
```
helm install my-release k8s/helm --set cache.host=kvs-redis-host,database.host=kvs-postgres-host
```


### Usage
`POST /{key}` with the request body as the `value`

`GET /{key}` to return the `value` in the response body

`DELETE /{key}` to delete the `key` `value` pair from the cache and database

Pass the `-grpc` to the binary on startup to enable the gRPC server instead of json. The gRPC spec is in `/pkg/kvs`.

Set the `PORT` environment variable to define the port the service listens on.

Set the `TTL` environment variable to define how long Redis will cache key values.

- If the Redis hostname isn't passed, the cache functionality is disabled.
- If `SQLITE_DB` environment variable is set, Sqlite3 is used instead of Postgres.

### Kubernetes and Helm
The `/k8s/helm` directory contains the Helm Chart to easily deploy to a Kuberenetes cluster.

### TODO
- [x] Add SQLITE_DB env var to use local sqlite db instead of postgres
- [x] Add env var to disable redis cache requirement
- [x] Update Helm chart to take database and cache parameters, but don't manage them
- [ ] Ignore cache if `Cache-Control: no-cache` header is passed
