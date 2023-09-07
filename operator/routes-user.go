package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/services"
	"github.com/mogenius/punq/utils"
)

func InitUserRoutes(router *gin.Engine) {

	userRoutes := router.Group("/user")
	{
		userRoutes.GET("/all", Auth(dtos.ADMIN), userList)
		userRoutes.GET("/:id", Auth(dtos.ADMIN), userGet)
		userRoutes.DELETE("/:id", Auth(dtos.ADMIN), userDelete)
		userRoutes.PATCH("/", Auth(dtos.ADMIN), userUpdate)
		userRoutes.POST("/", Auth(dtos.ADMIN), userAdd)
	}

}

// @Tags User
// @Produce json
// @Success 200 {array} dtos.PunqUser
// @Router /user/all [get]
// @Security Bearer
func userList(c *gin.Context) {
	users := services.ListUsers()

	utils.RespondForHttpResult(c, kubernetes.WorkloadResult(users, nil))
}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user/{id} [delete]
// @Param id path string false  "ID of the user"
// @Security Bearer
func userDelete(c *gin.Context) {
	userId := c.Param("id")

	utils.RespondForHttpResult(c, services.DeleteUser(userId))
}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user/{id} [get]
// @Param id path string false  "ID of the user"
// @Security Bearer
func userGet(c *gin.Context) {
	userId := c.Param("id")

	user := services.GetUser(userId)
	utils.RespondForHttpResult(c, kubernetes.WorkloadResult(user, nil))
}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user [patch]
// @Param body body dtos.PunqUser false "PunqUser"
// @Security Bearer
func userUpdate(c *gin.Context) {
	var data dtos.PunqUser
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.RespondForHttpResult(c, services.UpdateUser(data))
}

// @Tags User
// @Produce json
// @Success 200 {object} dtos.PunqUser
// @Router /user [post]
// @Param body body dtos.PunqUser false "PunqUser"
// @Security Bearer
func userAdd(c *gin.Context) {
	var data dtos.PunqUser
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		utils.MalformedMessage(c, err.Error())
		return
	}
	utils.RespondForHttpResult(c, services.AddUser(data))
}
