;通过调用system执行打开bilibili网站

_start:
    ; find kernel32.dll address
    xor eax, eax;
    mov ebx, fs: [eax + 30h] ; peb addr
    mov ebx, [ebx + 0ch]; ldr addr
    mov ebx, [ebx + 0ch]; inload结构体addr
    mov ebx, [ebx]; ntdll 模块
    mov ebx, [ebx]; kernel32 模块
    mov eax, [ebx + 18h]; kernel32 addr
    mov esi, eax;
    mov ebx, eax;
    ; find getprocess addr

    add esi, 3ch;
    add ebx, [esi]; pe addr
    lea ebx, [ebx + 0x78]; 18h--option - header, 60h--data - directory addr
    mov ebx, [ebx];
    add ebx, eax;
    mov esi, [ebx + 0x20]; 名称字段偏移量
    add esi, eax;

    xor ecx, ecx;

GetFunc:
    inc ecx;
    mov edx, [esi];
    add esi, 4;
    add edx, eax; edx位name数据实际地址
    mov edi, [edx];
    cmp edi, 0x50746547;
    jnz GetFunc;
    mov edi, [edx + 0x4];
    cmp edi, 0x41636f72;
    jnz GetFunc;
    mov edi, [edx + 0x8];
    cmp edi, 0x65726464;
    jnz GetFunc;


    mov esi, [ebx + 0x24]; num - offset
    add esi, eax; num - addr
    mov cx, [esi + ecx * 2];
    dec ecx;

    mov esi, [ebx + 0x1c]; addr - offset
    add esi, eax; addr - addr
    mov edx, [esi + ecx * 4];
    add edx, eax;

    ; find loadlibaddr
    xor ecx, ecx;
    push ebp;
    mov ebp, esp;
    push edx;
    push eax;
    push ecx;
    push 0x41797261;
    push 0x7262694c;
    push 0x64616f4c;
    push esp;
    push eax;
    call edx;

    add esp, 0xc;
    pop ecx;
    push eax; loadaddr
    push ecx;
    mov cx, 0x6c6c;
    push ecx;
    push 0x642e7472;
    push 0x6376736d;
    push esp;
    call eax;

    add esp, 0xc;
    mov cx, 0x6d65;
    push ecx;
    push 0x74737973;
    push esp;
    push eax;
    mov ebx, [esp + 0x1c]
    call ebx;

    add esp, 0x8;

    xor ecx, ecx;
    push 0x002f6d6f;
    push 0x632e696c;
    push 0x6962696c;
    push 0x69622e77;
    push 0x77772f2f;
    push 0x3a707474;
    push 0x68206578;
    push 0x652e6572;
    push 0x6f6c7078;
    push 0x65692074;
    push 0x72617473;
    push esp;
    call eax;