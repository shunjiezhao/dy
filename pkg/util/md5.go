package util

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/klauspost/compress/zstd"
	"io"
)

func EncodeMD5(r io.Reader) string {
	h := sha256.New()
	readAll, _ := io.ReadAll(r)
	h.Write(readAll)
	// 注意需要先将byte转换为16进制表示
	return base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(h.Sum(nil))))
}

func Compress(src []byte) ([]byte, error) {
	enc, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}

	return enc.EncodeAll(src, make([]byte, 0, len(src))), nil
}

func DeCompress(data []byte) ([]byte, error) {
	var decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
	return decoder.DecodeAll(data, nil)
}
