#include <stdio.h>
#include <stdlib.h>
#include <windows.h>

#define XQAPI(RetType, Name, ...)													\
	typedef RetType(__stdcall *Name##_Type)(unsigned char * authid, ##__VA_ARGS__); \
	Name##_Type Name##_Ptr;															\
	RetType Name(__VA_ARGS__);

#define LoadAPI(Name) Name##_Ptr = (Name##_Type)GetProcAddress(hmod, #Name)

unsigned char * authid;
int id;
XQAPI(void, S3_Api_SendMsg, char *, int, char *, char *, char *, int);
XQAPI(char * , S3_Api_GetGroupList, char *);
XQAPI(void, S3_Api_OutPutLog, char *);

// sender
XQAPI(char *, S3_Api_GetNick, char *, char *);
XQAPI(int, S3_Api_GetGender, char *, char *);
XQAPI(int, S3_Api_GetAge, char *, char *);

extern void __stdcall XQ_AuthId(int ID, int IMAddr){
	authid = (unsigned char *)malloc(sizeof(unsigned char)*16);
	*((int*)authid) = 1;
	*((int*)(authid + 4)) = 8;
	*((int*)(authid + 8)) = ID;
	*((int*)(authid + 12)) = IMAddr;
	authid += 8;
	HMODULE hmod = LoadLibraryA("xqapi.dll");
	LoadAPI(S3_Api_SendMsg);
	LoadAPI(S3_Api_OutPutLog);
	LoadAPI(S3_Api_GetGroupList);

	// sender
	LoadAPI(S3_Api_GetNick);
	LoadAPI(S3_Api_GetGender);
	LoadAPI(S3_Api_GetAge);

	id = ID;
	return;
}

void S3_Api_SendMsg(char * selfID, int messageType, char * groupID, char * userID, char * message, int bubble){
	S3_Api_SendMsg_Ptr(authid, selfID, messageType, groupID, userID, message, bubble);
	free(selfID);
	free(groupID);
	free(userID);
	free(message);
}

void S3_Api_OutPutLog(char * message){
	S3_Api_OutPutLog_Ptr(authid, message);
	free(message);
}

char * S3_Api_GetGroupList(char * selfID) {
	char * group_list = S3_Api_GetGroupList_Ptr(authid, selfID);
	free(selfID);
	return group_list;
}

// sender
char * S3_Api_GetNick(char * selfID, char * userID) {
	char * nickname = S3_Api_GetNick_Ptr(authid, selfID, userID);
	free(selfID);
	free(userID);
	return nickname;
}

int S3_Api_GetGender(char * selfID, char * userID) {
	int sex = S3_Api_GetGender_Ptr(authid, selfID, userID);
	free(selfID);
	free(userID);
	return sex;
}

int S3_Api_GetAge(char * selfID, char * userID) {
	int age = S3_Api_GetAge_Ptr(authid, selfID, userID);
	free(selfID);
	free(userID);
	return age;
}
