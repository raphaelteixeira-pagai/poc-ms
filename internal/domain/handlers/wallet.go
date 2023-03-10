package handlers

import (
	"github.com/raphaelteixeira-pagai/poc-ms/cmd/httphdl"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golangsugar/chatty"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/entities"
	"github.com/raphaelteixeira-pagai/poc-ms/internal/domain/services"
)

type wallethandlers struct {
	srv services.IWalletService
}

func RegisterWalletRoutes(srv services.IWalletService, router *gin.RouterGroup) {
	hdl := &wallethandlers{srv: srv}
	walletGroup := router.Group("wallet")
	walletGroup.
		GET("/:owner", httphdl.MetricMiddleware("wallet-get", "count"), hdl.FetchWallet).
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
	body := entities.Wallet{}
	if err := c.ShouldBindJSON(&body); err != nil {
		chatty.Error("internal error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": "invalid request"})
		return
	}

	err := w.srv.Create(c, body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Wallet Created")
}

func (w *wallethandlers) DepositWallet(c *gin.Context) {
	body := entities.Wallet{}
	if err := c.ShouldBindJSON(&body); err != nil {
		chatty.Error("internal error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": "invalid request"})
		return
	}

	err := w.srv.Deposit(c, body.Balance, body.Owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Wallet Deposited")
}

func (w *wallethandlers) WithdrawWallet(c *gin.Context) {
	body := entities.Wallet{}
	if err := c.ShouldBindJSON(&body); err != nil {
		chatty.Error("internal error")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "message": "invalid request"})
		return
	}

	err := w.srv.Withdraw(c, body.Balance, body.Owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Wallet Withdrawn")
}

func (w *wallethandlers) DeleteWallet(c *gin.Context) {
	owner := c.Param("owner")

	err := w.srv.Delete(c, owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Wallet Deleted")
}
