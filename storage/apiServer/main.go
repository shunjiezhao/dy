package main

import (
	"ch2/apiServer/heartbeat"
	"ch2/apiServer/locate"
	"ch2/apiServer/objects"
	"ch2/apiServer/temp"
	"ch2/apiServer/versions"
	"log"
	"net/http"
	"os"
)

var (
	port = ""
	Addr = ""
)

func main() {
	port = os.Getenv("RUN_PORTS")
	port = "9001"
	log.SetFlags(log.Llongfile)
	// api 层 需要了解下面的数据层 有哪些可用，那些不可用
	Addr := os.Getenv("IP") + ":" + port
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)

	println(Addr)
	os.MkdirAll("./objects", 0666)
	log.Fatal(http.ListenAndServe(Addr, nil))
}
