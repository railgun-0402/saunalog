package domain

type SaunaFacilityID string

type SaunaFacility struct {
	ID        SaunaFacilityID
	Name      string
	Address   string
	Price     int
	ImageURL  string
	SaunaInfo SaunaInfo
	CreatedAt int64
}

type SaunaInfo struct {
	Temperature  int
	Water        int
	HasMeal      bool
	HasRestArea  bool
	HasSleepRoom bool
}
