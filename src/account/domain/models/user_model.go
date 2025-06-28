package models

type User struct {
	ID       string  `json:"id" gorm:"primaryKey"`
	Email    string  `json:"email"`
	Password *string `json:"-" gorm:"column:password"`
	GoogleId *string `json:"googleId,omitempty"`
}

func CreateUserWithPassword(id string, email string, password string) *User {
	return &User{
		ID:       id,
		Email:    email,
		Password: &password,
	}
}

func CreateUserWithGoogleId(id string, email string, googleId string) *User {
	return &User{
		ID:       id,
		Email:    email,
		GoogleId: &googleId,
	}
}
