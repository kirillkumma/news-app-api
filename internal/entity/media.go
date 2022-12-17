package entity

type (
	Editor struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	Media struct {
		ID                 int64  `json:"id"`
		RegistrationNumber int64  `json:"registrationNumber"`
		Name               string `json:"name"`
		Email              string `json:"email"`
		Editor             Editor `json:"editor"`
		Password           string `json:"-"`
	}

	MediaListItem struct {
		ID                 int64  `json:"id"`
		RegistrationNumber int64  `json:"registrationNumber"`
		Name               string `json:"name"`
		Email              string `json:"email"`
		Editor             Editor `json:"editor"`
		SubscriptionCount  int64  `json:"subscriptionCount"`
	}
)
