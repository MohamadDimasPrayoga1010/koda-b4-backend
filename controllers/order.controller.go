package controllers

import (
	"context"
	"main/models"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionController struct {
	DB *pgxpool.Pool
}


// GetTransactions godoc
// @Summary Get all transactions
// @Description Mengambil daftar semua transaksi dengan pagination, search, dan sort (Admin Only)
// @Tags Transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param search query string false "Search by user fullname or invoice number"
// @Param sort query string false "Sort by column e.g. created_at, total" default("created_at")
// @Param order query string false "Sort order asc or desc" default("desc")
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/transactions [get]
func (tc *TransactionController) GetTransactions(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")
	sort := ctx.DefaultQuery("sort", "created_at")
	order := ctx.DefaultQuery("order", "desc")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	transactions, err := models.GetAllTransactions(tc.DB, search, sort, order, limit, offset)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to fetch transactions",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Transactions fetched successfully",
		"data":    transactions,
	})
}

// GetTransactionByID godoc
// @Summary Get transaction by ID
// @Description Mengambil detail transaksi berdasarkan ID (Admin Only)
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/transactions/{id} [get]
func (tc *TransactionController) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	transaction, err := models.GetTransactionByID(tc.DB, id)
	if err != nil {
		ctx.JSON(404, gin.H{
			"success": false,
			"message": "Transaction not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Transaction fetched successfully",
		"data":    transaction,
	})
}

// UpdateTransactionStatus godoc
// @Summary Update transaction status
// @Description Admin dapat mengupdate status transaksi berdasarkan ID
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Param body body map[string]int true "StatusID payload"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/transactions/{id}/status [patch]
func (tc *TransactionController) UpdateTransactionStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		StatusID int64 `json:"status_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Invalid request body",
			"data":    map[string]int{"code": 400},
		})
		return
	}

	var statusName string
	err := tc.DB.QueryRow(context.Background(), `SELECT name FROM status WHERE id=$1`, req.StatusID).Scan(&statusName)
	if err != nil {
		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Status not found",
			"data":    map[string]int{"code": 400},
		})
		return
	}

	query := `UPDATE transactions SET status=$1, updated_at=now() WHERE id=$2 RETURNING id`
	var transactionID int64
	err = tc.DB.QueryRow(context.Background(), query, statusName, id).Scan(&transactionID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to update transaction status",
			"data":    map[string]int{"code": 500},
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Transaction status updated successfully",
		"data": map[string]interface{}{
			"transaction_id": transactionID,
			"new_status":     statusName,
			"code":           200,
		},
	})
}

// DeleteTransaction godoc
// @Summary Delete transaction
// @Description Admin dapat menghapus transaksi beserta items terkait
// @Tags Transactions
// @Accept json
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/transactions/{id} [delete]
func (tc *TransactionController) DeleteTransaction(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := tc.DB.Exec(context.Background(), `DELETE FROM transaction_items WHERE transaction_id=$1`, id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete transaction items",
			"data":    err.Error(),
		})
		return
	}

	_, err = tc.DB.Exec(context.Background(), `DELETE FROM transactions WHERE id=$1`, id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete transaction",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Transaction deleted successfully",
	})
}


// GetHistoryTransactions godoc
// @Summary Get user transaction history
// @Description Fetch transaction history for authenticated user. Supports filter by status, month, pagination, and limit.
// @Tags Transactions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param status query string false "Transaction status filter" default(on progress)
// @Param month query int false "Month filter (1-12)" default(0)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(5)
// @Router /history [get]
func (tc *TransactionController) GetHistoryTransactions(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	monthStr := ctx.DefaultQuery("month", "0") 
	status := ctx.DefaultQuery("status", "OnProgress")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	month, _ := strconv.Atoi(monthStr)
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	histories, err := models.GetHistoryTransactions(tc.DB, userID, status, month, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "History transactions fetched successfully",
		"data":    histories,
	})
}


func (tc *TransactionController) GetHistoryDetailById(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	userID, ok := userIDValue.(int64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	transactionID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid transaction ID",
		})
		return
	}

	history, err := models.GetHistoryDetail(tc.DB, transactionID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Transaction detail fetched successfully",
		Data: history,
	})
}



