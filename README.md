![OneBot-YaYa](https://socialify.git.ci/Yiwen-Chan/OneBot-YaYa/image?description=1&descriptionEditable=OneBot%20base%20on%20XQ&font=Inter&logo=https%3A%2F%2Fgithub.com%2Fhowmanybots%2Fonebot%2Fraw%2Fmaster%2Fassets%2Flogo-256.png&owner=1&pattern=Circuit%20Board&theme=Light)

# 兼容性

### 接口
- [ ] HTTP API
- [ ] 正向Websocket
- [x] 反向Websocket

### 实现
<details>
<summary>已实现CQ码</summary>

- [CQ:image]
- [CQ:record]
- [CQ:emoji]
- [CQ:face]
- [CQ:at]
- [CQ:music]
- [CQ:json]
- [CQ:xml]


</details>

<details>
<summary>已实现API</summary>

##### 注意: 部分API实现与CQHTTP原版略有差异，请参考文档
| API                      | 功能                                                         |
| ------------------------ | ------------------------------------------------------------ |
| /send_private_msg        | [发送私聊消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_private_msg-发送私聊消息) |
| /send_group_msg          | [发送群消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_group_msg-发送群消息) |
| /send_msg                | [发送消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_msg-发送消息) |
| /send_like        | [发送好友赞](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_like-发送好友赞) |
| /set_group_kick        | [群组踢人](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_kick-群组踢人) |
| /set_group_ban        | [群组单人禁言](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_ban-群组单人禁言) |
| /set_group_whole_ban        | [群组全员禁言](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_whole_ban-群组全员禁言) |
| /set_group_anonymous        | [群组匿名](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous-群组匿名) |
| /set_group_card        | [设置群名片群备注](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_card-设置群名片群备注) |
| /set_group_leave        | [退出群组](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_leave-退出群组) |
| /get_login_info        | [获取登录号信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_login_info-获取登录号信息) |

</details>

<details>
<summary>已实现Event</summary>

#### 已实现Event
| 信息事件                     | 备注                                                         |
| ------------------------ | ------------------------------------------------------------ |
| [私聊信息](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md) |  |
| [群消息](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md) |  |

| 通知事件                    | 备注                                                         |
| ------------------------ | ------------------------------------------------------------ |
| [群文件上传](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [群管理员变动](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [群成员减少](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [群成员增加](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |
| [群禁言](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/notice.md) |  |

| 请求事件                     | 备注                                                         |
| ------------------------ | ------------------------------------------------------------ |
|  |  |

| 元事件                     | 备注                                                         |
| ------------------------ | ------------------------------------------------------------ |
| [生命周期](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md) |  |
| [心跳](https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md) |  |

</details>

