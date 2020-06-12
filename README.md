# terraform-varchecking-service

## Usage
1.Download binary executable file in binary directory and run it:

`./varchecking-service` 

2.Run the curl command below:

```

- Get all Datasource supported now:
curl http://127.0.0.1:10001/GetDatasource

- Get all Resource supported now:
curl http://127.0.0.1:10001/GetResource

- Query variable name which contains the specified one:
curl http://127.0.0.1:10001/QueryName -d '{"QueryKey":"vpc"}'

```

## Note:
1. Current tencentcloud provider version: `v1.36.1`
2. For the query function, it's case insensitive. But it only supports single word. 
For example, if you want to query `force_delete`, both `{"QueryKey":"force"}` and 
`{"QueryKey":"delete"}` are ok, but can't specify `{"QueryKey":"force_delete"}`.
3. For zsh shell, response contains an extra `%` symbol at the end of line.

