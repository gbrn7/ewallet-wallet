package cmd

import (
	"ewallet-wallet/external"
	"ewallet-wallet/helpers"
	healthHandler "ewallet-wallet/internal/handler/healthcheck"
	walletHandler "ewallet-wallet/internal/handler/wallet"
	"ewallet-wallet/internal/repository"
	"ewallet-wallet/internal/services"
	"ewallet-wallet/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func ServeHttp() {

	r := gin.Default()

	healthcheckSvc := &services.Healthcheck{}

	walletRepo := &repository.WalletRepo{
		DB: helpers.DB,
	}
	walletSvc := &services.WalletService{
		WalletRepo: walletRepo,
	}

	external := &external.External{}

	middleware := &middleware.ExternalDependency{
		External: external,
	}

	walletHandler := walletHandler.NewHandler(r, walletSvc, external, middleware)
	walletHandler.RegisterRoute()

	healthcheckHandler := healthHandler.NewHandler(r, healthcheckSvc)
	healthcheckHandler.RegisterRoute()

	err := r.Run(":" + helpers.GetEnv("PORT", ""))
	if err != nil {
		log.Fatal(err)
	}
}
