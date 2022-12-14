package dto

import (
	"gopkg.in/guregu/null.v3"
	"news-app-api/internal/entity"
)

type (
	GetFeedParams struct {
		UserID int64
		Since  null.Int `query:"since"`
		Limit  null.Int `query:"limit"`
		Offset null.Int `query:"offset"`
	}

	GetFeedResult struct {
		Total int64                 `json:"total"`
		Items []entity.NewsListItem `json:"items"`
	}
)
