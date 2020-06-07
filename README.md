# terraform-varchecking-service

## Usage
- Get all Datasource supported now:

`curl http://127.0.0.1:10001/GetDatasource`

- Get all Resource supported now:

`curl http://127.0.0.1:10001/GetResource`

- Query variable name which contains the specified one.

`curl http://127.0.0.1:10001/QueryName -d '{"QueryKey":"vpc"}'`

## Note:
1. For the query function, it's case insensitive. But it only supports single word. 
For example, if you want to query `force_delete`, both `{"QueryKey":"force"}` and 
`{"QueryKey":"delete"}` are ok, but can't specify `{"QueryKey":"force_delete"}`.
2. For zsh shell, response contains an extra `%` symbol at the end of line.
