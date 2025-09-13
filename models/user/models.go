package user

// Models
type GroceryCategory struct {
	ID        int    `json:"id"`
	Title     string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type GrocerySubcategory struct {
	ID         int    `json:"id"`
	CategoryID int    `json:"category_id"`
	Title      string `json:"name"`
	CreatedAt  string `json:"created_at"`
	ImageURL   string `json:"image_url"` // You'll need to add this field to your DB
}

type GroceryItem struct {
	ID             int     `json:"id"`
	Title          string  `json:"name"`
	Description    string  `json:"description"`
	PriceRetail    float64 `json:"price_retail"`
	Mrp            float64 `json:"mrp"`
	PriceWholesale float64 `json:"price_wholesale"`
	ImageURL1      string  `json:"image_url_1"`
}

type RecentSearch struct {
	Search    string `json:"search"`
	CreatedAt string `json:"created_at"`
}

type CategoryWithItems struct {
	ID    int           `json:"id"`
	Name  string        `json:"name"`
	Items []GroceryItem `json:"items"`
}

type CategoryWithSubcategories struct {
	ID            int                  `json:"id"`
	Name          string               `json:"name"`
	Subcategories []GrocerySubcategory `json:"subcategories"`
}

type DashboardNudge struct {
	NudgeTitle         string `json:"nudge_title"`
	NudgeBody          string `json:"nudge_body"`
	NudgeImage         string `json:"nudge_image"`
	NudgeNavigationURI string `json:"nudge_navigation_uri"`
}
