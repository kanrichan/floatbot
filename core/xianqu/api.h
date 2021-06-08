// from https://github.com/Tnze/CoolQ-Golang-SDK/blob/master/cqp/apis_native.go
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>
#include <tchar.h>

#define XQAPI(RetType, Name, ...)                                                   \
    typedef RetType(__stdcall *Name##_Type)(unsigned char *authid, ##__VA_ARGS__); \
    Name##_Type Name##_Ptr;                                                         \
    RetType Name(__VA_ARGS__);

#define LoadAPI(Name) Name##_Ptr = (Name##_Type)GetProcAddress(hmod, #Name)

unsigned char *authid;

XQAPI(int, S3_Api_ApiInit);
XQAPI(void, S3_Api_SetAuthId, int, int);
XQAPI(char *, S3_Api_GetFriendList, char *);
XQAPI(char *, S3_Api_GetOnLineList);
XQAPI(int, S3_Api_Getbotisonline, char *);
XQAPI(char *, S3_Api_GetGroupMemberList, char *, char *);
XQAPI(char *, S3_Api_GetGroupCard, char *, char *, char *);
XQAPI(void, S3_Api_SendMsg, char *, int, char *, char *, char *, int);
XQAPI(char *, S3_Api_UpLoadPic, char *, int, char *, char *);
XQAPI(char *, S3_Api_GetGroupAdmin, char *, char *);
XQAPI(void, S3_Api_ShutUP, char *, char *, char *, int);
XQAPI(int, S3_Api_SetGroupCard, char *, char *, char *, char *);
XQAPI(void, S3_Api_KickGroupMBR, char *, char *, char *, int);
XQAPI(char *, S3_Api_GetNotice, char *, char *);
XQAPI(int, S3_Api_IsShutUp, char *, char *, char *);
XQAPI(int, S3_Api_IfFriend, char *, char *);
XQAPI(void, S3_Api_SetRInf, char *, int, char *);
XQAPI(char *, S3_Api_GetGroupPsKey, char *);
XQAPI(char *, S3_Api_GetZonePsKey, char *);
XQAPI(char *, S3_Api_GetCookies, char *);
XQAPI(int, S3_Api_PBGroupNotic, char *, char *, char *, char *);
XQAPI(char *, S3_Api_WithdrawMsg, char *, char *, char *, char *);
XQAPI(void, S3_Api_OutPutLog, char *);
XQAPI(char *, S3_Api_OcrPic, char *, char *);
XQAPI(void, S3_Api_JoinGroup, char *, char *, char *);
XQAPI(char *, S3_Api_UpVote, char *, char *);
XQAPI(char *, S3_Api_UpVote_temp, char *, char *);
XQAPI(int, S3_Api_GetObjVote, char *, char *);
XQAPI(int, S3_Api_IsEnable);
XQAPI(void, S3_Api_HandleFriendEvent, char *, char *, int, char *);
XQAPI(void, S3_Api_HandleGroupEvent, char *, int, char *, char *, char *, int, char *);
XQAPI(char *, S3_Api_GetQQList);
XQAPI(char *, S3_Api_AddQQ, char *, char *, int);
XQAPI(void, S3_Api_LoginQQ, char *);
XQAPI(void, S3_Api_OffLineQQ, char *);
XQAPI(char *, S3_Api_DelQQ, char *);
XQAPI(int, S3_Api_DelFriend, char *, char *);
XQAPI(char *, S3_Api_GetNick, char *, char *);
XQAPI(char *, S3_Api_GetFriendsRemark, char *, char *);
XQAPI(char *, S3_Api_GetClientkey, char *);
XQAPI(char *, S3_Api_GetBkn, char *);
XQAPI(void, S3_Api_SetFriendsRemark, char *, char *, char *);
XQAPI(void, S3_Api_InviteGroup, char *, char *, char *);
XQAPI(int, S3_Api_InviteGroupMember, char *, char *, char *, char *);
XQAPI(char *, S3_Api_CreateDisGroup, char *);
XQAPI(char *, S3_Api_CreateGroup, char *, char *);
XQAPI(void, S3_Api_QuitGroup, char *, char *);
XQAPI(char *, S3_Api_GetGroupList, char *);
XQAPI(char *, S3_Api_GetGroupList_B, char *);
XQAPI(char *, S3_Api_GetFriendList_B, char *);
XQAPI(char *, S3_Api_GetQrcode, char *);
XQAPI(int, S3_Api_CheckQrcode, char *);
XQAPI(char *, S3_Api_GetGroupName, char *, char *);
XQAPI(char *, S3_Api_GetGroupMemberNum, char *, char *);
XQAPI(int, S3_Api_GetGroupLv, char *, char *);
XQAPI(void, S3_Api_SetShieldedGroup, char *, char *, int);
XQAPI(char *, S3_Api_GetGroupMemberList_B, char *, char *);
XQAPI(char *, S3_Api_GetGroupMemberList_C, char *, char *);
XQAPI(int, S3_Api_IsOnline, char *, char *);
XQAPI(char *, S3_Api_GetRInf, char *);
XQAPI(void, S3_Api_DelFriend_A, char *, char *, int);
XQAPI(void, S3_Api_Setcation, char *, int);
XQAPI(void, S3_Api_Setcation_problem_A, char *, char *, char *);
XQAPI(void, S3_Api_Setcation_problem_B, char *, char *, char *, char *);
XQAPI(int, S3_Api_AddFriend, char *, char *, char *, int);
XQAPI(void, S3_Api_DbgName, char *);
XQAPI(void, S3_Api_Mark, char *);
XQAPI(void, S3_Api_SendJSON, char *, int, int, char *, char *, char *);
XQAPI(void, S3_Api_SendXML, char *, int, int, char *, char *, char *, int);
XQAPI(char *, S3_Api_UpLoadVoice, char *, int, char *, char *);
XQAPI(void, S3_Api_SendMsgEX, char *, int, char *, char *, char *, int, int);
XQAPI(char *, S3_Api_GetVoiLink, char *, char *);
XQAPI(int, S3_Api_GetAnon, char *, char *);
XQAPI(int, S3_Api_SetAnon, char *, char *, int);
XQAPI(char *, S3_Api_GetLongClientkey, char *);
XQAPI(char *, S3_Api_GetBlogPsKey, char *);
XQAPI(char *, S3_Api_GetClassRoomPsKey, char *);
XQAPI(char *, S3_Api_GetRepPsKey, char *);
XQAPI(char *, S3_Api_GetTenPayPsKey, char *);
XQAPI(int, S3_Api_SetHeadPic, char *, char *);
XQAPI(char *, S3_Api_VoiToText, char *, char *, int, char *);
XQAPI(int, S3_Api_SignIn, char *, char *, char *, char *);
XQAPI(char *, S3_Api_GetPicLink, char *, int, char *, char *);
XQAPI(char *, S3_Api_GetVer);
XQAPI(int, S3_Api_GetAge, char *, char *);
XQAPI(int, S3_Api_GetGender, char *, char *);
XQAPI(int, S3_Api_ShakeWindow, char *, char *);
XQAPI(char *, S3_Api_SendMsgEX_V2, char *, int, char *, char *, char *, int, int, char *);
XQAPI(char *, S3_Api_WithdrawMsgEX, char *, int, char *, char *, char *, char *, char *);
XQAPI(int, S3_Api_Reload);
XQAPI(char *, S3_Api_GetPluginList);
XQAPI(int, S3_Api_GetWpa, char *);
XQAPI(int, S3_Api_Uninstall);

extern void __stdcall __declspec (dllexport) XQ_AuthId(int ID, int IMAddr){
    authid = (unsigned char *)malloc(sizeof(unsigned char)*16);
    *((int*)authid) = 1;
    *((int*)(authid + 4)) = 8;
    *((int*)(authid + 8)) = ID;
    *((int*)(authid + 12)) = IMAddr;
    authid += 8;
    HMODULE hmod = LoadLibraryA("xqapi.dll");
    LoadAPI(S3_Api_ApiInit);
    LoadAPI(S3_Api_SetAuthId);
    LoadAPI(S3_Api_GetFriendList);
    LoadAPI(S3_Api_GetOnLineList);
    LoadAPI(S3_Api_Getbotisonline);
    LoadAPI(S3_Api_GetGroupMemberList);
    LoadAPI(S3_Api_GetGroupCard);
    LoadAPI(S3_Api_SendMsg);
    LoadAPI(S3_Api_UpLoadPic);
    LoadAPI(S3_Api_GetGroupAdmin);
    LoadAPI(S3_Api_ShutUP);
    LoadAPI(S3_Api_SetGroupCard);
    LoadAPI(S3_Api_KickGroupMBR);
    LoadAPI(S3_Api_GetNotice);
    LoadAPI(S3_Api_IsShutUp);
    LoadAPI(S3_Api_IfFriend);
    LoadAPI(S3_Api_SetRInf);
    LoadAPI(S3_Api_GetGroupPsKey);
    LoadAPI(S3_Api_GetZonePsKey);
    LoadAPI(S3_Api_GetCookies);
    LoadAPI(S3_Api_PBGroupNotic);
    LoadAPI(S3_Api_WithdrawMsg);
    LoadAPI(S3_Api_OutPutLog);
    LoadAPI(S3_Api_OcrPic);
    LoadAPI(S3_Api_JoinGroup);
    LoadAPI(S3_Api_UpVote);
    LoadAPI(S3_Api_UpVote_temp);
    LoadAPI(S3_Api_GetObjVote);
    LoadAPI(S3_Api_IsEnable);
    LoadAPI(S3_Api_HandleFriendEvent);
    LoadAPI(S3_Api_HandleGroupEvent);
    LoadAPI(S3_Api_GetQQList);
    LoadAPI(S3_Api_AddQQ);
    LoadAPI(S3_Api_LoginQQ);
    LoadAPI(S3_Api_OffLineQQ);
    LoadAPI(S3_Api_DelQQ);
    LoadAPI(S3_Api_DelFriend);
    LoadAPI(S3_Api_GetNick);
    LoadAPI(S3_Api_GetFriendsRemark);
    LoadAPI(S3_Api_GetClientkey);
    LoadAPI(S3_Api_GetBkn);
    LoadAPI(S3_Api_SetFriendsRemark);
    LoadAPI(S3_Api_InviteGroup);
    LoadAPI(S3_Api_InviteGroupMember);
    LoadAPI(S3_Api_CreateDisGroup);
    LoadAPI(S3_Api_CreateGroup);
    LoadAPI(S3_Api_QuitGroup);
    LoadAPI(S3_Api_GetGroupList);
    LoadAPI(S3_Api_GetGroupList_B);
    LoadAPI(S3_Api_GetFriendList_B);
    LoadAPI(S3_Api_GetQrcode);
    LoadAPI(S3_Api_CheckQrcode);
    LoadAPI(S3_Api_GetGroupName);
    LoadAPI(S3_Api_GetGroupMemberNum);
    LoadAPI(S3_Api_GetGroupLv);
    LoadAPI(S3_Api_SetShieldedGroup);
    LoadAPI(S3_Api_GetGroupMemberList_B);
    LoadAPI(S3_Api_GetGroupMemberList_C);
    LoadAPI(S3_Api_IsOnline);
    LoadAPI(S3_Api_GetRInf);
    LoadAPI(S3_Api_DelFriend_A);
    LoadAPI(S3_Api_Setcation);
    LoadAPI(S3_Api_Setcation_problem_A);
    LoadAPI(S3_Api_Setcation_problem_B);
    LoadAPI(S3_Api_AddFriend);
    LoadAPI(S3_Api_DbgName);
    LoadAPI(S3_Api_Mark);
    LoadAPI(S3_Api_SendJSON);
    LoadAPI(S3_Api_SendXML);
    LoadAPI(S3_Api_UpLoadVoice);
    LoadAPI(S3_Api_SendMsgEX);
    LoadAPI(S3_Api_GetVoiLink);
    LoadAPI(S3_Api_GetAnon);
    LoadAPI(S3_Api_SetAnon);
    LoadAPI(S3_Api_GetLongClientkey);
    LoadAPI(S3_Api_GetBlogPsKey);
    LoadAPI(S3_Api_GetClassRoomPsKey);
    LoadAPI(S3_Api_GetRepPsKey);
    LoadAPI(S3_Api_GetTenPayPsKey);
    LoadAPI(S3_Api_SetHeadPic);
    LoadAPI(S3_Api_VoiToText);
    LoadAPI(S3_Api_SignIn);
    LoadAPI(S3_Api_GetPicLink);
    LoadAPI(S3_Api_GetVer);
    LoadAPI(S3_Api_GetAge);
    LoadAPI(S3_Api_GetGender);
    LoadAPI(S3_Api_ShakeWindow);
    LoadAPI(S3_Api_SendMsgEX_V2);
    LoadAPI(S3_Api_WithdrawMsgEX);
    LoadAPI(S3_Api_Reload);
    LoadAPI(S3_Api_GetPluginList);
    LoadAPI(S3_Api_GetWpa);
    LoadAPI(S3_Api_Uninstall);

    return;
}

char *str_same(char *a, char* b) {
    if (!a || !b) return NULL;
	int len = strlen(a);
    int i = 0;
    for (; i < len; ++i) if ((*(a+i) - *(b+i))) break;
	char *c = (char*)malloc(sizeof(char)*i);
	memcpy(c, a, i);
    c[i]='\0';
    return c;
}

char *fix_str(char *string) {
    if (!string) return NULL;
    return string+4;
}

int index_of(char *string, const char *dest) {
    if (!string || !dest ) return -1;
    int i = 0;
    int j = 0;
    while (string[i] != '\0') {
        if (string[i] != dest[0]) {
            i ++;
            continue;
        }
        j = 0;
        while (string[i+j] != '\0' && dest[j] != '\0') {
            if (string[i+j] != dest[j]) {
                break;
            }
            j ++;
        }
        if (dest[j] == '\0') return i+strlen(dest);
        i ++;
    }
    return -1;
}

char *cut_str(char *string,int index) {
    if (!string) return NULL;
    if (index>strlen(string)) return string;
	char *ret = (char*)malloc(sizeof(char)*index);
	memcpy(ret, string, index);
    ret[index]='\0';
    return ret;
}

// 在使用本类方法前必须调用本函数(返回框架版本号)
int S3_Api_ApiInit(){
    int ret = S3_Api_ApiInit_Ptr(authid);
    return ret;
}

// 在没置入此ID前所有接口调用都不会有效果(除了取框架版本与标记异常流程)
// id  id  整数型  无
// addr  addr  整数型  无
void S3_Api_SetAuthId(int id, int addr){
    S3_Api_SetAuthId_Ptr(authid, id, addr);

}

// 取得好友列表，返回获取到的原始JSON格式信息，需自行解析，http模式
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetFriendList(char *selfID){
    char *ret = S3_Api_GetFriendList_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取机器人在线账号列表
char *S3_Api_GetOnLineList(){
    char *ret = S3_Api_GetOnLineList_Ptr(authid);
    return fix_str(ret);
}

// 取机器人账号是否在线
// selfID  响应QQ  文本型  机器人QQ
int S3_Api_Getbotisonline(char *selfID){
    int ret = S3_Api_Getbotisonline_Ptr(authid, selfID);
    free(selfID);
    return ret;
}

// 取群员列表
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取群成员列表群号
char *S3_Api_GetGroupMemberList(char *selfID, char *groupID){
    char *ret = S3_Api_GetGroupMemberList_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return fix_str(ret);
}

// 取群成员名片
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取群成员列表群号
// userID  对象QQ  文本型  无
char *S3_Api_GetGroupCard(char *selfID, char *groupID, char *userID){
    char *ret = S3_Api_GetGroupCard_Ptr(authid, selfID, groupID, userID);
    free(selfID);
    free(groupID);
    free(userID);
    return fix_str(ret);
}

// 发送消息
// selfID  响应QQ  文本型  机器人QQ
// messageType  信息类型  整数型  0在线临时会话 1好友 2群 3讨论组 4群临时会话 5讨论组临时会话 7好友验证回复会话
// groupID  收信对象群_讨论组  文本型  发送群信息、讨论组、群或讨论组临时会话信息时填写，如发送对象为好友或信息类型是0时可空
// userID  收信QQ  文本型  收信对象QQ
// message  内容  文本型  信息内容
// bubble  气泡ID  整数型  已支持请自己测试
void S3_Api_SendMsg(char *selfID, int messageType, char *groupID, char *userID, char *message, int bubble){
    S3_Api_SendMsg_Ptr(authid, selfID, messageType, groupID, userID, message, bubble);
}

// 上传图片
// selfID  响应QQ  文本型  机器人QQ
// postType  上传类型  整数型  1好友、临时会话  2群、讨论组 Ps：好友临时会话用类型 1，群讨论组用类型 2；当填写错误时，图片GUID发送不会成功
// userID  参考对象  文本型  上传该图片所属的群号或QQ
// file  图片数据  字节集  图片字节集数据
char *S3_Api_UpLoadPic(char *selfID, int postType, char *userID, char *file){
    char *ret = S3_Api_UpLoadPic_Ptr(authid, selfID, postType, userID, file);
    free(selfID);
    free(userID);
    return fix_str(ret);
}

// 取群管理员列表
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取管理员列表群号
char *S3_Api_GetGroupAdmin(char *selfID, char *groupID){
    char *ret = S3_Api_GetGroupAdmin_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return fix_str(ret);
}

// 群禁言
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲操作的群号
// userID  对象QQ  文本型  欲禁言的对象，如留空且机器人QQ为管理员，将设置该群为全群禁言
// time  时间  整数型  0为解除禁言 （禁言单位为秒），如为全群禁言，参数为非0，解除全群禁言为0
void S3_Api_ShutUp(char *selfID, char *groupID, char *userID, int time){
    S3_Api_ShutUP_Ptr(authid, selfID, groupID, userID, time);
    free(selfID);
    free(groupID);
    free(userID);

}

// 全群禁言
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲操作的群号
// enable  是否全群禁言  逻辑型  
void S3_Api_ShutUpAll(char *selfID, char *groupID, int enable){
    if (enable) S3_Api_ShutUP_Ptr(authid, selfID, groupID, "", 1);
    else S3_Api_ShutUP_Ptr(authid, selfID, groupID, "", 0);
    free(selfID);
    free(groupID);
}

// 修改群成员昵称
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  对象所处群号
// userID  对象QQ  文本型  被修改名片人QQ
// card  名片  文本型  需要修改的名片
int S3_Api_SetGroupCard(char *selfID, char *groupID, char *userID, char *card){
    int ret = S3_Api_SetGroupCard_Ptr(authid, selfID, groupID, userID, card);
    free(selfID);
    free(groupID);
    free(userID);
    free(card);
    return ret;
}

// 群删除成员
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  被执行群号
// userID  对象QQ  文本型  被执行对象
// reject_add_request  不在允许  逻辑型  真为不再接收，假为接收
void S3_Api_KickGroupMBR(char *selfID, char *groupID, char *userID, int reject_add_request){
    S3_Api_KickGroupMBR_Ptr(authid, selfID, groupID, userID, reject_add_request);
    free(selfID);
    free(groupID);
    free(userID);

}

// 获取群通知
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取得公告的群号
char *S3_Api_GetNotice(char *selfID, char *groupID){
    char *ret = S3_Api_GetNotice_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return fix_str(ret);
}

//  取群成员禁言状态 -1失败 0未被禁言 1被单独禁言 2开启了全群禁言
// selfID  响应QQ  文本型  无
// groupID  群号  文本型  无
// userID  对象QQ  文本型  无
int S3_Api_IsShutUp(char *selfID, char *groupID, char *userID){
    int ret = S3_Api_IsShutUp_Ptr(authid, selfID, groupID, userID);
    free(selfID);
    free(groupID);
    free(userID);
    return ret;
}

// 取是否好友
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
int S3_Api_IfFriend(char *selfID, char *userID){
    int ret = S3_Api_IfFriend_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return ret;
}

// 修改QQ在线状态
// selfID  响应QQ  文本型  无
// type  类型  整数型  1、我在线上 2、Q我吧 3、离开 4、忙碌 5、请勿打扰 6、隐身 7、修改昵称 8、修改个性签名 9、修改性别 
// text  修改内容  文本型  类型为7和8时填写修改内容  类型9时“1”为男 “2”为女      其他填“”
void S3_Api_SetRInf(char *selfID, int type, char *text){
    S3_Api_SetRInf_Ptr(authid, selfID, type, text);
    free(selfID);
    free(text);

}

// 取得QQ群页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetGroupPsKey(char *selfID){
    char *ret = S3_Api_GetGroupPsKey_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取得QQ空间页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetZonePsKey(char *selfID){
    char *ret = S3_Api_GetZonePsKey_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取得机器人网页操作用的Cookies
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetCookies(char *selfID){
    char *a = S3_Api_GetCookies_Ptr(authid, selfID);
    char *fix = fix_str(a);
    int index = index_of(fix,"skey=");
    char *ret = cut_str(fix, index+10);
    free(selfID);
    return ret;
}

// 发布群公告
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲发布公告的群号
// title  标题  文本型  公告标题
// message  内容  文本型  公告内容
int S3_Api_PBGroupNotic(char *selfID, char *groupID, char *title, char *message){
    int ret = S3_Api_PBGroupNotic_Ptr(authid, selfID, groupID, title, message);
    free(selfID);
    free(groupID);
    free(title);
    free(message);
    return ret;
}

// 撤回群消息
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需撤回消息群号
// messageNum  消息序号  文本型  需撤回消息序号
// messageID  消息ID  文本型  需撤回消息ID
char *S3_Api_WithdrawMsg(char *selfID, char *groupID, char *messageNum, char *messageID){
    char *ret = S3_Api_WithdrawMsg_Ptr(authid, selfID, groupID, messageNum, messageID);
    free(selfID);
    free(groupID);
    free(messageNum);
    free(messageID);
    return fix_str(ret);
}

// 输出一行日志
// message  内容  文本型  任意想输出的文本格式信息
void S3_Api_OutPutLog(char *message){
    S3_Api_OutPutLog_Ptr(authid, message);
    free(message);

}

// 提取图片中的文字
// selfID  响应QQ  文本型  机器人QQ
// file  图片数据  字节集  图片数据
char *S3_Api_OcrPic(char *selfID, char *file){
    char *ret = S3_Api_OcrPic_Ptr(authid, selfID, file);
    free(selfID);
    return fix_str(ret);
}

// 主动加群
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲申请加入的群号
// reason  理由  文本型  附加理由，可留空（需回答正确问题时，请填写问题答案）
void S3_Api_JoinGroup(char *selfID, char *groupID, char *reason){
    S3_Api_JoinGroup_Ptr(authid, selfID, groupID, reason);
    free(selfID);
    free(groupID);
    free(reason);

}

// 点赞
// selfID  响应QQ  文本型  机器人QQ
// userID  被赞QQ  文本型  填写被赞人QQ
char *S3_Api_UpVote(char *selfID, char *userID){
    char *ret = S3_Api_UpVote_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return fix_str(ret);
}

// 通过列表或群临时通道点赞
// selfID  响应QQ  文本型  机器人QQ
// userID  被赞QQ  文本型  填写被赞人QQ
char *S3_Api_UpVote_temp(char *selfID, char *userID){
    char *ret = S3_Api_UpVote_temp_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return fix_str(ret);
}

// 获取赞数量
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
int S3_Api_GetObjVote(char *selfID, char *userID){
    int ret = S3_Api_GetObjVote_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return ret;
}

// 取插件是否启用
int S3_Api_IsEnable(){
    int ret = S3_Api_IsEnable_Ptr(authid);
    return ret;
}

// 置好友添加请求
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  申请入群 被邀请人 请求添加好友人的QQ （当请求类型为214时这里为邀请人QQ）
// approve  处理方式  整数型  1同意 0拒绝
// remark  附加信息  文本型  拒绝入群，拒绝添加好友 附加信息
void S3_Api_HandleFriendEvent(char *selfID, char *userID, int approve, char *remark){
    if (approve) S3_Api_HandleFriendEvent_Ptr(authid, selfID, userID, 10, remark);
    else S3_Api_HandleFriendEvent_Ptr(authid, selfID, userID, 20, remark);
    free(selfID);
    free(userID);
    free(remark);
}

// 置群请求
// selfID  响应QQ  文本型  机器人QQ
// sub_type  请求类型  整数型  213请求入群  214我被邀请加入某群  215某人被邀请加入群  101某人请求添加好友
// userID  对象QQ  文本型  申请入群 被邀请人 请求添加好友人的QQ （当请求类型为214时这里为邀请人QQ）
// groupID  群号  文本型  收到请求群号（好友添加时这里请为空）
// flag  seq  文本型  需要处理事件的seq
// approve  处理方式  整数型  10同意 20拒绝 30忽略
// remark  附加信息  文本型  拒绝入群，拒绝添加好友 附加信息
void S3_Api_HandleGroupEvent(char *selfID, int sub_type, char *userID, char *groupID, char *flag, int approve, char *remark){
    if (approve) S3_Api_HandleGroupEvent_Ptr(authid, selfID, sub_type, userID, groupID, flag, 10, remark);
    else S3_Api_HandleGroupEvent_Ptr(authid, selfID, sub_type, userID, groupID, flag, 20, remark);
    free(selfID);
    free(userID);
    free(groupID);
    free(flag);
    free(remark);
}

// 取所有QQ列表
char *S3_Api_GetQQList(){
    char *ret = S3_Api_GetQQList_Ptr(authid);
    return fix_str(ret);
}

// 向框架添加一个QQ
// account  帐号  文本型  无
// password  密码  文本型  无
// enable  自动登录  逻辑型  真 为自动登录
char *S3_Api_AddQQ(char *account, char *password, int enable){
    char *ret = S3_Api_AddQQ_Ptr(authid, account, password, enable);
    free(account);
    free(password);
    return fix_str(ret);
}

// 登录指定QQ
// userID  登录QQ  文本型  无
void S3_Api_LoginQQ(char *userID){
    S3_Api_LoginQQ_Ptr(authid, userID);
    free(userID);

}

// 离线指定QQ
// selfID  响应QQ  文本型  无
void S3_Api_OffLineQQ(char *selfID){
    S3_Api_OffLineQQ_Ptr(authid, selfID);
    free(selfID);

}

// 删除指定QQ
// selfID  响应QQ  文本型  无
char *S3_Api_DelQQ(char *selfID){
    char *ret = S3_Api_DelQQ_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 删除指定好友
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  被删除对象
int S3_Api_DelFriend(char *selfID, char *userID){
    int ret = S3_Api_DelFriend_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return ret;
}

// 取QQ昵称
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  欲取得的QQ的号码
char *S3_Api_GetNick(char *selfID, char *userID){
    char *a = S3_Api_GetNick_Ptr(authid, selfID, userID);
    char *b = S3_Api_GetNick_Ptr(authid, selfID, userID);
    char *ret = str_same(fix_str(a), fix_str(b));
    free(selfID);
    free(userID);
    return ret;
}

// 取好友备注姓名
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  需获取对象好友QQ
char *S3_Api_GetFriendsRemark(char *selfID, char *userID){
    char *ret = S3_Api_GetFriendsRemark_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return fix_str(ret);
}

// 取短Clientkey
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetClientkey(char *selfID){
    char *ret = S3_Api_GetClientkey_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取bkn
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetBkn(char *selfID){
    char *ret = S3_Api_GetBkn_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 修改好友备注名称
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  需获取对象好友QQ
// remark  备注  文本型  需要修改的备注姓名
void S3_Api_SetFriendsRemark(char *selfID, char *userID, char *remark){
    S3_Api_SetFriendsRemark_Ptr(authid, selfID, userID, remark);
    free(selfID);
    free(userID);
    free(remark);

}

// 邀请好友加入群
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  被邀请加入的群号
// userID  对象QQ  文本型  被邀请人QQ号码
void S3_Api_InviteGroup(char *selfID, char *groupID, char *userID){
    S3_Api_InviteGroup_Ptr(authid, selfID, groupID, userID);
    free(selfID);
    free(groupID);
    free(userID);

}

// 邀请群成员加入群
// selfID  响应QQ  文本型  机器人QQ
// targetID  群号  文本型  邀请到哪个群
// groupID  所在群  文本型  被邀请成员所在群
// userID  邀请QQ  文本型  被邀请人的QQ
int S3_Api_InviteGroupMember(char *selfID, char *targetID, char *groupID, char *userID){
    int ret = S3_Api_InviteGroupMember_Ptr(authid, selfID, targetID, groupID, userID);
    free(selfID);
    free(groupID);
    free(groupID);
    free(userID);
    return ret;
}

// 创建群 组包模式
// selfID  响应QQ  文本型  机器人
char *S3_Api_CreateDisGroup(char *selfID){
    char *ret = S3_Api_CreateDisGroup_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 创建群 群官网Http模式
// selfID  响应QQ  文本型  机器人QQ
// nickname  群昵称  文本型  预创建的群名称
char *S3_Api_CreateGroup(char *selfID, char *nickname){
    char *ret = S3_Api_CreateGroup_Ptr(authid, selfID, nickname);
    free(selfID);
    free(nickname);
    return fix_str(ret);
}

// 退出群
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲退出的群号
void S3_Api_QuitGroup(char *selfID, char *groupID){
    S3_Api_QuitGroup_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);

}

// 封包模式获取群号列表(最多可以取得999)
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetGroupList(char *selfID){
    char *ret = S3_Api_GetGroupList_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 封包模式获取群号列表(最多可以取得999)
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetGroupList_B(char *selfID){
    char *ret = S3_Api_GetGroupList_B_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 封包模式取好友列表(与封包模式取群列表同源)
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetFriendList_B(char *selfID){
    char *ret = S3_Api_GetFriendList_B_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取登录二维码base64
// key  key  字节集  无
char *S3_Api_GetQrcode(char *key){
    char *ret = S3_Api_GetQrcode_Ptr(authid, key);
    return fix_str(ret);
}

// 检查登录二维码状态
// key  key  字节集  返回数值
int S3_Api_CheckQrcode(char *key){
    int ret = S3_Api_CheckQrcode_Ptr(authid, key);
    return ret;
}

// 取指定的群名称
// selfID  响应QQ  文本型  无
// groupID  群号  文本型  无
char *S3_Api_GetGroupName(char *selfID, char *groupID){
    char *ret = S3_Api_GetGroupName_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return fix_str(ret);
}

// 取群人数上线与当前人数 换行符分隔
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需查询的群号
char *S3_Api_GetGroupMemberNum(char *selfID, char *groupID){
    char *ret = S3_Api_GetGroupMemberNum_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return fix_str(ret);
}

// 取群等级
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需查询的群号
int S3_Api_GetGroupLv(char *selfID, char *groupID){
    int ret = S3_Api_GetGroupLv_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return ret;
}

// 屏蔽或接收某群消息
// selfID  响应QQ  文本型  无
// groupID  群号  文本型  无
// type  类型  逻辑型  真 为屏蔽接收  假为接收并提醒
void S3_Api_SetShieldedGroup(char *selfID, char *groupID, int type){
    S3_Api_SetShieldedGroup_Ptr(authid, selfID, groupID, type);
    free(selfID);
    free(groupID);

}

// 取群成员列表
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取群成员列表群号
char *S3_Api_GetGroupMemberList_B(char *selfID, char *groupID){
    char *ret = S3_Api_GetGroupMemberList_B_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return fix_str(ret);
}

// 封包模式取群成员列表返回重组后的json文本
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  欲取群成员列表群号
char *S3_Api_GetGroupMemberList_C(char *selfID, char *groupID){
    char *ret = S3_Api_GetGroupMemberList_C_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return fix_str(ret);
}

// 检查指定QQ是否在线
// selfID  响应QQ  文本型  机器人QQ
// userID  对象QQ  文本型  需获取对象QQ
int S3_Api_IsOnline(char *selfID, char *userID){
    int ret = S3_Api_IsOnline_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return ret;
}

// 取机器人账号在线信息
// selfID  响应QQ  文本型  无
char *S3_Api_GetRInf(char *selfID){
    char *ret = S3_Api_GetRInf_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 多功能删除好友 可删除陌生人或者删除为单项好友
// selfID  响应QQ  文本型  机器人QQ
// userID  目标QQ  文本型  欲操作的目标
// sub_type  删除类型  整数型  1为在对方的列表删除我 2为在我的列表删除对方
void S3_Api_DelFriend_A(char *selfID, char *userID, int sub_type){
    S3_Api_DelFriend_A_Ptr(authid, selfID, userID, sub_type);
    free(selfID);
    free(userID);

}

// 设置机器人被添加好友时的验证方式
// selfID  响应QQ  文本型  机器人QQ
// sub_type  验证类型  整数型  0 允许任何人 1 需要验证消息 2不允许任何人 3需要回答问题 4需要回答问题并由我确认
void S3_Api_Setcation(char *selfID, int sub_type){
    S3_Api_Setcation_Ptr(authid, selfID, sub_type);
    free(selfID);

}

// 设置机器人被添加好友时的问题与答案
// selfID  响应QQ  文本型  机器人QQ
// question  设置问题  文本型  设置的问题
// answer  问题答案  文本型  设置的问题答案
void S3_Api_Setcation_problem_A(char *selfID, char *question, char *answer){
    S3_Api_Setcation_problem_A_Ptr(authid, selfID, question, answer);
    free(selfID);
    free(question);
    free(answer);

}

// 设置机器人被添加好友时的三个可选问题
// selfID  响应QQ  文本型  机器人QQ
// questionFrist  设置问题一  文本型  设置问题一
// questionSecond  设置问题二  文本型  设置问题二
// questionThird  设置问题三  文本型  设置问题三
void S3_Api_Setcation_problem_B(char *selfID, char *questionFrist, char *questionSecond, char *questionThird){
    S3_Api_Setcation_problem_B_Ptr(authid, selfID, questionFrist, questionSecond, questionThird);
    free(selfID);
    free(questionFrist);
    free(questionSecond);
    free(questionThird);

}

// 主动添加好友 请求成功返回真否则返回假
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
// remark  验证消息  文本型  无
// subType  来源信息  整数型  1QQ号码查找 2昵称查找 3条件查找 5临时会话 6QQ群 10QQ空间 11拍拍网 12最近联系人 14企业查找 其他的自己测试吧 1-255
int S3_Api_AddFriend(char *selfID, char *userID, char *remark, int subType){
    int ret = S3_Api_AddFriend_Ptr(authid, selfID, userID, remark, subType);
    free(selfID);
    free(userID);
    free(remark);
    return ret;
}

// 标记函数执行流程 debug时使用 每个函数内只需要调用一次
// text  标记内容  文本型  无
void S3_Api_DbgName(char *text){
    S3_Api_DbgName_Ptr(authid, text);
    free(text);

}

// 函数内标记附加信息 函数内可多次调用
// text  标记内容  文本型  无
void S3_Api_Mark(char *text){
    S3_Api_Mark_Ptr(authid, text);
    free(text);

}

// 发送json结构消息
// selfID  响应QQ  文本型  机器人QQ
// anonymous  发送方式  整数型  1普通 2匿名（匿名需要群开启）
// messageType  信息类型  整数型  0在线临时会话 1好友 2群 3讨论组 4群临时会话 5讨论组临时会话 7好友验证回复会话
// groupID  收信对象所属群_讨论组  文本型  发送群信息、讨论组、群或讨论组临时会话信息时填，如发送对象为好友或信息类型是0时可空
// userID  收信对象QQ  文本型  收信对象QQ
// jsonData  Json结构  文本型  Json结构内容
void S3_Api_SendJSON(char *selfID, int anonymous, int messageType, char *groupID, char *userID, char *jsonData){
    S3_Api_SendJSON_Ptr(authid, selfID, anonymous, messageType, groupID, userID, jsonData);
    // free(selfID);
    // free(groupID);
    // free(userID);
    // free(jsonData);

}

// 发送xml结构消息
// selfID  响应QQ  文本型  机器人QQ
// anonymous  发送方式  整数型  1普通 2匿名（匿名需要群开启）
// messageType  信息类型  整数型  0在线临时会话 1好友 2群 3讨论组 4群临时会话 5讨论组临时会话 7好友验证回复会话
// groupID  收信对象所属群_讨论组  文本型  发送群信息、讨论组、群或讨论组临时会话信息时填，如发送对象为好友或信息类型是0时可空
// userID  收信对象QQ  文本型  收信对象QQ
// xmlData  XML结构  文本型  Json结构内容
// nothing  NULL  整数型  无
void S3_Api_SendXML(char *selfID, int anonymous, int messageType, char *groupID, char *userID, char *xmlData, int nothing){
    S3_Api_SendXML_Ptr(authid, selfID, anonymous, messageType, groupID, userID, xmlData, nothing);
    // free(selfID);
    // free(groupID);
    // free(userID);
    // free(xmlData);

}

// 上传silk语音文件
// selfID  响应QQ  文本型  机器人QQ
// postType  上传类型  整数型  2、QQ群 讨论组
// groupID  接收群号  文本型  需上传的群号
// file  语音数据  字节集  语音字节集数据（AMR Silk编码）
char *S3_Api_UpLoadVoice(char *selfID, int postType, char *groupID, char *file){
    char *ret = S3_Api_UpLoadVoice_Ptr(authid, selfID, postType, groupID, file);
    free(selfID);
    free(groupID);
    return fix_str(ret);
}

// 发送普通消息支持群匿名方式
// selfID  响应QQ  文本型  机器人QQ
// messageType  信息类型  整数型  0在线临时会话 1好友 2群 3讨论组 4群临时会话 5讨论组临时会话 7好友验证回复会话
// groupID  收信对象群_讨论组  文本型  发送群信息、讨论组、群或讨论组临时会话信息时填写，如发送对象为好友或信息类型是0时可空
// userID  收信QQ  文本型  收信对象QQ
// message  内容  文本型  信息内容
// bubble  气泡ID  整数型  已支持请自己测试
// anonymous  群匿名  逻辑型  不需要匿名请填写假 可调用Api_GetAnon函数 查看群是否开启匿名如果群没有开启匿名发送消息会失 败
void S3_Api_SendMsgEX(char *selfID, int messageType, char *groupID, char *userID, char *message, int bubble, int anonymous){
    S3_Api_SendMsgEX_Ptr(authid, selfID, messageType, groupID, userID, message, bubble, anonymous);
    free(selfID);
    free(groupID);
    free(userID);
    free(message);

}

// 通过语音GUID获取语音文件下载连接
// selfID  响应QQ  文本型  机器人QQ
// recordUUID  语音GUID  文本型  [IR:Voi={xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxx}.amr]
char *S3_Api_GetVoiLink(char *selfID, char *recordUUID){
    char *ret = S3_Api_GetVoiLink_Ptr(authid, selfID, recordUUID);
    free(selfID);
    free(recordUUID);
    return fix_str(ret);
}

// 查询指定群是否允许匿名消息
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需开获取匿名功能开关的群号
int S3_Api_GetAnon(char *selfID, char *groupID){
    int ret = S3_Api_GetAnon_Ptr(authid, selfID, groupID);
    free(selfID);
    free(groupID);
    return ret;
}

// 开关群匿名功能
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  需开关群匿名功能群号
// enable  开关  逻辑型  真开    假关
int S3_Api_SetAnon(char *selfID, char *groupID, int enable){
    int ret = S3_Api_SetAnon_Ptr(authid, selfID, groupID, enable);
    free(selfID);
    free(groupID);
    return ret;
}

// 取得机器人网页操作用的长Clientkey
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetLongClientkey(char *selfID){
    char *ret = S3_Api_GetLongClientkey_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取得腾讯微博页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetBlogPsKey(char *selfID){
    char *ret = S3_Api_GetBlogPsKey_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取得腾讯课堂页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetClassRoomPsKey(char *selfID){
    char *ret = S3_Api_GetClassRoomPsKey_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取得QQ举报页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetRepPsKey(char *selfID){
    char *ret = S3_Api_GetRepPsKey_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 取得财付通页面操作用参数P_skey
// selfID  响应QQ  文本型  机器人QQ
char *S3_Api_GetTenPayPsKey(char *selfID){
    char *ret = S3_Api_GetTenPayPsKey_Ptr(authid, selfID);
    free(selfID);
    return fix_str(ret);
}

// 修改机器人自身头像
// selfID  响应QQ  文本型  机器人QQ
// data  数据  字节集  图像数据
int S3_Api_SetHeadPic(char *selfID, char *data){
    int ret = S3_Api_SetHeadPic_Ptr(authid, selfID, data);
    free(selfID);
    return ret;
}

// 语音GUID转换为文本内容
// selfID  响应QQ  文本型  无
// userID  参考对象  文本型  无
// subType  参考类型  整数型  无
// recordUUID  语音GUID  文本型  无
char *S3_Api_VoiToText(char *selfID, char *userID, int subType, char *recordUUID){
    char *ret = S3_Api_VoiToText_Ptr(authid, selfID, userID, subType, recordUUID);
    free(selfID);
    free(userID);
    free(recordUUID);
    return fix_str(ret);
}

// 群签到
// selfID  响应QQ  文本型  机器人QQ
// groupID  群号  文本型  QQ群号
// place  地名  文本型  签到地名
// message  内容  文本型  想发表的内容
int S3_Api_SignIn(char *selfID, char *groupID, char *place, char *message){
    int ret = S3_Api_SignIn_Ptr(authid, selfID, groupID, place, message);
    free(selfID);
    free(groupID);
    free(place);
    free(message);
    return ret;
}

// 通过图片GUID获取图片下注链接
// selfID  响应QQ  文本型  机器人QQ
// imageType  图片类型  整数型  2群 讨论组  1临时会话和好友
// userID  参考对象  文本型  图片所属对应的群号（可随意乱填写，只有群图片需要填写）
// imageUUID  图片GUID  文本型  例如[IR:pic={xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}.jpg]
char *S3_Api_GetPicLink(char *selfID, int imageType, char *userID, char *imageUUID){
    char *ret = S3_Api_GetPicLink_Ptr(authid, selfID, imageType, userID, imageUUID);
    free(selfID);
    free(userID);
    free(imageUUID);
    return fix_str(ret);
}

// 获取框架版本号
char *S3_Api_GetVer(){
    char *ret = S3_Api_GetVer_Ptr(authid);
    return fix_str(ret);
}

// 获取指定QQ个人资料的年龄
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
int S3_Api_GetAge(char *selfID, char *userID){
    int ret = S3_Api_GetAge_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return ret;
}

// 获取QQ个人资料的性别
// selfID  响应QQ  文本型  无
// userID  对象QQ  文本型  无
int S3_Api_GetGender(char *selfID, char *userID){
    int ret = S3_Api_GetGender_Ptr(authid, selfID, userID);
    if (ret == -1) ret = 0;
    free(selfID);
    free(userID);
    return ret;
}

// 向好友发送窗口抖动消息
// selfID  响应QQ  文本型  机器人QQ
// userID  接收QQ  文本型  接收抖动消息的QQ
int S3_Api_ShakeWindow(char *selfID, char *userID){
    int ret = S3_Api_ShakeWindow_Ptr(authid, selfID, userID);
    free(selfID);
    free(userID);
    return ret;
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
char *S3_Api_SendMsgEX_V2(char *selfID, int messageType, char *groupID, char *userID, char *message, int bubble, int anonymous, char *jsonData){
    char *ret = S3_Api_SendMsgEX_V2_Ptr(authid, selfID, messageType, groupID, userID, message, bubble, anonymous, jsonData);
    // free(selfID);
    // free(groupID);
    // free(userID);
    // free(message);
    // free(jsonData);
    return fix_str(ret);
}

// 撤回群消息或者私聊消息
// selfID  响应QQ  文本型  机器人QQ
// sub_type  撤回类型  整数型  1好友 2群聊 4群临时会话
// groupID  参考来源  文本型  非临时会话时请留空 临时会话请填群号
// userID  参考对象  文本型  非私聊消息请留空 私聊消息请填写对方QQ号码
// messageNum  消息序号  文本型  需撤回消息序号
// messageID  消息ID  文本型  需撤回消息ID
// time  消息时间  文本型  私聊消息需要群聊时3可留空
char *S3_Api_WithdrawMsgEX(char *selfID, int sub_type, char *groupID, char *userID, char *messageNum, char *messageID, char *time){
    char *ret = S3_Api_WithdrawMsgEX_Ptr(authid, selfID, sub_type, groupID, userID, messageNum, messageID, time);
    free(selfID);
    free(groupID);
    free(userID);
    free(messageNum);
    free(messageID);
    free(time);
    return fix_str(ret);
}

// 重新从Plugin目录下载入本插件(一般用做自动更新)
int S3_Api_Reload(){
    int ret = S3_Api_Reload_Ptr(authid);
    return ret;
}

// 返回框架加载的所有插件列表(包含本插件)的json文本
char *S3_Api_GetPluginList(){
    char *ret = S3_Api_GetPluginList_Ptr(authid);
    return fix_str(ret);
}

// 查询指定对象是否允许发送在线状态临时会话 获取失败返回0 允许返回1 禁止返回2
// userID  对象QQ  文本型  无
int S3_Api_GetWpa(char *userID){
    int ret = S3_Api_GetWpa_Ptr(authid, userID);
    free(userID);
    return ret;
}

// 主动卸载插件自身
int S3_Api_Uninstall(){
    int ret = S3_Api_Uninstall_Ptr(authid);
    return ret;
}

void S3_Api_CallMessageBox(char *text) {
    TCHAR ch[1000];
    _stprintf(ch, TEXT("%s"), text);
    MessageBox(NULL, ch, TEXT("OneBot-YaYa"), 0);
    return;
}

int S3_Api_MessageBoxButton(char *text) {
    TCHAR ch[1000];
    _stprintf(ch, TEXT("%s"), text);
    return MessageBox(NULL, ch, TEXT("OneBot-YaYa"), MB_YESNO|MB_ICONQUESTION|MB_SYSTEMMODAL);
}
