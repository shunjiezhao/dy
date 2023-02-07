package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/klauspost/compress/zstd"
)

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
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
