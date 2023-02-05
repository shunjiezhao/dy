package constants

func init() {

}

const (
	// user service
	UserServiceName = "user"
	UserTableName   = "user_info"
	FollowTableName = "follow_list"

	// video service
	VideoServiceName        = "video"
	VideoTableName          = "video_info"
	FavouriteVideoTableName = "user_favourite_video"
	CommentTableName        = "comment_info"

	//  jwt
	SecretKey   = "secret key"
	IdentityKey = "user_id"

	// MYSQL DSN
	MySQLDefaultDSN = "dy:123456@tcp(localhost:3307)/dy?charset=utf8mb4&parseTime=True&loc=Local"

	EtcdAddress       = "127.0.0.1:2379"
	ApiServerAddress  = ":8888"
	UserServerAddress = "127.0.0.1:8889"

	UploadSavePath     = "./static"
	UploadImageMaxSize = 10
	UploadServerUrl    = ApiServerAddress
)

var (
	PrivateKeyFile = []byte(
		"-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQD0glQart3rf0Ap\nJvFrZfK" +
			"+eq5VgAkcZgDBl5AMQx2oLmUye4SAt7cVsGrWvxiUGs9petvM4q0rMflL\nH4HIcYu4T2dIX/ReYeSxlx7XBYYZoqkM6UjzGdtQ8MqJLmyjpY9joRkG2z1kBySR\n5hP6N/2RRZh8e2qs1hxu6gOo+VEUfB7pG8yRJmlQOg/PbHh0fS7ibo5lv77r/kgh\no8WrfX0Rrm0UtEmmA16j6m61kIRCeofXbWEr1TzGlGl8h39zxN6+Z4qqRYmHi8ZD\nWhBUYhyVuG1BaSwaaDXfbyI+/ih/KLEGjpiLyOmdnNIeawSlLrW3iYOh12MsZHKp\no9aSQ/R9AgMBAAECggEBANe7pjlk0KlYPWQR2DDKYsNtuyP1NBS6azBkadRn42Lg\njKleEir/7apVXe7b7PPANAD9RbIgzmmuTibaRch1ZrHYXWieQR6FgSKwE6XkWc2E\nl2Os8ZCM39Uqn4kqTPCWw01EdrB2AFSheMLCHh5ICJKEtWYf/p9AyxWRpGkSkVdf\n8tBn1znvoUVuQNEtzHpu3rpbGl4xyNtoGwtJEqpwdeS/r4whDrrAkQ1/gnE55iWH\nF3zsa9HzeniB3fT/lpEbJilkAj6llSQZGcZiVXIh2iFVmKqrUT1bOqfu28wuzT4y\nhu0li841CvzqAzkyNqVRkgU/5Iym8LxFz9+1X25HXIECgYEA/0YK6VEbKhka+KAX\nt0CuI5LquDM9FqgUtBlWdw5JHF3SG+/Qb74SB6IoqaFsPaHwjZXstGHDcZOpXo2f\nh1yTS6xSM7L4c+wKCasKpH2pFkbY+tjQvqsOegr/s0vpEyEcfP10BnIoQRQDxL/h\nKRMuWSXvUld+xPoKOZNB2su6Y0sCgYEA9TRxwc7Rd12WWaGYV57uzw85OdJnSJQq\nOP5d5EY0YFNJdKa6pcofwXs8z8WHxG8JfrNjl/iKcVOEzFre2dc4slwh5Z4DeCr0\n7MP2VBbL9E0+rTNzeCngW1ZyiJ0uwUx7QDiCKeJfPIElcqZpoNBxK5dU3zULay5J\nKy3gf/BM4lcCgYBHkpTmm/X41LcqNIDRwZHRqZSj9sHPA2tin6QNl3TKPkf1y0Ru\nwCT//OhXv0nA8hGnMP0AClUpGBSpzR2Ib11hHzyhADIHFowt78X5Hr5034Jgur+0\nZfOWJlVRKRx9X5BEPy/zyrgcnwb7eC0iPh2Fo0w5kwyZH94UDISvWuW0hwKBgQCc\n0AVQJKvg4oEcoTOEFagz01CNoflbeSXnfQUez6b/U0ROzbHgBPt6CQ5C8dh5z2kL\nFj5DGjevcfIjnpmWRwWDS1iCOCOP3ij0of4OmOWmPyAuNBFMb7uDri1hIOSdygOo\ndnsHvjWZxB3mzHYQ2j0F26nzdUDwMpGog5ZnO45v0QKBgBvzrdbundsv3+TjXqz6\na06hlMhnLZhqSSy0J3mdRxSCZX1UYP0BZQCYGU2Icwr+5O6Zs17AR8ZXffywwDrJ\n/1LvXqPI5LA/FL6euyVrpwtXn2qNiV44VQ29mGUst1xhy6VAR95qtiMtYbePq5Sh\nook1LIIDvtrM0aPU9JAC8BXF\n-----END PRIVATE KEY-----\n")
	PublicKeyFile = []byte("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9IJUGq7d639AKSbxa2Xy\nvnquVYAJHGYAwZeQDEMdqC5lMnuEgLe3FbBq1r8YlBrPaXrbzOKtKzH5Sx+ByHGL\nuE9nSF/0XmHksZce1wWGGaKpDOlI8xnbUPDKiS5so6WPY6EZBts9ZAckkeYT+jf9\nkUWYfHtqrNYcbuoDqPlRFHwe6RvMkSZpUDoPz2x4dH0u4m6OZb++6/5IIaPFq319\nEa5tFLRJpgNeo+putZCEQnqH121hK9U8xpRpfId/c8TevmeKqkWJh4vGQ1oQVGIc\nlbhtQWksGmg1328iPv4ofyixBo6Yi8jpnZzSHmsEpS61t4mDoddjLGRyqaPWkkP0\nfQIDAQAB\n-----END PUBLIC KEY-----")
)
