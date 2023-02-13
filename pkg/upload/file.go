package upload

import (
	"bytes"
	"first/pkg/constants"
	"first/pkg/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

type FileType int

const (
	TypeImage FileType = iota + 1
	TypeExcel
	TypeTxt
)

//GetFileName 获取文件名称，先是通过获取文件后缀并筛出原始文件名进行 MD5 加密，最后返回经过加密处理后的文件名。
func GetFileName(name string) string {
	ext := GetFileExt(name)
	// end with ext ?
	fileName := strings.TrimSuffix(name, ext)
	file, err := os.Open(name)
	if err != nil {
		return ""
	}
	fileName = util.EncodeMD5(file)

	return fileName + ext
}

//GetFileExt 获取文件后缀，主要是通过调用 path.Ext 方法进行循环查找”.“符号，
//最后通过切片索引返回对应的文化后缀名称。
func GetFileExt(name string) string {
	return path.Ext(name)
}

//GetSavePath 获取文件保存地址，这里直接返回配置中的文件保存目录即可，也便于后续的调整。
func GetSavePath() string {
	return constants.UploadSavePath
}

func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsNotExist(err)
}

//CheckContainExt 检查文件扩展名是否合法
func CheckContainExt(t FileType, name string) bool {
	return true
}

// CheckMaxSize 检查文件大小是否超出最大大小限制
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)

	switch t {
	case TypeImage:
		if size >= constants.UploadImageMaxSize*1024*1024 {
			return true
		}
	}
	return false
}

// CheckPermission 检查是否拥有权限
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

// CreateSavePath 递归创建传递 perm 创建保护目录
func CreateSavePath(dst string, perm os.FileMode) error {
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

//SaveFile 保存所上传的文件，该方法主要是通过调用 os.Create 方法创建目标地址的文件，再通过 file.Open 方法打开源地址的文件，结合 io.Copy 方法实现两者之间的文件内容拷贝。
func SaveFile(file []byte, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	reader := bytes.NewReader(file)

	_, err = io.Copy(out, reader)
	return err
}
