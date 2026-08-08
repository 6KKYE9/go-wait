# go-wait

编解码这种小事，犯不着每次都跑在线工具网站溜一圈。

等某个条件满足再退出，常写在脚本或 CI 里：服务等起来再跑测试、等构建产物生成再打包。

零依赖，只用 Go 标准库。

## 用法

```
# 等端口可连
go-wait -port 127.0.0.1:8080

# 等文件出现
go-wait -file dist/app.zip

# 等一条命令跑成功（支持管道）
go-wait -cmd "curl -sf http://localhost/health"

# 带超时和轮询间隔
go-wait -port 127.0.0.1:8080 -timeout 60s -interval 1s
```

选项：
- `-port <地址>`：等待 `host:port` 可建立 TCP 连接
- `-file <路径>`：等待某个文件出现
- `-cmd <命令>`：等待某条命令成功退出（Windows 用 `cmd /C`，其他用 `sh -c`）
- `-timeout <时长>`：最长等待，默认 30s，例如 `2m`
- `-interval <间隔>`：轮询间隔，默认 500ms

三种条件必须且只能选一种。条件满足退出码 0，超时退出码 1。
