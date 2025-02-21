package repository

import (
	"context"
	"example/API_Gateway/internal/db"
	"example/API_Gateway/internal/models"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
    db *db.DB
}

func NewUserRepository(db *db.DB) *UserRepository {
    return &UserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("error hashing password: %v", err)
    }

    var user models.User
    err = r.db.Pool.QueryRow(ctx,
        `INSERT INTO users (username, password_hash, email, role)
         VALUES ($1, $2, $3, $4)
         RETURNING id, username, email, role, created_at, updated_at`,
        req.Username, string(hashedPassword), req.Email, req.Role,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        return nil, fmt.Errorf("error creating user: %v", err)
    }

    return &user, nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    var user models.User
    err := r.db.Pool.QueryRow(ctx,
        `SELECT id, username, email, role, created_at, updated_at
         FROM users WHERE id = $1`,
        id,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        return nil, fmt.Errorf("error getting user: %v", err)
    }

    return &user, nil
}

// UpdateUser updates a user's information
func (r *UserRepository) UpdateUser(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
    var user models.User
    
    // Build query dynamically based on provided fields
    query := "UPDATE users SET "
    params := []interface{}{}
    paramCount := 1
    updates := []string{}

    if req.Username != "" {
        updates = append(updates, fmt.Sprintf("username = $%d", paramCount))
        params = append(params, req.Username)
        paramCount++
    }
    if req.Email != "" {
        updates = append(updates, fmt.Sprintf("email = $%d", paramCount))
        params = append(params, req.Email)
        paramCount++
    }
    if req.Role != "" {
        updates = append(updates, fmt.Sprintf("role = $%d", paramCount))
        params = append(params, req.Role)
        paramCount++
    }

    if len(updates) == 0 {
        return nil, fmt.Errorf("no fields to update")
    }

    query += strings.Join(updates, ", ")
    query += fmt.Sprintf(" WHERE id = $%d RETURNING id, username, email, role, created_at, updated_at", paramCount)
    params = append(params, id)

    err := r.db.Pool.QueryRow(ctx, query, params...).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.Role,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err != nil {
        return nil, fmt.Errorf("error updating user: %v", err)
    }

    return &user, nil
}

// DeleteUser deletes a user by their ID
func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
    _, err := r.db.Pool.Exec(ctx,
        `DELETE FROM users WHERE id = $1`,
        id,
    )

    if err != nil {
        return fmt.Errorf("error deleting user: %v", err)
    }

    return nil
}

// GetAllUsers retrieves all users (for admin use)
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
    rows, err := r.db.Pool.Query(ctx,
        `SELECT id, username, email, role, created_at, updated_at
         FROM users`,
    )
    if err != nil {
        return nil, fmt.Errorf("error querying users: %v", err)
    }
    defer rows.Close()

    var users []*models.User
    for rows.Next() {
        var user models.User
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("error scanning user: %v", err)
        }
        users = append(users, &user)
    }

    return users, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
    var user models.User
    err := r.db.Pool.QueryRow(ctx,
        `SELECT id, username, email, role, created_at, updated_at
         FROM users WHERE username = $1`,
        username,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        return nil, fmt.Errorf("error getting user by username: %v", err)
    }

    return &user, nil
} 