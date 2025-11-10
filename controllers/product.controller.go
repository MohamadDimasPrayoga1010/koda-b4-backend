package controllers

import (
	"context"
	"fmt"
	"main/models"
	"strings"

	"os"
	"path/filepath"
	"strconv"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductController struct {
	DB *pgxpool.Pool
}

// @Summary Create a new product
// @Description Membuat produk baru beserta gambar dan ukuran
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Product Title"
// @Param description formData string false "Description"
// @Param base_price formData number true "Base Price"
// @Param stock formData int true "Stock"
// @Param category_id formData int false "Category ID"
// @Param variant_id formData int false "Variant ID"
// @Param sizes formData []int false "Size IDs (multiple allowed)"
// @Param images formData file false "Upload product image (repeat for multiple)"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /admin/products [post]
func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	var req models.ProductRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return
	}

	query := `
		INSERT INTO products (title, description, base_price, stock, category_id, variant_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, title, description, base_price, stock, category_id, variant_id, created_at
	`
	var product models.ProductResponse
	err := pc.DB.QueryRow(context.Background(), query,
		req.Title, req.Description, req.BasePrice, req.Stock, req.CategoryID, req.VariantID,
	).Scan(&product.ID, &product.Title, &product.Description, &product.BasePrice,
		&product.Stock, &product.CategoryID, &product.VariantID, &product.CreatedAt)

	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to create product",
			Data:    err.Error(),
		})
		return
	}

	uploadDir := "./uploads/products"
	os.MkdirAll(uploadDir, os.ModePerm)

	allowedExts := []string{".jpg", ".jpeg", ".png"}
	maxSize := int64(2 * 1024 * 1024)

	for _, file := range req.Images {
		if file.Size > maxSize {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "File too large: " + file.Filename,
			})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		valid := false
		for _, e := range allowedExts {
			if ext == e {
				valid = true
				break
			}
		}
		if !valid {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Invalid file type: " + file.Filename,
			})
			return
		}

		name := strings.TrimSuffix(file.Filename, ext)
		name = strings.ReplaceAll(name, " ", "_")
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + name + ext
		fullPath := filepath.Join(uploadDir, filename)

		if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Failed to save file: " + file.Filename,
				Data:    err.Error(),
			})
			return
		}

		_, err := pc.DB.Exec(context.Background(),
			`INSERT INTO product_images (product_id, image) VALUES ($1, $2)`,
			product.ID, filename,
		)
		if err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Failed to save image record: " + file.Filename,
				Data:    err.Error(),
			})
			return
		}

		product.Images = append(product.Images, models.ProductImage{
			ProductID: product.ID,
			Image:     filename, 
			UpdatedAt: time.Now(),
		})
	}


	for _, sizeID := range req.Sizes {
		var size models.Size

		err := pc.DB.QueryRow(context.Background(),
			`SELECT id, name, additional_price FROM sizes WHERE id = $1`, sizeID,
		).Scan(&size.ID, &size.Name, &size.AdditionalPrice)

		if err != nil {
			continue
		}

		_, err = pc.DB.Exec(context.Background(),
			`INSERT INTO product_sizes (product_id, size_id) VALUES ($1, $2)`,
			product.ID, sizeID,
		)
		if err == nil {
			product.Sizes = append(product.Sizes, size)
		}
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}




// GetProduct godoc
// @Summary Get list of products
// @Description Mengambil daftar products dengan pagination dan optional search
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param search query string false "Search by title or description"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/products [get]
func (pc *ProductController) GetProduct(ctx *gin.Context) {
	search := ctx.Query("search")
	limitStr := ctx.DefaultQuery("limit", "10")
	pageStr := ctx.DefaultQuery("page", "1")

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	offset := (page - 1) * limit

	query := `
		SELECT id, title, description, base_price, stock, category_id, variant_id, created_at
		FROM products
	`

	var limitation []interface{}
	if search != "" {
		query += " AND (LOWER(title) LIKE LOWER($1) OR LOWER(description) LIKE LOWER($1))"
		limitation = append(limitation, "%"+search+"%")
	}

	if search != "" {
		query += " ORDER BY id DESC LIMIT $2 OFFSET $3"
		limitation = append(limitation, limit, offset)
	} else {
		query += " ORDER BY id DESC LIMIT $1 OFFSET $2"
		limitation = append(limitation, limit, offset)
	}

	rows, err := pc.DB.Query(context.Background(), query, limitation...)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch product",
			Data:    err.Error(),
		})
		return
	}
	defer rows.Close()

	var products []models.ProductResponse

	for rows.Next() {
		var p models.ProductResponse
		err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.BasePrice, &p.Stock, &p.CategoryID, &p.VariantID, &p.CreatedAt)
		if err != nil {
			fmt.Println("Scan error:", err)
			continue
		}

		imgRows, _ := pc.DB.Query(context.Background(),
			`SELECT image, updated_at, deleted_at FROM product_images WHERE product_id=$1 AND deleted_at`, p.ID)
		for imgRows.Next() {
			var img models.ProductImage
			img.ProductID = p.ID
			imgRows.Scan(&img.Image, &img.UpdatedAt, &img.DeletedAt)
			p.Images = append(p.Images, img)
		}
		imgRows.Close()

		sizeRows, _ := pc.DB.Query(context.Background(), `
			SELECT s.id, s.name, s.additional_price 
			FROM sizes s
			JOIN product_sizes ps ON ps.size_id = s.id
			WHERE ps.product_id = $1
		`, p.ID)
		for sizeRows.Next() {
			var s models.Size
			sizeRows.Scan(&s.ID, &s.Name, &s.AdditionalPrice)
			p.Sizes = append(p.Sizes, s)
		}
		sizeRows.Close()

		products = append(products, p)
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Products fetched successfully",
		Data: gin.H{
			"page":     page,
			"limit":    limit,
			"products": products,
		},
	})
}

// GetProductByID godoc
// @Summary Get a product by ID
// @Description Mengambil product berdasarkan ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /admin/products/{id} [get]
func (pc *ProductController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
		SELECT id, title, description, base_price, stock, category_id, variant_id, created_at
		FROM products
		WHERE id = $1 
	`

	var p models.ProductResponse
	err := pc.DB.QueryRow(context.Background(), query, id).
		Scan(&p.ID, &p.Title, &p.Description, &p.BasePrice, &p.Stock, &p.CategoryID, &p.VariantID, &p.CreatedAt)

	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Product not found",
			Data:    nil,
		})
		return
	}

	imgRows, _ := pc.DB.Query(context.Background(),
		`SELECT image, updated_at, deleted_at FROM product_images WHERE product_id=$1`, p.ID)
	for imgRows.Next() {
		var img models.ProductImage
		img.ProductID = p.ID
		imgRows.Scan(&img.Image, &img.UpdatedAt, &img.DeletedAt)
		p.Images = append(p.Images, img)
	}
	imgRows.Close()

	sizeRows, _ := pc.DB.Query(context.Background(), `
		SELECT s.id, s.name, s.additional_price
		FROM sizes s
		JOIN product_sizes ps ON ps.size_id = s.id
		WHERE ps.product_id = $1
	`, p.ID)
	for sizeRows.Next() {
		var s models.Size
		sizeRows.Scan(&s.ID, &s.Name, &s.AdditionalPrice)
		p.Sizes = append(p.Sizes, s)
	}
	sizeRows.Close()

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product fetched successfully",
		Data:    p,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Mengupdate product beserta images dan sizes berdasarkan ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param body body models.ProductRequest true "Product request body"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/products/{id} [patch]
func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("id")

	var req models.ProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return 
	}

	query := `
		UPDATE products
		SET title = $1, description = $2, base_price = $3, stock = $4,
			category_id = $5, variant_id = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, title, description, base_price, stock, category_id, variant_id, created_at
	`
	var product models.ProductResponse
	err := pc.DB.QueryRow(context.Background(), query,
		req.Title, req.Description, req.BasePrice, req.Stock,
		req.CategoryID, req.VariantID, productID,
	).Scan(&product.ID, &product.Title, &product.Description, &product.BasePrice,
		&product.Stock, &product.CategoryID, &product.VariantID, &product.CreatedAt)

	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to update product",
			Data:    err.Error(),
		})
		return
	}

	pc.DB.Exec(context.Background(), `DELETE FROM product_images WHERE product_id = $1`, product.ID)
	for _, img := range req.Images {
		_, err := pc.DB.Exec(context.Background(),
			`INSERT INTO product_images (product_id, image) VALUES ($1, $2)`,
			product.ID, img,
		)
		if err == nil {
			product.Images = append(product.Images, models.ProductImage{ProductID: product.ID, Image: img.Filename})
		}
	}

	pc.DB.Exec(context.Background(), `DELETE FROM product_sizes WHERE product_id = $1`, product.ID)
	for _, sizeID := range req.Sizes {
		_, err := pc.DB.Exec(context.Background(),
			`INSERT INTO product_sizes (product_id, size_id) VALUES ($1, $2)`,
			product.ID, sizeID,
		)
		if err == nil {
			product.Sizes = append(product.Sizes, models.Size{ID: sizeID})
		}
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product updated successfully",
		Data:    product,
	})
}


// DeleteProduct godoc
// @Summary Delete a product
// @Description Menghapus product berdasarkan ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/products/{id} [delete]
func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `DELETE FROM products WHERE id = $1`

	result, err := pc.DB.Exec(context.Background(), query, id)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete product",
			Data:    err.Error(),
		})
		return
	}

	if result.RowsAffected() == 0 {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Product not found",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product deleted successfully",
	})
}


