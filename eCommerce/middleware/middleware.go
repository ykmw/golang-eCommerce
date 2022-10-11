package middleware

import (
	token "example.com/m/tokens"
	"net/http"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context){
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"No authorization header provided"})
			c.About()
			return
		}

		claims, err := token.ValicateToken(ClientToken)
		if err != ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}