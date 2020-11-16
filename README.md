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

</details>

<details>
<summary>已实现API</summary>

##### 注意: 部分API实现与CQHTTP原版略有差异，请参考文档
| API                      | 功能                                                         |
| ------------------------ | ------------------------------------------------------------ |
| /send_msg                | [发送消息](https://cqhttp.cc/docs/4.15/#/API?id=send_msg-发送消息) |
| /send_group_msg          | [发送群消息](https://cqhttp.cc/docs/4.15/#/API?id=send_group_msg-发送群消息) |
| /send_private_msg        | [发送私聊消息](https://cqhttp.cc/docs/4.15/#/API?id=send_private_msg-发送私聊消息) |

</details>

<details>
<summary>已实现Event</summary>

##### 注意: 部分Event数据与CQHTTP原版略有差异，请参考文档
| Event                                                        |
| ------------------------------------------------------------ |
| [私聊信息](https://cqhttp.cc/docs/4.15/#/Post?id=私聊消息)   |
| [群消息](https://cqhttp.cc/docs/4.15/#/Post?id=群消息)       |

</details>
