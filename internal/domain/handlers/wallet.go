package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golangsugar/chatty"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/services"
	"net/http"
)

type wallethandlers struct {
	srv services.IWalletService
}

func RegisterWalletRoutes(srv services.IWalletService, router *gin.RouterGroup) {
	hdl := &wallethandlers{srv: srv}
	walletGroup := router.Group("wallet")
	walletGroup.
		GET("/:owner", hdl.FetchWallet).
		POST("/", hdl.CreateWallet).
		PUT("/withdraw/:owner", hdl.WithdrawWallet).
		PUT("/deposit/:owner", hdl.DepositWallet).
		DELETE("/:owner", hdl.DeleteWallet)
}

func (w *wallethandlers) FetchWallet(c *gin.Context) {
	owner := c.Param("owner")
	if owner == "" {
		chatty.Error("owner empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "field violation", "message": "owner empty"})
		return
	}

	wallet, err := w.srv.Fetch(c, owner)
	if err != nil {
		chatty.Error("internal error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "internal", "message": "could not retrieve wallet"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (w *wallethandlers) CreateWallet(c *gin.Context) {

}

func (w *wallethandlers) DepositWallet(c *gin.Context) {

}

func (w *wallethandlers) WithdrawWallet(c *gin.Context) {

}

func (w *wallethandlers) DeleteWallet(c *gin.Context) {

}
