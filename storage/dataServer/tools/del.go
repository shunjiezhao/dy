package tools

import (
	"ch2/lib/es"
	"ch2/lib/utils"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	dataPath string
	addr     string
)

func SetAddr(addr string) {
	addr = addr
	_, dataPath = utils.SetAddr(addr)
	strings.TrimSuffix(dataPath, "/")
}

// 删除对象
func Work() {
	files, _ := filepath.Glob(dataPath + "/*")
	for _, file := range files {
		hash := strings.Split(filepath.Base(file), ".")[0]
		hashInMetadata, e := es.HasHash(hash)
		if e != nil {
			log.Println(e)
			return
		}
		if !hashInMetadata {
			del(hash)
		}
	}
}

func del(hash string) {
	log.Println("delete", hash)
	url := "http://" + addr + "/objects/" + hash
	request, _ := http.NewRequest("DELETE", url, nil)
	client := http.Client{}
	client.Do(request)
}
