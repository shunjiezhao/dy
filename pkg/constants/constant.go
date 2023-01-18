package constants

const (
	UserTableName   = "user"
	UserServiceName = "user"

	SecretKey = "secret key"

	MySQLDefaultDSN = "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
	EtcdAddress     = "127.0.0.1:2379"
)
