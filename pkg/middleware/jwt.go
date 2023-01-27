package middleware

import (
	"context"
	"first/pkg/constants"
	"first/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"log"
	"net/http"
	"time"
)

func JwtMiddle() *jwt.HertzJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Realm:            "test zone",
		Key:              []byte(constants.SecretKey),
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		IdentityKey:      constants.IdentityKey,
		SigningAlgorithm: "RS256",
		PubKeyBytes:      constants.PublicKeyFile,
		PrivKeyBytes:     constants.PrivateKeyFile,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					constants.IdentityKey: v,
				}
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, map[string]interface{}{
				"code":    http.StatusOK,
				"user_id": c.MustGet(constants.IdentityKey).(int64), // Authenticator 会先处理没有 uuid 的情况
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
			})
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    errno.AuthorizationFailedErrCode,
				"message": message,
			})
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			value, exists := c.Get(constants.IdentityKey)
			if !exists {
				return "", jwt.ErrMissingLoginValues
			}
			return value, nil
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	if err = authMiddleware.MiddlewareInit(); err != nil {
		log.Fatal("JWT Init Error:" + err.Error())
	}
	return authMiddleware
}
