package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/shoppinglist-api/internal/domain/auth"
	domainErr "github.com/tclutin/shoppinglist-api/internal/domain/errors"
	"github.com/tclutin/shoppinglist-api/pkg/response"
	"net/http"
	"strings"
)

func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				response.NewAPIError[string](http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
			return
		}

		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				response.NewAPIError[string](http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
			return
		}

		userID, err := authService.VerifyCredentials(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				response.NewAPIError[string](http.StatusUnauthorized, domainErr.ErrMissingCredentials.Error(), nil))
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
