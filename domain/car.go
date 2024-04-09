package domain

// Car represents the Car entity.
type Car struct {
	ID     string `json:"id,omitempty"   db:"car_id"`
	RegNum string `json:"regNum"         db:"reg_num"`
	Mark   string `json:"mark"           db:"mark"`
	Model  string `json:"model"          db:"model"`
	Year   int    `json:"year,omitempty" db:"year"`
}

// FetchCarsFilter represents the filter to fetch many cars.
type FetchCarsFilter struct {
	IDs     []string
	RegNums []string
	Marks   []string
	Models  []string
	Years   []int
	Offset  int
	Limit   int
}
