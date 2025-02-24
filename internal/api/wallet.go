package api

import (
	"ewallet-wallet/constants"
	"ewallet-wallet/helpers"
	"ewallet-wallet/internal/interfaces/i_service"
	"ewallet-wallet/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletAPI struct {
	WalletService i_service.IWalletService
}

func (api *WalletAPI) Create(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.Wallet
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse request,", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if req.UserID == 0 {
		log.Error("user id is empty")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	err := api.WalletService.Create(c.Request.Context(), &req)
	if err != nil {
		log.Error("failed to create wallet:", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, req)
}

func (api *WalletAPI) CreditBalance(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.TransactionRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		log.Error("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		log.Error("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := api.WalletService.CreditBalance(c.Request.Context(), tokenData.UserID, req)
	if err != nil {
		log.Error("failed to debit balance of wallet:", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (api *WalletAPI) DebitBalance(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.TransactionRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		log.Error("failed to validate request, ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		log.Error("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		log.Error("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := api.WalletService.DebitBalance(c.Request.Context(), tokenData.UserID, req)
	if err != nil {
		log.Error("failed to debit of wallet, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (api *WalletAPI) GetBalance(c *gin.Context) {
	var (
		log = helpers.Logger
	)

	token, ok := c.Get("token")
	if !ok {
		log.Error("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		log.Error("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := api.WalletService.GetBalance(c.Request.Context(), tokenData.UserID)
	if err != nil {
		log.Error("failed to get balance of wallet, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (api *WalletAPI) GetWalletHistory(c *gin.Context) {
	var (
		log   = helpers.Logger
		param models.WalletHistoryParam
	)

	if err := c.ShouldBindQuery(&param); err != nil {
		log.Error("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrFailedBadRequest, nil)
		return
	}

	if param.WalletTransactionType != "" {
		if param.WalletTransactionType != "CREDIT" && param.WalletTransactionType != "DEBIT" {
			log.Error("invalid wallet transaction_type")
			helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrFailedBadRequest, nil)
			return
		}
	}

	token, ok := c.Get("token")
	if !ok {
		log.Error("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		log.Error("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := api.WalletService.GetWalletHistory(c.Request.Context(), tokenData.UserID, param)
	if err != nil {
		log.Error("failed to get balance of wallet, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (api *WalletAPI) CreateWalletLink(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.WalletLink
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse query req: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrFailedBadRequest, nil)
		return
	}

	clientID, ok := c.Get("client_id")
	if !ok {
		log.Error("failed to get client id")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	clientSource, ok := clientID.(string)
	if !ok {
		log.Error("failed to parse client id")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := api.WalletService.CreateWalletLink(c.Request.Context(), clientSource, &req)
	if err != nil {
		log.Error("failed to create wallet link, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, &resp)
}

func (api *WalletAPI) WalletLinkConfirmation(c *gin.Context) {
	var (
		log = helpers.Logger
		req *models.WalletLink
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse query req: ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	walletIDs := c.Param("wallet_id")
	if walletIDs == "" {
		log.Error("failed to get wallet id: ")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		log.Error("failed to parse wallet id to int : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	clientID, ok := c.Get("client_id")
	if !ok {
		log.Error("failed to get client id")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	clientSource, ok := clientID.(string)
	if !ok {
		log.Error("failed to parse client id")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	err = api.WalletService.WalletLinkConfirmation(c.Request.Context(), walletID, clientSource, req.OTP)
	if err != nil {
		log.Error("failed to confirm wallet link, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)
}

func (api *WalletAPI) WalletUnlink(c *gin.Context) {
	var (
		log = helpers.Logger
	)

	walletIDs := c.Param("wallet_id")
	if walletIDs == "" {
		log.Error("failed to get wallet id: ")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		log.Error("failed to parse wallet id to int : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	clientID, ok := c.Get("client_id")
	if !ok {
		log.Error("failed to get client id")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	clientSource, ok := clientID.(string)
	if !ok {
		log.Error("failed to parse client id")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	err = api.WalletService.WalletUnlink(c.Request.Context(), walletID, clientSource)
	if err != nil {
		log.Error("failed to unlink wallet: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)
}

func (api *WalletAPI) ExGetBalance(c *gin.Context) {
	var (
		log = helpers.Logger
	)

	walletIDs := c.Param("wallet_id")
	if walletIDs == "" {
		log.Error("failed to get wallet id: ")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		log.Error("failed to parse wallet id to int : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	resp, err := api.WalletService.ExGetBalance(c.Request.Context(), walletID)
	if err != nil {
		log.Error("failed to unlink wallet: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (api *WalletAPI) ExternalTransaction(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.ExternalTransactionRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error("failed to parse query req: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrFailedBadRequest, nil)
		return
	}

	resp, err := api.WalletService.ExternalTransaction(c.Request.Context(), req)
	if err != nil {
		log.Error("failed to create external transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}
