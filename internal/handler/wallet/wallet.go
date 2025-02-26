package wallet

import (
	"ewallet-wallet/constants"
	"ewallet-wallet/helpers"
	"ewallet-wallet/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Create(c *gin.Context) {
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

	err := h.Service.Create(c.Request.Context(), &req)
	if err != nil {
		fmt.Printf("failed to created wallet: %v\n", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}
	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, req)
}

func (h *Handler) CreditBalance(c *gin.Context) {
	var (
		log = helpers.Logger
		req models.TransactionRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("failed to parse request, %v\n", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		fmt.Printf("failed to validate request, %v\n", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		log.Error("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := h.Service.CreditBalance(c.Request.Context(), tokenData.UserID, req)

	if err != nil {
		fmt.Printf("failed to debit balance of wallet: %v", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (h *Handler) DebitBalance(c *gin.Context) {
	var (
		req models.TransactionRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("failed to parse request, %v\n", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		fmt.Printf("failed to validate request, %v\n", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	token, ok := c.Get("token")
	if !ok {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		fmt.Println("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := h.Service.DebitBalance(c.Request.Context(), tokenData.UserID, req)
	if err != nil {
		fmt.Printf("failed to debit balance of wallet: %v", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (h *Handler) GetBalance(c *gin.Context) {

	token, ok := c.Get("token")
	if !ok {
		fmt.Print("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		fmt.Print("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := h.Service.GetBalance(c.Request.Context(), tokenData.UserID)
	if err != nil {
		fmt.Printf("failed to get balance of wallet, %v\n", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (h *Handler) GetWalletHistory(c *gin.Context) {
	var (
		param models.WalletHistoryParam
	)

	if err := c.ShouldBindQuery(&param); err != nil {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrFailedBadRequest, nil)
		return
	}

	if param.WalletTransactionType != "" {
		if param.WalletTransactionType != "CREDIT" && param.WalletTransactionType != "DEBIT" {
			fmt.Println("invalid wallet transaction_type")
			helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrFailedBadRequest, nil)
			return
		}
	}

	token, ok := c.Get("token")
	if !ok {
		fmt.Println("failed to get token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	tokenData, ok := token.(models.TokenData)
	if !ok {
		fmt.Println("failed to parse token data")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := h.Service.GetWalletHistory(c.Request.Context(), tokenData.UserID, param)
	if err != nil {
		fmt.Printf("failed to get balance of wallet, %v\n", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (h *Handler) CreateWalletLink(c *gin.Context) {
	var (
		req models.WalletLink
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("failed to parse query req: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrFailedBadRequest, nil)
		return
	}

	clientID, ok := c.Get("client_id")
	if !ok {
		fmt.Println("fsdfsdf")
		fmt.Println("fsdfsdf")
		fmt.Println("fsdfsdf")
		fmt.Println("failed to get client id")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	clientSource, ok := clientID.(string)
	if !ok {
		fmt.Println("failed to parse client id")
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	resp, err := h.Service.CreateWalletLink(c.Request.Context(), clientSource, &req)
	if err != nil {
		fmt.Println("failed to create wallet link, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, &resp)
}

func (h *Handler) WalletLinkConfirmation(c *gin.Context) {
	var (
		req *models.WalletLink
	)
	fmt.Println("fdfsf")

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("failed to parse query req: ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	walletIDs := c.Param("wallet_id")
	if walletIDs == "" {
		fmt.Println("failed to get wallet id: ")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		fmt.Println("failed to parse wallet id to int : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	clientID, ok := c.Get("client_id")
	if !ok {
		fmt.Println("failed to get client id")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	clientSource, ok := clientID.(string)
	if !ok {
		fmt.Println("failed to parse client id")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	err = h.Service.WalletLinkConfirmation(c.Request.Context(), walletID, clientSource, req.OTP)
	if err != nil {
		fmt.Println("failed to confirm wallet link, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)
}

func (h *Handler) WalletUnlink(c *gin.Context) {

	walletIDs := c.Param("wallet_id")
	if walletIDs == "" {
		fmt.Println("failed to get wallet id: ")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		fmt.Println("failed to parse wallet id to int : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	clientID, ok := c.Get("client_id")
	if !ok {
		fmt.Println("failed to get client id")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	clientSource, ok := clientID.(string)
	if !ok {
		fmt.Println("failed to parse client id")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrServerError, nil)
		return
	}

	err = h.Service.WalletUnlink(c.Request.Context(), walletID, clientSource)
	if err != nil {
		fmt.Println("failed to unlink wallet: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, nil)
}

func (h *Handler) ExGetBalance(c *gin.Context) {

	walletIDs := c.Param("wallet_id")
	if walletIDs == "" {
		fmt.Println("failed to get wallet id: ")
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		fmt.Println("failed to parse wallet id to int : ", err)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFailedBadRequest, nil)
		return
	}

	resp, err := h.Service.ExGetBalance(c.Request.Context(), walletID)
	if err != nil {
		fmt.Println("failed to unlink wallet: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}

func (h *Handler) ExternalTransaction(c *gin.Context) {
	var (
		req models.ExternalTransactionRequest
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("failed to parse query req: ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrFailedBadRequest, nil)
		return
	}

	resp, err := h.Service.ExternalTransaction(c.Request.Context(), req)
	if err != nil {
		fmt.Println("failed to create external transaction, ", err)
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, constants.ErrServerError, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.SuccessMessage, resp)
}
