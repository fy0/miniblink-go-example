/*
   Copyright 2018 fy

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"os"
	"path/filepath"
)

/*
#cgo CFLAGS: -g -std=c99

#include <stdio.h>
#include <stdbool.h>
#include <stdint.h>
#include <windows.h>
#include "wkedefine.h"

typedef struct {
    wkeWebView window;
    wchar_t url[MAX_PATH + 1];
} Application;

wchar_t* u8towcs(const char* szU8) {
    //预转换，得到所需空间的大小;
    int wcsLen = MultiByteToWideChar(CP_UTF8, 0, szU8, strlen(szU8), NULL, 0);
    //分配空间要给'\0'留个空间，MultiByteToWideChar不会给'\0'空间
    wchar_t* wszString = (wchar_t*)malloc(sizeof(wchar_t) * (wcsLen + 1));
    //转换
    MultiByteToWideChar(CP_UTF8, 0, szU8, strlen(szU8), wszString, wcsLen);
    //最后加上'\0'
    wszString[wcsLen] = '\0';
    return wszString;
}

// 回调：窗口已销毁
void HandleWindowDestroy(wkeWebView webWindow, void* param) {
    Application* app = (Application*)param;
    app->window = NULL;
    PostQuitMessage(0);
}

// 回调：文档加载成功
void HandleDocumentReady(wkeWebView webWindow, void* param) {
    wkeShowWindow(webWindow, true);
}

// 创建主页面窗口
bool CreateWebWindow(Application* app, int width, int height) {
	// WKE_WINDOW_TYPE_TRANSPARENT
    app->window = wkeCreateWebWindow(WKE_WINDOW_TYPE_POPUP, NULL, 0, 0, width, width);

    if (!app->window)
        return FALSE;

    // wkeOnWindowClosing(app->window, HandleWindowClosing, app); // 弹窗
    wkeOnWindowDestroy(app->window, HandleWindowDestroy, app);
    wkeOnDocumentReady(app->window, HandleDocumentReady, app);
    // wkeOnTitleChanged(app->window, HandleTitleChanged, app);
    // wkeOnCreateView(app->window, HandleCreateView, app);
    // wkeOnLoadUrlBegin(app->window, HandleLoadUrlBegin, app);
    // wkeOnLoadUrlEnd(app->window, HandleLoadUrlEnd, app);

    wkeMoveToCenter(app->window);
    wkeLoadURLW(app->window, app->url);

    return true;
}

void PrintHelpAndQuit(Application* app) {
    PostQuitMessage(0);
}

void RunMessageLoop(Application* app) {
    MSG msg = { 0 };
    while (GetMessageW(&msg, NULL, 0, 0)) {
        TranslateMessage(&msg);
        DispatchMessageW(&msg);
    }
}

void QuitApplication(Application* app) {
    if (app->window) {
        wkeDestroyWebWindow(app->window);
        app->window = NULL;
    }
}

void* application_new(char *u8url, int width, int height) {
    Application *app = malloc(sizeof(Application));
    memset(app, 0, sizeof(Application));

	wchar_t *url = u8towcs(u8url);
#ifdef _MSC_VER
    wcsncpy_s(app->url, MAX_PATH, url, MAX_PATH);
#else
    wcsncpy(app->url, url, MAX_PATH);
#endif
	free(url);

    if (!CreateWebWindow(app, width, height)) {
        PrintHelpAndQuit(app);
        return NULL;
    }

    return (void*)app;
}

void application_run(void *app) {
    RunMessageLoop((Application*)app);
}

void wke_init() {
	wkeInitialize();
}

void wke_final() {
	wkeFinalize();
}

jsValue JS_CALL js_msgBox(jsExecState es) {
	const wchar_t* text = jsToStringW(es, jsArg(es, 0));
	const wchar_t* title = jsToStringW(es, jsArg(es, 1));

	MessageBoxW(NULL, text, title, 0);

	//return jsUndefined();
	//return jsInt(1234);
	return jsStringW(es, L"C++返回字符串");
}

void set_devtools(void *app, char *u8path) {
	wkeSetDebugConfig(((Application*)app)->window, "showDevTools", u8path);
}

void test_bind() {
	jsBindFunction("msgBox", &js_msgBox, 2);
}
*/
import "C"

func getDevtoolsPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "devtools/inspector.html")
}

func main() {
	C.wke_init()
	app := C.application_new(C.CString("http://www.baidu.com"), 800, 600)
	C.set_devtools(app, C.CString(getDevtoolsPath()))
	C.test_bind()
	C.application_run(app)
	C.wke_final()
}
