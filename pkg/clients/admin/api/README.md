
# How to generate HS API code

Run:

```
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.16.2
oapi-codegen -generate client,types -package api pkg/hyperstore/api/spec.yaml > pkg/hyperstore/api/api.go
```
