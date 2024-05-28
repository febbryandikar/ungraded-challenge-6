package entity

type Recipe struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CookTime    int     `json:"cook_time"`
	Rating      float64 `json:"rating"`
}
