package wallet

import "github.com/gin-gonic/gin"

//go:generate mockgen -source=middleware.go -destination=middleware_mock_test.go -package=wallet
type Middleware interface {
	MiddlewareValidateToken(c *gin.Context)
	MiddlewareSignatureValidation(c *gin.Context)
}
