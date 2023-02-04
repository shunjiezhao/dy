package middleware

import (
	"first/pkg/constants"
	"first/pkg/errno"
	"first/service/api/handlers"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	middleware *jwt.GinJWTMiddleware
	token      gin.HandlerFunc
	once       sync.Once
)

func JwtMiddle() (*jwt.GinJWTMiddleware, gin.HandlerFunc) {
	once.Do(func() {
		var err error
		middleware, err = jwt.New(&jwt.GinJWTMiddleware{
			Realm:            "test zone",
			Key:              []byte(constants.SecretKey),
			Timeout:          time.Hour,
			MaxRefresh:       time.Hour,
			IdentityKey:      constants.IdentityKey,
			SigningAlgorithm: "RS256",
			Authorizator:     nil,
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				if v, ok := data.(int64); ok {
					return jwt.MapClaims{
						constants.IdentityKey: v,
					}
				}
				return jwt.MapClaims{}
			},
			LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
				if code == http.StatusOK {
					c.JSON(http.StatusOK, map[string]interface{}{
						"status_code": errno.SuccessCode,
						"user_id":     c.MustGet(constants.IdentityKey).(int64),
						// Authenticator 会先处理没有 uuid 的情况
						"token": message,
					})
					return
				}
				c.JSON(http.StatusOK, errno.AuthorizationFailedErr)
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.JSON(code, map[string]interface{}{
					"code":    errno.AuthorizationFailedErrCode,
					"message": message,
				})
			},
			Authenticator: func(c *gin.Context) (interface{}, error) {
				value, exists := c.Get(constants.IdentityKey)
				if !exists {
					return "", jwt.ErrMissingLoginValues
				}
				return value, nil
			},
			PubKeyBytes:   constants.PublicKeyFile,
			PrivKeyBytes:  constants.PrivateKeyFile,
			TokenLookup:   "header: Authorization, form: token, cookie: jwt, query: token",
			TokenHeadName: "Bearer",
			TimeFunc:      time.Now,
		})
		if err != nil {
			log.Fatal("JWT Error:" + err.Error())
		}
		// When you use jwt.New(), the function is already automatically called for checking,
		// which means you don't need to call it again.
		if err = middleware.MiddlewareInit(); err != nil {
			log.Fatal("JWT Init Error:" + err.Error())
		}
		token = gin.HandlerFunc(func(c *gin.Context) {
			fromJWT, err := middleware.GetClaimsFromJWT(c)
			if err != nil {
				handlers.SendResponse(c, errno.AuthorizationFailedErr)
				c.Abort()
				return
			}
			var curUserId int64
			tmp, ok := fromJWT[constants.IdentityKey].(float64)
			if ok {
				curUserId = int64(tmp)
			} else {
				curUserId = fromJWT[constants.IdentityKey].(int64)
			}
			c.Set(constants.IdentityKey, curUserId)
			c.Next()
		})
	})

	return middleware, token
}
