package controllers

import (
	"coffeeder-backend/libs"
	"coffeeder-backend/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
               Message: "Invalid file type(only png and jpeg): " + file.Filename,
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
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch products",
			Data:    err.Error(),
		})
		return
	}

	dataJSON, _ := json.Marshal(products)
	libs.RedisClient.Set(libs.Ctx, cacheKey, dataJSON, 10*time.Minute)

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
			COALESCE(pi.image, '') AS image,
			COALESCE(v.name, '') AS variant_name,
			COALESCE(json_agg(DISTINCT s.name) FILTER (WHERE s.name IS NOT NULL), '[]') AS sizes
		FROM products p
		LEFT JOIN product_images pi ON pi.product_id = p.id
		LEFT JOIN product_sizes ps ON ps.product_id = p.id
		LEFT JOIN sizes s ON s.id = ps.size_id
		LEFT JOIN variants v ON v.id = p.variant_id
		WHERE 1=1
	`

	var pFilter []interface{}
	filterIndex := 1

	if len(filter.Categories) > 0 {
		query += fmt.Sprintf(" AND p.category_id = ANY($%d)", filterIndex)
		pFilter = append(pFilter, filter.Categories)
		filterIndex++
	}

	if filter.IsFavorite != nil {
		query += fmt.Sprintf(" AND p.is_favorite = $%d", filterIndex)
		pFilter = append(pFilter, *filter.IsFavorite)
		filterIndex++
	}

	if filter.PriceMin != nil {
		query += fmt.Sprintf(" AND p.base_price >= $%d", filterIndex)
		pFilter = append(pFilter, *filter.PriceMin)
		filterIndex++
	}

	if filter.PriceMax != nil {
		query += fmt.Sprintf(" AND p.base_price <= $%d", filterIndex)
		pFilter = append(pFilter, *filter.PriceMax)
		filterIndex++
	}

	if searchQuery != "" {
		query += fmt.Sprintf(" AND (LOWER(p.title) LIKE LOWER($%d) OR LOWER(p.description) LIKE LOWER($%d))", filterIndex, filterIndex)
		pFilter = append(pFilter, "%"+searchQuery+"%")
		filterIndex++
	}

	query += " GROUP BY p.id, pi.image, v.name"

	switch filter.SortBy {
	case "name":
		query += " ORDER BY p.title ASC"
	case "baseprice":
		query += " ORDER BY p.base_price ASC"
	default:
		query += " ORDER BY p.title ASC"
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := pc.DB.Query(context.Background(), query, pFilter...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menjalankan query filter produk",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	var products []map[string]interface{}
	for rows.Next() {
		var (
			id          int64
			title       string
			desc        string
			price       float64
			image       string
			variantName string
			sizesRaw    []byte
		)

		if err := rows.Scan(&id, &title, &desc, &price, &image, &variantName, &sizesRaw); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal membaca data produk",
				"error":   err.Error(),
			})
			return
		}

		var sizes []string
		json.Unmarshal(sizesRaw, &sizes)

		products = append(products, map[string]interface{}{
			"id":          id,
			"title":       title,
			"description": desc,
			"base_price":  price,
			"image":       image,
			"variant":     variantName,
			"sizes":       sizes,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data produk berhasil difilter",
		"page":    page,
		"limit":   limit,
		"data":    products,
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
	var variantID *int64
	var categoryID int64

	query := `
		SELECT id, title, description, base_price, stock, category_id, variant_id
		FROM products
		WHERE id=$1
	`
	err = pc.DB.QueryRow(context.Background(), query, productID).Scan(
		&product.ID, &product.Title, &product.Description,
		&product.BasePrice, &product.Stock, &categoryID, &variantID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Product not found",
		})
		return
	}
	product.CategoryID = categoryID

	if variantID != nil {
		var v models.Variant
		if err := pc.DB.QueryRow(context.Background(), `SELECT id, name FROM variants WHERE id=$1`, *variantID).
			Scan(&v.ID, &v.Name); err == nil {
			product.Variant = &v
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
	var recQuery string
	var recArgs []interface{}
	if variantID != nil {
		recQuery = `
			SELECT p.id, p.title, p.description, p.base_price, p.stock, p.category_id, p.variant_id
			FROM recommended_products rp
			JOIN products p ON rp.recommended_id = p.id
			WHERE rp.product_id=$1 AND p.variant_id=$2
		`
		recArgs = append(recArgs, product.ID, *variantID)
	} else {
		recQuery = `
			SELECT p.id, p.title, p.description, p.base_price, p.stock, p.category_id, p.variant_id
			FROM recommended_products rp
			JOIN products p ON rp.recommended_id = p.id
			WHERE rp.product_id=$1
		`
		recArgs = append(recArgs, product.ID)
	}

	rowsRec, err := pc.DB.Query(context.Background(), recQuery, recArgs...)
	if err == nil {
		for rowsRec.Next() {
			var rec models.RecommendedProductInfo
			var recVariantID *int64
			if err := rowsRec.Scan(&rec.ID, &rec.Title, &rec.Description, &rec.BasePrice, &rec.Stock, &rec.CategoryID, &recVariantID); err != nil {
				continue
			}

			if recVariantID != nil {
				var v models.Variant
				if err := pc.DB.QueryRow(context.Background(), `SELECT id, name FROM variants WHERE id=$1`, *recVariantID).Scan(&v.ID, &v.Name); err == nil {
					rec.Variant = &v
				}
			}

			rec.Images = []models.ProductImage{}
			imgRows, err := pc.DB.Query(context.Background(),
				`SELECT image, updated_at FROM product_images WHERE product_id=$1 ORDER BY updated_at ASC`,
				rec.ID,
			)
			if err == nil {
				for imgRows.Next() {
					var img models.ProductImage
					if err := imgRows.Scan(&img.Image, &img.UpdatedAt); err == nil {
						img.ProductID = rec.ID
						rec.Images = append(rec.Images, img)
					}
				}
				imgRows.Close()
			}

			rec.Sizes = []models.Size{}
			if rec.Variant == nil || rec.Variant.Name != "Food" {
				sizeRows, err := pc.DB.Query(context.Background(),
					`SELECT s.id, s.name, s.additional_price FROM product_sizes ps JOIN sizes s ON ps.size_id=s.id WHERE ps.product_id=$1`,
					rec.ID,
				)
				if err == nil {
					for sizeRows.Next() {
						var s models.Size
						if err := sizeRows.Scan(&s.ID, &s.Name, &s.AdditionalPrice); err == nil {
							rec.Sizes = append(rec.Sizes, s)
						}
					}
					sizeRows.Close()
				}
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

func (pc *ProductController) AddToCart(ctx *gin.Context) {
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
			"message": "Invalid user id",
		})
		return
	}

	var carts []models.Cart
	if err := ctx.ShouldBindJSON(&carts); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	var results []models.CartItemResponse
	for _, c := range carts {
		item, err := models.AddOrUpdateCart(pc.DB, userID, c.ProductID, c.SizeID, c.VariantID, c.Quantity)
		if err != nil {
			if err.Error() == "quantity exceeds available stock" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": err.Error(),
				})
				return
			}
			// Error lainnya
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		results = append(results, item)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Items added successfully",
		"result":  results,
	})
}

func (pc *ProductController) GetCart(ctx *gin.Context) {
	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized"})
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid user id"})
		return
	}

	cartItems, err := models.GetCartByUser(pc.DB, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error()})
		return
	}

	var total float64
	for _, item := range cartItems {
		total += item.Subtotal
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Cart fetched successfully",
		"result":  cartItems,
		"total":   total,
	})
}

func (pc *ProductController) CreateTransaction(ctx *gin.Context) {

	userIDValue, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	var userID int64
	switch v := userIDValue.(type) {
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Invalid user ID type",
		})
		return
	}

	var req models.OrderTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	req.UserID = userID

	var userEmail string
	err := pc.DB.QueryRow(ctx.Request.Context(), "SELECT email FROM users WHERE id=$1", userID).Scan(&userEmail)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch user email",
			"error":   err.Error(),
		})
		return
	}
	if req.Email == "" {
		req.Email = userEmail
	}

	var profilePhone, profileAddress *string
	err = pc.DB.QueryRow(ctx.Request.Context(),
		"SELECT phone, address FROM profile WHERE user_id=$1", userID).
		Scan(&profilePhone, &profileAddress)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch user profile",
			"error":   err.Error(),
		})
		return
	}

	if req.Phone == "" {
		if profilePhone != nil {
			req.Phone = *profilePhone
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Phone must be provided",
			})
			return
		}
	}

	if req.Address == "" {
		if profileAddress != nil {
			req.Address = *profileAddress
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Address must be provided",
			})
			return
		}
	}

	order, err := models.CreateOrderTransaction(pc.DB, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   order,
	})
}
