package models

// AdvertisementResponse структура ответа при создании/получении объявления
type AdvertisementResponse struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Quantity    int             `json:"quantity"`
	UserID      string          `json:"user_id"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Photos      []PhotoResponse `json:"photos"`
}

// PhotoResponse структура фотографии
type PhotoResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at"`
}

// MyAdvertisementsResponse структура ответа GET /my/advertisements
type MyAdvertisementsResponse struct {
	Items      []AdvertisementResponse `json:"items"`
	Total      int                     `json:"total"`
	HasNext    bool                    `json:"has_next"`
	NextCursor string                  `json:"next_cursor"`
}

// ErrorResponse структура ошибки
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
