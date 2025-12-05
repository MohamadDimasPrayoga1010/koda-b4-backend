package controllers

import (
	"coffeeder-backend/libs"
	"coffeeder-backend/models"
	"context"
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
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
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

	var total int
	totalQuery := "SELECT COUNT(*) FROM transactions t LEFT JOIN users u ON u.id = t.user_id WHERE 1=1"
	var args []interface{}
	if search != "" {
		totalQuery += " AND (LOWER(u.fullname) LIKE LOWER($1) OR LOWER(t.invoice_number) LIKE LOWER($1))"
		args = append(args, "%"+search+"%")
	}

	err = tc.DB.QueryRow(context.Background(), totalQuery, args...).Scan(&total)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to count transactions",
			"data":    err.Error(),
		})
		return
	}

	pagination, links := libs.BuildHateoasGlobal("/transactions", page, limit, total, ctx.Request.URL.Query())

	response := models.ProductListResponse{
		Success:    true,
		Message:    "Transactions fetched successfully",
		Pagination: pagination,
		Links:      links,
		Data:       transactions,
	}

	ctx.JSON(200, response)
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
		StatusID int64 `json:"statusId" binding:"required"`
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
			"transactionId": transactionID,
			"newStatus":     statusName,
			"code":          200,
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
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	var userID int64
	switch v := userIDValue.(type) {
	case int64:
		userID = v
	case int:
		userID = int64(v)
	case float64:
		userID = int64(v)
	case string:
		tmp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Invalid user ID",
			})
			return
		}
		userID = tmp
	default:
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	monthStr := ctx.DefaultQuery("month", "0")
	status := ctx.DefaultQuery("status", "")
	pageStr := ctx.DefaultQuery("page", "1")

	month, _ := strconv.Atoi(monthStr)
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit := 5

	histories, total, err := models.GetHistoryTransactions(tc.DB, userID, status, month, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch history transactions",
			Data:    err.Error(),
		})
		return
	}

	pagination, links := libs.BuildHateoasGlobal("/transactions/history", page, limit, total, ctx.Request.URL.Query())

	response := models.ProductListResponse{
		Success:    true,
		Message:    "History transactions fetched successfully",
		Pagination: pagination,
		Links:      links,
		Data:       histories,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetHistoryDetailById godoc
// @Summary      Get detail of a user's transaction by ID
// @Description  Fetch detailed information of a single transaction including items for the authenticated user
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Transaction ID"
// @Success      200  {object}  models.Response{data=models.HistoryDetail} "Transaction detail fetched successfully"
// @Failure      400  {object}  models.Response "Invalid user ID or transaction ID"
// @Failure      401  {object}  models.Response "Unauthorized"
// @Failure      500  {object}  models.Response "Internal server error"
// @Security     ApiKeyAuth
// @Router       /history/{id} [get]
func (tc *TransactionController) GetHistoryDetailById(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	var userID int64
	switch v := userIDValue.(type) {
	case int64:
		userID = v
	case int:
		userID = int64(v)
	case float64:
		userID = int64(v)
	case string:
		tmp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Invalid user ID",
			})
			return
		}
		userID = tmp
	default:
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
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Transaction detail fetched successfully",
		Data:    history,
	})
}

func (tc *TransactionController) GetShippingMethods(ctx *gin.Context) {
	if _, exists := ctx.Get("userID"); !exists {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	shippings := []models.ShippingMethod{}

	rows, err := tc.DB.Query(ctx, `SELECT id, name FROM shippings ORDER BY id ASC`)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch shipping methods",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s models.ShippingMethod
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			continue 
		}
		shippings = append(shippings, s)
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Shipping methods fetched successfully",
		Data:    shippings,
	})
}


func (tc *TransactionController) GetPaymentMethods(ctx *gin.Context) {
	if _, exists := ctx.Get("userID"); !exists {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	payments := []models.PaymentMethod{}

	rows, err := tc.DB.Query(ctx, `SELECT id, name, image FROM payment_methods ORDER BY id ASC`)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch payment methods",
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PaymentMethod
		if err := rows.Scan(&p.ID, &p.Name, &p.Image); err != nil {
			continue 
		}
		payments = append(payments, p)
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Payment methods fetched successfully",
		Data:    payments,
	})
}
