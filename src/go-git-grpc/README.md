# go-git-grpc

支持GRPC的go-git

- [x] 通过grpc远程调用 go-git 读取仓库信息
- [x] 通过grpc远程调用 receive-pack、upload-pack 命令完成推拉操作
- [x] growerlab/hulk 为 hooks 目录下的钩子程序（将提供推拉操作产生的事件、分支保护、文件保护等功能）

#### 测试

- 执行 `test/init.sh` 初始化测试仓库
- 执行 `test/test.go` 测试grpc的go-git
- 执行 `test/test_door` 测试grpc的git推、拉

#### 性能

待测

#### 生成 proto

打开 https://github.com/protocolbuffers/protobuf/releases
下载protoc编译器：protoc-xxx-osx-x86_64.zip
将 bin/protoc 移动 $GOPATH/bin 目录下。

```
$ go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc

$ protoc --go_out=$GOPATH/src --go-grpc_out=$GOPATH/src pb/storer.proto --plugin=grpc
```
