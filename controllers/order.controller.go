package controllers

import (
	"context"
	"fmt"
	"main/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderController struct {
	DB *pgxpool.Pool
}

func (oc *OrderController) GetOrders(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	query := `
		SELECT 
			o.id, o.no_orders, o.total, o.created_at,
			s.id AS status_id, s.name AS status_name,
			u.id AS user_id, u.fullname AS user_name,
			pm.id AS payment_method_id, pm.name AS payment_method,
			sh.name AS shipping_name
		FROM orders o
		JOIN users u ON u.id = o.user_id
		JOIN status s ON s.id = o.status_id
		JOIN payment_methods pm ON pm.id = o.payment_method
		JOIN shippings sh ON sh.id = o.shipping_id
		WHERE 1=1
	`

	limitation := []interface{}{}
	if search != "" {
		query += " AND (LOWER(u.fullname) LIKE LOWER($1) OR LOWER(o.no_orders) LIKE LOWER($1))"
		limitation = append(limitation, "%"+search+"%")
		query += " ORDER BY o.id DESC LIMIT $2 OFFSET $3"
		limitation = append(limitation, limit, offset)
	} else {
		query += " ORDER BY o.id DESC LIMIT $1 OFFSET $2"
		limitation = append(limitation, limit, offset)
	}

	rows, err := oc.DB.Query(context.Background(), query, limitation...)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch orders",
			Data:    err.Error(),
		})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID, &o.NoOrders, &o.Total, &o.CreatedAt,
			&o.StatusID, &o.StatusName,
			&o.UserID, &o.UserName,
			&o.PaymentMethodID, &o.PaymentMethod,
			&o.ShippingName,
		); err != nil {
			fmt.Println("scan error:", err)
			continue
		}

		// Get order items
		itemRows, err := oc.DB.Query(context.Background(), `
			SELECT p.title, s.name AS size, op.qty
			FROM orders_products op
			JOIN products p ON p.id = op.product_id
			JOIN sizes s ON s.id = op.size_id
			WHERE op.orderd_id = $1
		`, o.ID)
		if err != nil {
			fmt.Println("order items query error:", err)
		} else {
			for itemRows.Next() {
				var item models.OrderItemResp
				itemRows.Scan(&item.Title, &item.Size, &item.Qty)
				o.OrderItems = append(o.OrderItems, item)
			}
			itemRows.Close()
		}

		orders = append(orders, o)
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Orders fetched successfully",
		Data:    orders,
	})
}

func (oc *OrderController) UpdateOrderStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		StatusID int64 `json:"status_id" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	query := `
		UPDATE orders SET status_id=$1, updated_at=now()
		WHERE id=$2 RETURNING id
	`
	var orderID int64
	err := oc.DB.QueryRow(context.Background(), query, req.StatusID, id).Scan(&orderID)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to update order status",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Order status updated successfully",
		Data:    map[string]interface{}{"order_id": orderID, "new_status_id": req.StatusID},
	})
}

func (oc *OrderController) DeleteOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := oc.DB.Exec(context.Background(), `DELETE FROM orders_products WHERE orderd_id = $1`, id)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete order items",
			Data:    err.Error(),
		})
		return
	}

	_, err = oc.DB.Exec(context.Background(), `DELETE FROM orders WHERE id = $1`, id)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete order",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Order deleted successfully",
		Data:    nil,
	})
}
