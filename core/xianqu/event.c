#include <stdint.h>
#include <windows.h>

extern char *GoCreate(char *version);
extern int GoEvent(char *self_id, int message_type, int sub_type, char *group_id, char *user_id, char *notice_id, char *message, char *message_num, char *message_id, char *raw_message, char *time, int ret);
extern int GoSetUp();
extern int GoDestroyPlugin();

//export XQ_Create
char* __stdcall __declspec (dllexport) XQ_Create(char *version)
{
	return GoCreate(version);
}

// self_id		机器人QQ, 多Q版用于判定哪个QQ接收到该消息
// message_type	消息类型, 接收到消息类型，该类型可在常量表中查询具体定义，此处仅列举： -1 未定义事件 0,在线状态临时会话 1,好友信息 2,群信息 3,讨论组信息 4,群临时会话 5,讨论组临时会话 6,财付通转账 7,好友验证回复会话
// sub_type		消息子类型, 此参数在不同消息类型下，有不同的定义，暂定：接收财付通转账时 1为好友 4为群临时会话 5为讨论组临时会话    有人请求入群时，不良成员这里为1
// group_id		消息来源, 此消息的来源，如：群号、讨论组ID、临时会话QQ、好友QQ等
// user_id		触发对象_主动, 主动发送这条消息的QQ，踢人时为踢人管理员QQ
// notice_id	触发对象_被动, 被动触发的QQ，如某人被踢出群，则此参数为被踢出人QQ
// message		消息内容, 此参数有多重含义，常见为：对方发送的消息内容，但当消息类型为 某人申请入群，则为入群申请理由
// message_num	消息序号, 此参数暂定用于消息回复，消息撤回
// message_id	消息ID, 此参数暂定用于消息回复，消息撤回
// raw_message	原始信息, UDP收到的原始信息，特殊情况下会返回JSON结构（入群事件时，这里为该事件seq）
// time			消息时间戳, 接受到消息的时间戳
// ret			回传文本指针, 此参数用于插件加载拒绝理由

//export XQ_Event
int __stdcall __declspec (dllexport) XQ_Event(char *self_id, int message_type, int sub_type, char *group_id, char *user_id, char *notice_id, char *message, char *message_num, char *message_id, char *raw_message, char *time, int ret)
{
	return GoEvent(self_id,message_type,sub_type, group_id,user_id,notice_id,message,message_num,message_id,raw_message,time,ret);
}

//export XQ_SetUp
int __stdcall __declspec (dllexport) XQ_SetUp()
{
	return GoSetUp();
}

//export XQ_DestroyPlugin
int __stdcall __declspec (dllexport) XQ_DestroyPlugin()
{
	return GoDestroyPlugin();
}