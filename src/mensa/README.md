# mensa

Rethinking the git transport protocol

### 环境变量

`ENV=local`      本地环境
`ENV=dev`        测试环境
`ENV=production` 生产环境

`NOAUTH=true` 关闭验证（主要用于测试推拉等功能）

### TODO

- [x] HTTP Server
- [x] SSH Server
- [ ] Archive
- [x] Authenticate

#### 流程图

Mensa

```mermaid
graph LR
    User(用户请求) --推拉代码--> Mensa[Mensa]
    subgraph m1 [Mensa 集合仓库相关的操作]
        Mensa -- git receive --> Hook{Hook}
        Hook --> Event(创建events)
        Hook --> Prod(分支保护)
        Hook --> OPs(其他hook相关功能)
        Mensa --启动--> Svc(SVC:listen 9000)
        Svc --> Repo{仓库操作}
        Repo --> CreateRepo(创建仓库)
        Repo --> DelRepo(删除仓库)
        Repo --> Merge(合并仓库)
        Repo --> RepoOther(仓库其他的操作)
    end

    Backend(网站前端请求) --仓库相关操作--> Mensa

```
