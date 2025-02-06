package types

type Carttypes struct {
	ID       int `gorm:"primaryKey"`
	UserID   int
	BaseID   int
	SizeID   int
	FlavorID int
	ExtraIDs string // String chứa danh sách ID của Extras (có thể là JSON)
	Quantity int
	Price    float64
	Base     Basestypes
	Size     SizesTypes
	Flavor   Flavorstypes
	Extras   []Extrastypes
}
