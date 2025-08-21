package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-blockchain-receipt/internal/services"
)

type Handler struct {
	receiptService *services.ReceiptService
}

func NewHandler(receiptService *services.ReceiptService) *Handler {
	return &Handler{
		receiptService: receiptService,
	}
}

// HTTPError represents an error response
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CreateReceipt godoc
// @Summary Create a new receipt
// @Description Create a new receipt from transaction payload
// @Tags receipts
// @Accept json
// @Produce json
// @Param payload body map[string]interface{} true "Transaction payload"
// @Success 201 {object} models.CreateReceiptResponse
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /receipts [post]
func (h *Handler) CreateReceipt(c echo.Context) error {
	var payload map[string]interface{}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.receiptService.CreateReceipt(c.Request().Context(), payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

// VerifyReceipt godoc
// @Summary Verify a receipt
// @Description Verify a receipt and its on-chain status
// @Tags receipts
// @Accept json
// @Produce json
// @Param rid query string true "Receipt ID"
// @Param jws query string true "JWS token"
// @Success 200 {object} models.VerifyResponse
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /verify [get]
func (h *Handler) VerifyReceipt(c echo.Context) error {
	rid := c.QueryParam("rid")
	jws := c.QueryParam("jws")

	if rid == "" || jws == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing rid or jws")
	}

	resp, err := h.receiptService.VerifyReceipt(c.Request().Context(), rid, jws)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// GetJWKS godoc
// @Summary Get JWKS
// @Description Get public JWKS for signature verification
// @Tags jwks
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /jwks.json [get]
func (h *Handler) GetJWKS(c echo.Context) error {
	return c.JSON(http.StatusOK, h.receiptService.GetJWKS())
}

// HealthCheck godoc
// @Summary Health check
// @Description Check health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "healthy",
	})
}