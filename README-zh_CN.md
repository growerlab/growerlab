![Services](https://github.com/growerlab/growerlab/workflows/Services/badge.svg)

## Growerlab

代码托管平台

### 运行环境

```
ENV=test         # 测试环境
ENV=production   # 生产环境
ENV=dev || null  # ENV是dev或空值则默认为开发环境
```

### 安装教程

主要介绍 growerlab 在 ubuntu 的安装，你应拥有 sudo 权限

- 基于 ubuntu19.10
- 未来将打包成 docker 镜像

需要安装运行的服务：

- mensa - SSH/HTTP git 仓库服务
  - 端口 8022、8080
- svc - 仓库相关的操作服务（文件列表、分支、tag 等）
  - 端口 9000
- frontend - web 网站前端
- backend - web 网站后端
  - 端口 8081

#### 依赖

```shell
$ sudo apt install -y golang-go
$ sudo apt install docker.io -y
$ sudo apt install rsync openssh-client -yq
$ sudo curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
$ sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
```

```shell
$ cd ~
$ wget https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
$ tar -C /usr/local -xzf go1.14.1.linux-amd64.tar.gz
$ echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
$ source ~/.profile
```

#### 创建 Git 账号

创建系统账号，便于 ssh 访问以及托管网站服务。

- `git` 用于 ssh 访问、仓库操作等进程服务
- `growerlab` 用于部署网站服务

```shell
$ sudo adduser growerlab --system --group --disabled-password
$ sudo usermod -aG docker growerlab
```

#### 配置 SSH

如果需要 ssh 操作仓库时使用 22 端口，则需要修改 /etc/ssh/sshd_config 中的端口设置为非 22 端口（系统中默认可能是`Port 22`）。
并修改`mensa`服务的端口为 22 端口。

#### 依赖的 docker 镜像

```
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
eqalpha/keydb       latest              4e94005a2d74        10 days ago         180MB
postgres            latest              73119b8892f9        10 days ago         314MB
nginx               latest              6678c7c2e56c        10 days ago         127MB
ubuntu              latest              72300a873c2c        3 weeks ago         64.2MB
```
