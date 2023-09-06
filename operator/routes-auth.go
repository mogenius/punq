package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/services"
	"net/http"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func InitAuthRoutes(router *gin.Engine) {
	router.POST("/auth/login", login)
	router.GET("/auth/authenticate", Auth(dtos.READER), authenticate)
}

// @Tags Auth
// @Produce json
// @Success 200 {object} dtos.PunqToken
// @Router /auth/login [post]
// @Param body body LoginInput true "LoginInput"
func login(c *gin.Context) {
	input := LoginInput{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	user := services.GetUserByEmail(input.Email)

	_, err = user.PasswordCheck(input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "username or password is incorrect."})
		return
	}

	token, _ := services.GenerateToken(user)

	c.JSON(http.StatusOK, token)
}

// @Tags Auth
// @Produce json
// @Success 200 {object} dtos.PunqToken
// @Router /auth/authenticate [get]
// @Security Bearer
func authenticate(c *gin.Context) {
	if temp, exists := c.Get("user"); exists {
		user, ok := temp.(dtos.PunqUser)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}
		c.JSON(http.StatusOK, user)
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{})
}
