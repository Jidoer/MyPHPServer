package main

import (
	"fmt"
	"github.com/kangc666/file/file"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cgi"
	"strings"
)

func main() {
	app := iris.New()
	//PHP Server
	app.Get("/{path:path}", func(ctx iris.Context) {
		//ctx.Request()
		Root := "/home/wwwroot/www/"
		path := ctx.Params().Get("path")

		log.Print(path)

		if strings.Contains(path, ".") {
			if strings.Contains(path, "./") {
				// ../../../home/user/
				ctx.HTML("Baned!")
			} else {
				lastname := string(path[strings.LastIndex(path, "."):])
				if strings.Contains(lastname, "/") {
					// c/com.co.kang/cc
					//ctx.HTML("Baned!")
					//return
					lastname = "/"
				}
				//good file!
				if file.Exists(Root + path) {
					//!
					if lastname == ".php" {
						loadphp(ctx, Root, path)
					} else if lastname == ".html" || lastname == ".css" || lastname == ".js" {
						ctx.Write([]byte(file.Reader(Root + path)))
					} else if lastname == ".zip" || lastname == ".mp3" || lastname == ".png" || lastname == ".gif" || lastname ==".mp4" {
						log.Print("File:" + Root + path)
						// /c.mp3 or cc/c.mp3
						if strings.Contains(path, "/") {
							ctx.SendFile(Root+path, string(path[strings.LastIndex(path, "/"):])) //file full name
						} else {
							ctx.SendFile(Root+path, path) //file full name
						}
					} else {
						ctx.JSON(iris.Map{"error": false, "error_message": "no func"})
					}
				} else {
					//404
					ctx.HTML("404!")
				}

			}
		} else {
			//No file or cc/bb/index
			if file.Exists(Root + path + "/index.php") {
				loadphp(ctx, Root, path+"/index.php")
			} else if file.Exists(Root + path + "/index.html") {
				ctx.Write([]byte(file.Reader(Root + path + "/index.html")))
			} else if file.Exists(Root + path + "/index.htm") {
				ctx.Write([]byte(file.Reader(Root + path + "/index.htm")))
			} else if file.Exists(Root + path + "/index.c") {
				ctx.Write([]byte(file.Reader(Root + path + "/index.c")))
			} else {
				//403!
				ctx.HTML(ListFile(Root,path))
			}
		}

	})

	//Add page Proxy & File Proxy From http/https
	app.Get("/proxy/{apiCall:path}", func(ctx iris.Context) {
		endpoint := "http://127.0.0.1:8896/"
		apiCall := ctx.Params().Get("apiCall")
		if len(apiCall) > 0 {
			endpoint = fmt.Sprintf("%s%s", endpoint, apiCall)
		}
		rawQuery := ctx.Request().URL.RawQuery
		if len(rawQuery) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, rawQuery)
		}
		var mc string
		if strings.Contains(apiCall, "./") {
			ctx.HTML("error")
		} else {
			if strings.Contains(apiCall, ".") {
				mc = string(apiCall[strings.LastIndex(apiCall, "."):])
				if strings.Contains(mc, "/") {
					//nb mc = mc[:len(mc)-1]
					//Such as: c/com.co.kang/cc
					mc = "/"
				}
			} else {
				mc = "/"
			}
			log.Print(mc)
			if mc == ".html" || mc == ".htm" || mc == ".php" || mc == "/" || mc == ".js" || mc == ".css" {
				//http&https proxy
				resp, err := http.Get(endpoint)
				if err != nil {
					ctx.JSON(iris.Map{"success": false, "error_message": err.Error()})
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode == 200 || resp.StatusCode == 201 {
					body, _ := ioutil.ReadAll(resp.Body)
					ctx.ContentType(context.ContentJSONHeaderValue)
					//ctx.Write(body)
					ctx.HTML(string(body))
				} else {
					ctx.JSON(iris.Map{"success": false, "error_message": "target not ok"})
				}
			} else if mc == ".zip" || mc == ".mp3" {
				if file.Exists("./res/" + apiCall) {
					ctx.SendFile("./res/"+apiCall, "index"+mc)
				} else {
					ctx.HTML("No File!")
				}
			} else {
				ctx.JSON(iris.Map{"success": false, "error_message": "target not ok"})
			}

		}
	})

	err := app.Run(
		iris.Addr(":9000"),
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)

	if err != nil {
		log.Println(err.Error())
	}
}

func ListFile(ROOT, Path string) string {
	dir_list, e := ioutil.ReadDir(ROOT + Path)
	if e != nil {
		fmt.Println("read dir error")
		return "error"
	}
	   /*
	   for i, v := range dir_list {
	       fmt.Println(i, "=", v.Name())
	   }*/
	tatol := "<hr>"
	for i, v := range dir_list {
		var myname string
		if file.IsDir(ROOT + Path + "/" + v.Name()) {
			myname = v.Name() + "/"
		} else {
			myname = v.Name()
		}

		tatol = tatol + "<a href='/" + Path + "/" + v.Name() + "'>" + myname /*v.Name()*/ + "</a><hr>"
		//log.Print(i,string(rune(i)))
		fmt.Println(i, v.Name()+string(rune(i)))
	}
	//log.Print(tatol)
	return "<h1>Index Of</h1>" + tatol + "<p align='center'>Powered By: FlyKO</p>"
}

func loadphp(ctx iris.Context /*path string*/, Root string, Path string) bool {
	PHPcgi(ctx.ResponseWriter(), ctx.Request() /*cgibin*/, "/usr/lib/cgi-bin/php7.4", Root, Path)
	return true
}

func PHPcgi(w http.ResponseWriter, r *http.Request, cgiBin string /*scriptFile string*/, Root string, Path string) {
	//scriptFile:PHPFilePath
	scriptFile := Root + Path
	handler := new(cgi.Handler)
	handler.Path = cgiBin
	handler.Env = append(handler.Env, "REDIRECT_STATUS=CGI")
	handler.Env = append(handler.Env, "SCRIPT_FILENAME="+scriptFile)
	//new
	handler.Env = append(handler.Env, "HTTP_HOST="+r.Host)
	handler.Env = append(handler.Env, "SERVER_NAME="+r.Host)

	var DOCUMENT_ROOT string
	var SCRIPT_NAME string
	/*GetDOCUMENT_ROOT*/
	if strings.Contains(Path, "/") {
		DOCUMENT_ROOT = string(scriptFile[:strings.LastIndex(scriptFile, "/")])
		SCRIPT_NAME = string(scriptFile[strings.LastIndex(scriptFile, "/"):])
	} else {
		// xxx.php index.php
		DOCUMENT_ROOT = Root
		SCRIPT_NAME = Path
	}

	handler.Env = append(handler.Env, "SCRIPT_NAME="+SCRIPT_NAME)
	handler.Env = append(handler.Env, "DOCUMENT_ROOT="+DOCUMENT_ROOT)

	print("DOCUMENT_ROOT=" + DOCUMENT_ROOT + "\n")
	//log.Print(r.Host)
	handler.ServeHTTP(w, r)
}
