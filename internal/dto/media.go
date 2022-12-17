package dto

import (
	"gopkg.in/guregu/null.v3"
	"news-app-api/internal/entity"
)

type (
	RegisterMediaParams struct {
		RegistrationNumber int64         `json:"registrationNumber"`
		Name               string        `json:"name"`
		Email              string        `json:"email"`
		Editor             entity.Editor `json:"editor"`
		Password           string        `json:"password"`
	}

	LoginMediaParams struct {
		RegistrationNumber int64  `json:"registrationNumber"`
		Password           string `json:"password"`
	}

	GetMediaListParams struct {
		Limit  null.Int `query:"limit"`
		Offset null.Int `query:"offset"`
	}

	GetMediaListResult struct {
		Total int64                  `json:"total"`
		Items []entity.MediaListItem `json:"items"`
	}

	ToggleSubscriptionParams struct {
		MediaID int64 `params:"media_id"`
		UserID  int64 `params:"-"`
	}

	ToggleSubscriptionResult struct {
		IsSubscribed bool `json:"isSubscribed"`
	}
)

func (p *RegisterMediaParams) Validate() error {
	if len(p.Name) > 32 {
		return &AppError{
			Message: "Максимальная длина названия - 32 символа",
			Code:    ErrCodeBadRequest,
		}
	} else if len(p.Email) > 16 {
		return &AppError{
			Message: "Максимальная длина email - 16 символа",
			Code:    ErrCodeBadRequest,
		}
	} else if len(p.Editor.FirstName) > 16 {
		return &AppError{
			Message: "Максимальная длина имени - 16 символа",
			Code:    ErrCodeBadRequest,
		}
	} else if len(p.Editor.LastName) > 16 {
		return &AppError{
			Message: "Максимальная длина фамилии - 16 символа",
			Code:    ErrCodeBadRequest,
		}
	} else if len(p.Password) > 32 {
		return &AppError{
			Message: "Максимальная длина пароля - 16 символа",
			Code:    ErrCodeBadRequest,
		}
	}
	return nil
}
