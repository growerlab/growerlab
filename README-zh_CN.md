## Growerlab 

代码托管平台

### 安装教程

主要介绍growerlab在ubuntu的安装，你应拥有sudo权限

- 基于ubuntu19.10
- 未来将打包成docker镜像

需要安装运行的服务：

- mensa - SSH/HTTP git仓库服务
- svc - 仓库相关的操作服务（文件列表、分支、tag等）
- frontend - web网站前端
- backend - web网站后端

#### 依赖

```shell
$ sudo apt install docker.io -y
$ sudo apt install rsync openssh-client -yq
```

#### 创建Git账号

创建系统账号，便于ssh访问以及托管网站服务。
- `git` 用于ssh访问、仓库操作等进程服务
- `growerlab` 用于部署网站服务

```shell
$ sudo adduser git --system --shell /bin/bash --gecos 'git version control' --group --disabled-password --home /home/git
$ sudo adduser growerlab --disabled-login --gecos 'growerlab system'
$ sudo usermod -a -G git growerlab
$ sudo usermod -a -G growerlab git

$ sudo -s -u git
$ cd
$ mkdir /home/git/repositories
$ chmod -R g+rwX /home/git/repositories/
$ chown -R git:git /home/git/repositories/
```

#### 配置SSH

如果需要ssh操作仓库时使用22端口，则需要修改 /etc/ssh/sshd_config 中的端口设置为非22端口（系统中默认可能是`Port 22`）。
并修改`mensa`服务的端口为22端口。


#### 依赖的docker镜像

```
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
eqalpha/keydb       latest              4e94005a2d74        10 days ago         180MB
postgres            latest              73119b8892f9        10 days ago         314MB
nginx               latest              6678c7c2e56c        10 days ago         127MB
ubuntu              latest              72300a873c2c        3 weeks ago         64.2MB
```
