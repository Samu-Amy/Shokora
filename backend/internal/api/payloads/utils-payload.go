package payloads

import "github.com/Samu-Amy/Shokora/internal/store/user"

func CreateUserResPayload(user *user.User) UserResPayload {
	return UserResPayload{
		Id:         user.Id,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		ImageUrl:   user.ImageUrl,
		BirthDate:  user.BirthDate,
		IsVerified: user.IsVerified,
		Role:       user.Role,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
