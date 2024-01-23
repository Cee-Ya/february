# February

> 常规精简golang 应用快速脚手架



| 包         | 链接                            | 版本                | 作用           |
| ---------- | ------------------------------- | ------------------- | -------------- |
| gin        | github.com/gin-gonic/gin        | v1.9.1              | web请求框架    |
| uuid       | github.com/google/uuid          | v1.5.0              | trace id       |
| copier     | github.com/jinzhu/copier        | v0.4.0              | 结构体拷贝     |
| lumberjack | github.com/natefinch/lumberjack | v2.0.0+incompatible | 日志切割       |
| errors     | github.com/pkg/errors           | v0.9.1              | 更高级的errors |
| go-redis   | github.com/redis/go-redis/v9    | v9.4.0              | redis          |
| viper      | github.com/spf13/viper          | v1.18.2             | 配置读取       |
| zap        | go.uber.org/zap                 | v1.26.0             | 系统日志       |
| crypto     | golang.org/x/crypto             | v0.18.0             | 加密           |
| gorm       | gorm.io/gorm                    | v1.25.5             | orm软件        |


## 打包

### windows + cmd
```shell
set GOOS=windows
set GOARCH=amd64
go build -o ./bin/february.exe main.go
```

### windows + powershell
```shell
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -o ./bin/february.exe main.go
```

### linux
```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/february main.go
```

### mac
```shell
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/february main.go
```