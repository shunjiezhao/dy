package constants

const (
	UserTableName   = "user_info"
	UserServiceName = "user"

	// jwt
	SecretKey      = "secret key"
	IdentityKey    = "id"
	PublicKeyFile  = "E:\\DY2023\\pkg\\constants\\pub.key"
	PrivateKeyFile = "E:\\DY2023\\pkg\\constants\\private.key"

	MySQLDefaultDSN   = "zsj:az123.@tcp(localhost:3306)/dy?charset=utf8mb4&parseTime=True&loc=Local"
	EtcdAddress       = "127.0.0.1:2379"
	ApiServerAddress  = "127.0.0.1:8888"
	UserServerAddress = "127.0.0.1:8889"
)
