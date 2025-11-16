package controllers

import (
	"coffeeder-backend/libs"
	"coffeeder-backend/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductController struct {
	DB *pgxpool.Pool
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with multiple variants, sizes, and images
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Product Title"
// @Param description formData string false "Product Description"
// @Param base_price formData number true "Base Price"
// @Param stock formData int true "Stock"
// @Param category_id formData int true "Category ID"
// @Param variant_ids formData string false "Comma-separated Variant IDs (example: 1,2,3)"
// @Param sizes formData string false "Comma-separated Size IDs (example: 1,3)"
// @Param images formData file false "Upload product image (repeat for multiple)"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/products [post]
func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	var req models.ProductRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid form-data",
			Data:    err.Error(),
		})
		return
	}

	if len(req.VariantID) == 0 {
		raw := ctx.PostForm("variant_id")
		if raw != "" {
			for _, v := range strings.Split(raw, ",") {
				v = strings.TrimSpace(v)
				if v == "" {
					continue
				}
				if id, err := strconv.ParseInt(v, 10, 64); err == nil {
					req.VariantID = append(req.VariantID, id)
				}
			}
		}
	}

	if len(req.Sizes) == 0 {
		raw := ctx.PostForm("sizes")
		if raw != "" {
			for _, s := range strings.Split(raw, ",") {
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

	contains := func(arr []string, s string) bool {
		for _, v := range arr {
			if v == s {
				return true
			}
		}
		return false
	}

	maxSize := int64(2 * 1024 * 1024)
	allowedExts := []string{".jpg", ".jpeg", ".png"}

	uploadDir := "./uploads/products"
	os.MkdirAll(uploadDir, os.ModePerm)

	savedFiles := []string{}

	files := req.Images
	for _, file := range files {

		if file.Size > maxSize {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "File too large(max 2mb): " + file.Filename,
			})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !contains(allowedExts, ext) {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Invalid file type(only png, jpg, jpeg): " + file.Filename,
			})
			return
		}

		name := strings.TrimSuffix(file.Filename, ext)
		name = strings.ReplaceAll(name, " ", "_")

		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), name, ext)
		fullPath := filepath.Join(uploadDir, filename)

		if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Failed to save file",
				Data:    err.Error(),
			})
			return
		}

		savedFiles = append(savedFiles, filename)
	}

	product, err := models.CreateProduct(pc.DB, req, savedFiles)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to create product",
			Data:    err.Error(),
		})
		return
	}

	iter := libs.RedisClient.Scan(libs.Ctx, 0, "products:*", 0).Iterator()
	for iter.Next(libs.Ctx) {
		libs.RedisClient.Del(libs.Ctx, iter.Val())
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Product created successfully",
		Data:    product,
	})
}

// GetProduct godoc
// @Summary Get list of products
// @Description Get paginated products list with search, sorting, and filtering options. Includes images, sizes, and variants.
// @Tags Products
// @Accept json
// @Produce json
// @Param search query string false "Search by title or description"
// @Param limit query int false "Number of products per page" default(10)
// @Param page query int false "Page number" default(1)
// @Param sort_by query string false "Sort field (id, title, base_price, created_at)" default(created_at)
// @Param order query string false "Sort order (ASC or DESC)" default(DESC)
// @Success 200 {object} models.Response{data=map[string]interface{}}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/products [get]
func (pc *ProductController) GetProducts(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	order := strings.ToUpper(ctx.DefaultQuery("order", "DESC"))

	cacheKey := fmt.Sprintf("products:page:%d:limit:%d:search:%s:sort:%s:order:%s",
		page, limit, search, sortBy, order)

	cached, err := libs.RedisClient.Get(libs.Ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var products []models.ProductResponse
		json.Unmarshal([]byte(cached), &products)
		ctx.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Products fetched from cache",
			Data: gin.H{
				"page":     page,
				"limit":    limit,
				"products": products,
			},
		})
		return
	}

	products, err := models.GetProducts(pc.DB, page, limit, search, sortBy, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch products",
			Data:    err.Error(),
		})
		return
	}

	dataJSON, _ := json.Marshal(products)
	libs.RedisClient.Set(libs.Ctx, cacheKey, dataJSON, 10*time.Minute)

	ctx.JSON(http.StatusOK, models.Response{
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
// @Description Get detailed information of a product, including its variants, sizes, and images.
// @Tags Products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response{data=models.ProductResponse} "Product fetched successfully"
// @Failure 400 {object} models.Response "Invalid product ID"
// @Failure 404 {object} models.Response "Product not found"
// @Failure 500 {object} models.Response "Internal server error"
// @Router /products/{id} [get]
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

	product, err := models.GetProductByID(pc.DB, productID)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Product not found",
			Data:    nil,
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Product fetched successfully",
		Data:    product,
	})
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update product details including title, description, price, stock, category, variants, sizes, and images
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param title formData string true "Product title"
// @Param description formData string false "Product description"
// @Param base_price formData number true "Product base price"
// @Param stock formData int true "Product stock"
// @Param category_id formData int true "Category ID"
// @Param variant_id formData []int false "Variant IDs (array)"
// @Param sizes formData []int false "Size IDs (array)"
// @Param images formData file false "Product images"
// @Success 200 {object} models.Response{data=models.ProductResponse} "Product updated successfully"
// @Failure 400 {object} models.Response "Invalid request or product ID"
// @Failure 500 {object} models.Response "Failed to update product"
// @Router /admin/products/{id} [patch]
func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid product ID"})
		return
	}

	var req models.ProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return
	}

	uploadDir := "./uploads/products"
	os.MkdirAll(uploadDir, os.ModePerm)

	savedFiles := []string{}
	maxSize := int64(2 * 1024 * 1024)
	allowedExts := []string{".jpg", ".jpeg", ".png"}

	contains := func(arr []string, s string) bool {
		for _, v := range arr {
			if v == s {
				return true
			}
		}
		return false
	}

	for _, file := range req.Images {
		if file.Size > maxSize {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "File too large(max 2mb): " + file.Filename,
			})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !contains(allowedExts, ext) {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Invalid file type(only png, jpg, jpeg): " + file.Filename})
			return
		}

		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.TrimSuffix(file.Filename, ext), ext)
		fullPath := filepath.Join(uploadDir, filename)
		if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
			continue
		}
		savedFiles = append(savedFiles, filename)
	}

	product, err := models.UpdateProduct(pc.DB, productID, req, savedFiles)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to update product",
			Data:    err.Error(),
		})
		return
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

// GetFavoriteProducts godoc
// @Summary Get favorite products
// @Description Menampilkan daftar produk favorit (bisa diatur limit-nya lewat query param ?limit=5)
// @Tags Products
// @Accept json
// @Produce json
// @Param limit query int false "Limit produk favorit yang ditampilkan" default(10)
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /favorite-product [get]
func (pc *ProductController) GetFavoriteProducts(ctx *gin.Context) {
	limitParam := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	query := `
		SELECT p.id, p.title, p.description, p.base_price, pi.image
		FROM products p
		LEFT JOIN product_images pi ON pi.product_id = p.id
		WHERE p.is_favorite = true
		GROUP BY p.id, pi.image
		ORDER BY p.updated_at DESC
		LIMIT $1
	`

	rows, err := pc.DB.Query(context.Background(), query, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil produk favorit",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	var favorites []map[string]interface{}
	for rows.Next() {
		var id int64
		var title, desc, image string
		var price float64

		if err := rows.Scan(&id, &title, &desc, &price, &image); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal membaca data produk",
				"error":   err.Error(),
			})
			return
		}

		favorites = append(favorites, map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": desc,
			"base_price":  price,
			"image":       image,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data produk favorit berhasil diambil",
		"data":    favorites,
	})
}

// FilterProducts godoc
// @Summary      Filter dan ambil daftar produk
// @Description  Endpoint ini mengambil daftar produk berdasarkan kategori, favorit, rentang harga, dan urutan (sort by). Semua parameter opsional.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        cat         query []int  false  "ID kategori, bisa lebih dari satu" collectionFormat(multi)  example(1,3)
// @Param        favorite    query bool   false  "Filter produk favorit"  example(true)
// @Param        price_min   query number false  "Batas harga minimum"  example(10000)
// @Param        price_max   query number false  "Batas harga maksimum"  example(50000)
// @Param        sortby      query string false "Urutkan hasil: name=A-Z, baseprice=termurah ke termahal"  Enums(name, baseprice)  example(baseprice)
// @Success      200  {object} map[string]interface{} "Data produk berhasil difilter"
// @Failure      500  {object} map[string]interface{} "Terjadi kesalahan server"
// @Router       /products [get]
func (pc *ProductController) FilterProducts(ctx *gin.Context) {
	var filter models.ProductFilter

	if cats := ctx.QueryArray("cat"); len(cats) > 0 {
		for _, c := range cats {
			if id, err := strconv.ParseInt(c, 10, 64); err == nil {
				filter.Categories = append(filter.Categories, id)
			}
		}
	}

	if fav := ctx.Query("favorite"); fav != "" {
		b := fav == "true"
		filter.IsFavorite = &b
	}

	if pmin := ctx.Query("price_min"); pmin != "" {
		if f, err := strconv.ParseFloat(pmin, 64); err == nil {
			filter.PriceMin = &f
		}
	}
	if pmax := ctx.Query("price_max"); pmax != "" {
		if f, err := strconv.ParseFloat(pmax, 64); err == nil {
			filter.PriceMax = &f
		}
	}

	filter.SortBy = ctx.DefaultQuery("sortby", "name")
	searchQuery := ctx.Query("q")

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	query := `
		SELECT 
			p.id,
			p.title,
			p.description,
			p.base_price,
			p.stock,
			p.category_id,
			p.created_at,
			p.updated_at,
			COALESCE(
				(SELECT pi.image FROM product_images pi WHERE pi.product_id = p.id LIMIT 1),
				''
			) AS image,
			COALESCE((
				SELECT json_agg(s.name)
				FROM product_sizes ps
				JOIN sizes s ON ps.size_id = s.id
				WHERE ps.product_id = p.id
			), '[]') AS sizes,
			COALESCE((
				SELECT json_agg(
					json_build_object(
						'id', v.id,
						'name', v.name,
						'additional_price', v.additional_price
					)
				)
				FROM product_variants pv
				JOIN variants v ON pv.variant_id = v.id
				WHERE pv.product_id = p.id
			), '[]') AS variants

		FROM products p
		WHERE 1=1
	`

	var args []interface{}
	argIndex := 1

	if len(filter.Categories) > 0 {
		query += fmt.Sprintf(" AND p.category_id = ANY($%d)", argIndex)
		args = append(args, filter.Categories)
		argIndex++
	}

	if filter.IsFavorite != nil {
		query += fmt.Sprintf(" AND p.is_favorite = $%d", argIndex)
		args = append(args, *filter.IsFavorite)
		argIndex++
	}

	if filter.PriceMin != nil {
		query += fmt.Sprintf(" AND p.base_price >= $%d", argIndex)
		args = append(args, *filter.PriceMin)
		argIndex++
	}

	if filter.PriceMax != nil {
		query += fmt.Sprintf(" AND p.base_price <= $%d", argIndex)
		args = append(args, *filter.PriceMax)
		argIndex++
	}

	if searchQuery != "" {
		query += fmt.Sprintf(" AND (LOWER(p.title) LIKE LOWER($%d) OR LOWER(p.description) LIKE LOWER($%d))", argIndex, argIndex+1)
		args = append(args, "%"+searchQuery+"%")
		args = append(args, "%"+searchQuery+"%")
		argIndex += 2
	}

	switch filter.SortBy {
	case "baseprice":
		query += " ORDER BY p.base_price ASC"
	default:
		query += " ORDER BY p.title ASC"
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := pc.DB.Query(context.Background(), query, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menjalankan query filter produk",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	var products []models.ProductResponseFilter

	for rows.Next() {
		var (
			id          int64
			title       string
			desc        string
			price       float64
			stock       int
			categoryID  int64
			createdAt   time.Time
			updatedAt   time.Time
			image       string
			sizesJSON   []byte
			variantsJSON []byte
		)

		err := rows.Scan(
			&id, &title, &desc, &price, &stock, &categoryID,
			&createdAt, &updatedAt,
			&image,
			&sizesJSON,
			&variantsJSON,
		)
		if err != nil {
			continue
		}

		var sizes []string
		var variants []map[string]interface{}

		json.Unmarshal(sizesJSON, &sizes)
		json.Unmarshal(variantsJSON, &variants)

		products = append(products, models.ProductResponseFilter{
			ID:          id,
			Title:       title,
			Description: desc,
			BasePrice:   price,
			Stock:       stock,
			CategoryID:  categoryID,
			Image:       image,
			Sizes:       sizes,
			Variants:    variants,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	nextPage := page + 1
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "Data produk berhasil difilter",
		"page":     page,
		"limit":    limit,
		"next":     nextPage,
		"previous": prevPage,
		"data":     products,
	})
}


// GetProductDetail godoc
// @Summary Get product detail by ID
// @Description Returns detailed information of a product including images, sizes, and recommended products
// @Tags Products
// @Param id path int true "Product ID"
// @Produce json
// @Success 200 {object} models.ProductDetail
// @Failure 400 {object} map[string]interface{} "Invalid product ID"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Router /products/{id} [get]
func (pc *ProductController) GetProductDetail(ctx *gin.Context) {

	idParam := ctx.Param("id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid product ID",
		})
		return
	}

	var product models.ProductDetail
	var categoryID int64

	query := `
		SELECT id, title, description, base_price, stock, category_id
		FROM products
		WHERE id=$1
	`
	err = pc.DB.QueryRow(context.Background(), query, productID).Scan(
		&product.ID, &product.Title, &product.Description,
		&product.BasePrice, &product.Stock, &categoryID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Product not found",
		})
		return
	}
	product.CategoryID = categoryID

	product.Variant = nil
	variantRows, err := pc.DB.Query(context.Background(),
		`SELECT v.id, v.name, v.additional_price
		 FROM variants v
		 JOIN product_variants pv ON pv.variant_id = v.id
		 WHERE pv.product_id=$1`, product.ID)
	if err == nil {
		var variants []models.Variant
		for variantRows.Next() {
			var v models.Variant
			if err := variantRows.Scan(&v.ID, &v.Name, &v.AdditionalPrice); err == nil {
				variants = append(variants, v)
			}
		}
		variantRows.Close()
		if len(variants) == 1 {
			product.Variant = &variants[0]
		} else if len(variants) > 1 {
			product.Variant = &variants[0]
		}
	}

	product.Images = []models.ProductImage{}
	rowsImg, err := pc.DB.Query(context.Background(),
		`SELECT image, updated_at FROM product_images WHERE product_id=$1 ORDER BY updated_at ASC`,
		product.ID,
	)
	if err == nil {
		for rowsImg.Next() {
			var img models.ProductImage
			if err := rowsImg.Scan(&img.Image, &img.UpdatedAt); err == nil {
				img.ProductID = product.ID
				product.Images = append(product.Images, img)
			}
		}
		rowsImg.Close()
	}

	product.Sizes = []models.Size{}
	if product.Variant == nil || product.Variant.Name != "Food" {
		rowsSize, err := pc.DB.Query(context.Background(),
			`SELECT s.id, s.name, s.additional_price 
			 FROM product_sizes ps
			 JOIN sizes s ON ps.size_id = s.id
			 WHERE ps.product_id=$1`,
			product.ID,
		)
		if err == nil {
			for rowsSize.Next() {
				var s models.Size
				if err := rowsSize.Scan(&s.ID, &s.Name, &s.AdditionalPrice); err == nil {
					product.Sizes = append(product.Sizes, s)
				}
			}
			rowsSize.Close()
		}
	}

	product.Recommended = []models.RecommendedProductInfo{}
	recQuery := `
		SELECT p.id, p.title, p.description, p.base_price, p.stock, p.category_id
		FROM recommended_products rp
		JOIN products p ON rp.recommended_id = p.id
		WHERE rp.product_id=$1
	`
	rowsRec, err := pc.DB.Query(context.Background(), recQuery, product.ID)
	if err == nil {
		for rowsRec.Next() {
			var rec models.RecommendedProductInfo
			if err := rowsRec.Scan(&rec.ID, &rec.Title, &rec.Description, &rec.BasePrice, &rec.Stock, &rec.CategoryID); err != nil {
				continue
			}

			var recVariantRows, _ = pc.DB.Query(context.Background(),
				`SELECT v.id, v.name, v.additional_price
				 FROM variants v
				 JOIN product_variants pv ON pv.variant_id = v.id
				 WHERE pv.product_id=$1`, rec.ID)
			var recVariants []models.Variant
			for recVariantRows.Next() {
				var v models.Variant
				if err := recVariantRows.Scan(&v.ID, &v.Name, &v.AdditionalPrice); err == nil {
					recVariants = append(recVariants, v)
				}
			}
			recVariantRows.Close()
			if len(recVariants) > 0 {
				rec.Variant = &recVariants[0]
			}

			rec.Images = []models.ProductImage{}
			imgRows, _ := pc.DB.Query(context.Background(),
				`SELECT image, updated_at FROM product_images WHERE product_id=$1 ORDER BY updated_at ASC`, rec.ID)
			for imgRows.Next() {
				var img models.ProductImage
				if err := imgRows.Scan(&img.Image, &img.UpdatedAt); err == nil {
					img.ProductID = rec.ID
					rec.Images = append(rec.Images, img)
				}
			}
			imgRows.Close()

			rec.Sizes = []models.Size{}
			if rec.Variant == nil || rec.Variant.Name != "Food" {
				sizeRows, _ := pc.DB.Query(context.Background(),
					`SELECT s.id, s.name, s.additional_price 
					 FROM product_sizes ps
					 JOIN sizes s ON ps.size_id = s.id
					 WHERE ps.product_id=$1`, rec.ID)
				for sizeRows.Next() {
					var s models.Size
					if err := sizeRows.Scan(&s.ID, &s.Name, &s.AdditionalPrice); err == nil {
						rec.Sizes = append(rec.Sizes, s)
					}
				}
				sizeRows.Close()
			}

			product.Recommended = append(product.Recommended, rec)
		}
		rowsRec.Close()
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Product detail fetched successfully",
		"data":    product,
	})
}

// AddToCart godoc
// @Summary Add items to cart
// @Description Add or update multiple items in user's cart. If item exists, quantity will be updated.
// @Tags Cart
// @Accept json
// @Produce json
// @Param carts body []models.Cart true "List of cart items"
// @Success 200 {object} models.Response{data=[]models.CartItemResponse} "Items added successfully"
// @Failure 400 {object} models.Response "Invalid request or stock not enough"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /cart [post]
func (pc *ProductController) AddToCart(ctx *gin.Context) {

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
		uid, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Invalid user ID",
			})
			return
		}
		userID = uid
	default:
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var carts []models.Cart
	if err := ctx.ShouldBindJSON(&carts); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return
	}

	var results []models.CartItemResponse
	for _, c := range carts {
		item, err := models.AddOrUpdateCart(pc.DB, userID, c.ProductID, c.SizeID, c.VariantID, c.Quantity)
		if err != nil {
			if err.Error() == "stock not enough for product" {
				ctx.JSON(http.StatusBadRequest, models.Response{
					Success: false,
					Message: err.Error(),
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		results = append(results, item)
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Items added successfully",
		Data:    results,
	})
}

// GetCart godoc
// @Summary Get user's cart
// @Description Retrieve all items in the authenticated user's cart.
// @Tags Cart
// @Produce json
// @Success 200 {object} models.Response{data=[]models.CartItemResponse} "Cart fetched successfully"
// @Failure 400 {object} models.Response "Invalid user ID"
// @Failure 401 {object} models.Response "Unauthorized"
// @Failure 500 {object} models.Response "Failed to fetch cart"
// @Security ApiKeyAuth
// @Router /cart [get]
func (pc *ProductController) GetCart(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(401, models.Response{
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
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Invalid user ID",
			})
			return
		}
		userID = tmp
	default:
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	cartResp, err := models.GetCartByUser(pc.DB, userID)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch cart",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Cart fetched successfully",
		Data:    cartResp,
	})
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a transaction for the authenticated user's cart items. User profile fields are used if not provided in request.
// @Tags Transaction
// @Accept json
// @Produce json
// @Param request body models.OrderTransactionRequest true "Transaction request body"
// @Success 201 {object} models.Response{data=models.OrderTransaction} "Transaction created successfully"
// @Failure 400 {object} models.Response "Invalid request or missing user info"
// @Failure 401 {object} models.Response "User not authenticated"
// @Failure 500 {object} models.Response "Failed to create transaction"
// @Security ApiKeyAuth
// @Router /transactions [post]
func (pc *ProductController) CreateTransaction(ctx *gin.Context) {
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

	var req models.OrderTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
		return
	}

	req.UserID = userID

	var dbFullname, dbEmail, dbPhone, dbAddress *string
	err := pc.DB.QueryRow(ctx, `
		SELECT u.fullname, u.email, p.phone, p.address
		FROM users u
		LEFT JOIN profile p ON p.user_id = u.id
		WHERE u.id=$1
	`, userID).Scan(&dbFullname, &dbEmail, &dbPhone, &dbAddress)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch user info",
			Data:    err.Error(),
		})
		return
	}

	if req.Fullname == "" {
		if dbFullname != nil && *dbFullname != "" {
			req.Fullname = *dbFullname
		} else {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Fullname must be provided",
			})
			return
		}
	}
	if req.Email == "" {
		if dbEmail != nil && *dbEmail != "" {
			req.Email = *dbEmail
		} else {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Email must be provided",
			})
			return
		}
	}
	if req.Phone == "" {
		if dbPhone != nil && *dbPhone != "" {
			req.Phone = *dbPhone
		} else {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Phone must be provided",
			})
			return
		}
	}
	if req.Address == "" {
		if dbAddress != nil && *dbAddress != "" {
			req.Address = *dbAddress
		} else {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Address must be provided",
			})
			return
		}
	}

	order, err := models.CreateOrderTransaction(pc.DB, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to create transaction",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Transaction created successfully",
		Data:    order,
	})
}
