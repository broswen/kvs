# KVS (key-value store)

A key value webservice that is cached with Redis and backed by Postgres.


### Usage
`POST /{key}` with the request body as the `value`

`GET /{key}` to return the saved item

```json
  "key": "key here",
  "value": "value here"
```