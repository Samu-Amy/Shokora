package payloads

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// TODO: se dovesse servire UserBase (quindi solo una parte dei dati, fai Embed in User e magari anche metodo ToBaseUserPayload (da mettere eventualmente dentro ToBaseUserPayload per creare UserBase))
// type UserBase struct {
// 	Id        int64     `json:"id"`
// 	FirstName string    `json:"first_name"`
// 	LastName  string    `json:"last_name"`
// 	Email     string    `json:"email"`
// 	ImageUrl  string    `json:"image_url"`
// 	BirthDate time.Time `json:"birth_date"`
// }

type UserRes struct {
	Id         int64     `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	ImageUrl   string    `json:"image_url"`
	BirthDate  time.Time `json:"birth_date"`
	IsVerified bool      `json:"is_verified"`
	Role       user.Role `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
	UpdateAt   time.Time `json:"updated_at"`
}

func ToUserRes(user user.User) UserRes {
	return UserRes{
		Id:         user.Id,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		ImageUrl:   user.ImageUrl,
		BirthDate:  user.BirthDate,
		IsVerified: user.IsVerified,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		UpdateAt:   user.UpdatedAt,
	}
}
