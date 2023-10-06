package operator

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/logger"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"
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
		return nil, errors.New("headerStr value is empty")
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

func AuthByParameter(requiredAccessLevel dtos.AccessLevel) gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuthorized, err := HasSufficientAccessByParameter(c, requiredAccessLevel)
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
	token := utils.GetRequiredHeader(c, "authorization")
	if token == "" {
		return nil, fmt.Errorf("missing header 'authorization'")
	}

	authorization, err := parseAuthHeader(token)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}

	claims, err := services.ValidationToken(authorization.Value)
	if err != nil {
		return nil, err
	}
	userId := claims.UserID

	// updateLocalUserStore()

	getGinContextUser := services.GetGinContextUser(c)
	if getGinContextUser != nil && getGinContextUser.Id == userId {
		return getGinContextUser, nil
	}

	user, err := services.GetUser(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CheckUserAuthorizationByParameter(c *gin.Context) (*dtos.PunqUser, error) {
	token, tokenOk := c.GetQuery("token")
	if !tokenOk || token == "" {
		return nil, fmt.Errorf("missing header 'authorization'")
	}

	claims, err := services.ValidationToken(token)
	if err != nil {
		return nil, err
	}
	userId := claims.UserID

	getGinContextUser := services.GetGinContextUser(c)
	if getGinContextUser != nil && getGinContextUser.Id == userId {
		return getGinContextUser, nil
	}

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
	return false, fmt.Errorf(errStr)
}

func HasSufficientAccessByParameter(c *gin.Context, requiredAccessLevel dtos.AccessLevel) (bool, error) {
	user, err := CheckUserAuthorizationByParameter(c)
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
	return false, fmt.Errorf(errStr)
}
