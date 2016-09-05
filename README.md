
# AutoDeploy 
自定部署服务，自动化出来

使用Go语言，基于Beego框架实现的API服务器程序。

目前仅支持 github 的 push 通知

# 简单使用
在配置文件中根据自行情况配置
配置格式：
```
HTTPAddr = WEB运行地址
httpport = WEB运行端口
runmode = 运行模式（prod,dev）
[github项目完整地址]
secret= github密钥
sh = shell脚本地址
work = 工作目录，在脚本中可将此此配置项作为环境变量`$WORK`使用
```
执行命令运行程序：
```sh
go run main.go 
//or 
bee run
```

# 实例
## 配置
```
HTTPAddr = http://yushuangqi.com
httpport = 8080
runmode = prod

[https://github.com/ysqi/ysqi.github.io]
secret = ******
sh = $HOME/blogauto/gitpull.sh
work = $HOME/www/yushuangqi.com
```

## Shell脚本
$HOME/blogauto/gitpull.sh ：
```sh
#!/bin/sh  
cd $work
git clean -fd && git checkout --force && git pull origin master
```