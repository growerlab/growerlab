## 依赖

### 数据库

MySql v8.0.x

KeyDB 5.x

#### 初始化数据库

创建数据库名称，数据库用户，用户密码均为 growerlab 的数据库

```
CREATE USER 'growerlab'@'localhost' IDENTIFIED BY 'growerlab';
GRANT DELETE, EXECUTE, SELECT, CREATE ROUTINE, ALTER ROUTINE, CREATE VIEW, GRANT OPTION, REFERENCES, TRIGGER, UPDATE, DROP, CREATE, LOCK TABLES, EVENT, INDEX, ALTER, SHOW VIEW, INSERT, CREATE TEMPORARY TABLES 
    ON `growerlab`.* TO 'growerlab'@'localhost';
UPDATE mysql.user SET max_questions = 0, max_updates = 0, max_connections = 0 WHERE User = 'growerlab' AND Host = 'localhost';
CREATE DATABASE `growerlab`;
USE `growerlab`;
```


#### 初始化数据库表结构

使用 `db/growerlab.sql` 文件初始化表结构

如果有种子数据，应该放入 `db/seed.sql` 文件中

```
初始账号 admin@growerlab.com
初始密码 growerlabadmin
```

#### GraphQL

基于 `gqlgen` 通过自动化生成GraphQL的基础代码

如果有修改 `*.graphql` 应该使用 `gqlgen` 工具生成代码
