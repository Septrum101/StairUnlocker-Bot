### Demo： [@stairunlock_test_bot](https://t.me/stairunlock_test_bot)

### 配置说明

### 配置说明
````yaml
converterAPI: https://api.dler.io   # subconverter API 地址
maxConn: 256                        # 节点检测同时最大连接数  
maxOnline: 10                       # 最大同时在线用户 
log_level: info                     # 日志等级: debug/info/warning/error/silent
internal: 5                         # 检测间隔时间(单位: 分钟)
telegramToken:                      # 机器人token 
````

### 命令参数

````bash
Usage of StairUnlocker-Bot:
  -f	[config Path] specify configuration file
  -h	this help
  -v	show current version of StairUnlock
````