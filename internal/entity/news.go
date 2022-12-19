package entity

type (
	News struct {
		ID                      int64  `json:"id"`
		MediaRegistrationNumber int64  `json:"mediaRegistrationNumber"`
		Title                   string `json:"title"`
		Text                    string `json:"text"`
		CreatedAt               int64  `json:"createdAt"`
	}

	NewsListItem struct {
		ID         int64         `json:"id"`
		Media      MediaListItem `json:"media"`
		Title      string        `json:"title"`
		Text       string        `json:"text"`
		IsFavorite bool          `json:"isFavorite"`
		CreatedAt  int64         `json:"createdAt"`
	}
)
