package svc

import "github.com/rs/zerolog/log"

// Tap is the representation of a beer on tap
type Tap struct {
	ID          string  `json:"id,omitempty"`
	Name        string  `json:"name"`
	Tag         string  `json:"tag"`
	Description string  `json:"description"`
	TapNumber   int     `json:"tapNumber"`
	Gravity     float64 `json:"originalGravity"`
	Color       string  `json:"color"`
	IBUs        float64 `json:"ibu"`
	Calories    int     `json:"calorires"`
	ABV         float64 `json:"abv"`
}

// NewTap creates a new instance of a Tap
func NewTap() *Tap {
	return &Tap{}
}

// TapList is a list of taps
type TapList []Tap

// TapWriter writes the list of taps to something.
type TapWriter interface {
	Write([]Tap) error
}

// TapReader reads the list of taps
type TapReader interface {
	Read() ([]Tap, error)
}

// TapServicer is the service interface for the Taps
type TapServicer interface {
	FindAll() ([]Tap, error)
	Save(Tap) (Tap, error)
}

// TapService is the default implementation for the TapServicer interface
type TapService struct {
	reader TapReader
	writer TapWriter
}

// FindAll returns the list of beers that are currently on tap or
// an error if something unexpected occurred.
func (svc *TapService) FindAll() ([]Tap, error) {
	return svc.reader.Read()
}

// Save allows you to save a tap
func (svc *TapService) Save(t Tap) (Tap, error) {
	// Load the taps from the service
	taps, err := svc.FindAll()
	log.Debug().Msgf("Found % taps", len(taps))
	if err != nil {
		return t, err
	}
	if len(taps) == 0 {
		log.Warn().Msg("Found no taps, creating new list")
		taps = []Tap{t}
	} else {
		log.Debug().Msgf("Found %s taps, adding %", len(taps), t.Name)
		taps = append(taps, t)
	}

	err = svc.writer.Write(taps)
	return t, err
}

// NewTapService returns a new instance of the TapService
func NewTapService(w TapWriter, r TapReader) TapServicer {
	return &TapService{
		writer: w,
		reader: r,
	}
}
