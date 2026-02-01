package controller

import (
	"Authentication_Service/internal/dto/common"
	"Authentication_Service/internal/dto/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
}

func (a AuthenticationController) RegisterController(c *gin.Context) {
	var userRequestDTO request.RegisterReqDto
	err := c.ShouldBindJSON(&userRequestDTO)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.ApiResponse[any]{
			Code:    http.StatusBadRequest,
			Message: "Bad Request in DTO",
		})
	}
}

func (a AuthenticationController) LoginController(c *gin.Context) {
}
