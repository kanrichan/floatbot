# 消息事件

<details>
<summary>目录</summary>
<p>

- [私聊消息](#私聊消息)
- [群消息](#群消息)
- [群文件上传](#群文件上传)
- [群管理员变动](#群管理员变动)
- [群成员减少](#群成员减少)
- [群成员增加](#群成员增加)
- [群禁言](#群禁言)
- [好友添加](#好友添加)
- [群消息撤回](#群消息撤回)
- [好友消息撤回](#好友消息撤回)

</p>
</details>

## 私聊消息

### 事件数据

| 字段名 | 数据类型 | 可能的值 | 说明 |
| ----- | ------- | ------- | ---- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type` | string | `message` | 上报类型 |
| `message_type` | string | `private` | 消息类型 |
| `sub_type` | string | `friend`、`group`、`other` | 消息子类型，如果是好友则是 `friend`，如果是群临时会话则是 `group` |
| `message_id` | number (int32) | - | 先驱暂时不实现 |
| `user_id` | number (int64) | - | 发送者 QQ 号 |
| `message` | message | - | 消息内容 |
| `raw_message` | string | - | 原始消息内容 |
| `font` | number (int32) | - | 先驱暂时不实现 |
| `sender` | object | - | 发送人信息 |

其中 `sender` 字段的内容如下：

| 字段名 | 数据类型 | 说明 |
| ----- | ------ | ---- |
| `user_id` | number (int64) | 发送者 QQ 号 |
| `nickname` | string | 先驱暂时不实现 |
| `sex` | string | 先驱暂时不实现 |
| `age` | number (int32) | 先驱暂时不实现 |

需要注意的是，`sender` 中的各字段是尽最大努力提供的，也就是说，不保证每个字段都一定存在，也不保证存在的字段都是完全正确的（缓存可能过期）。

## 群消息

### 事件数据

| 字段名 | 数据类型 | 可能的值 | 说明 |
| ----- | ------- | ------- | --- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type` | string | `message` | 上报类型 |
| `message_type` | string | `group` | 消息类型 |
| `sub_type` | string | `normal`、`anonymous`、`notice` | 消息子类型，正常消息是 `normal`，匿名消息是 `anonymous`，系统提示（如「管理员已禁止群内匿名聊天」）是 `notice` |
| `message_id` | number (int32) | - | 先驱暂时不实现 |
| `group_id` | number (int64) | - | 群号 |
| `user_id` | number (int64) | - | 发送者 QQ 号 |
| `anonymous` | object | - | 先驱暂时不实现 |
| `message` | message | - | 消息内容 |
| `raw_message` | string | - | 原始消息内容 |
| `font` | number (int32) | - | 先驱暂时不实现 |
| `sender` | object | - | 发送人信息 |

其中 `anonymous` 字段的内容如下：

| 字段名 | 数据类型 | 说明 |
| ----- | ------ | ---- |
| `id` | number (int64) | 先驱暂时不实现 |
| `name` | string | 先驱暂时不实现 |
| `flag` | string | 先驱暂时不实现 |

`sender` 字段的内容如下：

| 字段名 | 数据类型 | 说明 |
| ----- | ------ | ---- |
| `user_id` | number (int64) | 发送者 QQ 号 |
| `nickname` | string | 先驱暂时不实现 |
| `card` | string | 先驱暂时不实现 |
| `sex` | string | 先驱暂时不实现 |
| `age` | number (int32) | 先驱暂时不实现 |
| `area` | string | 先驱暂时不实现 |
| `level` | string | 先驱暂时不实现 |
| `role` | string | 先驱暂时不实现 |
| `title` | string | 先驱暂时不实现 |

需要注意的是，`sender` 中的各字段是尽最大努力提供的，也就是说，不保证每个字段都一定存在，也不保证存在的字段都是完全正确的（缓存可能过期）。尤其对于匿名消息，此字段不具有参考价值。

## 群文件上传

### 事件数据

| 字段名 | 数据类型 | 可能的值 | 说明 |
| ----- | ------ | ------- | ---- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type` | string | `notice` | 上报类型 |
| `notice_type` | string | `group_upload` | 通知类型 |
| `group_id` | number (int64) | - | 群号 |
| `user_id` | number (int64) | - | 发送者 QQ 号 |
| `file` | object | - | 文件信息 |

其中 `file` 字段的内容如下：

| 字段名 | 数据类型 | 说明 |
| ----- | ------ | ---- |
| `id` | string | 先驱暂时不实现 |
| `name` | string | 文件名 |
| `size` | number (int64) | 先驱暂时不实现 |
| `busid` | number (int64) | 先驱暂时不实现 |

## 群管理员变动

### 事件数据

| 字段名 | 数据类型 | 可能的值 | 说明 |
| ----- | ------ | -------- | --- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type` | string | `notice` | 上报类型 |
| `notice_type` | string | `group_admin` | 通知类型 |
| `sub_type` | string | `set`、`unset` | 事件子类型，分别表示设置和取消管理员 |
| `group_id` | number (int64) | - | 群号 |
| `user_id` | number (int64) | - | 管理员 QQ 号 |

## 群成员减少

### 事件数据

| 字段名 | 数据类型 | 可能的值 | 说明 |
| ----- | ------ | -------- | --- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type` | string | `notice` | 上报类型 |
| `notice_type` | string | `group_decrease` | 通知类型 |
| `sub_type` | string | `leave`、`kick`、`kick_me` | 事件子类型，分别表示主动退群、成员被踢、登录号被踢 |
| `group_id` | number (int64) | - | 群号 |
| `operator_id` | number (int64) | - | 操作者 QQ 号（如果是主动退群，则和 `user_id` 相同） |
| `user_id` | number (int64) | - | 离开者 QQ 号 |

## 群成员增加

### 事件数据

| 字段名 | 数据类型 | 可能的值 | 说明 |
| ----- | ------ | -------- | --- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type` | string | `notice` | 上报类型 |
| `notice_type` | string | `group_increase` | 通知类型 |
| `sub_type` | string | `approve`、`invite` | 事件子类型，分别表示管理员已同意入群、管理员邀请入群 |
| `group_id` | number (int64) | - | 群号 |
| `operator_id` | number (int64) | - | 操作者 QQ 号 |
| `user_id` | number (int64) | - | 加入者 QQ 号 |

## 群禁言

### 事件数据

| 字段名 | 数据类型 | 可能的值 | 说明 |
| ----- | ------ | -------- | --- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type` | string | `notice` | 上报类型 |
| `notice_type` | string | `group_ban` | 通知类型 |
| `sub_type` | string | `ban`、`lift_ban` | 事件子类型，分别表示禁言、解除禁言 |
| `group_id` | number (int64) | - | 群号 |
| `operator_id` | number (int64) | - | 操作者 QQ 号 |
| `user_id` | number (int64) | - | 被禁言 QQ 号 |
| `duration` | number (int64) | - | 先驱暂时不实现 |

## 好友添加

### 事件数据

| 字段名 | 数据类型 | 可能的值 | 说明 |
| ----- | ------ | -------- | --- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type` | string | `notice` | 上报类型 |
| `notice_type` | string | `friend_add` | 通知类型 |
| `user_id` | number (int64) | - | 新添加好友 QQ 号 |

## 群消息撤回

### 事件数据

| 字段名          | 数据类型   | 可能的值       | 说明           |
| ------------- | ------ | -------------- | -------------- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type`   | string | `notice`       | 上报类型       |
| `notice_type` | string | `group_recall` | 通知类型       |
| `group_id`    | number (int64)  |                | 群号           |
| `user_id`     | number (int64)  |                | 消息发送者 QQ 号   |
| `operator_id` | number (int64)  |                | 操作者 QQ 号  |
| `message_id`  | number (int64)  |                | 先驱暂时不实现 |

## 好友消息撤回

### 事件数据

| 字段名          | 数据类型   | 可能的值       | 说明           |
| ------------- | ------ | -------------- | -------------- |
| `time` | number (int64) | - | 事件发生的时间戳 |
| `self_id` | number (int64) | - | 收到事件的机器人 QQ 号 |
| `post_type`   | string | `notice`       | 上报类型       |
| `notice_type` | string | `friend_recall`| 通知类型       |
| `user_id`     | number (int64)  |                | 好友 QQ 号        |
| `message_id`  | number (int64)  |                | 先驱暂时不实现 |

