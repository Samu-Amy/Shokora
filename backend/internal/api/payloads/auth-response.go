package payloads

// type UserResPayload struct {
// 	Id         int64     `json:"id"`
// 	FirstName  string    `json:"first_name"`
// 	LastName   string    `json:"last_name"`
// 	Email      string    `json:"email"`
// 	ImageUrl   string    `json:"image_url"`
// 	BirthDate  time.Time `json:"birth_date"`
// 	IsVerified bool      `json:"is_verified"`
// 	Role       user.Role `json:"role"`
// 	CreatedAt  time.Time `json:"created_at"`
// 	UpdatedAt  time.Time `json:"updated_at"`
// }

type RegisterUserRes struct {
	User              UserRes `json:"user"`
	VerificationId    *int64  `json:"verification_id,omitempty"`
	VerificationError string  `json:"verification_error,omitempty"`
	AuthError         bool    `json:"auth_error,omitempty"`
}
