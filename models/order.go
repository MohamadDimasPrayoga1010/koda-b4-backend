package models

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionItem struct {
	Title   string `json:"title"`
	Qty     int    `json:"qty"`
	Size    string `json:"size,omitempty"`
	Image   string `json:"image,omitempty"`
	Variant string `json:"variant,omitempty"`
}

type Transaction struct {
	ID            int64             `json:"id"`
	NoOrders      string            `json:"noOrders"`
	CreatedAt     time.Time         `json:"createdAt"`
	StatusName    string            `json:"statusName"`
	Total         float64           `json:"total"`
	UserFullname  string            `json:"userFullname"`
	UserAddress   string            `json:"userAddress"`
	UserPhone     string            `json:"userPhone"`
	PaymentMethod string            `json:"paymentMethod"`
	ShippingName  string            `json:"shippingName"`
	VariantName   *string           `json:"variant,omitempty"`
	OrderItems    []TransactionItem `json:"orderItems"`
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
			fmt.Println("Scan transaction error:", err)
			continue
		}

		itemRows, err := db.Query(context.Background(), `
            SELECT 
                pr.title, 
                s.name AS size, 
                v.name AS variant, 
                ti.quantity AS qty, 
                COALESCE(pi.image, '') AS image
            FROM transaction_items ti
            JOIN products pr ON pr.id = ti.product_id
            LEFT JOIN sizes s ON s.id = ti.size_id
            LEFT JOIN variants v ON v.id = ti.variant_id
            LEFT JOIN LATERAL (
                SELECT image 
                FROM product_images 
                WHERE product_id = pr.id
                ORDER BY id ASC
                LIMIT 1
            ) pi ON true
            WHERE ti.transaction_id=$1
        `, t.ID)
		if err != nil {
			fmt.Println("Query order items error:", err)
			continue
		}

		for itemRows.Next() {
			var item TransactionItem
			var sizeName *string
			var variantName *string
			var imageName string

			if err := itemRows.Scan(&item.Title, &sizeName, &variantName, &item.Qty, &imageName); err != nil {
				fmt.Println("Scan order item error:", err)
				continue
			}

			if sizeName != nil {
				item.Size = *sizeName
			}
			if variantName != nil {
				item.Variant = *variantName
			}
			item.Image = imageName

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

	itemRows, err := db.Query(context.Background(), `
		SELECT 
			pr.title, 
			s.name AS size, 
			v.name AS variant, 
			ti.quantity AS qty, 
			COALESCE(pi.image, '') AS image
		FROM transaction_items ti
		JOIN products pr ON pr.id = ti.product_id
		LEFT JOIN sizes s ON s.id = ti.size_id
		LEFT JOIN variants v ON v.id = ti.variant_id
		LEFT JOIN LATERAL (
			SELECT image 
			FROM product_images 
			WHERE product_id = pr.id
			ORDER BY id ASC
			LIMIT 1
		) pi ON true
		WHERE ti.transaction_id=$1
	`, t.ID)
	if err != nil {
		return t, err
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var item TransactionItem
		var sizeName *string
		var variantName *string
		var imageName string

		if err := itemRows.Scan(&item.Title, &sizeName, &variantName, &item.Qty, &imageName); err != nil {
			fmt.Println("Scan order item error:", err)
			continue
		}

		if sizeName != nil {
			item.Size = *sizeName
		}
		if variantName != nil {
			item.Variant = *variantName
		}
		item.Image = imageName

		t.OrderItems = append(t.OrderItems, item)
	}

	return t, nil
}

type HistoryTransaction struct {
	ID            int64     `json:"id"`
	InvoiceNumber string    `json:"invoiceNumber"`
	Image         string    `json:"image"`
	Total         float64   `json:"total"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	ShippingName  string    `json:"shippingName"`  
	ShippingFee   float64   `json:"shippingFee"`
}

func GetHistoryTransactions(db *pgxpool.Pool, userID int64, status string, month, page, limit int) ([]HistoryTransaction, int, error) {
	ctx := context.Background()
	offset := (page - 1) * limit

	totalQuery := "SELECT COUNT(DISTINCT t.id) FROM transactions t WHERE t.user_id=$1"
	params := []interface{}{userID}

	if status != "" {
		totalQuery += " AND t.status ILIKE $" + strconv.Itoa(len(params)+1)
		params = append(params, status)
	}

	if month >= 1 && month <= 12 {
		totalQuery += " AND EXTRACT(MONTH FROM t.created_at) = $" + strconv.Itoa(len(params)+1)
		params = append(params, month)
	}

	var total int
	if err := db.QueryRow(ctx, totalQuery, params...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
	SELECT 
		t.id,
		t.invoice_number,
		t.total,
		t.status,
		t.created_at,
		COALESCE(s.name, '') AS shipping_name,
		COALESCE(s.additional_price, 0) AS shipping_fee,
		COALESCE(MAX(pi.image), '') AS image
	FROM transactions t
	LEFT JOIN shippings s ON s.id = t.shipping_id
	LEFT JOIN transaction_items ti ON ti.transaction_id = t.id
	LEFT JOIN products p ON p.id = ti.product_id
	LEFT JOIN product_images pi ON pi.product_id = p.id
	WHERE t.user_id = $1
	`
	params = []interface{}{userID}

	if status != "" {
		query += " AND t.status ILIKE $" + strconv.Itoa(len(params)+1)
		params = append(params, status)
	}
	if month >= 1 && month <= 12 {
		query += " AND EXTRACT(MONTH FROM t.created_at) = $" + strconv.Itoa(len(params)+1)
		params = append(params, month)
	}

	query += `
	GROUP BY t.id, s.name, s.additional_price
	ORDER BY t.created_at DESC
	LIMIT $` + strconv.Itoa(len(params)+1) + " OFFSET $" + strconv.Itoa(len(params)+2)
	params = append(params, limit, offset)

	rows, err := db.Query(ctx, query, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var histories []HistoryTransaction
	for rows.Next() {
		var h HistoryTransaction
		if err := rows.Scan(&h.ID, &h.InvoiceNumber, &h.Total, &h.Status, &h.CreatedAt, &h.ShippingName, &h.ShippingFee, &h.Image); err != nil {
			continue
		}
		histories = append(histories, h)
	}

	return histories, total, nil
}

type HistoryDetail struct {
	ID             int64                   `json:"id"`
	InvoiceNumber  string                  `json:"invoice"`
	CustName       string                  `json:"custName"`
	CustPhone      string                  `json:"custPhone"`
	CustEmail      string                  `json:"custEmail"`
	CustAddress    string                  `json:"custAddress"`
	PaymentMethod  string                  `json:"paymentMethod"`
	DeliveryMethod string                  `json:"deliveryMethod"`
	Status         string                  `json:"status"`
	Total          float64                 `json:"total"`
	CreatedAt      string                  `json:"createdAt"`
	ShippingPrice  float64                 `json:"shippingPrice"`
	Items          []TransactionItemDetail `json:"items"`
}

type TransactionItemDetail struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Image         string  `json:"image"`
	Size          *string `json:"size,omitempty"`
	BasePrice     float64 `json:"basePrice"`
	DiscountPrice float64 `json:"discountPrice"`
	Variant       *string `json:"variant,omitempty"`
	Quantity      int     `json:"quantity"`
	Subtotal      float64 `json:"subtotal"`
}

type ShippingMethod struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PaymentMethod struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func GetHistoryDetail(db *pgxpool.Pool, transactionID, userID int64) (*HistoryDetail, error) {
	ctx := context.Background()

	queryHeader := `
	SELECT 
		t.id,
		t.invoice_number,
		t.fullname,
		t.phone,
		t.email,
		t.address,
		pm.name AS payment_method,
		s.name AS delivery_method,
		t.status,
		t.total,
		TO_CHAR(t.created_at, 'YYYY-MM-DD HH24:MI:SS') AS created_at
	FROM transactions t
	LEFT JOIN payment_methods pm ON pm.id = t.payment_method_id
	LEFT JOIN shippings s ON s.id = t.shipping_id
	WHERE t.id = $1 AND t.user_id = $2
	`
	var header HistoryDetail
	err := db.QueryRow(ctx, queryHeader, transactionID, userID).Scan(
		&header.ID,
		&header.InvoiceNumber,
		&header.CustName,
		&header.CustPhone,
		&header.CustEmail,
		&header.CustAddress,
		&header.PaymentMethod,
		&header.DeliveryMethod,
		&header.Status,
		&header.Total,
		&header.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	queryItems := `
	SELECT 
		ti.id,
		p.title AS name,
		COALESCE((SELECT pi.image FROM product_images pi WHERE pi.product_id = p.id LIMIT 1), '') AS image,
		sz.name AS size,
		p.base_price,
		0 AS discount_price,
		v.name AS variant,
		ti.quantity,
		ti.subtotal
	FROM transaction_items ti
	LEFT JOIN products p ON p.id = ti.product_id
	LEFT JOIN sizes sz ON sz.id = ti.size_id
	LEFT JOIN variants v ON v.id = ti.variant_id
	WHERE ti.transaction_id = $1
	`

	rows, err := db.Query(ctx, queryItems, transactionID)
	if err != nil {
		return nil, err
	}
	shippingPrice := 0.0
	if header.DeliveryMethod == "DoorDelivery" {
		shippingPrice = 10000
	}

	header.ShippingPrice = shippingPrice
	header.Total += shippingPrice

	defer rows.Close()

	for rows.Next() {
		var item TransactionItemDetail
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Image,
			&item.Size,
			&item.BasePrice,
			&item.DiscountPrice,
			&item.Variant,
			&item.Quantity,
			&item.Subtotal,
		)
		if err != nil {
			return nil, err
		}
		header.Items = append(header.Items, item)
	}

	return &header, nil

}
