package middleware

import (
	"fmt"
	"os"
	"time"

	initializers "github.com/naqeeb8a/Coffee-shop/initializers"
	models "github.com/naqeeb8a/Coffee-shop/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Missing authorization token",
		})
	}
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["expiresIn"].(float64) {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "token expired",
			})
		}
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Unable to find user with this token",
			})
			return
		}
		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Invalid token",
		})
	}
}
