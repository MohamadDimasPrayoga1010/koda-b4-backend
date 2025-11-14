package controllers

import (
	"coffeeder-backend/libs"
	"coffeeder-backend/models"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
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
		baseQuery += " AND (LOWER(u.fullname) LIKE LOWER($" + strconv.Itoa(argIndex) + ") OR LOWER(u.email) LIKE LOWER($" + strconv.Itoa(argIndex) + "))"
		args = append(args, "%"+search+"%")
		argIndex++
	}

	orderClause := " ORDER BY u." + sortBy + " " + sortOrder

	limitClause := " LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
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
		var p models.Profile
		err := rows.Scan(
			&u.ID, &u.Fullname, &u.Email, &u.Role,
			&p.Image, &p.Phone, &p.Address,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			continue
		}
		u.Profile = &p
		users = append(users, u)
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Users fetched successfully",
		Data: gin.H{
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
	id := ctx.Param("id")

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
	var p models.Profile
	err := uc.DB.QueryRow(context.Background(), query, id).Scan(
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

	u.Profile = &p

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
	fullname := ctx.PostForm("fullname")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	role := ctx.PostForm("role")
	phone := ctx.PostForm("phone")
	address := ctx.PostForm("address")

	if fullname == "" || email == "" || password == "" || role == "" {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "fullname, email, password, dan role wajib diisi",
		})
		return
	}

	if len(password) < 6 {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Password minimal 6 karakter",
		})
		return
	}

	var imagePath string
	file, err := ctx.FormFile("image")
	if err == nil {
		const maxSize = 2 << 20
		if file.Size > maxSize {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Ukuran file melebihi 2MB",
			})
			return
		}

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
		}

		opened, _ := file.Open()
		defer opened.Close()

		buffer := make([]byte, 512)
		opened.Read(buffer)
		contentType := http.DetectContentType(buffer)
		if !allowedTypes[contentType] {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Format gambar tidak didukung (hanya JPG dan PNG)",
			})
			return
		}

		savePath := "uploads/" + file.Filename
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Gagal menyimpan file gambar",
				Data:    err.Error(),
			})
			return
		}
		imagePath = savePath
	}

	hashed, err := libs.HashPassword(password)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Gagal meng-hash password",
		})
		return
	}

	query := `
		INSERT INTO users (fullname, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, fullname, email, role
	`
	var userID int64
	var name, userEmail, userRole string
	err = auc.DB.QueryRow(context.Background(), query, fullname, email, hashed, role).
		Scan(&userID, &name, &userEmail, &userRole)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
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

	profileQuery := `
		INSERT INTO profile (user_id, image, phone, address)
		VALUES ($1, $2, $3, $4)
	`
	_, err = auc.DB.Exec(context.Background(), profileQuery, userID, imagePath, phone, address)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Gagal membuat profile user",
			Data:    err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "User berhasil dibuat",
		Data: map[string]interface{}{
			"id":       userID,
			"fullname": name,
			"email":    userEmail,
			"role":     userRole,
			"image":    imagePath,
		},
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
	id := ctx.Param("id")

	fullname := ctx.PostForm("fullname")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	phone := ctx.PostForm("phone")
	address := ctx.PostForm("address")

	if fullname == "" && email == "" && phone == "" && address == "" && password == "" {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Tidak ada data yang diubah",
		})
		return
	}

	if password != "" && len(password) < 6 {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Password minimal 6 karakter",
		})
		return
	}

	var imagePath string
	file, err := ctx.FormFile("image")
	if err == nil {
		const maxSize = 2 << 20
		if file.Size > maxSize {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Ukuran file melebihi 2MB",
			})
			return
		}

		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
		}

		opened, _ := file.Open()
		defer opened.Close()

		buffer := make([]byte, 512)
		_, _ = opened.Read(buffer)
		contentType := http.DetectContentType(buffer)
		if !allowedTypes[contentType] {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Format gambar tidak didukung (hanya JPG dan PNG)",
			})
			return
		}

		savePath := "uploads/" + file.Filename
		if err := ctx.SaveUploadedFile(file, savePath); err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Gagal menyimpan file gambar",
				Data:    err.Error(),
			})
			return
		}
		imagePath = savePath
	}

	updateFields := []string{}
	args := []interface{}{}
	argIdx := 1

	if fullname != "" {
		updateFields = append(updateFields, fmt.Sprintf("fullname=$%d", argIdx))
		args = append(args, fullname)
		argIdx++
	}
	if email != "" {
		updateFields = append(updateFields, fmt.Sprintf("email=$%d", argIdx))
		args = append(args, email)
		argIdx++
	}
	if password != "" {
		hashed, err := libs.HashPassword(password)
		if err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Gagal meng-hash password",
			})
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("password=$%d", argIdx))
		args = append(args, hashed)
		argIdx++
	}

	if len(updateFields) > 0 {
		query := fmt.Sprintf(`
			UPDATE users 
			SET %s, updated_at=now()
			WHERE id=$%d
		`, strings.Join(updateFields, ", "), argIdx)
		args = append(args, id)

		_, err := auc.DB.Exec(context.Background(), query, args...)
		if err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Gagal mengupdate data user",
				Data:    err.Error(),
			})
			return
		}
	}

	profileFields := []string{}
	profileArgs := []interface{}{}
	pIdx := 1

	if imagePath != "" {
		profileFields = append(profileFields, fmt.Sprintf("image=$%d", pIdx))
		profileArgs = append(profileArgs, imagePath)
		pIdx++
	}
	if phone != "" {
		profileFields = append(profileFields, fmt.Sprintf("phone=$%d", pIdx))
		profileArgs = append(profileArgs, phone)
		pIdx++
	}
	if address != "" {
		profileFields = append(profileFields, fmt.Sprintf("address=$%d", pIdx))
		profileArgs = append(profileArgs, address)
		pIdx++
	}

	if len(profileFields) > 0 {
		query := fmt.Sprintf(`
			UPDATE profile 
			SET %s, updated_at=now()
			WHERE user_id=$%d
		`, strings.Join(profileFields, ", "), pIdx)
		profileArgs = append(profileArgs, id)

		_, err := auc.DB.Exec(context.Background(), query, profileArgs...)
		if err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Gagal mengupdate profil user",
				Data:    err.Error(),
			})
			return
		}
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Data user berhasil diperbarui",
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}
	userID := userIDValue.(int64)

	phone := ctx.PostForm("phone")
	address := ctx.PostForm("address")
	file, _ := ctx.FormFile("image")

	profile, err := models.UpdateProfile(uc.DB, userID, phone, address, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": profile})
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User not authenticated"})
		return
	}
	userID := userIDValue.(int64)

	var profile models.ProfileUser
	err := uc.DB.QueryRow(ctx, `
        SELECT id, phone, address, image, user_id, created_at, updated_at
        FROM profile WHERE user_id=$1
    `, userID).Scan(&profile.ID, &profile.Phone, &profile.Address, &profile.Image, &profile.UserID, &profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": profile})
}
