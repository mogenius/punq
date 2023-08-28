package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mogenius/punq/dtos"
	"github.com/mogenius/punq/kubernetes"
	"github.com/mogenius/punq/services"
)

func InitUserRoutes(router *gin.Engine) {
	router.GET("/user/all", Auth(dtos.ADMIN), userList)
	router.GET("/user", Auth(dtos.ADMIN), userGet)
	router.DELETE("/user", Auth(dtos.ADMIN), userDelete)
	router.PATCH("/user", Auth(dtos.ADMIN), userUpdate)
	router.POST("/user", Auth(dtos.ADMIN), userAdd)
}

func userList(c *gin.Context) {
	users := services.ListUsers()

	RespondForWorkloadResult(c, kubernetes.WorkloadResult(users, nil))
}

func userDelete(c *gin.Context) {
	userId := c.Query("userId")

	RespondForWorkloadResult(c, services.DeleteUser(userId))
}

func userGet(c *gin.Context) {
	userId := c.Query("userId")

	user := services.GetUser(userId)
	RespondForWorkloadResult(c, kubernetes.WorkloadResult(user, nil))

}

func userUpdate(c *gin.Context) {
	var data dtos.PunqUser
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, services.UpdateUser(data))
}

func userAdd(c *gin.Context) {
	var data dtos.PunqUser
	err := c.MustBindWith(&data, binding.JSON)
	if err != nil {
		MalformedMessage(c, err.Error())
		return
	}
	RespondForWorkloadResult(c, services.AddUser(data))
}
