package core

//#include <xqapi.h>
import "C"

// 在使用本类方法前必须调用本函数(返回框架版本号)
func ApiInit() bool {
	return GoBool(C.S3_Api_ApiInit())
}

// 在没置入此ID前所有接口调用都不会有效果(除了取框架版本与标记异常流程)
// id  id  整数型  无
// addr  addr  整数型  无
func SetAuthId(id int64, addr int64) {
	C.S3_Api_SetAuthId(
		C.int(id), C.int(addr),
	)
}

// 取得好友列表，返回获取到的原始JSON格式信息，需自行解析，http模式
// selfID  响应QQ  文本型  机器人QQ
func GetFriendList(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetFriendList(
		GoInt2CStr(selfID),
	))
}

// 取机器人在线账号列表
func GetOnLineList() string {
	return CPtr2GoStr(C.S3_Api_GetOnLineList())
}

// 取机器人账号是否在线
// selfID  响应QQ  文本型  机器人QQ
func Getbotisonline(selfID int64) bool {
	return GoBool(C.S3_Api_Getbotisonline(
		GoInt2CStr(selfID),
	))
}

// 取群员列表
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取群成员列表群号
func GetGroupMemberList(selfID int64, groupID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupMemberList(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

// 取群成员名片
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取群成员列表群号
// userID  对象QQ  文本型  无
func GetGroupCard(selfID int64, groupID int64, userID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupCard(
		GoInt2CStr(selfID), GoInt2CStr(groupID), GoInt2CStr(userID),
	))
}

// 发送消息
// selfID  响应QQ  文本型  机器人QQ
// messageType  信息类型  整数型  0在线临时会话 1好友 2群 3讨论组 4群临时会话 5讨论组临时会话 7好友验证回复会话
// groupID  收信对象群_讨论组  文本型  发送群信息、讨论组、群或讨论组临时会话信息时填写，如发送对象为好友或信息类型是0时可空
// userID  收信QQ  文本型  收信对象QQ
// message  内容  文本型  信息内容
// bubble  气泡ID  整数型  已支持请自己测试
func SendMsg(selfID int64, messageType int64, groupID int64, userID int64, message string, bubble int64) {
	C.S3_Api_SendMsg(
		GoInt2CStr(selfID), C.int(messageType), GoInt2CStr(groupID), GoInt2CStr(userID), CString(message), C.int(bubble),
	)
}

// 上传图片
// selfID  响应QQ  文本型  机器人QQ
// postType  上传类型  整数型  1好友、临时会话  2群、讨论组 Ps：好友临时会话用类型 1，群讨论组用类型 2；当填写错误时，图片GUID发送不会成功
// userID  参考对象  文本型  上传该图片所属的群号或QQ
// file  图片数据  字节集  图片字节集数据
func UpLoadPic(selfID int64, postType int64, userID int64, file []byte) string {
	return CPtr2GoStr(C.S3_Api_UpLoadPic(
		GoInt2CStr(selfID), C.int(postType), GoInt2CStr(userID), CByte(file),
	))
}

// 取群管理员列表
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取管理员列表群号
func GetGroupAdmin(selfID int64, groupID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupAdmin(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

// 群禁言
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲操作的群号
// userID  对象QQ  文本型  欲禁言的对象，如留空且机器人QQ为管理员，将设置该群为全群禁言
// time  时间  整数型  0为解除禁言 （禁言单位为秒），如为全群禁言，参数为非0，解除全群禁言为0
func ShutUP(selfID int64, groupID int64, userID int64, time int64) {
	C.S3_Api_ShutUP(
		GoInt2CStr(selfID), GoInt2CStr(groupID), GoInt2CStr(userID), C.int(time),
	)
}

// 修改群成员昵称
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  对象所处群号
// userID  对象QQ  文本型  被修改名片人QQ
// card  名片  文本型  需要修改的名片
func SetGroupCard(selfID int64, groupID int64, userID int64, card string) bool {
	return GoBool(C.S3_Api_SetGroupCard(
		GoInt2CStr(selfID), GoInt2CStr(groupID), GoInt2CStr(userID), CString(card),
	))
}

// 群删除成员
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  被执行群号
// userID  对象QQ  文本型  被执行对象
// reject_add_request  不在允许  逻辑型  真为不再接收，假为接收
func KickGroupMBR(selfID int64, groupID int64, userID int64, reject_add_request bool) {
	C.S3_Api_KickGroupMBR(
		GoInt2CStr(selfID), GoInt2CStr(groupID), GoInt2CStr(userID), CBool(reject_add_request),
	)
}

// 获取群通知
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取得公告的群号
func GetNotice(selfID int64, groupID int64) string {
	return CPtr2GoStr(C.S3_Api_GetNotice(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

//  取群成员禁言状态 -1失败 0未被禁言 1被单独禁言 2开启了全群禁言
// selfID  响应QQ  文本型  无
// groupID  群号  文本型  无
// userID  对象QQ  文本型  无
func IsShutUp(selfID int64, groupID int64, userID int64) int64 {
	return int64(C.S3_Api_IsShutUp(
		GoInt2CStr(selfID), GoInt2CStr(groupID), GoInt2CStr(userID),
	))
}

// 取是否好友
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
func IfFriend(selfID int64, userID int64) bool {
	return GoBool(C.S3_Api_IfFriend(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 修改QQ在线状态
// selfID  响应QQ  文本型  无
// type  类型  整数型  1、我在线上 2、Q我吧 3、离开 4、忙碌 5、请勿打扰 6、隐身 7、修改昵称 8、修改个性签名 9、修改性别
// text  修改内容  文本型  类型为7和8时填写修改内容  类型9时“1”为男 “2”为女      其他填“”
func SetRInf(selfID int64, subType int64, text string) {
	C.S3_Api_SetRInf(
		GoInt2CStr(selfID), C.int(subType), CString(text),
	)
}

// 取得QQ群页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
func GetGroupPsKey(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupPsKey(
		GoInt2CStr(selfID),
	))
}

// 取得QQ空间页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
func GetZonePsKey(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetZonePsKey(
		GoInt2CStr(selfID),
	))
}

// 取得机器人网页操作用的Cookies
// selfID  响应QQ  文本型  机器人QQ
func GetCookies(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetCookies(
		GoInt2CStr(selfID),
	))
}

// 发布群公告
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲发布公告的群号
// title  标题  文本型  公告标题
// message  内容  文本型  公告内容
func PBGroupNotic(selfID int64, groupID int64, title string, message string) bool {
	return GoBool(C.S3_Api_PBGroupNotic(
		GoInt2CStr(selfID), GoInt2CStr(groupID), CString(title), CString(message),
	))
}

// 撤回群消息
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需撤回消息群号
// messageNum  消息序号  文本型  需撤回消息序号
// messageID  消息ID  文本型  需撤回消息ID
func WithdrawMsg(selfID int64, groupID int64, messageNum int64, messageID int64) string {
	return CPtr2GoStr(C.S3_Api_WithdrawMsg(
		GoInt2CStr(selfID), GoInt2CStr(groupID), GoInt2CStr(messageNum), GoInt2CStr(messageID),
	))
}

// 输出一行日志
// message  内容  文本型  任意想输出的文本格式信息
func OutPutLog(message string) {
	C.S3_Api_OutPutLog(
		CString(message),
	)
}

// 提取图片中的文字
// selfID  响应QQ  文本型  机器人QQ
// file  图片数据  字节集  图片数据
func OcrPic(selfID int64, file []byte) string {
	return CPtr2GoStr(C.S3_Api_OcrPic(
		GoInt2CStr(selfID), CByte(file),
	))
}

// 主动加群
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲申请加入的群号
// reason  理由  文本型  附加理由，可留空（需回答正确问题时，请填写问题答案）
func JoinGroup(selfID int64, groupID int64, reason string) {
	C.S3_Api_JoinGroup(
		GoInt2CStr(selfID), GoInt2CStr(groupID), CString(reason),
	)
}

// 点赞
// selfID  响应QQ  文本型  机器人QQ
// userID  被赞QQ  文本型  填写被赞人QQ
func UpVote(selfID int64, userID int64) string {
	return CPtr2GoStr(C.S3_Api_UpVote(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 通过列表或群临时通道点赞
// selfID  响应QQ  文本型  机器人QQ
// userID  被赞QQ  文本型  填写被赞人QQ
func UpVote_temp(selfID int64, userID int64) string {
	return CPtr2GoStr(C.S3_Api_UpVote_temp(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 获取赞数量
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
func GetObjVote(selfID int64, userID int64) int64 {
	return int64(C.S3_Api_GetObjVote(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 取插件是否启用
func IsEnable() bool {
	return GoBool(C.S3_Api_IsEnable())
}

// 置好友添加请求
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  申请入群 被邀请人 请求添加好友人的QQ （当请求类型为214时这里为邀请人QQ）
// approve  处理方式  整数型  10同意 20拒绝 30忽略 40同意单项好友的请求
// remark  附加信息  文本型  拒绝入群，拒绝添加好友 附加信息
func HandleFriendEvent(selfID int64, userID int64, approve int64, remark string) {
	C.S3_Api_HandleFriendEvent(
		GoInt2CStr(selfID), GoInt2CStr(userID), C.int(approve), CString(remark),
	)
}

// 置群请求
// selfID  响应QQ  文本型  机器人QQ
// sub_type  请求类型  整数型  213请求入群  214我被邀请加入某群  215某人被邀请加入群  101某人请求添加好友
// userID  对象QQ  文本型  申请入群 被邀请人 请求添加好友人的QQ （当请求类型为214时这里为邀请人QQ）
// groupID  群号  文本型  收到请求群号（好友添加时这里请为空）
// flag  seq  文本型  需要处理事件的seq
// approve  处理方式  整数型  10同意 20拒绝 30忽略
// remark  附加信息  文本型  拒绝入群，拒绝添加好友 附加信息
func HandleGroupEvent(selfID int64, sub_type int64, userID int64, groupID int64, flag int64, approve int64, remark string) {
	C.S3_Api_HandleGroupEvent(
		GoInt2CStr(selfID), C.int(sub_type), GoInt2CStr(userID), GoInt2CStr(groupID), GoInt2CStr(flag), C.int(approve), CString(remark),
	)
}

// 取所有QQ列表
func GetQQList() string {
	return CPtr2GoStr(C.S3_Api_GetQQList())
}

// 向框架添加一个QQ
// account  帐号  文本型  无
// password  密码  文本型  无
// enable  自动登录  逻辑型  真 为自动登录
func AddQQ(account string, password string, enable bool) string {
	return CPtr2GoStr(C.S3_Api_AddQQ(
		CString(account), CString(password), CBool(enable),
	))
}

// 登录指定QQ
// userID  登录QQ  文本型  无
func LoginQQ(userID string) {
	C.S3_Api_LoginQQ(
		CString(userID),
	)
}

// 离线指定QQ
// selfID  响应QQ  文本型  无
func OffLineQQ(selfID int64) {
	C.S3_Api_OffLineQQ(
		GoInt2CStr(selfID),
	)
}

// 删除指定QQ
// selfID  响应QQ  文本型  无
func DelQQ(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_DelQQ(
		GoInt2CStr(selfID),
	))
}

// 删除指定好友
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  被删除对象
func DelFriend(selfID int64, userID int64) bool {
	return GoBool(C.S3_Api_DelFriend(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 取QQ昵称
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  欲取得的QQ的号码
func GetNick(selfID int64, userID int64) string {
	return CPtr2GoStr(C.S3_Api_GetNick(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 取好友备注姓名
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  需获取对象好友QQ
func GetFriendsRemark(selfID int64, userID int64) string {
	return CPtr2GoStr(C.S3_Api_GetFriendsRemark(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 取短Clientkey
// selfID  响应QQ  文本型  机器人QQ
func GetClientkey(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetClientkey(
		GoInt2CStr(selfID),
	))
}

// 取bkn
// selfID  响应QQ  文本型  机器人QQ
func GetBkn(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetBkn(
		GoInt2CStr(selfID),
	))
}

// 修改好友备注名称
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  需获取对象好友QQ
// remark  备注  文本型  需要修改的备注姓名
func SetFriendsRemark(selfID int64, userID int64, remark string) {
	C.S3_Api_SetFriendsRemark(
		GoInt2CStr(selfID), GoInt2CStr(userID), CString(remark),
	)
}

// 邀请好友加入群
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  被邀请加入的群号
// userID  对象QQ  文本型  被邀请人QQ号码
func InviteGroup(selfID int64, groupID int64, userID int64) {
	C.S3_Api_InviteGroup(
		GoInt2CStr(selfID), GoInt2CStr(groupID), GoInt2CStr(userID),
	)
}

// 邀请群成员加入群
// selfID  响应QQ  文本型  机器人QQ
// targetID  群号  文本型  邀请到哪个群
// groupID  所在群  文本型  被邀请成员所在群
// userID  邀请QQ  文本型  被邀请人的QQ
func InviteGroupMember(selfID int64, groupID int64, targetID int64, userID int64) bool {
	return GoBool(C.S3_Api_InviteGroupMember(
		GoInt2CStr(selfID), GoInt2CStr(targetID), GoInt2CStr(groupID), GoInt2CStr(userID),
	))
}

// 创建群 组包模式
// selfID  响应QQ  文本型  机器人
func CreateDisGroup(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_CreateDisGroup(
		GoInt2CStr(selfID),
	))
}

// 创建群 群官网Http模式
// selfID  响应QQ  文本型  机器人QQ
// nickname  群昵称  文本型  预创建的群名称
func CreateGroup(selfID int64, nickname string) string {
	return CPtr2GoStr(C.S3_Api_CreateGroup(
		GoInt2CStr(selfID), CString(nickname),
	))
}

// 退出群
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲退出的群号
func QuitGroup(selfID int64, groupID int64) {
	C.S3_Api_QuitGroup(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	)
}

// 封包模式获取群号列表(最多可以取得999)
// selfID  响应QQ  文本型  机器人QQ
func GetGroupList(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupList(
		GoInt2CStr(selfID),
	))
}

// 封包模式获取群号列表(最多可以取得999)
// selfID  响应QQ  文本型  机器人QQ
func GetGroupList_B(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupList_B(
		GoInt2CStr(selfID),
	))
}

// 封包模式取好友列表(与封包模式取群列表同源)
// selfID  响应QQ  文本型  机器人QQ
func GetFriendList_B(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetFriendList_B(
		GoInt2CStr(selfID),
	))
}

// 取登录二维码base64
// key  key  字节集  无
func GetQrcode(key []byte) string {
	return CPtr2GoStr(C.S3_Api_GetQrcode(
		CByte(key),
	))
}

// 检查登录二维码状态
// key  key  字节集  返回数值
func CheckQrcode(key []byte) int64 {
	return int64(C.S3_Api_CheckQrcode(
		CByte(key),
	))
}

// 取指定的群名称
// selfID  响应QQ  文本型  无
// groupID  群号  文本型  无
func GetGroupName(selfID int64, groupID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupName(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

// 取群人数上线与当前人数 换行符分隔
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需查询的群号
func GetGroupMemberNum(selfID int64, groupID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupMemberNum(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

// 取群等级
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需查询的群号
func GetGroupLv(selfID int64, groupID int64) int64 {
	return int64(C.S3_Api_GetGroupLv(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

// 屏蔽或接收某群消息
// selfID  响应QQ  文本型  无
// groupID  群号  文本型  无
// enable  类型  逻辑型  真 为屏蔽接收  假为接收并提醒
func SetShieldedGroup(selfID int64, groupID int64, enable bool) {
	C.S3_Api_SetShieldedGroup(
		GoInt2CStr(selfID), GoInt2CStr(groupID), CBool(enable),
	)
}

// 取群成员列表
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取群成员列表群号
func GetGroupMemberList_B(selfID int64, groupID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupMemberList_B(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

// 封包模式取群成员列表返回重组后的json文本
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取群成员列表群号
func GetGroupMemberList_C(selfID int64, groupID int64) string {
	return CPtr2GoStr(C.S3_Api_GetGroupMemberList_C(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

// 检查指定QQ是否在线
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  需获取对象QQ
func IsOnline(selfID int64, userID int64) bool {
	return GoBool(C.S3_Api_IsOnline(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 取机器人账号在线信息
// selfID  响应QQ  文本型  无
func GetRInf(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetRInf(
		GoInt2CStr(selfID),
	))
}

// 多功能删除好友 可删除陌生人或者删除为单项好友
// selfID  响应QQ  文本型  机器人QQ
// userID  目标QQ  文本型  欲操作的目标
// sub_type  删除类型  整数型  1为在对方的列表删除我 2为在我的列表删除对方
func DelFriend_A(selfID int64, userID int64, sub_type int64) {
	C.S3_Api_DelFriend_A(
		GoInt2CStr(selfID), GoInt2CStr(userID), C.int(sub_type),
	)
}

// 设置机器人被添加好友时的验证方式
// selfID  响应QQ  文本型  机器人QQ
// sub_type  验证类型  整数型  0 允许任何人 1 需要验证消息 2不允许任何人 3需要回答问题 4需要回答问题并由我确认
func Setcation(selfID int64, sub_type int64) {
	C.S3_Api_Setcation(
		GoInt2CStr(selfID), C.int(sub_type),
	)
}

// 设置机器人被添加好友时的问题与答案
// selfID  响应QQ  文本型  机器人QQ
// question  设置问题  文本型  设置的问题
// answer  问题答案  文本型  设置的问题答案
func Setcation_problem_A(selfID int64, question string, answer string) {
	C.S3_Api_Setcation_problem_A(
		GoInt2CStr(selfID), CString(question), CString(answer),
	)
}

// 设置机器人被添加好友时的三个可选问题
// selfID  响应QQ  文本型  机器人QQ
// questionFrist  设置问题一  文本型  设置问题一
// questionSecond  设置问题二  文本型  设置问题二
// questionThird  设置问题三  文本型  设置问题三
func Setcation_problem_B(selfID int64, questionFrist string, questionSecond string, questionThird string) {
	C.S3_Api_Setcation_problem_B(
		GoInt2CStr(selfID), CString(questionFrist), CString(questionSecond), CString(questionThird),
	)
}

// 主动添加好友 请求成功返回真否则返回假
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
// remark  验证消息  文本型  无
// subType  来源信息  整数型  1QQ号码查找 2昵称查找 3条件查找 5临时会话 6QQ群 10QQ空间 11拍拍网 12最近联系人 14企业查找 其他的自己测试吧 1-255
func AddFriend(selfID int64, userID int64, remark string, subType int64) bool {
	return GoBool(C.S3_Api_AddFriend(
		GoInt2CStr(selfID), GoInt2CStr(userID), CString(remark), C.int(subType),
	))
}

// 标记函数执行流程 debug时使用 每个函数内只需要调用一次
// text  标记内容  文本型  无
func DbgName(text string) {
	C.S3_Api_DbgName(
		CString(text),
	)
}

// 函数内标记附加信息 函数内可多次调用
// text  标记内容  文本型  无
func Mark(text string) {
	C.S3_Api_Mark(
		CString(text),
	)
}

// 发送json结构消息
// selfID  响应QQ  文本型  机器人QQ
// anonymous  发送方式  整数型  1普通 2匿名（匿名需要群开启）
// messageType  信息类型  整数型  0在线临时会话 1好友 2群 3讨论组 4群临时会话 5讨论组临时会话 7好友验证回复会话
// groupID  收信对象所属群_讨论组  文本型  发送群信息、讨论组、群或讨论组临时会话信息时填，如发送对象为好友或信息类型是0时可空
// userID  收信对象QQ  文本型  收信对象QQ
// jsonData  Json结构  文本型  Json结构内容
func SendJSON(selfID int64, anonymous int64, messageType int64, groupID int64, userID int64, jsonData string) {
	C.S3_Api_SendJSON(
		GoInt2CStr(selfID), C.int(anonymous), C.int(messageType), GoInt2CStr(groupID), GoInt2CStr(userID), CString(jsonData),
	)
}

// 发送xml结构消息
// selfID  响应QQ  文本型  机器人QQ
// anonymous  发送方式  整数型  1普通 2匿名（匿名需要群开启）
// messageType  信息类型  整数型  0在线临时会话 1好友 2群 3讨论组 4群临时会话 5讨论组临时会话 7好友验证回复会话
// groupID  收信对象所属群_讨论组  文本型  发送群信息、讨论组、群或讨论组临时会话信息时填，如发送对象为好友或信息类型是0时可空
// userID  收信对象QQ  文本型  收信对象QQ
// xmlData  XML结构  文本型  Json结构内容
// nothing  NULL  整数型  无
func SendXML(selfID int64, anonymous int64, messageType int64, groupID int64, userID int64, xmlData string, nothing int64) {
	C.S3_Api_SendXML(
		GoInt2CStr(selfID), C.int(anonymous), C.int(messageType), GoInt2CStr(groupID), GoInt2CStr(userID), CString(xmlData), C.int(nothing),
	)
}

// 上传silk语音文件
// selfID  响应QQ  文本型  机器人QQ
// postType  上传类型  整数型  2、QQ群 讨论组
// groupID  接收群号  文本型  需上传的群号
// file  语音数据  字节集  语音字节集数据（AMR Silk编码）
func UpLoadVoice(selfID int64, postType int64, groupID int64, file []byte) string {
	return CPtr2GoStr(C.S3_Api_UpLoadVoice(
		GoInt2CStr(selfID), C.int(postType), GoInt2CStr(groupID), CByte(file),
	))
}

// 发送普通消息支持群匿名方式
// selfID  响应QQ  文本型  机器人QQ
// messageType  信息类型  整数型  0在线临时会话 1好友 2群 3讨论组 4群临时会话 5讨论组临时会话 7好友验证回复会话
// groupID  收信对象群_讨论组  文本型  发送群信息、讨论组、群或讨论组临时会话信息时填写，如发送对象为好友或信息类型是0时可空
// userID  收信QQ  文本型  收信对象QQ
// message  内容  文本型  信息内容
// bubble  气泡ID  整数型  已支持请自己测试
// anonymous  群匿名  逻辑型  不需要匿名请填写假 可调用Api_GetAnon函数 查看群是否开启匿名如果群没有开启匿名发送消息会失 败
func SendMsgEX(selfID int64, messageType int64, groupID int64, userID int64, message string, bubble int64, anonymous bool) {
	C.S3_Api_SendMsgEX(
		GoInt2CStr(selfID), C.int(messageType), GoInt2CStr(groupID), GoInt2CStr(userID), CString(message), C.int(bubble), CBool(anonymous),
	)
}

// 通过语音GUID获取语音文件下载连接
// selfID  响应QQ  文本型  机器人QQ
// recordUUID  语音GUID  文本型  [IR:Voi={xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx}.amr]
func GetVoiLink(selfID int64, recordUUID string) string {
	return CPtr2GoStr(C.S3_Api_GetVoiLink(
		GoInt2CStr(selfID), CString(recordUUID),
	))
}

// 查询指定群是否允许匿名消息
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需开获取匿名功能开关的群号
func GetAnon(selfID int64, groupID int64) bool {
	return GoBool(C.S3_Api_GetAnon(
		GoInt2CStr(selfID), GoInt2CStr(groupID),
	))
}

// 开关群匿名功能
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需开关群匿名功能群号
// enable  开关  逻辑型  真开    假关
func SetAnon(selfID int64, groupID int64, enable bool) bool {
	return GoBool(C.S3_Api_SetAnon(
		GoInt2CStr(selfID), GoInt2CStr(groupID), CBool(enable),
	))
}

// 取得机器人网页操作用的长Clientkey
// selfID  响应QQ  文本型  机器人QQ
func GetLongClientkey(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetLongClientkey(
		GoInt2CStr(selfID),
	))
}

// 取得腾讯微博页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
func GetBlogPsKey(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetBlogPsKey(
		GoInt2CStr(selfID),
	))
}

// 取得腾讯课堂页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
func GetClassRoomPsKey(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetClassRoomPsKey(
		GoInt2CStr(selfID),
	))
}

// 取得QQ举报页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
func GetRepPsKey(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetRepPsKey(
		GoInt2CStr(selfID),
	))
}

// 取得财付通页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
func GetTenPayPsKey(selfID int64) string {
	return CPtr2GoStr(C.S3_Api_GetTenPayPsKey(
		GoInt2CStr(selfID),
	))
}

// 修改机器人自身头像
// selfID  响应QQ  文本型  机器人QQ
// data  数据  字节集  图像数据
func SetHeadPic(selfID int64, data []byte) bool {
	return GoBool(C.S3_Api_SetHeadPic(
		GoInt2CStr(selfID), CByte(data),
	))
}

// 语音GUID转换为文本内容
// selfID  响应QQ  文本型  无
// userID  参考对象  文本型  无
// subType  参考类型  整数型  无
// recordUUID  语音GUID  文本型  无
func VoiToText(selfID int64, userID int64, subType int64, recordUUID string) string {
	return CPtr2GoStr(C.S3_Api_VoiToText(
		GoInt2CStr(selfID), GoInt2CStr(userID), C.int(subType), CString(recordUUID),
	))
}

// 群签到
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  QQ群号
// place  地名  文本型  签到地名
// message  内容  文本型  想发表的内容
func SignIn(selfID int64, groupID int64, place string, message string) bool {
	return GoBool(C.S3_Api_SignIn(
		GoInt2CStr(selfID), GoInt2CStr(groupID), CString(place), CString(message),
	))
}

// 通过图片GUID获取图片下注链接
// selfID  响应QQ  文本型  机器人QQ
// imageType  图片类型  整数型  2群 讨论组  1临时会话和好友
// userID  参考对象  文本型  图片所属对应的群号（可随意乱填写，只有群图片需要填写）
// imageUUID  图片GUID  文本型  例如[IR:pic={xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}.jpg]
func GetPicLink(selfID int64, imageType int64, userID int64, imageUUID string) string {
	return CPtr2GoStr(C.S3_Api_GetPicLink(
		GoInt2CStr(selfID), C.int(imageType), GoInt2CStr(userID), CString(imageUUID),
	))
}

// 获取框架版本号
func GetVer() string {
	return CPtr2GoStr(C.S3_Api_GetVer())
}

// 获取指定QQ个人资料的年龄
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
func GetAge(selfID int64, userID int64) int64 {
	return int64(C.S3_Api_GetAge(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 获取QQ个人资料的性别
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
func GetGender(selfID int64, userID int64) int64 {
	return int64(C.S3_Api_GetGender(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 向好友发送窗口抖动消息
// selfID  响应QQ  文本型  机器人QQ
// userID  接收QQ  文本型  接收抖动消息的QQ
func ShakeWindow(selfID int64, userID int64) bool {
	return GoBool(C.S3_Api_ShakeWindow(
		GoInt2CStr(selfID), GoInt2CStr(userID),
	))
}

// 同步发送消息 有返回值可以用来撤回机器人自己发送的消息
// selfID  响应QQ  文本型  机器人QQ
// messageType  信息类型  整数型  1好友 2群 3讨论组 4群临时会话
// groupID  收信群_讨论组  文本型  发送群信息、讨论组信息、群临时会话信息、讨论组临时会话信息时填写
// userID  收信对象  文本型  最终接收这条信息的对象QQ
// message  内容  文本型  信息内容
// bubble  气泡ID  整数型  -2强制不使用奇葩 -1随机使用气泡 0跟随框架的设置
// anonymous  匿名模式  逻辑型  是否使用匿名模式
// jsonData  附加JSON参数  文本型  以后信息发送参数增加都是依靠这个json文本
func SendMsgEX_V2(selfID int64, messageType int64, groupID int64, userID int64, message string, bubble int64, anonymous bool, jsonData string) string {
	return CPtr2GoStr(C.S3_Api_SendMsgEX_V2(
		GoInt2CStr(selfID), C.int(messageType), GoInt2CStr(groupID), GoInt2CStr(userID), CString(message), C.int(bubble), CBool(anonymous), CString(jsonData),
	))
}

// 撤回群消息或者私聊消息
// selfID  响应QQ  文本型  机器人QQ
// sub_type  撤回类型  整数型  1好友 2群聊 4群临时会话
// groupID  参考来源  文本型  非临时会话时请留空 临时会话请填群号
// userID  参考对象  文本型  非私聊消息请留空 私聊消息请填写对方QQ号码
// messageNum  消息序号  文本型  需撤回消息序号
// messageID  消息ID  文本型  需撤回消息ID
// time  消息时间  文本型  私聊消息需要群聊时3可留空
func WithdrawMsgEX(selfID int64, sub_type int64, groupID int64, userID int64, messageNum int64, messageID int64, time int64) string {
	return CPtr2GoStr(C.S3_Api_WithdrawMsgEX(
		GoInt2CStr(selfID), C.int(sub_type), GoInt2CStr(groupID), GoInt2CStr(userID), GoInt2CStr(messageNum), GoInt2CStr(messageID), GoInt2CStr(time),
	))
}

// 重新从Plugin目录下载入本插件(一般用做自动更新)
func Reload() bool {
	return GoBool(C.S3_Api_Reload())
}

// 返回框架加载的所有插件列表(包含本插件)的json文本
func GetPluginList() string {
	return CPtr2GoStr(C.S3_Api_GetPluginList())
}

// 查询指定对象是否允许发送在线状态临时会话 获取失败返回0 允许返回1 禁止返回2
// userID  对象QQ  文本型  无
func GetWpa(userID int64) int64 {
	return int64(C.S3_Api_GetWpa(
		GoInt2CStr(userID),
	))
}

// 主动卸载插件自身
func Uninstall() bool {
	return GoBool(C.S3_Api_Uninstall())
}
