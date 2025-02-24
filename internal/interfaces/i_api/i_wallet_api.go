package i_api

import "github.com/gin-gonic/gin"

type IWalletAPI interface {
	Create(*gin.Context)
	CreditBalance(c *gin.Context)
	DebitBalance(c *gin.Context)
	GetBalance(c *gin.Context)
	GetWalletHistory(c *gin.Context)
	ExGetBalance(c *gin.Context)

	CreateWalletLink(c *gin.Context)
	WalletLinkConfirmation(c *gin.Context)
	WalletUnlink(c *gin.Context)
	ExternalTransaction(c *gin.Context)
}
