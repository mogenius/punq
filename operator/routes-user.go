package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"
	"net/http"
)

func InitUserRoutes(router *gin.Engine) {

	userRoutes := router.Group("/user", Auth(dtos.ADMIN))
	{
		userRoutes.GET("/all", userList)
		userRoutes.GET("/", currentUserGet)
		userRoutes.GET("/:id", validateParam("id"), userGet)
		userRoutes.DELETE("/:id", validateParam("id"), userDelete)
		userRoutes.PATCH("/", userUpdate)
		userRoutes.POST("/", userAdd)
	}

}

// @Tags User
// @Produce json
// @Success 200 {array} dtos.PunqUser
// @Router /user/all [get]
// @Security Bearer
func userList(c *gin.Context) {
	users := services.ListUsers()
	utils.HttpRespondForWorkloadResult(c, kubernetes.WorkloadResult(users, nil))
}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user/{id} [delete]
// @Param id path string false  "ID of the user"
// @Security Bearer
func userDelete(c *gin.Context) {
	userId := c.Param("id")
	c.JSON(http.StatusOK, services.DeleteUser(userId))
}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user [get]
// @Security Bearer
func currentUserGet(c *gin.Context) {
	user := services.GetGinContextUser(c)
	if user != nil {
		c.JSON(http.StatusOK, user)
		return
	}
	utils.Unauthorized(c, "Unauthorized")
	return

}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user/{id} [get]
// @Param id path string false  "ID of the user"
// @Security Bearer
func userGet(c *gin.Context) {
	userId := c.Param("id")

	user, err := services.GetUser(userId)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user [patch]
// @Param body body dtos.PunqUserUpdateInput false "PunqUserUpdateInput"
// @Security Bearer
func userUpdate(c *gin.Context) {
	var data dtos.PunqUserUpdateInput
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	user, err := services.UpdateUser(data)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user [post]
// @Param body body dtos.PunqUserCreateInput false "PunqUserCreateInput"
// @Security Bearer
func userAdd(c *gin.Context) {
	var data dtos.PunqUserCreateInput
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	user, err := services.AddUser(data)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}
