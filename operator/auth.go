package operator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/services"
	"net/http"
)

// var users []dtos.PunqUser = []dtos.PunqUser{}
// var nextUpdate time.Time = time.Now().Add(-1 * time.Minute) // trigger first update instant

func Auth(requiredAccessLevel dtos.AccessLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuthorized, err := HasSufficientAccess(c, requiredAccessLevel)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if isAuthorized {
			c.Next()
		}
	}
}

// func updateLocalUserStore() {
// 	if time.Now().After(nextUpdate) {
// 		users = services.ListUsers()
// 		nextUpdate = time.Now().Add(1 * time.Minute) // wait a minute for next update
// 	}
// }

func CheckUserAuthorization(c *gin.Context) (*dtos.PunqUser, error) {
	authorization := GetRequiredHeader(c, "authorization")
	if authorization == "" {
		return nil, fmt.Errorf("MalformedRequest")
	}
	claims := services.ValidationToken(authorization)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid token"})
		return nil, fmt.Errorf("MalformedRequest")
	}
	userId := claims.UserID

	// updateLocalUserStore()

	user := services.GetUser(userId)
	if user != nil {
		return user, nil
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"err": "Authorization failed.",
	})

	return nil, fmt.Errorf("UserNotFound")
}

func HasSufficientAccess(c *gin.Context, requiredAccessLevel dtos.AccessLevel) (bool, error) {
	user, err := CheckUserAuthorization(c)
	if err != nil {
		return false, err
	}
	if user != nil {
		if user.AccessLevel >= requiredAccessLevel {
			c.Set("user", *user)
			return true, err
		}
	}
	errStr := fmt.Sprintf("AccessLevel is insufficient (Current:%d - Required:%d).", user.AccessLevel, requiredAccessLevel)
	c.JSON(http.StatusNotFound, gin.H{
		"err": errStr,
	})
	return false, fmt.Errorf(errStr)
}
