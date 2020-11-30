![OneBot-YaYa](https://socialify.git.ci/Yiwen-Chan/OneBot-YaYa/image?description=1&descriptionEditable=OneBot%20(CQHTTP)%20%E5%85%88%E9%A9%B1%E5%B9%B3%E5%8F%B0%E7%9A%84%E5%AE%9E%E7%8E%B0&font=Raleway&forks=1&issues=1&logo=https%3A%2F%2Fgithub.com%2Fhowmanybots%2Fonebot%2Fraw%2Fmaster%2Fassets%2Flogo-256.png&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)

# OneBot-YaYa

OneBot-YaYa是基于GO和C语言混合编程开发的QQ机器人HTTP API，[OneBot标准](https://github.com/howmanybots/onebot)的先驱QQ机器人平台实现

![Badge](https://img.shields.io/badge/OneBot-v11-black)
[![License](https://img.shields.io/github/license/Yiwen-Chan/OneBot-YaYa.svg)](https://raw.githubusercontent.com/Yiwen-Chan/OneBot-YaYa/master/LICENSE)
[![qq 群](https://img.shields.io/badge/qq%E7%BE%A4-1048452984-green.svg)](https://jq.qq.com/?_wv=1027&k=QMb7x1mM)



### 开始使用

1. 下载 [先驱框架(测试版)](http://api.xianqubot.com/index.php?newver=beta) 与 [OneBot-YaYa](https://github.com/Yiwen-Chan/OneBot-YaYa/releases) 并解压

2. 解压先驱框架到文件夹，并运行[先驱.exe](http://api.xianqubot.com/index.php?newver=beta)，待相关文件生成完毕后，将`OneBot-YaYa.XQ.dll`放入`.\Plugin`，重启先驱框架

3. 切到`账号管理`界面登录机器人账号

<details>
<summary>4. 详细的配置文件说明 .\OneBot\config.yml</summary>

```
# 版本
version: 1.0.5
# 主人QQ号
master: 12345678
# 是否开启DEBUG日志
debug: true
# 心跳设置，默认不动
heratbeat:
  enable: true
  interval: 10000
# 缓存设置，暂未实现
cache:
  database: false
  image: false
  record: false
  video: false
# 不同姬气人的设置，注意yaml中 "-" 代表一个父节点有多个子节点
bots:
# 被设置的姬气人QQ
- bot: 0
  # 正向WS
  websocket:
  # 连接到的服务的名字，自己起
  - name: WSS EXAMPLE
    # 是否启动该服务的连接，连接为 true
    enable: false
    # OneBot建立服务器的HOST，无特殊需求一般为 127.0.0.1
    host: 127.0.0.1
    # OneBot建立服务器的PORT，与插件的端口要对应
    port: 6700
    # OneBot服务器 Token ,一般不动
    access_token: ""
    # OneBot上报格式，可为 string 或 array ，一般不动
    post_message_format: string
  # 反向WS
  websocket_reverse:
  # 连接到的服务的名字，自己起
  - name: WSC EXAMPLE
    # 是否启动该服务的连接，连接为 true
    enable: false
    # 插件服务器的地址，一般只需要改端口
    url: ws://127.0.0.1:8080/ws
    # 暂未实现
    api_url: ws://127.0.0.1:8080/api
    # 暂未实现
    event_url: ws://127.0.0.1:8080/event
    # 暂未实现
    use_universal_client: true
    # 插件填了 Token 这里也要填
    access_token: ""
    # OneBot上报格式，可为 string 或 array ，一般不动
    post_message_format: string
    # 掉线重连的时间间隔，单位毫秒
    reconnect_interval: 3000
  # HTTP 和 HTTP POST
  http:
  # 连接到的服务的名字，自己起
  - name: HTTP EXAMPLE
    # 是否启动该服务的连接，连接为 true
    enable: true
    # OneBot建立服务器的HOST，无特殊需求一般为 127.0.0.1
    host: 127.0.0.1
    # OneBot建立服务器的PORT，与插件的端口要对应
    port: 5700
    # OneBot服务器 Token ,一般不动
    token: ""
    # OneBot 上报的地址，即插件服务器地址
    post_url: 
    # OneBot 上报的 Secret，一般不填
    secret: ""
    # 等待响应时间，一般不动
    time_out: 0
    # OneBot上报格式，可为 string 或 array ，一般不动
    post_message_format: string
```

</details>

5. 每个姬气人都可以设置多个 正向WS 反向WS HTTP 服务，实在不懂[加群](https://jq.qq.com/?_wv=1027&k=PVW9Ol8b)问或者提 [issue](https://github.com/Yiwen-Chan/OneBot-YaYa/issues)

6. 再次重启先驱框架（热重载什么的咕了）

- 注：不要使用`重载插件`功能，否则会导致框架闪退，此为框架与go不兼容问题

### 支持的标准

##### 通信方式

- [x] HTTP
- [x] HTTP POST
- [x] 正向Websocket
- [x] 反向Websocket

<details>
<summary>消息段类型</summary>


- [纯文本](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#纯文本)

  ```
  纯文本内容
  ```

- [QQ表情](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#qq-表情)

  ```
  [CQ:face,id=123]
  ```

- [图片](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#图片)

  ```
  [CQ:image,file=http://baidu.com/1.jpg]
  ```

- [语音](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#语音)

  ```
  [CQ:record,file=http://baidu.com/1.mp3]
  ```

- [@某人](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#@某人)

  ```
  [CQ:at,qq=10001000]
  ```

- [窗口抖动](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#窗口抖动)

  ```
  [CQ:shake]
  ```

- [自定义音乐分享](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#自定义音乐分享)

  ```
  [CQ:music,type=custom,url=http://baidu.com,audio=http://baidu.com/1.mp3,title=音乐标题]
  ```

- [XML消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#xml-消息)

  ```
  [CQ:xml,data=<?xml ...]
  ```

- [JSON消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#json-消息)

  ```
  [CQ:json,data={"app": ...]
  ```

</details>

<details>
<summary>API</summary>



| API                      | 功能                                                         | 备注                                                       |
| ------------------------ | ------------------------------------------------------------ | ------------------------ |
| /send_private_msg        | [发送私聊消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_private_msg-发送私聊消息) |  |
| /send_group_msg          | [发送群消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_group_msg-发送群消息) |  |
| /send_msg                | [发送消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_msg-发送消息) |  |
| /delete_msg | [撤回信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#delete_msg-撤回消息) | 暂未实现 |
| /get_msg | [获取消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_msg-获取消息) | 暂未实现 |
| /get_forward_msg | [获取合并转发消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_forward_msg-获取合并转发消息) | 暂未实现 |
| /send_like | [发送好友赞](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_like-发送好友赞) |  |
| /set_group_kick | [群组踢人](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_kick-群组踢人) |  |
| /set_group_ban | [群组单人禁言](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_ban-群组单人禁言) |  |
| /set_group_anonymous_ban | [群组匿名用户禁言](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous_ban-群组匿名用户禁言) | 暂未实现 |
| /set_group_whole_ban | [群组全员禁言](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_whole_ban-群组全员禁言) |  |
| /set_group_admin         | [群组设置管理员](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_admin-群组设置管理员) | 先驱不支持 |
| /set_group_anonymous     | [群组匿名](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous-群组匿名) |  |
| /set_group_card          | [设置群名片群备注](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_card-设置群名片群备注) |  |
| /set_group_name          | [设置群名](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_name-设置群名) | 先驱不支持 |
| /set_group_leave         | [退出群组](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_leave-退出群组) |  |
| /set_group_special_title | [设置群组专属头衔](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_special_title-设置群组专属头衔) | 先驱不支持 |
| /set_friend_add_request  | [处理加好友请求](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_friend_add_request-处理加好友请求) |  |
| /set_group_add_request   | [处理加群请求/邀请](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_add_request-处理加群请求邀请) |            |
| /get_login_info | [获取登录号信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_login_info-获取登录号信息) |  |
| /get_stranger_info | [获取陌生人信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_stranger_info-获取陌生人信息) | 暂未实现 |
| /get_friend_list         | [获取好友列表](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_friend_list-获取好友列表) |  |
| /get_group_info | [获取群信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_info-获取群信息) | 暂未实现 |
| /get_group_list | [获取群列表](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_list-获取群列表) |  |
| /get_group_member_info | [获取群成员信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_info-获取群成员信息) | 暂未实现 |
| /get_group_member_list | [获取群成员列表](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_list-获取群成员列表) |  |
| /get_group_honor_info | [获取群荣誉信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_honor_info-获取群荣誉信息) | 先驱不支持 |
| /get_cookies | [获取 Cookies](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_cookies-获取-cookies) | 暂未实现 |
| /get_csrf_token | [获取 CSRF Token](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_csrf_token-获取-csrf-token) | 暂未实现 |
| /get_credentials | [获取 QQ 相关接口凭证](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_credentials-获取-qq-相关接口凭证) | 暂未实现 |
| /get_record | [获取语音](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_image-获取语音) | 暂未实现 |
| /get_image | [获取图片](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_image-获取图片) | 暂未实现 |
| /can_send_image | [检查是否可以发送图片](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#can_send_image-检查是否可以发送图片) |  |
| /can_send_record | [检查是否可以发送语音](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#can_send_record-检查是否可以发送语音) |  |
| /get_status | [获取运行状态](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_status-获取运行状态) |  |
| /get_version_info | [获取版本信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_version_info-获取版本信息) |  |
| /set_restart | [重启 onebot 实现](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_restart-重启-onebot-实现) | 暂未实现 |
| /clean_cache | [清理缓存](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#clean_cache-清理缓存) | 暂未实现 |

</details>

<details>
<summary>Event</summary>



| 信息事件                                                     | 备注                  |
| ------------------------------------------------------------ | --------------------- |
| [私聊信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md) | `sender` 字段暂未实现 |
| [群消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md) | `sender` 字段暂未实现 |

| 通知事件                    | 备注                                                         |
| ------------------------ | ------------------------------------------------------------ |
| [群文件上传](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [群管理员变动](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [群成员减少](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [群成员增加](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [群禁言](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [好友添加](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md#好友添加) | |
| [群消息撤回](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md#群消息撤回) | |
| [好友消息撤回](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md#好友消息撤回) | |
| [群内戳一戳](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md#群内戳一戳) | 先驱不支持 |
| [群红包运气王](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md#群红包运气王) | 先驱不支持 |
| [群成员荣誉变更](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md#群成员荣誉变更) | 先驱不支持 |

| 请求事件                     | 备注                                                         |
| ------------------------ | ------------------------------------------------------------ |
| [加好友请求](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/request.md#加好友请求) |  |
| [加群请求/邀请](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/request.md#加群请求邀请) | |

| 元事件                     | 备注                                                         |
| ------------------------ | ------------------------------------------------------------ |
| [生命周期](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md) |  |
| [心跳](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md) |  |

</details>

### OneBot 生态环境

#### SDK／开发框架

| 语言               | 通信方式               | 地址                                                         | 核心作者               |
| ------------------ | ---------------------- | ------------------------------------------------------------ | ---------------------- |
| Python             | HTTP, 反向 WS          | [nonebot/nonebot](https://github.com/nonebot/nonebot)        | richardchien yanyongyu |
| Go                 | 正向 WS                | [wdvxdr1123/ZeroBot](https://github.com/wdvxdr1123/ZeroBot)  | wdvxdr1123             |
| Node.js            | HTTP, 正向 WS, 反向 WS | [koishijs/koishi](https://github.com/koishijs/koishi)        | Shigma                 |
| PHP                | 反向 WS                | [zhamao-robot/zhamao-framework](https://github.com/zhamao-robot/zhamao-framework) | crazywhalecc           |
| C#                 | HTTP, 正向 WS, 反向 WS | [frank-bots/cqhttp.Cyan](https://github.com/frank-bots/cqhttp.Cyan) | frankli0324            |
| Java Kotlin Groovy | 反向 WS                | [lz1998/Spring-CQ](https://github.com/lz1998/Spring-CQ)（[教程](https://www.bilibili.com/video/av89649630/)） | lz1998                 |

[More](https://github.com/Yiwen-Chan/OneBot-YaYa/blob/master/docs/sdk.md)

#### 应用案例

| 项目地址                                                     | 简介或功能                                                   | 依赖                                          | 核心作者                                      |
| ------------------------------------------------------------ | ------------------------------------------------------------ | --------------------------------------------- | --------------------------------------------- |
| [Kyomotoi/ATRI](https://github.com/Kyomotoi/ATRI/tree/master) | 名为 [ATRI](https://atri-mdm.com/) 的BOT                     | [nonebot](https://github.com/nonebot/nonebot) | [Kyomotoi](https://github.com/Kyomotoi)       |
| [fz6m/nonebot-plugin](https://github.com/fz6m/nonebot-plugin) | 各种即开即用、良好兼容的插件                                 | [nonebot](https://github.com/nonebot/nonebot) | [fz6m](https://github.com/fz6m)               |
| [mnixry/coolQPythonBot](https://github.com/mnixry/coolQPythonBot) | 识图识番搜图涩图 番剧查询 B站视频解析 RSS 维基百科 广播 欢迎 一言 嘴臭 身份生成 | [nonebot](https://github.com/nonebot/nonebot) | [mnixry](https://github.com/mnixry)           |
| [cleoold/sendo-erika](https://github.com/cleoold/sendo-erika) | 自定义回复 点歌 签到 运势 谷歌搜索 百度热搜 碧蓝建造         | [nonebot](https://github.com/nonebot/nonebot) | [cleoold](https://github.com/cleoold)         |
| [Quan666/ELF_RSS](https://github.com/Quan666/ELF_RSS)        | RSS                                                          | [nonebot](https://github.com/nonebot/nonebot) | [Quan666](https://github.com/Quan666)         |
| [Bluefissure/OtterBot](https://github.com/Bluefissure/OtterBot) | FF14 玩家BOT                                                 | 云BOT                                         | [Bluefissure](https://github.com/Bluefissure) |

[More](https://github.com/Yiwen-Chan/OneBot-YaYa/blob/master/docs/plugin.md)

### 特别感谢

- [OneBot标准](https://github.com/howmanybots/onebot)
- [CQHTTP](https://github.com/richardchien/coolq-http-api) - [LICENSE](https://github.com/richardchien/coolq-http-api/blob/master/LICENSE)
- [go-cqhttp](https://github.com/Mrs4s/go-cqhttp) - [LICENSE](https://github.com/Mrs4s/go-cqhttp/blob/master/LICENSE)