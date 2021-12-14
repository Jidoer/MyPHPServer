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
	app.Get("/{path:path}", func(ctx iris.Context) {
		//ctx.Request()
		Root := "./"//"/home/wwwroot/www/"
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
					if lastname == ".php" {
						loadphp(ctx, Root+path)
					} else if lastname == ".html" || lastname == ".css" || lastname == ".js" {
						ctx.HTML(file.Reader(Root + path))
					} else if lastname == ".zip" || lastname == ".mp3" {
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
				loadphp(ctx, Root+path+"/index.php")
			} else if file.Exists(Root + path + "/index.html") {
				ctx.HTML(file.Reader(Root + path + "/index.html"))
			} else if file.Exists(Root + path + "/index.htm") {
				ctx.HTML(file.Reader(Root + path + "/index.htm"))
			} else if file.Exists(Root + path + "/index.c") {
				ctx.HTML(file.Reader(Root + path + "/index.c"))
			} else {
				//403!
				ctx.Header("HTTP/1.1", "403")
				ctx.HTML("403 ERROR!")
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

func loadphp(ctx iris.Context, path string) bool {
	PHPcgi(ctx.ResponseWriter(), ctx.Request() /*cgibin*/, "/usr/lib/cgi-bin/php7.4", path)
	return true
}

func PHPcgi(w http.ResponseWriter, r *http.Request, cgiBin string, scriptFile string) {
	handler := new(cgi.Handler)
	handler.Path = cgiBin
	handler.Env = append(handler.Env, "REDIRECT_STATUS=CGI")
	handler.Env = append(handler.Env, "SCRIPT_FILENAME="+scriptFile)
	handler.ServeHTTP(w, r)
}
