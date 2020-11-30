[主页](https://github.com/howmanybots/onebot/blob/master/README.md)　[生态](https://github.com/howmanybots/onebot/blob/master/ecosystem.md)　[更新日志](https://github.com/howmanybots/onebot/blob/master/changelog.md)　标准版本：[v11](https://github.com/howmanybots/onebot/blob/master/v11/specs/README.md)　[v10](https://github.com/howmanybots/onebot/blob/master/v10/specs/README.md)

![OneBot-YaYa](https://socialify.git.ci/Yiwen-Chan/OneBot-YaYa/image?description=1&descriptionEditable=OneBot%20(CQHTTP)%20%E5%85%88%E9%A9%B1%E5%B9%B3%E5%8F%B0%E7%9A%84%E5%AE%9E%E7%8E%B0&font=Raleway&forks=1&issues=1&logo=https%3A%2F%2Fgithub.com%2Fhowmanybots%2Fonebot%2Fraw%2Fmaster%2Fassets%2Flogo-256.png&owner=1&pattern=Plus&pulls=1&stargazers=1&theme=Light)

# OneBot-YaYa

OneBot-YaYa是基于GO和C语言混合编程开发的QQ机器人HTTP API，在先驱QQ机器人平台实现了OneBot协议标准

![Badge](https://img.shields.io/badge/OneBot-v11-black)
[![License](https://img.shields.io/github/license/Yiwen-Chan/OneBot-YaYa.svg)](https://raw.githubusercontent.com/Yiwen-Chan/OneBot-YaYa/master/LICENSE)
[![QQ 群](https://img.shields.io/badge/QQ %E7%BE%A4-1048452984-green.svg)](https://jq.qq.com/?_wv=1027&k=QMb7x1mM)

### 支持的通信方式
- [x] HTTP
- [x] HTTP POST
- [x] 正向Websocket
- [x] 反向Websocket

### 支持的协议标准
<details>
<summary>消息段类型</summary>


- 纯文本

  ```
  纯文本内容
  ```

- QQ表情

  ```
  [CQ:face,id=123]
  ```

- 图片

  ```
  [CQ:image,file=http://baidu.com/1.jpg]
  ```

- 语音

  ```
  [CQ:record,file=http://baidu.com/1.mp3]
  ```

- @某人

  ```
  [CQ:at,qq=10001000]
  ```

- 窗口抖动

  ```
  [CQ:shake]
  ```

- 自定义音乐分享

  ```
  [CQ:music,type=custom,url=http://baidu.com,audio=http://baidu.com/1.mp3,title=音乐标题]
  ```

- XML消息

  ```
  [CQ:xml,data=<?xml ...]
  ```

- JSON消息

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

### 特别感谢

- [OneBot标准](https://github.com/howmanybots/onebot)
- [CQHTTP](https://github.com/richardchien/coolq-http-api) - [LICENSE](https://github.com/richardchien/coolq-http-api/blob/master/LICENSE)
- [go-cqhttp](https://github.com/Mrs4s/go-cqhttp) - [LICENSE](https://github.com/Mrs4s/go-cqhttp/blob/master/LICENSE)