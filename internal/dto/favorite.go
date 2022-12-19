package dto

import (
	"gopkg.in/guregu/null.v3"
	"news-app-api/internal/entity"
)

type (
	GetFavoriteListParams struct {
		UserID int64
		Limit  null.Int `query:"limit"`
		Offset null.Int `query:"offset"`
	}

	GetFavoriteListResult struct {
		Total int64                 `json:"total"`
		Items []entity.NewsListItem `json:"items"`
	}
)
