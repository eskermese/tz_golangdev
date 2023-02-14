package restHandler

import (
	"github.com/eskermese/tz_golangdev/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initTransactionRoutes(api fiber.Router) {
	transactions := api.Group("/transactions")
	{
		transactions.Get("/maximum-change", h.getMaximumChange)
	}
}

// @Summary Get Order total costs
// @Tags transactions
// @Description Get the address of the account which balance changed the most(also provides the receiver address) over the last 100 blocks
// @ModuleID orderTotalCosts
// @Accept  json
// @Produce  json
// @Success 200 {object} core.Transaction
// @Failure 400 {object} core.Error
// @Router /v1/transactions/maximum-change [get]
func (h *Handler) getMaximumChange(c *fiber.Ctx) error {
	transaction, err := h.service.GetMaxTransactionChange(c.Context())
	if err != nil {
		h.logger.Error("getMaximumChange: internal server error", logger.Error(err))

		return c.Status(fiber.StatusInternalServerError).JSON(response{Detail: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(transaction)
}
