#include "Windows.h"
#include "TlHelp32.h"
#include "Psapi.h"
#include "stdio.h"
#include "resource.h"
#include "easyx.h"

// 控件ID
#define IDC_BUTTON 109
#define IDC_BUTTON_D 110
#define IDC_BUTTON_coin 113
#define IDC_BUTTON_gold 115
#define IDC_EDIT_SUN 102
#define IDC_EDIT_Dmd 111
#define IDC_EDIT_coin 112
#define IDC_EDIT_gold 116
#define IDC_STATIC 103
#define IDC_STATIC_dimand 105
#define IDC_STATIC_coin 114
#define IDC_STATIC_gold 116
#define IDC_EDIT_CD 104
#define ModuleName L"PlantsVsZombies.exe"


// 全局变量 修改阳光的文本框窗口句柄、字体句柄、进程句柄、游戏进程id、图片句柄
HWND hEditSun;
HWND hEditDmd;
HWND hEditcoin;
HWND hEditgold;
HFONT hFont;
HANDLE hProcess = NULL;
DWORD gamePID = 0;
HBITMAP hBitmap;  

// 获取游戏程序imagebase基址，通过进程id和进程名称
LPVOID GetModuleBaseAddress(DWORD processID, LPCWSTR moduleName) {
	LPVOID lpBaseAddress = NULL;
	
	
	HANDLE hPrcoessb = OpenProcess(PROCESS_ALL_ACCESS, false, processID);
	if (hPrcoessb != NULL) {
		HMODULE hMods[1024];
		DWORD cbNeeded;
		if (EnumProcessModules(hPrcoessb, hMods, sizeof(hMods), &cbNeeded)) {
			DWORD dwModuleCount = cbNeeded / sizeof(HMODULE);	// 计算模块的数量

			for (DWORD i = 0; i < dwModuleCount; i++)
			{
				TCHAR szModName[MAX_PATH];
				if (GetModuleFileNameEx(hPrcoessb, hMods[i], szModName, MAX_PATH)) {
					if (wcsstr(szModName, moduleName)) {
						MODULEINFO modInfo = { 0 };
						if (GetModuleInformation(hPrcoessb, hMods[i], &modInfo, sizeof(MODULEINFO))) {
							lpBaseAddress = modInfo.lpBaseOfDll;	// 模块基地址
							break;
						}
					}
				}
			}

		}
		CloseHandle(hPrcoessb);
	}
	return lpBaseAddress;
}

// 修改金币
void UpdateGold() {
	if (hProcess == NULL) {
		MessageBox(NULL, L"游戏进程未启动", L"errorr", MB_OK | MB_ICONERROR);
		return;
	}

	wchar_t buffer[10];
	GetWindowText(hEditcoin, buffer, 10);
	int goldvalue = _wtoi(buffer);


	DWORD dwProtect;

	uintptr_t baseAddress = (uintptr_t)GetModuleBaseAddress(gamePID, ModuleName);

	// 计算第一个偏移
	uintptr_t oneNext = baseAddress + 0x2A9EC0;

	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);


	// 对基址进行第二个偏移
	oneNext += 0x82c;

	// 读取内存内容，获取下一个地址
	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);

	// 进行第三个偏移
	oneNext += 0x20c;




	// 读取内存内容，获取最终的地址
	uintptr_t goldAddress = oneNext;
	byte newbuffer[4];
	memcpy(newbuffer, &goldvalue, sizeof(goldvalue));

	// 修改内存保护权限并写入新值
	VirtualProtectEx(hProcess, (LPVOID)goldAddress, 4, PAGE_READWRITE, &dwProtect);
	WriteProcessMemory(hProcess, (LPVOID)goldAddress, newbuffer, 4, NULL);
	VirtualProtectEx(hProcess, (LPVOID)goldAddress, 4, dwProtect, &dwProtect);
}

// 修改银币
void UpdateCoin() {
	if (hProcess == NULL) {
		MessageBox(NULL, L"游戏进程未启动", L"errorr", MB_OK | MB_ICONERROR);
		return;
	}

	wchar_t buffer[10];
	GetWindowText(hEditcoin, buffer, 10);
	int coinvalue = _wtoi(buffer);


	DWORD dwProtect;

	uintptr_t baseAddress = (uintptr_t)GetModuleBaseAddress(gamePID, ModuleName);

	// 计算第一个偏移
	uintptr_t oneNext = baseAddress + 0x2A9EC0;

	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);


	// 对基址进行第二个偏移
	oneNext += 0x82c;

	// 读取内存内容，获取下一个地址
	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);

	// 进行第三个偏移
	oneNext += 0x208;




	// 读取内存内容，获取最终的地址
	uintptr_t coinAddress = oneNext;
	byte newbuffer[4];
	memcpy(newbuffer, &coinvalue, sizeof(coinvalue));

	// 修改内存保护权限并写入新值
	VirtualProtectEx(hProcess, (LPVOID)coinAddress, 4, PAGE_READWRITE, &dwProtect);
	WriteProcessMemory(hProcess, (LPVOID)coinAddress, newbuffer, 4, NULL);
	VirtualProtectEx(hProcess, (LPVOID)coinAddress, 4, dwProtect, &dwProtect);
}

// 修改钻石
void UpdateDmd() {
	if (hProcess == NULL) {
		MessageBox(NULL, L"游戏进程未启动", L"errorr", MB_OK | MB_ICONERROR);
		return;
	}

	wchar_t buffer[10];
	GetWindowText(hEditDmd, buffer, 10);
	int DmdValue = _wtoi(buffer);


	DWORD dwProtect;

	uintptr_t baseAddress = (uintptr_t)GetModuleBaseAddress(gamePID, ModuleName);

	// 计算第一个偏移
	uintptr_t oneNext = baseAddress + 0x2A9EC0;

	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);


	// 对基址进行第二个偏移
	oneNext += 0x82c;

	// 读取内存内容，获取下一个地址
	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);

	// 进行第三个偏移
	oneNext += 0x210;




	// 读取内存内容，获取最终的地址
	uintptr_t DmdAddress = oneNext;
	byte newbuffer[4];
	memcpy(newbuffer, &DmdValue, sizeof(DmdValue));

	// 修改内存保护权限并写入新值
	VirtualProtectEx(hProcess, (LPVOID)DmdAddress, 4, PAGE_READWRITE, &dwProtect);
	WriteProcessMemory(hProcess, (LPVOID)DmdAddress, newbuffer, 4, NULL);
	VirtualProtectEx(hProcess, (LPVOID)DmdAddress, 4, dwProtect, &dwProtect);
}

// 修改阳光
void UpdateSunlight() {
	if (hProcess == NULL) {
		MessageBox(NULL, L"游戏进程未启动", L"errorr", MB_OK | MB_ICONERROR);
		return;
	}

	wchar_t buffer[10];
	GetWindowText(hEditSun, buffer, 10);
	int sunValue = _wtoi(buffer);


	DWORD dwProtect;

	uintptr_t baseAddress = (uintptr_t)GetModuleBaseAddress(gamePID, ModuleName);

	// 计算第一个偏移
	uintptr_t oneNext = baseAddress + 0x2A9EC0;

	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);


	// 对基址进行第二个偏移
	oneNext += 0x768;

	// 读取内存内容，获取下一个地址
	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);

	// 进行第三个偏移
	oneNext += 0x5560;




	// 读取内存内容，获取最终的地址
	uintptr_t sunAddress = oneNext;
	byte newbuffer[4];
	memcpy(newbuffer, &sunValue, sizeof(sunValue));

	// 修改内存保护权限并写入新值
	VirtualProtectEx(hProcess, (LPVOID)sunAddress, 4, PAGE_READWRITE, &dwProtect);
	WriteProcessMemory(hProcess, (LPVOID)sunAddress, newbuffer, 4, NULL);
	VirtualProtectEx(hProcess, (LPVOID)sunAddress, 4, dwProtect, &dwProtect);
}

// cd汇编注入
void hook_code(bool bEnable) {
	DWORD dwProtect;
	uintptr_t baseAddress = (uintptr_t)GetModuleBaseAddress(gamePID, ModuleName);
	uintptr_t CodeAddr = baseAddress + 0x88E73;

	byte buffer[4] = { 0xC6,0x45,0x48,0x01 };
	byte oldbuf[4] = { 0xC6,0x45,0x48,0x00 };
	VirtualProtectEx(hProcess, (LPVOID)CodeAddr, 4, PAGE_EXECUTE_READWRITE, &dwProtect);
	// ReadProcessMemory(hProcess, (LPVOID)CodeAddr, oldbuf, 4, NULL);
	if (bEnable)
		WriteProcessMemory(hProcess, (LPVOID)CodeAddr, buffer, 4, NULL);
	if (!bEnable)
		WriteProcessMemory(hProcess, (LPVOID)CodeAddr, oldbuf, 4, NULL);
	VirtualProtectEx(hProcess, (LPVOID)CodeAddr, 4, dwProtect, &dwProtect);
}

// 修改卡槽cd
void EnableCDModification(bool bEnable) {
	if (hProcess == NULL) {
		MessageBox(NULL, L"游戏进程未启动", L"errorr", MB_OK | MB_ICONERROR);
		return;
	}


	DWORD dwProtect;

	uintptr_t baseAddress = (uintptr_t)GetModuleBaseAddress(gamePID, ModuleName);

	// 计算第一个偏移
	uintptr_t oneNext = baseAddress + 0x2A9EC0;

	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);


	// 对基址进行第二个偏移
	oneNext += 0x768;

	// 读取内存内容，获取下一个地址
	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);

	// 进行第三个偏移
	oneNext += 0x144;

	ReadProcessMemory(hProcess, (LPCVOID)oneNext, &oneNext, sizeof(uintptr_t), NULL);

	oneNext += 0x70;
	// oneNext -= 0x243550;

	// 读取内存内容，获取最终的地址
	uintptr_t sunAddress = oneNext;
	byte newbuffer[4];
	byte oldbuffer[4];

	ReadProcessMemory(hProcess, (LPCVOID)sunAddress, oldbuffer, sizeof(oldbuffer), NULL);


	if (bEnable) {

		for (int i = 0; i < 12; i++)
		{
			int nCD = 1;
			memcpy(newbuffer, &nCD, sizeof(nCD));

			VirtualProtectEx(hProcess, (LPVOID)sunAddress, 4, PAGE_READWRITE, &dwProtect);
			WriteProcessMemory(hProcess, (LPVOID)sunAddress, newbuffer, sizeof(newbuffer), NULL);
			VirtualProtectEx(hProcess, (LPVOID)sunAddress, 4, dwProtect, &dwProtect);
			sunAddress += 0x50;
		}

		hook_code(bEnable);

	}
	else
	{

		memcpy(newbuffer, oldbuffer, sizeof(oldbuffer));

		VirtualProtectEx(hProcess, (LPVOID)sunAddress, 4, PAGE_READWRITE, &dwProtect);
		WriteProcessMemory(hProcess, (LPVOID)sunAddress, newbuffer, sizeof(newbuffer), NULL);
		VirtualProtectEx(hProcess, (LPVOID)sunAddress, 4, dwProtect, &dwProtect);

		hook_code(bEnable);
	}





}


// 找到游戏进程id
DWORD FindGameProcessPID(const wchar_t* processName) {
	DWORD pid = 0;
	HANDLE hSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
	if (hSnapshot != INVALID_HANDLE_VALUE) {
		PROCESSENTRY32 pe32;
		pe32.dwSize = sizeof(PROCESSENTRY32);

		if (Process32First(hSnapshot, &pe32)) {
			do {
				if (_wcsicmp(pe32.szExeFile, processName) == 0) {
					pid = pe32.th32ProcessID;
					break;
				}
			} while (Process32Next(hSnapshot, &pe32));
		}
		CloseHandle(hSnapshot);
	}

	return pid;
}


LRESULT CALLBACK WindowProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam) {
	switch (uMsg)
	{
	case WM_CREATE:
	{

		// 字体修改
		hFont = CreateFont(
			-MulDiv(14, GetDeviceCaps(GetDC(hwnd), LOGPIXELSY), 72), // 字体高度
			0, 0, 0, FW_NORMAL, FALSE, FALSE, FALSE, DEFAULT_CHARSET,
			OUT_DEFAULT_PRECIS, CLIP_DEFAULT_PRECIS, CLEARTYPE_QUALITY, // 使用 ClearType
			DEFAULT_PITCH | FF_SWISS, L"微软雅黑");


		// 阳光控件
		HWND hStatic = CreateWindowEx(0, L"STATIC", L"阳光值", WS_CHILD | WS_VISIBLE, 10, 50, 50, 25, hwnd, (HMENU)IDC_STATIC, GetModuleHandle(NULL), NULL);
		// 将指定消息发送到一个或多个窗口，参数1窗口，参数2消息类型，参数3消息内容，
		SendMessage(hStatic, WM_SETFONT, (WPARAM)hFont, TRUE);

		hEditSun = CreateWindowEx(0, L"EDIT", L"999", WS_CHILD | WS_VISIBLE | WS_BORDER | ES_NUMBER,
			80, 50, 100, 25, hwnd, (HMENU)IDC_EDIT_SUN, GetModuleHandle(NULL), NULL
		);
		SendMessage(hEditSun, WM_SETFONT, (WPARAM)hFont, TRUE);

		HWND hButton = CreateWindowEx(0, L"button", L"修改", WS_CHILD | WS_VISIBLE | BS_PUSHBUTTON,
			200, 50, 50, 25, hwnd, (HMENU)IDC_BUTTON, GetModuleHandle(NULL), NULL
		);
		SendMessage(hButton, WM_SETFONT, (WPARAM)hFont, TRUE);



		// 钻石控件
		HWND hStatic_dimand = CreateWindowEx(0, L"STATIC", L"钻石数量", WS_CHILD | WS_VISIBLE, 10, 150, 50, 25, hwnd, (HMENU)IDC_STATIC_dimand, GetModuleHandle(NULL), NULL);
		SendMessage(hStatic_dimand, WM_SETFONT, (WPARAM)hFont, TRUE);

		hEditDmd = CreateWindowEx(0, L"EDIT", L"999", WS_CHILD | WS_VISIBLE | WS_BORDER | ES_NUMBER,
			80, 150, 100, 25, hwnd, (HMENU)IDC_EDIT_Dmd, GetModuleHandle(NULL), NULL
		);
		SendMessage(hEditDmd,WM_SETFONT, (WPARAM)hFont, TRUE);

		HWND hButton1 = CreateWindowEx(0, L"button", L"修改", WS_CHILD | WS_VISIBLE | BS_PUSHBUTTON,
			200, 150, 50, 25, hwnd, (HMENU)IDC_BUTTON_D, GetModuleHandle(NULL), NULL
		);
		SendMessage(hButton1, WM_SETFONT, (WPARAM)hFont, TRUE);

		// 银币控件
		HWND hStatic_coin = CreateWindowEx(0, L"STATIC", L"银币数量", WS_CHILD | WS_VISIBLE, 10, 200, 50, 25, hwnd, (HMENU)IDC_STATIC_coin, GetModuleHandle(NULL), NULL);
		SendMessage(hStatic_coin, WM_SETFONT, (WPARAM)hFont, TRUE);

		hEditcoin = CreateWindowEx(0, L"EDIT", L"999", WS_CHILD | WS_VISIBLE | WS_BORDER | ES_NUMBER,
			80, 200, 100, 25, hwnd, (HMENU)IDC_EDIT_coin, GetModuleHandle(NULL), NULL
		);
		SendMessage(hEditcoin, WM_SETFONT, (WPARAM)hFont, TRUE);

		HWND hButton2 = CreateWindowEx(0, L"button", L"修改", WS_CHILD | WS_VISIBLE | BS_PUSHBUTTON,
			200, 200, 50, 25, hwnd, (HMENU)IDC_BUTTON_coin, GetModuleHandle(NULL), NULL
		);
		SendMessage(hButton2, WM_SETFONT, (WPARAM)hFont, TRUE);

		// 金币控件
		HWND hStatic_gold = CreateWindowEx(0, L"STATIC", L"金币数量", WS_CHILD | WS_VISIBLE, 10, 250, 50, 25, hwnd, (HMENU)IDC_STATIC_gold, GetModuleHandle(NULL), NULL);
		SendMessage(hStatic_gold, WM_SETFONT, (WPARAM)hFont, TRUE);

		hEditgold = CreateWindowEx(0, L"EDIT", L"999", WS_CHILD | WS_VISIBLE | WS_BORDER | ES_NUMBER,
			80, 250, 100, 25, hwnd, (HMENU)IDC_EDIT_gold, GetModuleHandle(NULL), NULL
		);
		SendMessage(hEditgold, WM_SETFONT, (WPARAM)hFont, TRUE);

		HWND hButton3 = CreateWindowEx(0, L"button", L"修改", WS_CHILD | WS_VISIBLE | BS_PUSHBUTTON,
			200, 250, 50, 25, hwnd, (HMENU)IDC_BUTTON_gold, GetModuleHandle(NULL), NULL
		);
		SendMessage(hButton3, WM_SETFONT, (WPARAM)hFont, TRUE);


		// 卡槽cd控件
		HWND hCheckbox = CreateWindowEx(0, L"BUTTON", L"启用植物卡槽CD修改",
			WS_CHILD | WS_VISIBLE | BS_AUTOCHECKBOX,
			10, 100, 200, 25, hwnd, (HMENU)IDC_EDIT_CD,
			GetModuleHandle(NULL), NULL);
		SendMessage(hCheckbox, WM_SETFONT, (WPARAM)hFont, TRUE);

		gamePID = FindGameProcessPID(L"PlantsVsZombies.exe");

		if (gamePID != 0) {
			hProcess = OpenProcess(PROCESS_ALL_ACCESS, false, gamePID);
		}

		if (hProcess == NULL) {
			wchar_t debug[256];
			wsprintf(debug, L"this is gamepid %d\n", gamePID);
			OutputDebugString(debug);
		}
	}
	break;
	// 当用户从菜单中选择命令项、控件将通知消息发送到其父窗口或转换加速键时发送WM_COMMAND。
	case WM_COMMAND:
	{
		// loword从指定值检索低字序，如果是wmcommand消息，则wparam的低字为菜单的标识符id
		switch (LOWORD(wParam))
		{
		case IDC_BUTTON:
			UpdateSunlight();
			break;

		case IDC_BUTTON_D:
			UpdateDmd();
			break;
		case IDC_BUTTON_coin:
			UpdateCoin();
			break;
		case IDC_BUTTON_gold:
			UpdateGold();
			break;

		case IDC_EDIT_CD:
			if (SendMessage((HWND)lParam, BM_GETCHECK, 0, 0) == BST_CHECKED) {
				// 复选框被选中，启用CD修改

				EnableCDModification(true);
			}
			else {
				// 复选框未选中，禁用CD修改

				EnableCDModification(false);
			}
			break;
		}
		break;
	}

	case WM_PAINT:
	{
		PAINTSTRUCT ps;
		HDC hdc = BeginPaint(hwnd, &ps);

		// 创建兼容DC
		HDC hdcMem = CreateCompatibleDC(hdc);
		HBITMAP hbmOld = (HBITMAP)SelectObject(hdcMem, hBitmap);

		// 获取窗口大小
		RECT rect;
		GetClientRect(hwnd, &rect);

		// 绘制背景图
		BitBlt(hdc, 0, 0, rect.right, rect.bottom, hdcMem, 0, 0, SRCCOPY);

		// 清理
		SelectObject(hdcMem, hbmOld);
		DeleteDC(hdcMem);

		EndPaint(hwnd, &ps);
	}
	break;

	case WM_DESTROY:
		if (hProcess != NULL) {
			CloseHandle(hProcess);
		}
		PostQuitMessage(0);
		break;

	default:
		return DefWindowProc(hwnd, uMsg, wParam, lParam);
	}

	return 0;
}

int WINAPI wWinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, PWSTR pCmdLine, int nCmdShow) {
	const wchar_t title[] = L"PVZ MOD";

	WNDCLASS wc = {};

	wc.lpfnWndProc = WindowProc;
	wc.lpszClassName = title;
	wc.hInstance = hInstance;
	wc.hbrBackground = (HBRUSH)(COLOR_WINDOW + 1);

	RegisterClass(&wc);

	
	int screenWidth = GetSystemMetrics(SM_CXSCREEN);
	int screenHeight = GetSystemMetrics(SM_CYSCREEN);

	
	int windowWidth = 800;
	int windowHeight = 500;

	
	int xPos = (screenWidth - windowWidth) / 2;
	int yPos = (screenHeight - windowHeight) / 2;

	HWND hwnd = CreateWindowEx(
		0, title, L"pve", WS_OVERLAPPEDWINDOW,
		xPos, yPos, windowWidth, windowHeight, // 居中显示
		NULL, NULL, hInstance, NULL
	);

	if (hwnd == NULL) {
		return 0;
	}
	hBitmap = LoadBitmap(hInstance, MAKEINTRESOURCE(IDB_BITMAP1)); 

	// 当调用showwindow时会发送wm_print消息
	ShowWindow(hwnd, nCmdShow);

	

	MSG msg = {};
	while (GetMessage(&msg, NULL, 0, 0)) {
		TranslateMessage(&msg);
		DispatchMessage(&msg);
	}

	// initgraph(680, 480);

	return 0;
}
