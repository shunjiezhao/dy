package objects

import (
	"ch2/dataServer/locate"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// 第一次防止，防止竞争 -> 也就是利用es和http不是一个事务内，中间插空进行查询。
// 二次检验：
func del(w http.ResponseWriter, r *http.Request) {
	hash := strings.Split(r.URL.EscapedPath(), "/")[2]
	files, _ := filepath.Glob(dataPath + hash + ".*")
	if len(files) != 1 {
		return
	}
	locate.Del(hash)
	os.Rename(files[0], filepath.Dir(dataPath)+"/garbage/"+filepath.Base(files[0])) // 移入垃圾堆
}
