package entity

type (
	User struct {
		ID       int64  `json:"id"`
		Login    string `json:"login"`
		Password string `json:"-"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}
)
