package controllers

import (
	"context"
	"main/models"
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
// @Param search query string false "Search by user fullname or order number"
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

	args := []interface{}{}
	argIdx := 1

	query := `
		SELECT 
			o.id, o.no_orders, o.created_at,
			s.name AS status_name,
			u.fullname AS user_fullname, p.address AS user_address, p.phone AS user_phone,
			pm.name AS payment_method, sh.name AS shipping_name,
			COALESCE(SUM(pr.base_price * op.qty), 0) AS total
		FROM orders o
		JOIN users u ON u.id = o.user_id
		LEFT JOIN profile p ON p.user_id = u.id
		JOIN status s ON s.id = o.status_id
		JOIN payment_methods pm ON pm.id = o.payment_method
		JOIN shippings sh ON sh.id = o.shipping_id
		JOIN orders_products op ON op.order_id = o.id
		JOIN products pr ON pr.id = op.product_id
		WHERE 1=1
	`

	if search != "" {
		query += " AND (LOWER(u.fullname) LIKE LOWER($" + strconv.Itoa(argIdx) + ") OR LOWER(o.no_orders) LIKE LOWER($" + strconv.Itoa(argIdx) + "))"
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += " GROUP BY o.id, s.name, u.fullname, p.address, p.phone, pm.name, sh.name"
	query += " ORDER BY " + sort + " " + order
	query += " LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	args = append(args, limit, offset)

	rows, err := tc.DB.Query(context.Background(), query, args...)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch transactions",
			Data:    err.Error(),
		})
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(
			&t.ID, &t.NoOrders, &t.CreatedAt, &t.StatusName,
			&t.UserFullname, &t.UserAddress, &t.UserPhone,
			&t.PaymentMethod, &t.ShippingName, &t.Total,
		); err != nil {
			continue
		}
		
		itemRows, err := tc.DB.Query(context.Background(), `
			SELECT pr.title, s.name AS size, op.qty
			FROM orders_products op
			JOIN products pr ON pr.id = op.product_id
			LEFT JOIN sizes s ON s.id = op.size_id
			WHERE op.order_id=$1
		`, t.ID)
		if err == nil {
			for itemRows.Next() {
				var item models.TransactionItem
				itemRows.Scan(&item.Title, &item.Size, &item.Qty)
				t.OrderItems = append(t.OrderItems, item)
			}
			itemRows.Close()
		}

		transactions = append(transactions, t)
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Transactions fetched successfully",
		Data:    transactions,
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
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/transactions/{id} [get]
func (tc *TransactionController) GetTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
		SELECT 
			o.id, o.no_orders, o.created_at,
			s.name AS status_name,
			u.fullname AS user_fullname, p.address AS user_address, p.phone AS user_phone,
			pm.name AS payment_method, sh.name AS shipping_name,
			COALESCE(SUM(pr.base_price * op.qty), 0) AS total
		FROM orders o
		JOIN users u ON u.id = o.user_id
		LEFT JOIN profile p ON p.user_id = u.id
		JOIN status s ON s.id = o.status_id
		JOIN payment_methods pm ON pm.id = o.payment_method
		JOIN shippings sh ON sh.id = o.shipping_id
		JOIN orders_products op ON op.order_id = o.id
		JOIN products pr ON pr.id = op.product_id
		WHERE o.id=$1
		GROUP BY o.id, s.name, u.fullname, p.address, p.phone, pm.name, sh.name
	`

	var t models.Transaction
	err := tc.DB.QueryRow(context.Background(), query, id).Scan(
		&t.ID, &t.NoOrders, &t.CreatedAt, &t.StatusName,
		&t.UserFullname, &t.UserAddress, &t.UserPhone,
		&t.PaymentMethod, &t.ShippingName, &t.Total,
	)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	itemRows, err := tc.DB.Query(context.Background(), `
		SELECT pr.title, s.name AS size, op.qty
		FROM orders_products op
		JOIN products pr ON pr.id = op.product_id
		LEFT JOIN sizes s ON s.id = op.size_id
		WHERE op.order_id=$1
	`, t.ID)
	if err == nil {
		for itemRows.Next() {
			var item models.TransactionItem
			itemRows.Scan(&item.Title, &item.Size, &item.Qty)
			t.OrderItems = append(t.OrderItems, item)
		}
		itemRows.Close()
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Transaction fetched successfully",
		Data:    t,
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
		ctx.JSON(200, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    map[string]int{"code": 400},
		})
		return
	}

	query := `UPDATE orders SET status_id=$1, updated_at=now() WHERE id=$2 RETURNING id`
	var transactionID int64
	err := tc.DB.QueryRow(context.Background(), query, req.StatusID, id).Scan(&transactionID)
	if err != nil {
		ctx.JSON(200, models.Response{
			Success: false,
			Message: "Failed to update transaction status",
			Data:    map[string]int{"code": 500},
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Transaction status updated successfully",
		Data:    map[string]interface{}{"transaction_id": transactionID, "new_status_id": req.StatusID, "code": 200},
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

	_, err := tc.DB.Exec(context.Background(), `DELETE FROM orders_products WHERE order_id=$1`, id)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete transaction items",
			Data:    err.Error(),
		})
		return
	}

	_, err = tc.DB.Exec(context.Background(), `DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete transaction",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Transaction deleted successfully",
	})
}
