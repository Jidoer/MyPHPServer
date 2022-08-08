package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"strconv"
	"strings"

	"MyPHPServer/tool/file"

	"github.com/kataras/iris/v12"
)

type Server struct {
	app    *iris.Application
	Config ServerConfig
}
type ServerConfig struct {
	Root    string `json:"root"`
	CgiBin  string `json:"cgibin"`
	UseSSL  bool   `json:"usessl"`
	CrtPath string `json:"crtpath"`
	KeyPath string `json:"keypath"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}

func NewServer() *Server {
	var s Server
	if file.Exists("config.json") {
		file, err := os.Open("config.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		conf, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}
		er := json.Unmarshal(conf, &s.Config)
		if er != nil {
			panic(er)
		}
	}
	s.app = iris.New()
	return &s
}

func (s *Server) Start() {
	s.GetStaticHandler()
	s.GetPHPHandler()

	if s.Config.UseSSL {
		s.app.Run(iris.TLS(s.Config.Host+":"+strconv.Itoa(s.Config.Port), s.Config.CrtPath, s.Config.KeyPath))
	} else {
		s.app.Run(iris.Addr(s.Config.Host + ":" + strconv.Itoa(s.Config.Port)))
	}
}

func (s *Server) loadphp(ctx iris.Context /*path string*/, Root string, Path string) bool {
	PHPcgi(ctx.ResponseWriter(), ctx.Request() /*cgibin*/, s.Config.CgiBin /*"/usr/lib/cgi-bin/php7.4"*/, Root, Path)
	return true
}

func PHPcgi(w http.ResponseWriter, r *http.Request, cgiBin string /*scriptFile string*/, Root string, Path string) {
	//scriptFile:PHPFilePath
	scriptFile := Root + Path
	handler := new(cgi.Handler)
	handler.Path = cgiBin
	handler.Env = append(handler.Env, "REDIRECT_STATUS=CGI")
	handler.Env = append(handler.Env, "SCRIPT_FILENAME="+scriptFile)
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

func (s *Server) GetPHPHandler() {
	s.app.Get("/{path:path}", func(ctx iris.Context) {
		path := ctx.Params().Get("path")
		log.Println("Load: ", s.Config.Root+path)
		if path == "" {
			if file.Exists(s.Config.Root + "index.htm") {
				path = "index.htm"
			}
			if file.Exists(s.Config.Root + "index.html") {
				path = "index.html"
			}
			if file.Exists(s.Config.Root + "index.php") {
				path = "index.php"
			}
			//ctx.StatusCode(404)
			//return
		}
		if !file.Exists(s.Config.Root + path) {
			log.Println("Not Found: ", s.Config.Root+path)
			ctx.StatusCode(http.StatusNotFound)
			return
		}
		name, filetype := s.GetFileType(path)
		log.Println("FileType: ", name, filetype)
		switch GetSwitch(filetype) {
		case "run":
			s.loadphp(ctx, s.Config.Root, path)
		case "download":
			ctx.SendFile(s.Config.Root+path, name)
		case "view":
			ctx.View(s.Config.Root + path)
		case "null":
			ctx.StatusCode(http.StatusNotFound)
		case "dir":
			{
				if file.Exists(s.Config.Root + path + "/index.php") {
					s.loadphp(ctx, s.Config.Root, path+"/index.php")
					return
				}
				if file.Exists(s.Config.Root + path + "/index.html") {
					ctx.Write(file.Reader2Byte(s.Config.Root + path + "/index.html"))
					return
				}
				ctx.StatusCode(http.StatusForbidden)
				page := ListFile(s.Config.Root,path)
				ctx.HTML(page)
				//ctx.StatusCode(http.StatusNotFound)
				//403
				//ctx.StatusCode(http.StatusForbidden)
				break
			}
		case "source":
			ctx.Write(file.Reader2Byte(s.Config.Root + path))
		default:
			ctx.StatusCode(http.StatusNotFound)
		}
	})

}

//load js css
func (s *Server) GetStaticHandler() {
	s.app.HandleDir("/static", "./static")
}

func (s *Server) GetConfig() ServerConfig {
	return s.Config
}

//return filename,type
func (s *Server) GetFileType(path string) (string, string) { //full path
	if !file.Exists(s.Config.Root + path) {
		return "", "notfound"
	}
	if strings.Contains(path, "/") {
		path = path[strings.LastIndex(path, "/")+1:] //file name
	}
	if !strings.Contains(path, ".") {
		if file.IsDir(s.Config.Root+ path) {
			return path, "dir"
		}
		return path, "file"
	}
	fileType := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
	return path, fileType
}

func GetSwitch(filetype string) string {
	switch filetype {
	case "php":
		return "run"
	case "html", "htm", "js", "css", "json", "txt", "xml", "c":
		return "source"
	case "jpg", "jpeg", "png", "gif", "bmp", "ico", "zip", "rar", "7z", "gz", "bz2", "tar", "tgz", "tbz2", "mp4", "mp3", "avi", "flv", "wmv", "mkv", "mov", "mpeg", "mpg", "m4v", "3gp", "3g2", "wav", "wma", "flac", "aac", "m4a":
		return "download"
	case "dir":
		return "dir"
	case "tpl":
		return "view"
	default:
		return "null"
	}
}
