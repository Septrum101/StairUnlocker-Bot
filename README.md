### Demo： [@stairunlock_test_bot](https://t.me/stairunlock_test_bot)

### 配置说明

````yaml
converterAPI: https://api.dler.io   # subconverter API 地址
maxConn: 256                        # 节点检测同时最大连接数  
maxOnline: 10                       # 最大同时检查任务数 
log_level: info                     # 日志等级: debug/info/warning/error/silent
internal: 60                        # 检测间隔时间(单位: 秒)
telegramToken: YOUR_BOT_TOKEN       # 机器人token 
````

### 命令参数

````bash
Usage of StairUnlocker-Bot:
  -f	[config Path] specify configuration file
  -h	this help
  -v	show current version of StairUnlock
````

### 性能测试

实测601个节点，耗时29.876s

```bash
StairUnlocker Bot 3.1.0 Bulletin:
Total 601 nodes, Duration: 29.876s
<Connectivity>: 476
Abema: 5
Bahamut: 2
Disney Plus: 456
HBO: 247
Netflix: 3
TVB: 0
Youtube Premium: 393
Timestamp: 2022-10-17T09:55:18Z
-------------------------
@stairunlock_test_bot
Project: https://git.io/Jyl5l
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/thank243/StairUnlocker-Bot.svg)](https://starchart.cc/thank243/StairUnlocker-Bot)
