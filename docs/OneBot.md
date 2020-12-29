# OneBot
OneBot标准 是从原 CKYU 平台的 CQHTTP 插件接口修改而来的通用聊天机器人应用接口标准。符合 OneBot标准 的插件能在任意支持 OneBot标准 机器人平台上运行。
## 通信方式
OneBot 标准可通过四种方式与 OneBot插件 进行交互，OneBot-YaYa 中每个 BOT 都需要单独配置通信方式，每个 BOT 可配置多个通信方式。

| 通信方式                           | 对应配置               | 优点                   | 缺点                                      |
| ---------------------------------- | ---------------------- | ---------------------- | ----------------------------------------- |
| 正向WS (OneBot端为WebSocket服务端) | `websocket`            | 高性能                 | 不适合云插件开发                          |
| 反向WS (OneBot端为WebSocket客户端) | `websocket_reverse`    | 高性能，适合云插件开发 | 每个插件服务都需要开启一个WebSocket服务端 |
| HTTP                               | `http`                 | 比较简单               | 同时需要建立监听以及发起请求              |
| HTTP (快速回复)                    | `http` 中的 `post_url` | 最易上手               | 不支持主动发起调用                        |

### 正向WS

| 配置项                          | 默认值        | 说明                   |
| ------------------------------- | ------------- | ---------------------- |
| `websocket.name`                | `WSS EXAMPLE` | 正向WS服务 的 名字     |
| `websocket.enable`              | `false`       | 正向WS服务 的 开关     |
| `websocket.host`                | `127.0.0.1`   | 正向WS服务 的 监听IP   |
| `websocket.port`                | `6700`        | 正向WS服务 的 监听端口 |
| `websocket.access_token`        |               | 正向WS服务 的 Token    |
| `websocket.post_message_format` | `string`      | 正向WS服务 的 上报格式 |

### 反向WS

| 配置项                                   | 默认值                      | 说明                   |
| ---------------------------------------- | --------------------------- | ---------------------- |
| `websocket_reverse.name`                 | `WSC EXAMPLE`               | 反向WS服务 的 名字     |
| `websocket_reverse.enable`               | `false`                     | 反向WS服务 的 开关     |
| `websocket_reverse.url`                  | `ws://127.0.0.1:8080/ws`    | 反向WS服务 的 连接地址 |
| `websocket_reverse.api_url`              | `ws://127.0.0.1:8080/api`   | 暂未实现               |
| `websocket_reverse.event_url`            | `ws://127.0.0.1:8080/event` | 暂未实现               |
| `websocket_reverse.use_universal_client` | `true`                      | 暂未实现               |
| `websocket_reverse.access_token`         |                             | 反向WS服务 的 Token    |
| `websocket_reverse.post_message_format`  | `string`                    | 反向WS服务 的 上报格式 |
| `websocket_reverse.reconnect_interval`   | `3000`                      | 反向WS服务 的 重连间隔 |

### HTTP

| 配置项                          | 默认值        | 说明                    |
| ------------------------------- | ------------- | ----------------------- |
| `http.name`                     | `WSS EXAMPLE` | HTTP服务 的 名字        |
| `http.enable`                   | `false`       | HTTP服务 的 开关        |
| `http.host`                     | `127.0.0.1`   | HTTP服务 的 监听IP      |
| `http.port`                     | `5700`        | HTTP服务 的 监听端口    |
| `http.access_token`             |               | HTTP服务 的 Token       |
| `http.post_url`                 |               | HTTP服务 的 上报地址    |
| `http.secret`                   |               | HTTP服务 的 上报 Secret |
| `http.time_out`                 | `0`           | HTTP服务 的 上报超时    |
| `websocket.post_message_format` | `string`      | HTTP服务 的 上报格式    |

### HTTP (快速回复)

| 配置项                          | 默认值        | 说明                    |
| ------------------------------- | ------------- | ----------------------- |
| `http.name`                     | `WSS EXAMPLE` | HTTP服务 的 名字        |
| `http.enable`                   | `false`       | HTTP服务 的 开关        |
| `http.host`                     | `127.0.0.1`   | 不填                    |
| `http.port`                     | `5700`        | 不填                    |
| `http.access_token`             |               | 不填                    |
| `http.post_url`                 |               | HTTP服务 的 上报地址    |
| `http.secret`                   |               | HTTP服务 的 上报 Secret |
| `http.time_out`                 | `0`           | HTTP服务 的 上报超时    |
| `websocket.post_message_format` | `string`      | HTTP服务 的 上报格式    |

