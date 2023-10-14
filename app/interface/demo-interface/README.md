# demo-interface

# 项目框架

```
demo-interface # 项目名
    api # 文档与对外提供的数据结构
        constant # 常量
            enum # 枚举，按照业务模块划分
        docs # swagger文档
        http # HTTP服务
            request  # 请求结构体，按照业务模块划分
            response # 返回结构体，按照业务模块划分
        translate # 多语言
            locales # 本地化
                en-GB # 英语
                    messages.gotext.json # 翻译专用
                zh-CN # 中文
    boot  # 引导
    build # 编译
    cmd   # 启动
        command        # 命令行模式
        demo-interface # 项目启动
    configs    # 配置
        amqp.yaml     # AMQP配置
        business.yaml # 业务配置
        cron.yaml     # 定时配置
        database.yaml # 数据库配置
        http.yaml     # HTTP配置
        logger.yaml   # 日志配置
        redis.yaml    # Redis配置
        rpc.yaml      # RPC配置
    internal # 业务
        conf       # 配置
        consumer   # AMQP业务消费，按照业务模块划分
        controller # 控制器，按照业务模块划分
        cron       # 定时任务，按照业务模块划分
        repository # 数据持久层，按照业务模块划分
        router     # 路由，按照业务模块划分
        service    # 服务，按照业务模块划分
    migrations # 迁移
    mocks      # 模拟
    test-results   # 测试
    .gitignore     # git忽略
    .gitlab-ci.yml # CI
    .golangci.yaml # golang CI
    LICENSE   # 版本
    Makefile  # MF
    README.md # 说明
```

## 工具

### 帮助
```
make help
```

### 项目工程化
pkg/cmd/generate/README.md
```
init-module -t 数据库表名 -tn 模块名称 -d 项目绝对路径 -p 项目名称
init-module -t ip_address -tn IP地址 -d /Users/thooh/Projects/github.com/thoohv5/person/app/interface/demo-interface -p demo-interface
```

### 配置拷贝
从示例配置到配置文本，影响：项目/configs
```
make yaml
```

### 序列化
影响：枚举（项目/api/constant/enum）,多语言（项目/api/translate），生成耗时（建议在文件中，点击//go:generate生成）
```
make gen
```

### 代码格式化
```
make format
```

### wire依赖
在模块注册时，需要运行，比如运行init-module
```
make wire
```

### CI检查
提交代码时必须无报错
```
make lint
```

### 配置结构化
配置文件对应的GolangStruct的生成，需要调整配置时使用
```
make config
```


### swagger文档
```
make swag
```

### mock依赖
```
make mock
```

### 测试
```
make test
```

### 测试覆盖率
```
make cover
```

### swaggerHTTP测试
需要启动服务
```
make swagger-ci
```

### 项目编译
```
make build
```

### 项目运行
```
make run
```

### 数据库迁移工具
在项目中可以运行 `go run cmd/demo-interface/main.go migration COMMAND`
关联文件夹 项目/api/migrations，支持使用model或者原生SQL，SQL的文件夹的命名规范为 NUM_OP_DESC.tx.[up/down].sql。如：2_add_data.tx.up.sql
```
migration COMMAND
COMMAND:
  - init                  # 初始化
  - up                    # 升级到最新版本
  - up [target]           # 升级到指定版本
  - down                  # 降级到上一个版本
  - reset                 # 还原初始化
  - version               # 当前版本
  - set_version [version] # 设置版本
```



