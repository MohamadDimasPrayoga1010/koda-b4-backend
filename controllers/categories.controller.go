package controllers

import (
	"context"
	"main/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryController struct {
	DB *pgxpool.Pool
}

// GetCategories godoc
// @Summary Get all categories
// @Description Mengambil daftar semua kategori
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/categories [get]
func (cc *CategoryController) GetCategories(ctx *gin.Context) {
	categories, err := models.GetAllCategories(cc.DB)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch categories",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Categories fetched successfully",
		Data:    categories,
	})
}


// GetCategoryByID godoc
// @Summary Get category by ID
// @Description Mengambil kategori berdasarkan ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/categories/{id} [get]
func (cc *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	catID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid category ID",
		})
		return
	}

	category, err := models.GetCategoryByID(cc.DB, catID)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Category not found",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Category fetched successfully",
		Data:    category,
	})
}


// CreateCategory godoc
// @Summary Create category
// @Description Menambahkan kategori baru
// @Tags Categories
// @Accept json
// @Produce json
// @Param name body string true "Category name"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/categories [post]
func (cc *CategoryController) CreateCategory(ctx *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Name is required",
			Data:    err.Error(),
		})
		return
	}

	category, err := models.CreateCategory(cc.DB, req.Name)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to create category",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Category created successfully",
		Data: gin.H{
			"id":         category.ID,
			"name":       category.Name,
			"created_at": category.CreatedAt,
			"updated_at": category.UpdatedAt,
		},
	})
}


// UpdateCategory godoc
// @Summary Update category
// @Description Mengupdate kategori berdasarkan ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param name body string true "Category name"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/categories/{id} [patch]
func (cc *CategoryController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	catID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid category ID",
		})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Name is required",
			Data:    err.Error(),
		})
		return
	}

	result, err := cc.DB.Exec(context.Background(),
		`UPDATE categories SET name=$1, updated_at=now() WHERE id=$2`, req.Name, catID)

	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to update category",
			Data:    err.Error(),
		})
		return
	}

	if result.RowsAffected() == 0 {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Category not found",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Category updated successfully",
		Data: gin.H{
			"id":   catID,
			"name": req.Name,
		},
	})
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Menghapus kategori berdasarkan ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/categories/{id} [delete]
func (cc *CategoryController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("id")
	catID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid category ID",
		})
		return
	}

	result, err := cc.DB.Exec(context.Background(),
		`DELETE FROM categories WHERE id=$1`, catID)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete category",
			Data:    err.Error(),
		})
		return
	}

	if result.RowsAffected() == 0 {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Category not found",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Category deleted successfully",
	})
}
