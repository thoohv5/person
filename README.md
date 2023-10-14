# ipam

项目结构和make命令详见： app/interface/demo-interface/README.md

项目启动步骤：

1. demo目录下执行
```
go mod download
```

2. 创建数据库



3. app/interface/demo-interface目录下执行
```
make yaml # 生成配置文件，可以修改配置文件，
make run # 其中，make swag Windows环境无法执行，可以注释 app/interface/demo-interface/Makefile:20
```