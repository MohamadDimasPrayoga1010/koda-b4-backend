package models

import (
	"context"
	"strconv"
	"strings"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionItem struct {
	Title string `json:"title"`
	Qty   int    `json:"qty"`
	Size  string `json:"size,omitempty"`
}

type Transaction struct {
	ID            int64             `json:"id"`
	NoOrders      string            `json:"no_orders"`
	CreatedAt     time.Time         `json:"created_at"`
	StatusName    string            `json:"status_name"`
	Total         float64           `json:"total"`
	UserFullname  string            `json:"user_fullname"`
	UserAddress   string            `json:"user_address"`
	UserPhone     string            `json:"user_phone"`
	PaymentMethod string            `json:"payment_method"`
	ShippingName  string            `json:"shipping_name"`
	OrderItems    []TransactionItem `json:"order_items"`
}

func GetAllTransactions(db *pgxpool.Pool, search, sort, order string, limit, offset int) ([]Transaction, error) {
	queriParams := make([]interface{}, 0)
	argIdx := 1

	allowedSort := map[string]bool{"created_at": true, "total": true, "invoice_number": true}
	if !allowedSort[sort] {
		sort = "created_at"
	}
	if strings.ToLower(order) != "asc" && strings.ToLower(order) != "desc" {
		order = "desc"
	}

	query := `
		SELECT 
			t.id,
			t.invoice_number AS no_orders,
			t.created_at,
			t.status AS status_name,
			u.fullname AS user_fullname,
			t.address AS user_address,
			t.phone AS user_phone,
			pm.name AS payment_method,
			sh.name AS shipping_name,
			COALESCE(SUM(pr.base_price * ti.quantity), 0) AS total
		FROM transactions t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN payment_methods pm ON pm.id = t.payment_method_id
		LEFT JOIN shippings sh ON sh.id = t.shipping_id
		LEFT JOIN transaction_items ti ON ti.transaction_id = t.id
		LEFT JOIN products pr ON pr.id = ti.product_id
		WHERE 1=1
	`

	if search != "" {
		query += " AND (LOWER(u.fullname) LIKE LOWER($" + strconv.Itoa(argIdx) + ") OR LOWER(t.invoice_number) LIKE LOWER($" + strconv.Itoa(argIdx) + "))"
		queriParams = append(queriParams, "%"+search+"%")
		argIdx++
	}

	query += " GROUP BY t.id, u.fullname, t.address, t.phone, pm.name, sh.name"

	query += " ORDER BY " + sort + " " + order

	query += " LIMIT $" + strconv.Itoa(argIdx) + " OFFSET $" + strconv.Itoa(argIdx+1)
	queriParams = append(queriParams, limit, offset)

	rows, err := db.Query(context.Background(), query, queriParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		t.OrderItems = make([]TransactionItem, 0)

		if err := rows.Scan(
			&t.ID, &t.NoOrders, &t.CreatedAt, &t.StatusName,
			&t.UserFullname, &t.UserAddress, &t.UserPhone,
			&t.PaymentMethod, &t.ShippingName, &t.Total,
		); err != nil {
			continue
		}

		itemRows, _ := db.Query(context.Background(), `
			SELECT pr.title, s.name AS size, ti.quantity AS qty
			FROM transaction_items ti
			JOIN products pr ON pr.id = ti.product_id
			LEFT JOIN sizes s ON s.id = ti.size_id
			WHERE ti.transaction_id=$1
		`, t.ID)

		for itemRows.Next() {
			var item TransactionItem
			itemRows.Scan(&item.Title, &item.Size, &item.Qty)
			t.OrderItems = append(t.OrderItems, item)
		}
		itemRows.Close()

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func GetTransactionByID(db *pgxpool.Pool, id string) (Transaction, error) {
	var t Transaction
	t.OrderItems = make([]TransactionItem, 0)

	query := `
		SELECT 
			t.id,
			t.invoice_number AS no_orders,
			t.created_at,
			t.status AS status_name,
			u.fullname AS user_fullname,
			t.address AS user_address,
			t.phone AS user_phone,
			pm.name AS payment_method,
			sh.name AS shipping_name,
			COALESCE(SUM(pr.base_price * ti.quantity), 0) AS total
		FROM transactions t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN payment_methods pm ON pm.id = t.payment_method_id
		LEFT JOIN shippings sh ON sh.id = t.shipping_id
		LEFT JOIN transaction_items ti ON ti.transaction_id = t.id
		LEFT JOIN products pr ON pr.id = ti.product_id
		WHERE t.id=$1
		GROUP BY t.id, u.fullname, t.address, t.phone, pm.name, sh.name
	`

	err := db.QueryRow(context.Background(), query, id).Scan(
		&t.ID, &t.NoOrders, &t.CreatedAt, &t.StatusName,
		&t.UserFullname, &t.UserAddress, &t.UserPhone,
		&t.PaymentMethod, &t.ShippingName, &t.Total,
	)
	if err != nil {
		return t, err
	}

	itemRows, _ := db.Query(context.Background(), `
		SELECT pr.title, s.name AS size, ti.quantity AS qty
		FROM transaction_items ti
		JOIN products pr ON pr.id = ti.product_id
		LEFT JOIN sizes s ON s.id = ti.size_id
		WHERE ti.transaction_id=$1
	`, t.ID)
	for itemRows.Next() {
		var item TransactionItem
		itemRows.Scan(&item.Title, &item.Size, &item.Qty)
		t.OrderItems = append(t.OrderItems, item)
	}
	itemRows.Close()

	return t, nil
}


type HistoryTransaction struct {
	ID            int64     `json:"id"`
	InvoiceNumber string    `json:"invoice_number"`
	Image         string    `json:"image"`
	Total         float64   `json:"total"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}


func GetHistoryTransactions(db *pgxpool.Pool, userID int64, status string, month, page, limit int) ([]HistoryTransaction, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	query := `
	SELECT 
		t.id,
		t.invoice_number,
		t.total,
		t.status,
		t.created_at,
		COALESCE(MAX(pi.image), '') AS image
	FROM transactions t
	LEFT JOIN transaction_items ti ON ti.transaction_id = t.id
	LEFT JOIN products p ON p.id = ti.product_id
	LEFT JOIN product_images pi ON pi.product_id = p.id
	WHERE t.user_id = $1
	  AND t.status ILIKE $2
`

	queriParams := []interface{}{userID, status}

	if month >= 1 && month <= 12 {
		query += " AND EXTRACT(MONTH FROM t.created_at) = $" + strconv.Itoa(len(queriParams)+1)
		queriParams = append(queriParams, month)
	}

	query += `
		GROUP BY t.id, t.invoice_number, t.total, t.status, t.created_at
		ORDER BY t.created_at DESC
		LIMIT $` + strconv.Itoa(len(queriParams)+1) + " OFFSET $" + strconv.Itoa(len(queriParams)+2)
	queriParams = append(queriParams, limit, offset)

	rows, err := db.Query(ctx, query, queriParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []HistoryTransaction
	for rows.Next() {
		var h HistoryTransaction
		if err := rows.Scan(&h.ID, &h.InvoiceNumber, &h.Total, &h.Status, &h.CreatedAt, &h.Image); err != nil {
			return nil, err
		}
		histories = append(histories, h)
	}

	return histories, nil
}