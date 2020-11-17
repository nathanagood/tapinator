package svc

// Tap is the representation of a beer on tap
type Tap struct {
	ID          string  `json:"id,omitempty"`
	Name        string  `json:"name"`
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
	TapNumber   int     `json:"tapNumber"`
	Gravity     float64 `json:"originalGravity"`
	Color       float64 `json:"color"`
	IBUs        float64 `json:"ibu"`
	Calories    int     `json:"calorires"`
	ABV         int     `json:"abv"`
}

// NewTap creates a new instance of a Tap
func NewTap() *Tap {
	return &Tap{}
}

// TapList is a list of taps
type TapList []Tap

// TapServicer is the service interface for the Taps
type TapServicer interface {
	GetTaps() ([]Tap, error)
}

// TapService is the default implementation for the TapServicer interface
type TapService struct {
}

// GetTaps returns the list of beers that are currently on tap or
// an error if something unexpected occurred.
func (svc *TapService) GetTaps() ([]Tap, error) {
	tap := NewTap()
	tapList := []Tap{*tap}
	return tapList, nil
}

// NewTapService returns a new instance of the TapService
func NewTapService() TapServicer {
	return &TapService{}
}
