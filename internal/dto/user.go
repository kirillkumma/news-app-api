package dto

import (
	"gopkg.in/guregu/null.v3"
	"news-app-api/internal/entity"
)

type (
	RegisterUserParams struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Email    string `json:"email"`
	}

	LoginUserParams struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	GetSubscriptionListParams struct {
		UserID int64    `params:"user_id"`
		Limit  null.Int `query:"limit"`
		Offset null.Int `query:"offset"`
	}

	GetSubscriptionListResult struct {
		Total int64                  `json:"total"`
		Items []entity.MediaListItem `json:"items"`
	}
)

func (p *RegisterUserParams) Validate() error {
	if len(p.Login) > 32 {
		return &AppError{
			Message: "Максимальная длина логина - 32 символа",
			Code:    ErrCodeBadRequest,
		}
	} else if len(p.Password) > 32 {
		return &AppError{
			Message: "Максимальная длина пароля - 32 символа",
			Code:    ErrCodeBadRequest,
		}
	} else if len(p.Name) > 255 {
		return &AppError{
			Message: "Максимальная длина ФИО - 255 символов",
			Code:    ErrCodeBadRequest,
		}
	} else if len(p.Email) > 16 {
		return &AppError{
			Message: "Максимальная длина email - 16 символов",
			Code:    ErrCodeBadRequest,
		}
	}
	return nil
}
