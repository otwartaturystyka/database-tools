package models

// Datafile represents structure of data.json file.
type Datafile struct {
	Meta     Meta      `json:"meta"`
	Sections []Section `json:"sections"`
	Tracks   []Track   `json:"tracks"`
	Stories  []Story   `json:"stories"`
}

// AllPlaces returns places from all sections.
func (d *Datafile) AllPlaces() []Place {
	places := make([]Place, 0)
	for _, section := range d.Sections {
		places = append(places, section.Places...)
	}

	return places
}
