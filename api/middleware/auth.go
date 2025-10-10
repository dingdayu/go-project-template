package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/dingdayu/singbox-adapter/model/entity"
	"github.com/dingdayu/singbox-adapter/model/entity/contextkey"
	"github.com/dingdayu/singbox-adapter/pkg/jwt"
	"github.com/dingdayu/singbox-adapter/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Authorization() gin.HandlerFunc {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		panic("JWT secret is not configured")
	}

	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ErrAuthForbidden)
			return
		}
		token := strings.Split(authorization, " ")
		if token[0] != "Bearer" || len(token) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ErrAuthForbidden)
			return
		}
		// 检查 jwt token 是否有效
		user, err := jwt.ParseJwt(token[1], secret)
		if err != nil {
			logger.Logger().WarnContext(c.Request.Context(), "JWT parsing failed",
				"error", err,
				"ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent())
			c.AbortWithStatusJSON(http.StatusUnauthorized, entity.ErrAuthTokenInvalid)
			return
		}

		// 创建增强的 context
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, contextkey.Email, user.Email)
		ctx = context.WithValue(ctx, contextkey.RealName, user.RealName)
		ctx = context.WithValue(ctx, contextkey.Role, user.Roles)
		ctx = context.WithValue(ctx, contextkey.UserName, user.Username)
		ctx = context.WithValue(ctx, contextkey.IP, c.ClientIP())

		// 更新请求上下文
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
