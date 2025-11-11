package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"main/libs"
	"main/models"
	"net/http"
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
// @Param sizes formData string false "List of size IDs (example: 1,2,3)"
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

	if len(req.Sizes) == 0 {
		sizesStr := ctx.PostForm("sizes")
		if sizesStr != "" {
			for _, s := range strings.Split(sizesStr, ",") {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				if id, err := strconv.ParseInt(s, 10, 64); err == nil {
					req.Sizes = append(req.Sizes, id)
				}
			}
		}
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
		err := pc.DB.QueryRow(
			context.Background(),
			`SELECT id, name, additional_price FROM sizes WHERE id = $1`,
			sizeID,
		).Scan(&size.ID, &size.Name, &size.AdditionalPrice)

		if err != nil {
			continue
		}

		_, err = pc.DB.Exec(
			context.Background(),
			`INSERT INTO product_sizes (product_id, size_id) VALUES ($1, $2)`,
			product.ID, sizeID,
		)
		if err != nil {
			continue
		}

		product.Sizes = append(product.Sizes, size)
	}

	redis := libs.RedisClient.Scan(libs.Ctx, 0, "products:*", 0).Iterator()
	for redis.Next(libs.Ctx) {
		key := redis.Val()
		libs.RedisClient.Del(libs.Ctx, key)
	}

	if err := redis.Err(); err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to clear Redis cache",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// GetProduct godoc
// @Summary Get list of products
// @Description Mengambil daftar products dengan pagination, optional search, dan sorting
// @Tags Products
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param search query string false "Search by title or description"
// @Param sort_by query string false "Sort by field (id, title, base_price, created_at)" default(created_at)
// @Param order query string false "Sort order (asc/desc)" default(desc)
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/products [get]
func (pc *ProductController) GetProduct(ctx *gin.Context) {
	search := ctx.Query("search")
	limitStr := ctx.DefaultQuery("limit", "10")
	pageStr := ctx.DefaultQuery("page", "1")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	order := strings.ToUpper(ctx.DefaultQuery("order", "DESC"))

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	offset := (page - 1) * limit

	allowedSortFields := map[string]bool{
		"id":         true,
		"title":      true,
		"base_price": true,
		"created_at": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	cacheKey := fmt.Sprintf("products:page:%d:",
		page)

	cache, err := libs.RedisClient.Get(libs.Ctx, cacheKey).Result()
	if err == nil {
		var cachedData []models.ProductResponse
		json.Unmarshal([]byte(cache), &cachedData)

		ctx.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Product fetch from cache",
			Data: gin.H{
				"page": page,
				// "limit":    limit,
				// "sort_by":  sortBy,
				// "order":    order,
				"products": cachedData,
			},
		})
		return
	}

	query := `
		SELECT id, title, description, base_price, stock, category_id, variant_id, created_at
		FROM products
	`
	var args []interface{}
	argIndex := 1

	if search != "" {
		query += fmt.Sprintf(" WHERE LOWER(title) LIKE LOWER($%d) OR LOWER(description) LIKE LOWER($%d)", argIndex, argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortBy, order, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := pc.DB.Query(context.Background(), query, args...)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch products",
			Data:    err.Error(),
		})
		return
	}
	defer rows.Close()

	var products []models.ProductResponse

	for rows.Next() {
		var p models.ProductResponse
		err := rows.Scan(
			&p.ID, &p.Title, &p.Description,
			&p.BasePrice, &p.Stock, &p.CategoryID, &p.VariantID, &p.CreatedAt,
		)
		if err != nil {
			ctx.JSON(404, models.Response{
				Success: false,
				Message: "Product not found",
				Data:    nil,
			})
			return
		}

		imgRows, _ := pc.DB.Query(context.Background(),
			`SELECT image, updated_at, deleted_at 
			 FROM product_images 
			 WHERE product_id = $1 AND deleted_at IS NULL`, p.ID)
		for imgRows.Next() {
			var img models.ProductImage
			img.ProductID = p.ID
			imgRows.Scan(&img.Image, &img.UpdatedAt, &img.DeletedAt)
			p.Images = append(p.Images, img)
		}
		imgRows.Close()

		sizeRows, _ := pc.DB.Query(context.Background(),
			`SELECT s.id, s.name, s.additional_price 
			 FROM sizes s
			 JOIN product_sizes ps ON ps.size_id = s.id
			 WHERE ps.product_id = $1`,
			p.ID)
		for sizeRows.Next() {
			var s models.Size
			sizeRows.Scan(&s.ID, &s.Name, &s.AdditionalPrice)
			p.Sizes = append(p.Sizes, s)
		}
		sizeRows.Close()

		products = append(products, p)
	}

	if len(products) == 0 {
		ctx.JSON(200, models.Response{
			Success: true,
			Message: "No products found",
			Data: gin.H{
				"products": []models.ProductResponse{},
				"page":     page,
				"limit":    limit,
			},
		})
		return
	}

	jsonData, _ := json.Marshal(products)
	err = libs.RedisClient.Set(libs.Ctx, cacheKey, jsonData, 10*time.Minute).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Gagal menyimpan cache Redis",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Products fetched successfully",
		Data: gin.H{
			"page":     page,
			"limit":    limit,
			"sort_by":  sortBy,
			"order":    order,
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
// @Failure 500 {object} models.Response
// @Router /admin/products/{id} [get]
func (pc *ProductController) GetProductByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product ID",
			Data:    nil,
		})
		return
	}

	query := `
		SELECT id, title, description, base_price, stock, category_id, variant_id, created_at
		FROM products
		WHERE id = $1
	`

	var p models.ProductResponse
	err = pc.DB.QueryRow(context.Background(), query, productID).
		Scan(&p.ID, &p.Title, &p.Description, &p.BasePrice, &p.Stock,
			&p.CategoryID, &p.VariantID, &p.CreatedAt)

	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Product not found",
			Data:    nil,
		})
		return
	}

	imgRows, err := pc.DB.Query(context.Background(),
		`SELECT image, updated_at, deleted_at FROM product_images WHERE product_id=$1 AND deleted_at IS NULL`, p.ID)
	if err == nil {
		for imgRows.Next() {
			var img models.ProductImage
			img.ProductID = p.ID
			if err := imgRows.Scan(&img.Image, &img.UpdatedAt, &img.DeletedAt); err == nil {
				p.Images = append(p.Images, img)
			}
		}
		imgRows.Close()
	}

	sizeRows, err := pc.DB.Query(context.Background(),
		`SELECT s.id, s.name, s.additional_price
		 FROM sizes s
		 JOIN product_sizes ps ON ps.size_id = s.id
		 WHERE ps.product_id = $1`,
		p.ID)
	if err == nil {
		for sizeRows.Next() {
			var s models.Size
			if err := sizeRows.Scan(&s.ID, &s.Name, &s.AdditionalPrice); err == nil {
				p.Sizes = append(p.Sizes, s)
			}
		}
		sizeRows.Close()
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product fetched successfully",
		Data:    p,
	})
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product data by ID, including multiple sizes and images
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param title formData string true "Product title"
// @Param description formData string false "Product description"
// @Param base_price formData number true "Base price of the product"
// @Param stock formData int true "Stock quantity"
// @Param category_id formData int false "Category ID"
// @Param variant_id formData int false "Variant ID"
// @Param sizes formData int false "List of size IDs (send multiple 'sizes' params for multiple sizes)"
// @Param images formData file false "Product images (send multiple 'images' files for multiple uploads)"
// @Success 200 {object} models.Response "Product updated successfully"
// @Failure 400 {object} models.Response "Invalid request body"
// @Failure 404 {object} models.Response "Product not found"
// @Router /admin/products/{id} [patch]
func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("id")

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
		UPDATE products
		SET title = $1, description = $2, base_price = $3, stock = $4,
			category_id = $5, variant_id = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING id, title, description, base_price, stock, category_id, variant_id, created_at, updated_at
	`
	var product models.ProductResponse
	err := pc.DB.QueryRow(context.Background(), query,
		req.Title, req.Description, req.BasePrice, req.Stock,
		req.CategoryID, req.VariantID, productID,
	).Scan(
		&product.ID, &product.Title, &product.Description,
		&product.BasePrice, &product.Stock,
		&product.CategoryID, &product.VariantID,
		&product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to update product",
			Data:    err.Error(),
		})
		return
	}

	uploadDir := "./uploads/products"
	os.MkdirAll(uploadDir, os.ModePerm)

	pc.DB.Exec(context.Background(), `DELETE FROM product_images WHERE product_id = $1`, product.ID)

	for _, file := range req.Images {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.TrimSuffix(file.Filename, ext), ext)
		fullPath := filepath.Join(uploadDir, filename)

		if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
			continue
		}

		pc.DB.Exec(context.Background(),
			`INSERT INTO product_images (product_id, image) VALUES ($1, $2)`,
			product.ID, filename,
		)

		product.Images = append(product.Images, models.ProductImage{
			ProductID: product.ID,
			Image:     filename,
			UpdatedAt: time.Now(),
		})
	}

	pc.DB.Exec(context.Background(), `DELETE FROM product_sizes WHERE product_id = $1`, product.ID)
	for _, sizeID := range req.Sizes {
		pc.DB.Exec(context.Background(),
			`INSERT INTO product_sizes (product_id, size_id) VALUES ($1, $2)`,
			product.ID, sizeID,
		)
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

// GetProductImages godoc
// @Summary Get all images of a product
// @Description Mengambil semua gambar dari product berdasarkan product ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /product/{id}/images [get]
func (pc *ProductController) GetProductImages(ctx *gin.Context) {
	productIDParam := ctx.Param("id")
	productID, err := strconv.ParseInt(productIDParam, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product ID",
		})
		return
	}

	rows, err := pc.DB.Query(context.Background(),
		`SELECT row_number() OVER (), image, updated_at, deleted_at
		 FROM product_images 
		 WHERE product_id=$1 AND deleted_at IS NULL`, productID)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch product images",
			Data:    err.Error(),
		})
		return
	}
	defer rows.Close()

	var images []models.ProductImage
	for rows.Next() {
		var img models.ProductImage
		var idx int
		if err := rows.Scan(&idx, &img.Image, &img.UpdatedAt, &img.DeletedAt); err != nil {
			continue
		}
		img.ProductID = productID
		images = append(images, img)
	}

	if len(images) == 0 {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "No images found for this product",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product images fetched successfully",
		Data:    images,
	})
}

// GetProductImageByID godoc
// @Summary Get a single image of a product
// @Description Mengambil satu gambar dari product berdasarkan product ID dan image ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /product/{id}/images/{image_id} [get]
func (pc *ProductController) GetProductImageByID(ctx *gin.Context) {
	productIDParam := ctx.Param("id")
	imageIDParam := ctx.Param("image_id")

	productID, err1 := strconv.ParseInt(productIDParam, 10, 64)
	imageID, err2 := strconv.ParseInt(imageIDParam, 10, 64)
	if err1 != nil || err2 != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product ID or image ID",
		})
		return
	}

	var img models.ProductImage
	err := pc.DB.QueryRow(context.Background(),
		`SELECT image, updated_at, deleted_at
		 FROM product_images 
		 WHERE product_id=$1 AND id=$2 AND deleted_at IS NULL`,
		productID, imageID).Scan(&img.Image, &img.UpdatedAt, &img.DeletedAt)

	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Image not found",
		})
		return
	}

	img.ProductID = productID

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product image fetched successfully",
		Data:    img,
	})
}

// DeleteProductImage godoc
// @Summary Delete a product image
// @Description Menghapus gambar product berdasarkan product ID dan image ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /product/{id}/images/{image_id} [delete]
func (pc *ProductController) DeleteProductImage(ctx *gin.Context) {
	productIDParam := ctx.Param("id")
	imageIDParam := ctx.Param("image_id")

	productID, err1 := strconv.ParseInt(productIDParam, 10, 64)
	imageID, err2 := strconv.ParseInt(imageIDParam, 10, 64)
	if err1 != nil || err2 != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product ID or image ID",
		})
		return
	}

	query := `UPDATE product_images SET deleted_at=NOW() WHERE product_id=$1 AND id=$2 AND deleted_at IS NULL`
	result, err := pc.DB.Exec(context.Background(), query, productID, imageID)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete product image",
			Data:    err.Error(),
		})
		return
	}

	if result.RowsAffected() == 0 {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Image not found or already deleted",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product image deleted successfully",
	})
}

// UpdateProductImage godoc
// @Summary Update a product image
// @Description Mengupdate / mengganti file gambar product berdasarkan product ID dan image ID
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param image_id path int true "Image ID"
// @Param image formData file true "New product image file (jpg/png, max 2MB)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /product/{id}/images/{image_id} [patch]
func (pc *ProductController) UpdateProductImage(ctx *gin.Context) {
	productIDParam := ctx.Param("id")
	imageIDParam := ctx.Param("image_id")

	productID, err1 := strconv.ParseInt(productIDParam, 10, 64)
	imageID, err2 := strconv.ParseInt(imageIDParam, 10, 64)
	if err1 != nil || err2 != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product ID or image ID",
		})
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Image file is required",
		})
		return
	}

	const maxSize = 2 * 1024 * 1024
	if file.Size > maxSize {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "File size exceeds 2MB",
		})
		return
	}

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid file type",
		})
		return
	}

	uploadDir := "./uploads/products"
	os.MkdirAll(uploadDir, os.ModePerm)
	newFilename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.TrimSuffix(file.Filename, ext), ext)
	fullPath := filepath.Join(uploadDir, newFilename)

	if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to save image file",
			Data:    err.Error(),
		})
		return
	}

	query := `UPDATE product_images SET image=$1, updated_at=NOW() WHERE product_id=$2 AND id=$3 RETURNING image, updated_at`
	var updatedImage string
	var updatedAt time.Time
	err = pc.DB.QueryRow(context.Background(), query, newFilename, productID, imageID).Scan(&updatedImage, &updatedAt)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Image not found",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product image updated successfully",
		Data: models.ProductImage{
			ProductID: productID,
			Image:     updatedImage,
			UpdatedAt: updatedAt,
		},
	})
}
