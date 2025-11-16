package controllers

import (
	"coffeeder-backend/models"
	"context"
	"fmt"
	"net/http"

	"path/filepath"
	"strconv"
	"strings"


	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserController struct {
	DB *pgxpool.Pool
}

// GetUsersList godoc
// @Summary Get list of users
// @Description Mengambil daftar users dengan pagination, optional search, dan sort
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param search query string false "Search by fullname or email"
// @Param sort_by query string false "Sort by field (fullname, email, created_at)" default(created_at)
// @Param sort_order query string false "Sort order (asc or desc)" default(desc)
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/users [get]
func (uc *UserController) GetUsersList(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortOrder := strings.ToLower(ctx.DefaultQuery("sort_order", "desc"))

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	allowedSortBy := map[string]bool{
		"fullname":   true,
		"email":      true,
		"created_at": true,
	}
	if !allowedSortBy[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	baseQuery := `
		SELECT 
			u.id, u.fullname, u.email, u.role,
			p.image, p.phone, p.address,
			u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profile p ON p.user_id = u.id
		WHERE 1=1
	`

	var args []interface{}
	argIndex := 1

	if search != "" {
		baseQuery += fmt.Sprintf(
			" AND (LOWER(u.fullname) LIKE LOWER($%d) OR LOWER(u.email) LIKE LOWER($%d))",
			argIndex, argIndex,
		)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	orderClause := fmt.Sprintf(" ORDER BY u.%s %s", sortBy, sortOrder)
	limitClause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	query := baseQuery + orderClause + limitClause

	rows, err := uc.DB.Query(context.Background(), query, args...)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to fetch users",
			Data:    err.Error(),
		})
		return
	}
	defer rows.Close()

	var users []models.UserList
	for rows.Next() {
		var u models.UserList
		p := &models.Profile{}

		var img, phone, address *string

		err := rows.Scan(
			&u.ID, &u.Fullname, &u.Email, &u.Role,
			&img, &phone, &address,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if img != nil {
			p.Image = img
		}
		if phone != nil {
			p.Phone = phone
		}
		if address != nil {
			p.Address = address
		}

		u.Profile = p 
		users = append(users, u)
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Users fetched successfully",
		Data: map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"sort_by":    sortBy,
			"sort_order": sortOrder,
			"users":      users,
		},
	})
}




// GetUserByID godoc
// @Summary Get user by ID
// @Description Mengambil detail user beserta profile berdasarkan ID (Admin Only)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/users/{id} [get]
func (uc *UserController) GetUserByID(ctx *gin.Context) {
	
	idParam := ctx.Param("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid user ID",
			Data:    nil,
		})
		return
	}

	query := `
		SELECT 
			u.id, u.fullname, u.email, u.role,
			p.image, p.phone, p.address,
			u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profile p ON p.user_id = u.id
		WHERE u.id=$1
	`

	var u models.UserList
	p := &models.Profile{} 

	err = uc.DB.QueryRow(context.Background(), query, userID).Scan(
		&u.ID, &u.Fullname, &u.Email, &u.Role,
		&p.Image, &p.Phone, &p.Address,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "User not found",
			Data:    nil,
		})
		return
	}

	u.Profile = p

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "User fetched successfully",
		Data:    u,
	})
}


// AddUser godoc
// @Summary Create a new user
// @Description Menambahkan user baru beserta profile dengan upload gambar
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param fullname formData string true "Full Name"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param role formData string true "User Role (admin/user)"
// @Param phone formData string false "Phone number"
// @Param address formData string false "Address"
// @Param image formData file false "Profile Image (max 2MB, jpg/png only)"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 409 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/users [post]
func (auc *UserController) AddUser(ctx *gin.Context) {
	var req models.AdminUserRequest

	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Request tidak valid",
			Data:    err.Error(),
		})
		return
	}

	if req.Image != nil {
		const maxSize = 2 * 1024 * 1024 
		if req.Image.Size > maxSize {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "File terlalu besar (max 2MB): " + req.Image.Filename,
			})
			return
		}

		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		ext := strings.ToLower(filepath.Ext(req.Image.Filename))
		if !allowedExts[ext] {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Format file tidak didukung (jpg/jpeg/png): " + req.Image.Filename,
			})
			return
		}
	}

	user, profile, err := models.AddUser(auc.DB,
		req.Fullname, req.Email, req.Password, req.Role,
		req.Phone, req.Address, req.Image,
	)
	if err != nil {
		if err.Error() == "email already exists" {
			ctx.JSON(409, models.Response{
				Success: false,
				Message: "Email sudah digunakan",
			})
			return
		}
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Gagal membuat user",
			Data:    err.Error(),
		})
		return
	}

	type UserResponse struct {
		ID        int64           `json:"id"`
		Fullname  string          `json:"fullname"`
		Email     string          `json:"email"`
		Role      string          `json:"role"`
		Profile   *models.Profile `json:"profile,omitempty"`
	}

	resp := UserResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Email:    user.Email,
		Role:     user.Role,
		Profile:  profile,
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "User berhasil dibuat",
		Data:    resp,
	})
}




// EditUser godoc
// @Summary Update user data
// @Description Mengupdate data user dan profil berdasarkan ID (support upload image)
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "User ID"
// @Param fullname formData string false "Full Name"
// @Param email formData string false "Email"
// @Param password formData string false "Password (optional)"
// @Param phone formData string false "Phone number"
// @Param address formData string false "Address"
// @Param image formData file false "Profile image (max 2MB, jpg/png only)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/users/{id} [patch]
func (auc *UserController) EditUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	userID, _ := strconv.ParseInt(idStr, 10, 64)

	fullname := ctx.PostForm("fullname")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")
	address := ctx.PostForm("address")

	file, _ := ctx.FormFile("image")

	if file != nil {
		const maxSize = 2 * 1024 * 1024
		if file.Size > maxSize {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "File terlalu besar (max 2MB): " + file.Filename,
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
				Message: "Format file tidak didukung (jpg/jpeg/png): " + file.Filename,
			})
			return
		}
	}

	user, profile, err := models.UpdateUser(auc.DB, userID, fullname, email, password, phone, address, file)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Gagal update user",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "User berhasil diperbarui",
		Data: map[string]interface{}{
			"user":    user,
			"profile": profile,
		},
	})
}


// DeleteUser godoc
// @Summary Delete a user
// @Description Menghapus user beserta profile berdasarkan ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /admin/users/{id} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := uc.DB.Exec(context.Background(), `DELETE FROM profile WHERE user_id=$1`, id)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete user profile",
			Data:    err.Error(),
		})
		return
	}

	_, err = uc.DB.Exec(context.Background(), `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to delete user",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "User deleted successfully",
		Data:    nil,
	})
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update phone, address, and profile image
// @Tags         Profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        Authorization header string true "Bearer <JWT token>"
// @Param        phone formData string false "Phone number"
// @Param        address formData string false "Address"
// @Param        image formData file false "Profile image (jpg, jpeg, png, max 2MB)"
// @Success      200 {object} models.ProfileUser
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /profile [patch]
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
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
	default:
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Invalid user ID type",
		})
		return
	}

	phone := ctx.PostForm("phone")
	address := ctx.PostForm("address")
	fullname := ctx.PostForm("fullname")
	email := ctx.PostForm("email")
	file, _ := ctx.FormFile("image")

	if file != nil {
		if file.Size > 2*1024*1024 {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "File too large (max 2MB)",
			})
			return
		}

		ext := file.Filename
		ext = filepath.Ext(ext)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "Invalid file type (only jpg, jpeg, png)",
			})
			return
		}
	}

	profileResp, err := models.UpdateProfile(uc.DB, userID, phone, address, fullname, email, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to update profile",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Profile updated successfully",
		Data:    profileResp,
	})
}


// GetProfile godoc
// @Summary Get current user's profile
// @Description Get profile info for the logged-in user
// @Tags Profile
// @Accept json
// @Produce json
// @Success 200 {object} models.ProfileUser
// @Failure 401 {object} map[string]interface{}
// @Router /profile [get]
// @Security BearerAuth
func (uc *UserController) GetProfile(ctx *gin.Context) {
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
	default:
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Invalid user ID type",
		})
		return
	}

	var profile models.ProfileUser
	err := uc.DB.QueryRow(ctx, `
        SELECT id, phone, address, image, user_id, created_at, updated_at
        FROM profile WHERE user_id=$1
    `, userID).Scan(
		&profile.ID, &profile.Phone, &profile.Address, &profile.Image,
		&profile.UserID, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var fullname, email string
	err = uc.DB.QueryRow(ctx, `
        SELECT fullname, email FROM users WHERE id=$1
    `, userID).Scan(&fullname, &email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	resp := models.ProfileResponse{
		ID:        profile.ID,
		Fullname:  fullname,
		Email:     email,
		Image:     profile.Image,
		Phone:     profile.Phone,
		Address:   profile.Address,
		UserID:    profile.UserID,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Profile fetched successfully",
		Data:    resp,
	})
}

