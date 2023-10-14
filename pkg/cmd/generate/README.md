# 代码生成模块

为了生成指定代码模板，功能简单，不是很通用

推荐使用：init-module 命令 => 初始化 orm、controller、service 和 router

## 介绍

使用方式：在命令目录下执行 `go install` 进行编译安装，然后就可以直接使用了

## generate命令介绍

单纯的初始化数据，不会注册到对应文件中。相对来说比较通用

查看命名介绍指令：generate-xxx -h

所有命令：

```txt
generate-ctrl   #初始化controller
generate-dao   #初始化orm和dao层
generate-serv  #初始化service
generate-req   #初始化request
generate-resp   #初始化response
```

## register命令介绍

单纯注册到指定文件中，相对定制化。里面注册变量都是定死的了。后期可以考虑优化

查看命名介绍指令：register-xxx -h

所有命令：

```txt
register-ctrl       #注册controller
register-serv       #注册service
register-dao       #注册orm和dao层
register-route       #注册router
```

## init命令介绍

包含了 generate 和 register 命令，建议这样使用

init命令都有一个 config.yaml 文件，这个就是配置相关文件路径信息的

查看命名介绍指令：init-xxx -h

所有命令：

```txt
#初始化一个完成的module模板并注册
#包含：controller、service、request、response、router、dao（orm）
init-module        

#初始化orm和dao，并注册
init-dao        

#初始化逻辑层并注册，包含：controller、service、request、response、router
init-logic        

```

指令例子：init-module

参数解释：
```txt
Usage of init-module:
  -c string
        c: config file path：input absolute path (default "./config.yaml")
  -d string
        d: dir, project absolute path. If it is empty or '.', the configured default address will be used
  -m string
        m: project mode (default "ipam")
  -p string
        p: project name (default "ipam-interface")
  -t string
        t: table name, generate tables
  -tn string
        tn: table name (default "{xxx}")
```

当前执行指令路径: xxx/ipam/app/interface/ipam-interface

配置路径：xxx/ipam/app/interface/ipam-interface/configs/config/config.yaml

执行指令：init-module -t hello -tn 你好 -c configs/config.yaml

配置介绍：

```yaml
# ipam-interface config

# 项目名称
project_name: "ipam-interface"

# 项目绝对路径
project_abs_path: "/home/code/go/thoohv5/person/app/interface/ipam-interface"

# 各个模块配置
config:
  - name: "request"        # 模块名称
    registered_file: ""    # 需要被注册的文件路径
    path: "/api/http/request" # 模块在项目中的路径
  - name: "response"
    registered_file: ""
    path: "/api/http/response"
  - name: "controller"
    registered_file: "controller.go"
    path: "/internal/controller"
  - name: "service"
    registered_file: "service.go"
    path: "/internal/service"
  - name: "repository"
    registered_file: "repository.go"
    path: "/internal/repository"
  - name: "router"
    registered_file: "router.go"
    path: "/internal/router"
    parent_name: "ripam"  # 父级路由变量名称
```