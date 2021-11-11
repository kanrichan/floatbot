#include <stdint.h>
#include <windows.h>

extern char* __stdcall __declspec (dllexport) XQ_Create(char *version);
extern int __stdcall __declspec (dllexport) XQ_Event(char *self_id, int message_type, int sub_type, char *group_id, char *user_id, char *notice_id, char *message, char *message_num, char *message_id, char *raw_message, char *time, int ret);
extern int __stdcall __declspec (dllexport) XQ_SetUp();
extern int __stdcall __declspec (dllexport) XQ_DestroyPlugin();