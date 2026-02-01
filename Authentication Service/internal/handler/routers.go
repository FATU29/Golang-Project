package handler

import (
	"Authentication_Service/internal/constant"
	"Authentication_Service/internal/controller"

	"github.com/gin-gonic/gin"
)

func DefineRouters(r *gin.Engine) {
	root := r.Group(constant.Root)

	NewAuthenticationHandler(controller.AuthenticationController{}).RouterList(root)

}
