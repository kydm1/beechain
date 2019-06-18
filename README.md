[![pipeline status](https://api.travis-ci.org/33cn/plugin.svg?branch=master)](https://travis-ci.org/33cn/plugin/)
[![Go Report Card](https://goreportcard.com/badge/github.com/33cn/plugin?branch=master)](https://goreportcard.com/report/github.com/33cn/plugin)


# 基于 chain33 区块链开发 框架 开发的 beechain公有链系统


### 编译

```
git clone https://github.com/kydm1/beechain $GOPATH/src/github.com/kydm1/beechain
cd $GOPATH/src/github.com/kydm1/beechain
go build -i -o beechain
go build -i -o beechain-cli github.com/kydm1/beechain/cli
```

### 运行
拷贝编译好的beechain, beechain-cli, beechain.toml这三个文件置于同一个文件夹下，执行：
```
./beechain -f beechain.toml
```
