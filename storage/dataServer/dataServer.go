package main

import (
	"ch2/dataServer/heartbeat"
	"ch2/dataServer/locate"
	"ch2/dataServer/objects"
	"ch2/dataServer/temp"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	path2 "path"
	"strings"
)

var (
	port = ""
	Addr = ""
)

func main() {
	//port = os.Getenv("RUN_PORTS")
	flag.StringVar(&port, "port", "9002", "port")
	flag.Parse()
	log.SetFlags(log.Llongfile)
	// api 层 需要了解下面的数据层 有哪些可用，那些不可用
	Addr := os.Getenv("IP") + ":" + port

	path := "./store" + fmt.Sprintf("/objects_%s/", strings.Split(Addr, ":")[1])
	objects.SetAddr(path)
	temp.SetAddr(path)

	objPath := path2.Join(path, "objects")
	go heartbeat.StartHeartbeat(Addr)
	go locate.StartLocate(Addr)
	go locate.Collections(objPath)
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(Addr, nil))
}
