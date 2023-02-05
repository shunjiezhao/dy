package util

import "github.com/GUAIK-ORG/go-snowflake/snowflake"

var ss *snowflake.Snowflake

func init() {
	var err error
	ss, err = snowflake.NewSnowflake(int64(0), int64(1))
	if err != nil {
		panic(err)
	}
}
func NextVal() int64 {
	return ss.NextVal()
}
