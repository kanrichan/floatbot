package onebot

import "yaya/core"

// 初始化顺序
// core.C -> core.Go -> onebot.main(常量变量) -> onebot.main.init() -> onebot.Main -> app.go

// OneBot-YaYa 插件 生命周期
// onebot.onStart -> onebot.onEvent -> onebot.onDisable

// onebot.onStart 启动 生命周期
// 全局变量定义 -> 配置文件加载 -> go http|wss|wsc 挥手 -> go 监听服务 -> go 上报服务

// onebot.onEvent 事件 生命周期
// XQ -> CORE -> XQEvent() -> event.onEvent() -> event.Push() -> http|wss|wsc -> onebot标准插件 -> http|wss|wsc -> api.CallApi() -> api.sendMsg() -> CORE -> XQ

// 将package onebot里的函数绑定到package core
func init() {
	core.Create = XQCreate
	core.Event = XQEvent
	core.DestroyPlugin = XQDestroyPlugin
	core.SetUp = XQSetUp
}

// Main DLL加载到内存后，go初始化完毕，调用此函数
func Main() {
}

// AppInfoJson 全局插件信息
var AppInfoJson string

// XQCreate core里的XQCreate，XQ加载插件信息调用此函数
func XQCreate(version string) string {
	return AppInfoJson
}

// XQDestroyPlugin core里的XQDestroyPlugin，XQ界面点击卸载调用此函数
func XQDestroyPlugin() int64 {
	return 0
}

// XQSetUp core里的XQSetUp，XQ界面点击设置调用此函数
func XQSetUp() int64 {
	return 0
}

// XEvent 先驱事件的结构体
type XEvent struct {
	ID          int64  `db:"id"`           // ID 消息ID 唯一对应
	SelfID      int64  `db:"self_id"`      // SelfID 机器人QQ, 多Q版用于判定哪个QQ接收到该消息
	MessageType int64  `db:"message_type"` // MessageType 消息类型, 接收到消息类型，该类型可在常量表中查询具体定义，此处仅列举： -1 未定义事件 0,在线状态临时会话 1,好友信息 2,群信息 3,讨论组信息 4,群临时会话 5,讨论组临时会话 6,财付通转账 7,好友验证回复会话
	SubType     int64  `db:"sub_type"`     // SubType 消息子类型, 此参数在不同消息类型下，有不同的定义，暂定：接收财付通转账时 1为好友 4为群临时会话 5为讨论组临时会话    有人请求入群时，不良成员这里为1
	GroupID     int64  `db:"group_id"`     // GroupID 消息来源, 此消息的来源，如：群号、讨论组ID、临时会话QQ、好友QQ等
	UserID      int64  `db:"user_id"`      // UserID 触发对象_主动, 主动发送这条消息的QQ，踢人时为踢人管理员QQ
	NoticeID    int64  `db:"notice_id"`    // NoticeID 触发对象_被动, 被动触发的QQ，如某人被踢出群，则此参数为被踢出人QQ
	Message     string `db:"message"`      // Message 消息内容, 此参数有多重含义，常见为：对方发送的消息内容，但当消息类型为 某人申请入群，则为入群申请理由
	MessageNum  int64  `db:"message_num"`  // MessageNum 消息序号, 此参数暂定用于消息回复，消息撤回
	MessageID   int64  `db:"message_id"`   // MessageID 消息ID, 此参数暂定用于消息回复，消息撤回
	RawMessage  string `db:"raw_message"`  // RawMessage 原始信息, UDP收到的原始信息，特殊情况下会返回JSON结构（入群事件时，这里为该事件seq）
	Time        int64  `db:"time"`         // Time 消息时间戳, 接受到消息的时间戳
	Ret         int64  `db:"ret"`          // Ret 回传文本指针, 此参数用于插件加载拒绝理由
}

func XQEvent(selfID, messageType, subType, groupID, userID, noticeID int64, message string, messageNum, messageID int64, rawMessage string, time, ret int64) int64 {
	switch messageType {
	default:
		// TODO 初始化完毕后再处理消息等事件
		if !FirstStart {
			go ProtectRun(func() {
				onEvent(&XEvent{
					ID:          0,
					SelfID:      selfID,
					MessageType: messageType,
					SubType:     subType,
					GroupID:     groupID,
					UserID:      userID,
					NoticeID:    noticeID,
					Message:     message,
					MessageNum:  messageNum,
					MessageID:   messageID,
					RawMessage:  rawMessage,
					Time:        time,
					Ret:         ret,
				})
			}, "onEvent()")
		}
	case 12001:
		go ProtectRun(func() { onStart() }, "onStart()")
	case 12002:
		go ProtectRun(func() { onDisable() }, "onDisable()")
	}
	return 0
}

var (
	FirstStart bool = true // 首次启动标志

	// 各种媒体的路径
	XQPath     = PathExecute()
	AppPath    = XQPath + "OneBot/"
	ImagePath  = XQPath + "OneBot/image/"
	RecordPath = XQPath + "OneBot/record/"
	VideoPath  = XQPath + "OneBot/video/"
	CachePath  = XQPath + "OneBot/cache/"

	PicPool = PicsCache{Max: 1000} // 已经上传过的图片池子
)

// onStart 插件加载完毕，接收到框架初始化命令后调用此函数
func onStart() {
	if FirstStart {
		// 注册 XQApi ，建立API名与函数的映射关系
		apiMap.Register(&apiMap.this)
		// TODO 创建各种媒体的路径
		CreatePath(AppPath)
		CreatePath(ImagePath)
		CreatePath(RecordPath)
		CreatePath(VideoPath)

		INFO("夜夜は世界一かわいい")
		// TODO 加载配置文件并初始化
		Conf = Load(AppPath + "config.yml")
		if Conf == nil {
			ERROR("初始化失败，夜夜去睡觉了，晚安~")
			return
		}
		// TODO 启动数据库
		go Conf.runDB()
		// TODO 启动 http|wss|wsc
		go Conf.runOnebot()
		// 启动心跳
		go Conf.heartBeat()
	}
	FirstStart = false
}

// onDisable 插件加载完毕，XQ界面点击停止调用此函数
func onDisable() {
}
