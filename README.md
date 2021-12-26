# Goblog

一个简易的基于 go 语言开发的博客项目。

## 创建数据库

```shell
CREATE DATABASE goblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 编译运行

```shell

# 创建 tmp 目录
mkdir tmp

# 生成二进制文件 goblog
go build -o tmp/goblog-service

# 复制环境配置信息
cp .env tmp/

# 进入 tmp 目录，然后运行我们的应用程序
cd tmp
# 运行
./goblog-service

```