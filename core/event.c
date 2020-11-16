#include <stdint.h>
extern char* GO_Create(char* version);
extern int GO_Event(char * selfID, int mseeageType, int subType, char * groupID, char * userID, char * noticID, char * message, char * messageNum, char* messageID, char* rawMessage, char* time, char* ret);
extern int GO_SetUp();
extern int GO_DestroyPlugin();

extern char* __stdcall XQ_Create(char* version);
extern int __stdcall XQ_Event(char * selfID, int mseeageType, int subType, char * groupID, char * userID, char * noticID, char * message, char * messageNum, char* messageID, char* rawMessage, char* time, char* ret);
extern int __stdcall XQ_SetUp();
extern int __stdcall XQ_DestroyPlugin();

char* _stdcall XQ_Create(char* version)
{
	return GO_Create(version);
}

int _stdcall XQ_Event(char * selfID, int mseeageType, int subType, char * groupID, char * userID, char * noticID, char * message, char * messageNum, char* messageID, char* rawMessage, char* time, char* ret)
{
	return GO_Event(selfID, mseeageType, subType, groupID, userID, noticID, message, messageNum, messageID, rawMessage, time, ret);
}

int _stdcall XQ_SetUp()
{
	return GO_DestroyPlugin();
}

int _stdcall XQ_DestroyPlugin()
{
	return GO_SetUp();
}