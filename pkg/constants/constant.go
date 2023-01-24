package constants

const (
	// user service
	UserServiceName = "user"
	UserTableName   = "user_info"
	FollowTableName = "follow_list"

	// jwt
	SecretKey      = "secret key"
	IdentityKey    = "user_id"
	PublicKeyFile  = "E:\\DY2023\\pkg\\constants\\pub.key"
	PrivateKeyFile = ".\\pkg\\constants\\private.key"

	MySQLDefaultDSN   = "dy:123456@tcp(localhost:3307)/dy?charset=utf8mb4&parseTime=True&loc=Local"
	EtcdAddress       = "127.0.0.1:2379"
	ApiServerAddress  = "127.0.0.1:8888"
	UserServerAddress = "127.0.0.1:8889"
)
