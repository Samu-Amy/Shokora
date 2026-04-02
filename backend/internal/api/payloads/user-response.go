package payloads

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/store/user"
)

type UserRes struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	// ImageUrl    string          `json:"image_url"`
	Birthday    time.Time       `json:"birthday"`
	IsVerified  bool            `json:"is_verified"`
	Role        user.Role       `json:"role"`
	Permissions user.Permission `json:"permissions"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdateAt    time.Time       `json:"updated_at"`
}

func ToUserRes(user user.User) *UserRes {
	return &UserRes{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		// ImageUrl:    user.ImageUrl,
		Birthday:    user.Birthday,
		IsVerified:  user.IsVerified,
		Role:        user.Role,
		Permissions: user.Permissions,
		CreatedAt:   user.CreatedAt,
		UpdateAt:    user.UpdatedAt,
	}
}
