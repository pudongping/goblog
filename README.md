# Goblog

一个简易的基于 go 语言开发的博客项目。

## 运行代码

### 下载代码

```shell
git clone https://github.com/pudongping/goblog.git
```

### 创建数据库

```shell
CREATE DATABASE goblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 配置环境变量

```shell
cd goblog

# 复制配置环境变量文件，按照自己的需求填写环境变量信息
cp .env.example .env
```

### 下载项目依赖包

```shell
go mod tidy
```

### 运行代码

```shell
go run main.go
```

### 打包编译

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