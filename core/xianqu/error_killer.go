package xianqu

/*
#include <string.h>
#include <windows.h>
#include <tchar.h>
void kill(char *text) {
	TCHAR ch[100];
    _stprintf(ch, TEXT("%s"), text);
    SendMessage(FindWindow(NULL,ch),WM_CLOSE,0,0);
    return;
}
*/
import "C"
import "time"

func errorkiller() {
	for {
		time.Sleep(time.Second * 5)
		C.kill(cString("先驱框架运行时出现异常！"))
	}
}
