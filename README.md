# KVS (key-value store)

A key value webservice that is cached with Redis and backed by Postgres.

![diagram](kvs.png)


### Usage
`POST /{key}` with the request body as the `value`

`GET /{key}` to return the `value` in the response body
