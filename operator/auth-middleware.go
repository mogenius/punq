package operator

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"
	"regexp"
)

// var users []dtos.PunqUser = []dtos.PunqUser{}
// var nextUpdate time.Time = time.Now().Add(-1 * time.Minute) // trigger first update instant

type AuthHeader struct {
	Scheme string
	Value  string
}

func parseAuthHeader(headerStr string) (*AuthHeader, error) {
	re := regexp.MustCompile(`(\S+)\s+(\S+)`)

	if headerStr == "" {
		return nil, errors.New("headerStr value ist empty")
	}

	matches := re.FindStringSubmatch(headerStr)
	if matches == nil {
		return nil, errors.New("invalid authorization token")
	}

	return &AuthHeader{
		Scheme: matches[1],
		Value:  matches[2],
	}, nil
}

func Auth(requiredAccessLevel dtos.AccessLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuthorized, err := HasSufficientAccess(c, requiredAccessLevel)
		if err != nil {
			utils.Unauthorized(c, err.Error())
			c.Abort()
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
	authorization, err := parseAuthHeader(utils.GetRequiredHeader(c, "authorization"))

	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	claims := services.ValidationToken(authorization.Value)
	if claims == nil {
		msg := "claims failed"
		logger.Log.Error(msg)
		return nil, errors.New(msg)
	}
	userId := claims.UserID

	// updateLocalUserStore()

	user, err := services.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
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
	// c.JSON(http.StatusNotFound, gin.H{
	// 	"err": errStr,
	// })
	return false, fmt.Errorf(errStr)
}
