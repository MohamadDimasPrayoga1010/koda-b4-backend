package models

import (
	"context"
	"errors"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Product struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	BasePrice   float64        `json:"base_price"`
	Stock       int            `json:"stock"`
	CategoryID  int64          `json:"category_id"`
	VariantID   int64          `json:"variant_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
	Images      []ProductImage `json:"images,omitempty"`
	Sizes       []Size         `json:"sizes,omitempty"`
}

type ProductResponse struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	BasePrice   float64        `json:"base_price"`
	Stock       int            `json:"stock"`
	CategoryID  int64          `json:"category_id"`
	VariantID   int64          `json:"variant_id"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	Images      []ProductImage `json:"images,omitempty"`
	Sizes       []Size         `json:"sizes,omitempty"`
}

type ProductImage struct {
	ProductID int64      `json:"product_id"`
	Image     string     `json:"image"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Size struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name" form:"name"`
	AdditionalPrice float64 `json:"additional_price" form:"additional_price"`
}

type ProductRequest struct {
	Title       string                  `form:"title" binding:"required"`
	Description string                  `form:"description"`
	BasePrice   float64                 `form:"base_price" binding:"required"`
	Stock       int                     `form:"stock" binding:"required"`
	CategoryID  int64                   `form:"category_id"`
	VariantID   int64                   `form:"variant_id"`
	Sizes       []int64                 `form:"sizes"`
	Images      []*multipart.FileHeader `form:"images"`
}

type ProductFilter struct {
	Categories []int64  `json:"categories" form:"categories"`
	IsFavorite *bool    `json:"is_favorite" form:"is_favorite"`
	SortBy     string   `json:"sortby" form:"sortby"`
	PriceMin   *float64 `json:"price_min" form:"price_min"`
	PriceMax   *float64 `json:"price_max" from:"price_max"`
}

type ProductDetail struct {
	ID          int64                    `json:"id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	BasePrice   float64                  `json:"base_price"`
	Stock       int                      `json:"stock"`
	CategoryID  int64                    `json:"category_id"`
	Variant     *Variant                 `json:"variant,omitempty"`
	Sizes       []Size                   `json:"sizes"`
	Images      []ProductImage           `json:"images"`
	Recommended []RecommendedProductInfo `json:"recommended"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

type RecommendedProductInfo struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	BasePrice   float64        `json:"base_price"`
	Stock       int            `json:"stock"`
	CategoryID  int64          `json:"category_id"`
	Variant     *Variant       `json:"variant,omitempty"`
	Sizes       []Size         `json:"sizes"`
	Images      []ProductImage `json:"images"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type Variant struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type RecommendedProduct struct {
	ProductID     int64 `json:"product_id"`
	RecommendedID int64 `json:"recommended_id"`
}

type Cart struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	SizeID    *int64    `json:"size_id,omitempty"`
	VariantID *int64    `json:"variant_id,omitempty"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItemResponse struct {
	ProductID int64   `json:"product_id"`
	Title     string  `json:"title"`
	BasePrice float64 `json:"base_price"`
	Image     string  `json:"image"`
	Size      string  `json:"size,omitempty"`
	Variant   string  `json:"variant,omitempty"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

func AddOrUpdateCart(db *pgxpool.Pool, userID, productID int64, sizeID, variantID *int64, quantity int) (CartItemResponse, error) {
	ctx := context.Background()
	var cartID int64

	// Cek stok produk
	var stock int
	err := db.QueryRow(ctx, `
		SELECT stock FROM products 
		WHERE id=$1
	`, productID).Scan(&stock)
	if err != nil {
		return CartItemResponse{}, err
	}

	// Cek apakah sudah ada di cart
	var existingQty int
	checkQuery := `
		SELECT id, quantity FROM carts 
		WHERE user_id=$1 AND product_id=$2 AND 
		      COALESCE(size_id, 0) = COALESCE($3, 0) AND 
		      COALESCE(variant_id, 0) = COALESCE($4, 0)
		LIMIT 1
	`
	err = db.QueryRow(ctx, checkQuery, userID, productID, sizeID, variantID).Scan(&cartID, &existingQty)
	if err == nil {
		newQty := existingQty + quantity
		if newQty > stock {
			return CartItemResponse{}, errors.New("stock not enough for product")
		}

		updateQuery := `UPDATE carts SET quantity=$1, updated_at=NOW() WHERE id=$2`
		_, err := db.Exec(ctx, updateQuery, newQty, cartID)
		if err != nil {
			return CartItemResponse{}, err
		}
	} else {
		if quantity > stock {
			return CartItemResponse{}, errors.New("stock not enough for product")
		}

		insertQuery := `
			INSERT INTO carts (user_id, product_id, size_id, variant_id, quantity, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			RETURNING id
		`
		err := db.QueryRow(ctx, insertQuery, userID, productID, sizeID, variantID, quantity).Scan(&cartID)
		if err != nil {
			return CartItemResponse{}, err
		}
	}

	// Ambil detail cart item termasuk subtotal
	detailQuery := `
		SELECT 
			c.product_id,
			p.title,
			p.base_price,
			COALESCE(pi.image,'') AS image,
			COALESCE(s.name,'') AS size_name,
			COALESCE(v.name,'') AS variant_name,
			c.quantity,
			(p.base_price + COALESCE(s.additional_price,0)) * c.quantity AS subtotal
		FROM carts c
		JOIN products p ON p.id = c.product_id
		LEFT JOIN product_images pi ON pi.product_id = p.id
		LEFT JOIN sizes s ON s.id = c.size_id
		LEFT JOIN variants v ON v.id = c.variant_id
		WHERE c.id=$1
	`

	var item CartItemResponse
	err = db.QueryRow(ctx, detailQuery, cartID).Scan(
		&item.ProductID, &item.Title, &item.BasePrice,
		&item.Image, &item.Size, &item.Variant, &item.Quantity,
		&item.Subtotal,
	)
	if err != nil {
		return CartItemResponse{}, err
	}

	return item, nil
}


func GetCartByUser(db *pgxpool.Pool, userID int64) ([]CartItemResponse, error) {
	query := `
		SELECT 
			c.product_id,
			p.title,
			p.base_price,
			COALESCE(pi.image, '') AS image,
			COALESCE(s.name, '') AS size,
			COALESCE(v.name, '') AS variant,
			c.quantity,
			( (p.base_price + COALESCE(s.additional_price, 0)) * c.quantity ) AS subtotal
		FROM carts c
		JOIN products p ON p.id = c.product_id
		LEFT JOIN LATERAL (
			SELECT image 
			FROM product_images 
			WHERE product_id = p.id 
			ORDER BY id ASC 
			LIMIT 1
		) pi ON true
		LEFT JOIN sizes s ON s.id = c.size_id
		LEFT JOIN variants v ON v.id = c.variant_id
		WHERE c.user_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItemResponse
	for rows.Next() {
		var item CartItemResponse
		if err := rows.Scan(
			&item.ProductID,
			&item.Title,
			&item.BasePrice,
			&item.Image,
			&item.Size,
			&item.Variant,
			&item.Quantity,
			&item.Subtotal,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}




type OrderTransactionRequest struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email,omitempty"`
	Phone           string `json:"phone,omitempty"`
	Address         string `json:"address,omitempty"`
	PaymentMethodID int64  `json:"payment_method_id"`
	ShippingID      int64  `json:"shipping_id"`
	UserID          int64  `json:"user_id"`
}

type OrderTransactionItem struct {
    ID          int64   `json:"id"`
    ProductID   int64   `json:"product_id"`
    ProductName string  `json:"product_name"`
    VariantName *string `json:"variant_name,omitempty"`
    SizeName    *string `json:"size_name,omitempty"`
    Quantity    int     `json:"quantity"`
    Subtotal    float64 `json:"subtotal"`
}

type OrderTransaction struct {
	ID                int64                  `json:"id"`
	UserID            int64                  `json:"user_id"`
	Fullname          string                 `json:"fullname"`
	Email             string                 `json:"email"`
	Phone             string                 `json:"phone"`
	Address           string                 `json:"address"`
	PaymentMethodName string                 `json:"payment_method_name"`
	ShippingName      string                 `json:"shipping_name"`
	InvoiceNumber     string                 `json:"invoice_number"`
	Total             float64                `json:"total"`
	Status            string                 `json:"status"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	Items             []OrderTransactionItem `json:"items,omitempty"`
}

func CreateOrderTransaction(db *pgxpool.Pool, req OrderTransactionRequest) (*OrderTransaction, error) {
	ctx := context.Background()


	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var userEmail string
	err = tx.QueryRow(ctx, `SELECT email FROM users WHERE id=$1`, req.UserID).Scan(&userEmail)
	if err != nil {
		return nil, errors.New("failed to fetch user email")
	}
	if req.Email == "" {
		req.Email = userEmail
	}

	var profilePhone, profileAddress *string
	err = tx.QueryRow(ctx, `SELECT phone, address FROM profile WHERE user_id=$1`, req.UserID).Scan(&profilePhone, &profileAddress)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("failed to fetch user profile")
	}

	if req.Phone == "" {
		if profilePhone != nil {
			req.Phone = *profilePhone
		} else {
			return nil, errors.New("phone must be provided")
		}
	}
	if req.Address == "" {
		if profileAddress != nil {
			req.Address = *profileAddress
		} else {
			return nil, errors.New("address must be provided")
		}
	}

	queryCart := `
		SELECT 
			c.product_id, 
			p.title AS product_name,
			s.name AS size_name,
			v.name AS variant_name,
			c.quantity,
			(p.base_price + COALESCE(s.additional_price,0)) * c.quantity AS subtotal
		FROM carts c
		JOIN products p ON p.id = c.product_id
		LEFT JOIN sizes s ON s.id = c.size_id
		LEFT JOIN variants v ON v.id = c.variant_id
		WHERE c.user_id=$1
	`
	rows, err := tx.Query(ctx, queryCart, req.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderTransactionItem
	var total float64
	for rows.Next() {
		var item OrderTransactionItem
		var sizeName, variantName *string
		if err := rows.Scan(&item.ProductID, &item.ProductName, &sizeName, &variantName, &item.Quantity, &item.Subtotal); err != nil {
			return nil, err
		}
		item.SizeName = sizeName
		item.VariantName = variantName
		items = append(items, item)
		total += item.Subtotal
	}

	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	for i := range items {
		item := &items[i]

		var currentStock int
		err := tx.QueryRow(ctx, `SELECT COALESCE(stock,0) FROM products WHERE id=$1 FOR UPDATE`, item.ProductID).Scan(&currentStock)
		if err != nil {
			return nil, errors.New("failed to fetch stock for product " + item.ProductName + ": " + err.Error())
		}
		if item.Quantity > currentStock {
			return nil, errors.New("product " + item.ProductName + " stock insufficient: available " + strconv.Itoa(currentStock) + ", requested " + strconv.Itoa(item.Quantity))
		}

		_, err = tx.Exec(ctx, `UPDATE products SET stock = stock - $1 WHERE id=$2`, item.Quantity, item.ProductID)
		if err != nil {
			return nil, errors.New("failed to update stock for product " + item.ProductName + ": " + err.Error())
		}
	}

	var paymentName, shippingName string
	err = tx.QueryRow(ctx, `SELECT name FROM payment_methods WHERE id=$1`, req.PaymentMethodID).Scan(&paymentName)
	if err != nil {
		return nil, errors.New("invalid payment method")
	}
	err = tx.QueryRow(ctx, `SELECT name FROM shippings WHERE id=$1`, req.ShippingID).Scan(&shippingName)
	if err != nil {
		return nil, errors.New("invalid shipping method")
	}

	invoice := "INV-" + time.Now().Format("20060102150405") + "-" + strconv.FormatInt(req.UserID, 10)
	var orderID int64
	var createdAt, updatedAt time.Time
	insertOrder := `
		INSERT INTO transactions
			(user_id, fullname, email, phone, address, payment_method_id, shipping_id, invoice_number, total, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,'Pending',NOW(),NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, insertOrder,
		req.UserID, req.Fullname, req.Email, req.Phone, req.Address,
		req.PaymentMethodID, req.ShippingID, invoice, total,
	).Scan(&orderID, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	for i := range items {
		item := &items[i]

		insertItem := `
			INSERT INTO transaction_items
				(transaction_id, product_id, variant_id, size_id, quantity, subtotal)
			VALUES ($1,$2,$3,$4,$5,$6)
			RETURNING id
		`
		err = tx.QueryRow(ctx, insertItem,
			orderID, item.ProductID, nil, nil, item.Quantity, item.Subtotal,
		).Scan(&item.ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	order := &OrderTransaction{
		ID:                orderID,
		UserID:            req.UserID,
		Fullname:          req.Fullname,
		Email:             req.Email,
		Phone:             req.Phone,
		Address:           req.Address,
		PaymentMethodName: paymentName,
		ShippingName:      shippingName,
		InvoiceNumber:     invoice,
		Total:             total,
		Status:            "Pending",
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
		Items:             items,
	}

	return order, nil
}



