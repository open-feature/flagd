# HTTP(S) Service


Why is my `int` response a `string`?
```
$ curl -X POST "localhost:8013/flags/myIntFlag/resolve/int"
{"value":"1","reason":"STATIC","variant":"one"}
```
When interacting directly with the flagD http(s) api and requesting an `int` the response type will be a `string`. This behaviour is introduced by [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway), which uses [proto3 json mapping](https://developers.google.com/protocol-buffers/docs/proto3#json) to build the response object. If a number value is required, and none of the provided SDK's can be used, then it is recommended to use the `float64` endpoint instead:
```
$ curl -X POST "localhost:8013/flags/myIntFlag/resolve/float"
{"value":1,"reason":"STATIC","variant":"one"}
```
