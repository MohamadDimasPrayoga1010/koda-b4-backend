package controllers

import (
	"context"
	"fmt"
	"main/libs"
	"main/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserController struct {
	DB *pgxpool.Pool
}

func (uc *UserController) GetUsersList(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	search := ctx.DefaultQuery("search", "")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	query := `
		SELECT u.id, u.fullname, u.email, u.role,
		       p.image, p.phone, p.address,
		       u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profile p ON p.user_id = u.id
		WHERE 1=1
	`

	var args []interface{}
	if search != "" {
		query += " AND (LOWER(u.fullname) LIKE LOWER($1) OR LOWER(u.email) LIKE LOWER($1))"
		args = append(args, "%"+search+"%")
	}

	if len(args) > 0 {
		query += fmt.Sprintf(" ORDER BY u.id DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
		args = append(args, limit, offset)
	} else {
		query += " ORDER BY u.id DESC LIMIT $1 OFFSET $2"
		args = append(args, limit, offset)
	}

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
			fmt.Println("Scan error:", err)
			continue
		}
		u.Profile = &p
		users = append(users, u)
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Users fetched successfully",
		Data:    users,
	})
}

func (auc *UserController) AddUser(ctx *gin.Context) {
	var req models.AdminUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	if req.Password == "" {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Password is required",
			Data:    nil,
		})
		return
	}

	hashed, err := libs.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to hash password",
			Data:    nil,
		})
		return
	}

	query := `
		INSERT INTO users (fullname, email, password, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, fullname, email, role
	`
	var userID int64
	var fullname, email, role string
	err = auc.DB.QueryRow(context.Background(), query, req.Fullname, req.Email, hashed, req.Role).
		Scan(&userID, &fullname, &email, &role)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(409, models.Response{
				Success: false,
				Message: "Email already exists",
				Data:    nil,
			})
			return
		}
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to create user",
			Data:    err.Error(),
		})
		return
	}


	profileQuery := `
		INSERT INTO profile (user_id, image, phone, address)
		VALUES ($1, $2, $3, $4)
	`
	_, err = auc.DB.Exec(context.Background(), profileQuery, userID, req.Image, req.Phone, req.Address)
	if err != nil {
		fmt.Println("Failed to insert profile:", err)
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "User created successfully",
		Data: map[string]interface{}{
			"id":       userID,
			"fullname": fullname,
			"email":    email,
			"role":     role,
		},
	})
}


func (auc *UserController) EditUser(ctx *gin.Context) {
	id := ctx.Param("id")
	var req models.AdminUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
			Data:    nil,
		})
		return
	}

	passwordQuery := ""
	var hashed string
	if req.Password != "" {
		h, err := libs.HashPassword(req.Password)
		if err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Failed to hash password",
				Data:    nil,
			})
			return
		}
		hashed = h
		passwordQuery = ", password = '" + hashed + "'"
	}

	query := fmt.Sprintf(`
		UPDATE users 
		SET fullname=$1, email=$2, role=$3%s, updated_at=now()
		WHERE id=$4
	`, passwordQuery)

	_, err := auc.DB.Exec(context.Background(), query, req.Fullname, req.Email, req.Role, id)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to update user",
			Data:    err.Error(),
		})
		return
	}

	profileQuery := `
		UPDATE profile
		SET image=$1, phone=$2, address=$3, updated_at=now()
		WHERE user_id=$4
	`
	_, err = auc.DB.Exec(context.Background(), profileQuery, req.Image, req.Phone, req.Address, id)
	if err != nil {
		fmt.Println("Failed to update profile:", err)
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "User updated successfully",
	})
}

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



