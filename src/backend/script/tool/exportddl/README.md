### export ddl

Export DDL to `db/growerlab.sql` from PostgreSQL

该工具主要是为了方便将数据库结构快速导出到 db/growerlab.sql

避免，修改一个字段或者相关数据库结构操作频繁的手动导出到 db/growerlab.sql，同时避免手动操作可能的失误

#### 使用方法

```shell
$ cd backend
$ go run script/tool/exportddl/main.go
```