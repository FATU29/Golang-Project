package _interface

import "github.com/gin-gonic/gin"

type IAuthenticationController interface {
	RegisterController(c *gin.Context)
	LoginController(c *gin.Context)
}
