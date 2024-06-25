new Vue({
    el: '#app',
    data: {
        url: '',
        key: '',
        commandText: '',
        path: '',
        current: '',
        lines: [],
        kines: [],
        errorMessage: '',
        errorsMessage: '',
        tureMessage: '',
        checkFile: '',
        currentPage: 1,
        pageSize: 10,
        contextMenuVisible: false,
        selectedFile: null,
        showDetails: false,
        linkStart: '',
        cmd: false,
        cmdResult: '',
        ter: false,
        isDirectory: false
    },

    computed: {
        totalPages() {
            return Math.ceil(this.lines.length / this.pageSize);
        },
        paginatedLines() {
            if (this.lines.length === 0) {
                return [{
                    name: '',
                    type: '',
                    time: '',
                    size: ''
                }];
            }
            const start = (this.currentPage - 1) * this.pageSize;
            const end = start + this.pageSize;
            return this.lines.slice(start, end);
        }
    },
    methods: {
        executeCommand() {
            this.showDetails = false
            this.errorsMessage = '';
            this.tureMessage = '';
            this.linkStart = "123"
            this.cmd = false
            this.ter = false
            this.cmdResult = ''
            fetch("/", {
                method: "POST",
                body: new URLSearchParams({
                    url: this.url,
                    key: this.key,
                    linkStart: this.linkStart
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            this.errorsMessage = text; // 显示具体的错误信息
                            throw new Error('Network response was not ok');
                        });
                    }
                    return response.text();
                })
                .then(data => {
                    const part = data.split('------')
                    this.tureMessage = part[1]
                    this.path = part[0]
                })
                .catch(error => {
                    // 这里处理网络错误或者由throw new Error抛出的错误
                    console.error('There has been a problem with your fetch operation:', error);

                });
        },

        cmdT() {
            this.errorMessage = '';
            this.errorsMessage = '';
            this.tureMessage = '';
            this.cmd = true;
            this.ter = false;
            this.showDetails = false;
        },

        // 命令执行主界面
        command() {
            this.errorMessage = '';
            this.errorsMessage = '';
            this.tureMessage = '';
            this.cmd = true
            this.checkFile = 'cmd'
            fetch("/", {
                method: "POST",
                body: new URLSearchParams({
                    url: this.url,
                    key: this.key,
                    command: this.commandText,
                    checkFile: this.checkFile
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            this.errorsMessage = text; // 显示具体的错误信息
                            throw new Error('Network response was not ok');
                        });
                    }
                    return response.text();
                })
                .then(data => {
                    this.cmdResult = data.split('\n');
                })
                .catch(error => {
                    // 这里处理网络错误或者由throw new Error抛出的错误
                    console.error('There has been a problem with your fetch operation:', error);
                });

        },

        // 文件主界面
        execute() {
            this.errorsMessage = '';
            this.tureMessage = '';
            this.errorMessage = '';
            this.checkFile = 'home';
            this.ter = false;
            this.path = '';
            this.showDetails = true;
            this.cmd = false;
            this.cmdResult = '';
            fetch("/", {
                method: "POST",
                body: new URLSearchParams({
                    url: this.url,
                    key: this.key,
                    checkFile: this.checkFile
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            this.errorMessage = text; // 显示具体的错误信息
                            this.lines = []; // 清除旧的文件信息
                            throw new Error('Network response was not ok');
                        });
                    }
                    return response.text();
                })
                .then(data => {

                    const parts = data.split('---SPLIT---');
                    this.path = parts[0];
                    this.getPath(this.path);
                    this.lines = this.parseLines(parts[1]);
                    this.currentPage = 1; // Reset to first page after loading new data
                })
                .catch(error => {
                    // 这里处理网络错误或者由throw new Error抛出的错误
                    this.lines = []; // 清除旧的文件信息
                    console.error('There has been a problem with your fetch operation:', error);

                });
        },

        // 处理点击文件操作
        handleClick(item) {
            const fileExtension = item.name.split('.').pop().toLowerCase();
            const downloadExtensions = ['zip', 'exe', 'bat', 'cmd', 'go', 'jpg', 'png'];
            if (item.type === '文件' && downloadExtensions.includes(fileExtension)) {
                // 如果文件后缀为.zip、.exe、.bat或.cmd，则执行下载操作
                this.downloadFile(item.name);
            } else {
                // 否则执行默认的预览操作
                if (item.type === '文件') {
                    this.handleFile(item.name);
                } else {
                    this.errorMessage = '';
                    fetch("/", {
                        method: "POST",
                        body: new URLSearchParams({
                            clickedWord: item.name,
                            url: this.url,
                            key: this.key,
                            path: this.path,
                            current: this.current
                        }),
                        headers: {
                            "Content-Type": "application/x-www-form-urlencoded"
                        }
                    })
                        .then(response => {
                            if (!response.ok) {
                                return response.text().then(text => {
                                    this.errorMessage = text; // 显示具体的错误信息
                                    this.lines = []; // 清除旧的文件信息
                                    throw new Error('Network response was not ok');
                                });
                            }
                            return response.text();
                        })
                        .then(data => {
                            const parts = data.split('---SPLIT---');
                            this.path = parts[0];
                            this.getPath(this.path);
                            this.lines = this.parseLines(parts[1]);
                            this.currentPage = 1; // Reset to first page after loading new data
                        })
                        .catch(error => {
                            this.lines = [];
                            console.error("Error:", error);
                        });
                }
            }
        },
        handleInput(name) {
            this.errorMessage = '';
            fetch("/", {
                method: "POST",
                body: new URLSearchParams({
                    clickedWord: name,
                    url: this.url,
                    key: this.key,
                    path: this.path,
                    current: this.current
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            this.errorMessage = text; // 显示具体的错误信息
                            this.lines = []; // 清除旧的文件信息
                            throw new Error('Network response was not ok');
                        });
                    }
                    return response.text();
                })
                .then(data => {
                    const parts = data.split('---SPLIT---');
                    this.path = parts[0];
                    this.getPath(this.path);
                    this.lines = this.parseLines(parts[1]);
                    this.currentPage = 1; // Reset to first page after loading new data
                })
                .catch(error => {
                    this.lines = []; // 清除旧的文件信息
                    console.error("Error:", error);
                });
        },

        handleRootPath(name, root) {
            this.errorMessage = '';
            fetch("/", {
                method: "POST",
                body: new URLSearchParams({
                    clickedWord: name,
                    url: this.url,
                    key: this.key,
                    path: this.path,
                    rootPath: root
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            this.errorMessage = text; // 显示具体的错误信息
                            this.lines = []; // 清除旧的文件信息
                            throw new Error('Network response was not ok');
                        });
                    }
                    return response.text();
                })
                .then(data => {
                    const parts = data.split('---SPLIT---');
                    this.path = parts[0];
                    this.getPath(this.path);
                    this.lines = this.parseLines(parts[1]);
                    this.currentPage = 1; // Reset to first page after loading new data
                })
                .catch(error => {
                    this.lines = []; // 清除旧的文件信息
                    console.error("Error:", error);
                });
        },

        handleFile(name) {
            this.checkFile = 'view'
            fetch("/", {
                method: "POST",
                body: new URLSearchParams({
                    clickedWord: name,
                    url: this.url,
                    key: this.key,
                    path: this.path,
                    checkFile: this.checkFile
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => response.text())
                .then(data => {
                    const parts = data.split('---SPLIT---');

                    const escapeHtml = (unsafe) => {
                        return unsafe
                            .replace(/&/g, "&amp;")
                            .replace(/</g, "&lt;")
                            .replace(/>/g, "&gt;")
                            .replace(/"/g, "&quot;")
                            .replace(/'/g, "&#039;");
                    };


                    const newWindow = window.open("", "_blank");
                    newWindow.document.write("<pre>" + escapeHtml(parts[1]) + "</pre>");
                    newWindow.document.title = name;
                })
                .catch(error => {
                    console.error("Error:", error);
                });
        },

        homePath() {
            if (this.path.startsWith('/')) {
                this.path = "/";
                this.handleRootPath(this.path, "home");
            }
            const drivePattern = /^[A-Za-z]:\\/;
            if (drivePattern.test(this.path)) {
                this.path = this.path.slice(0, 3);
                this.handleRootPath(this.path, "home");
            }
            this.handleRootPath(this.path, "home");
        },

        UpDir() {
            if (this.path.endsWith('/')) {
                this.path = this.path.slice(0, -1);
            }
            const isWindows = /^[A-Za-z]:/;
            if (isWindows.test(this.path)) {
                // Windows 路径
                const lastBackslashIndex = this.path.lastIndexOf('/');
                if (lastBackslashIndex > 2) {
                    this.path = this.path.slice(0, lastBackslashIndex);
                } else {
                    // 保持根目录（例如 C:\）
                    this.path += '/'
                    this.path = this.path.slice(0, 2);
                }
            } else {
                // Linux 路径
                const lastSlashIndex = this.path.lastIndexOf('/');
                if (lastSlashIndex > 0) {
                    this.path = this.path.slice(0, lastSlashIndex);
                } else {
                    // 保持根目录
                    this.path = "";
                }
            }
            this.path += '/'
            this.handleRootPath(this.path, "up");
        },

        getPath(path) {
            this.current = path;
        },

        changePath() {
            this.handleInput("reset")
        },

        parseLines(data) {
            const lines = data.split('\n').filter(line => line.trim() !== '' && !['./', '../'].includes(line.trim().split('\t')[0]));
            return lines.map(line => {
                const parts = line.split('\t');
                return {
                    name: parts[0],
                    type: parts[0].endsWith('/') ? '目录' : '文件',
                    time: parts[1] || '0',
                    size: parts[2] || '0'
                };
            });
        },

        submitUpload() {
            const input = document.getElementById('file');
            if (input.files.length === 0) {
                alert('请选择要上传的文件');
                return;
            }

            // Append files to FormData
            let formData = new FormData();
            for (const file of input.files) {
                formData.append('file', file);
            }

            // Append other data to FormData
            formData.append('url', this.url);
            formData.append('path', this.path);

            // Send FormData via fetch
            fetch("/upload", {
                method: "POST",
                body: formData
            })
                .then(response => response.text())
                .then(data => {
                    alert(data)
                    this.handleInput("reset")
                    console.log("Upload response:", data);
                    // Handle response if needed
                })
                .catch(error => {
                    console.error("Error uploading files:", error);
                    // Handle error if needed
                });
        },
        nextPage() {
            if (this.currentPage < this.totalPages) {
                this.currentPage++;
            }
        },
        prevPage() {
            if (this.currentPage > 1) {
                this.currentPage--;
            }
        },
        showContextMenu(event, file) {
            this.selectedFile = file;
            this.isDirectory = file.type === '目录';
            const contextMenu = this.$refs.contextMenu;
            contextMenu.style.display = 'block';
            contextMenu.style.left = event.clientX + 'px';
            contextMenu.style.top = event.clientY + 'px';
            this.contextMenuVisible = true;
            document.addEventListener('click', this.hideContextMenu);
        },
        hideContextMenu() {
            this.$refs.contextMenu.style.display = 'none';
            this.contextMenuVisible = false;
            document.removeEventListener('click', this.hideContextMenu);
        },
        downloadFile(name) {
            // const form = document.getElementById('downloadForm');
            // const formData = new FormData(form);
            this.checkFile = 'download'
            fetch('/', {
                method: "POST",
                body: new URLSearchParams({
                    clickedWord: name,
                    url: this.url,
                    key: this.key,
                    path: this.path,
                    checkFile: this.checkFile
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    const disposition = response.headers.get('Content-Disposition');
                    const filename = disposition ? disposition.match(/filename="(.+)"/)[1] : 'downloaded_file';

                    return response.blob().then(blob => {
                        const url = window.URL.createObjectURL(blob);
                        const a = document.createElement('a');
                        a.style.display = 'none';
                        a.href = url;
                        a.download = filename;
                        document.body.appendChild(a);
                        a.click();
                        window.URL.revokeObjectURL(url);
                    });
                })
                .catch(error => {
                    console.error('There was a problem with the fetch operation:', error);
                });
        },
        deleteFile(name) {
            this.checkFile = 'delete'
            fetch('/', {
                method: "POST",
                body: new URLSearchParams({
                    clickedWord: name,
                    url: this.url,
                    key: this.key,
                    path: this.path,
                    checkFile: this.checkFile
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            this.errorMessage = text; // 显示具体的错误信息
                            // this.lines = []; // 清除旧的文件信息
                            throw new Error('file delete was not ok');
                        });
                    }
                    return response.text();
                })
                .then(data => {
                    alert(data)
                    this.handleInput("reset")
                    console.log("Upload response:", data);
                    // Handle response if needed
                })
                .catch(error => {
                    console.error("Error uploading files:", error);
                    // Handle error if needed
                });
            this.hideContextMenu();
        },

        handleCommand(command) {
            if (command.toLowerCase() === 'clear') {
                document.getElementById('output').innerHTML = ''; // 清空显示的内容
                this.createInputLine(); // 创建新的输入行
                return;
            }
            if (command.toLowerCase() === 'exit') {
                this.ter = false; // 关闭终端窗口
                return;
            }
            this.checkFile = 'terminal'
            fetch("/", {
                method: "POST",
                body: new URLSearchParams({
                    url: this.url,
                    key: this.key,
                    command: command,
                    path: this.path,
                    checkFile: this.checkFile
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error('Network response was not ok');
                    }
                    return response.text();
                })
                .then(data => {
                    const part = data.split('------')
                    this.path = part[1].replace(/\s+/g, '');

                    const responseLine = document.createElement('div');
                    responseLine.textContent = part[0];
                    document.getElementById('output').appendChild(responseLine);
                    this.createInputLine();
                })
                .catch(error => {
                    console.error('There has been a problem with your fetch operation:', error);
                    const errorLine = document.createElement('div');
                    errorLine.textContent = error.message;
                    document.getElementById('output').appendChild(errorLine);
                    this.createInputLine();
                });
        },

        createInputLine() {
            this.errorMessage = '';
            this.errorsMessage = '';
            this.tureMessage = '';
            this.cmdResult = '';
            this.ter = true;
            this.cmd = false;
            this.showDetails = false;
            const output = document.getElementById('output');
            const inputLine = document.createElement('div');
            inputLine.className = 'input-line';

            const prompt = document.createElement('span');
            prompt.className = 'prompt';
            prompt.textContent = `${this.path} $`;

            const input = document.createElement('input');
            input.type = 'text';
            input.className = 'input';
            input.autofocus = true;

            input.addEventListener('keydown', (event) => {
                if (event.key === 'Enter') {
                    const command = input.value.trim();
                    inputLine.innerHTML = `<span class="prompt">${prompt.textContent}</span> ${command}`;
                    if (command) {
                        this.handleCommand(command);
                    } else {
                        this.createInputLine();
                    }
                }
            });

            inputLine.appendChild(prompt);
            inputLine.appendChild(input);
            output.appendChild(inputLine);

            input.focus();
        },

        resetPath() {
            this.linkStart = '123'
            fetch("/", {
                method: "POST",
                body: new URLSearchParams({
                    url: this.url,
                    key: this.key,
                    linkStart: this.linkStart
                }),
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                }
            })
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(text => {
                            this.errorsMessage = text; // 显示具体的错误信息
                            throw new Error('Network response was not ok');
                        });
                    }
                    return response.text();
                })
                .then(data => {
                    const part = data.split('------')
                    this.path = part[0]
                })
                .catch(error => {
                    // 这里处理网络错误或者由throw new Error抛出的错误
                    console.error('There has been a problem with your fetch operation:', error);
                });
        }
    },
    mounted() {
        document.addEventListener('contextmenu', (event) => {
            if (this.contextMenuVisible) {
                event.preventDefault();
            }
        });
    },
    beforeDestroy() {
        document.removeEventListener('contextmenu', (event) => {
            if (this.contextMenuVisible) {
                event.preventDefault();
            }
        });
    }
});
