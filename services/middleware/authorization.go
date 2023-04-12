package middleware

import (
	"errors"
	"github.com/Zhoangp/Api_Gateway-Courses/pkg/common"
	"github.com/Zhoangp/Api_Gateway-Courses/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func extractToken(token string) (string, error) {

	parts := strings.Split(token, " ")
	if parts[0] != "Bearer" || len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
		return "", utils.ErrInvalidToken
	}
	return parts[1], nil
}
func (m *middleareManager) RequireVerifyToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := (c.Param("token"))
		payload, err := utils.ValidateJWT(token, m.cfg)
		if err != nil {
			panic(err)
		}
		c.Set("emailUser", payload.Email)
		c.Set("password", payload.Password)
		c.Set("verified", payload.Verified)
		c.Next()

	}
}
func (m *middleareManager) RequiredAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := c.Request.Header.Get("Authorization")

		token, err := extractToken(s)
		if err != nil {
			panic(err)
		}

		payload, err := utils.ValidateJWT(token, m.cfg)
		if err != nil {
			panic(err)
		}

		//User-Service, err := m.userRepo.FindDataWithCondition(map[string]any{"email": payload.Email})
		//if err != nil {
		//	panic(err)
		//}

		if payload.Verified {
			panic(common.NewCustomError(errors.New("This account has not been verified!"),"This account has not been verified!"))
		}
		c.Set("emailUser", payload.Email)

		c.Next()
	}
}