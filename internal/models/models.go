package models

// League represents a football league
type League struct {
	ID         int    `json:"id"`
	SportID    int    `json:"sport_id"`
	CountryID  int    `json:"country_id"`
	Name       string `json:"name"`
	Active     bool   `json:"active"`
	ShortCode  string `json:"short_code"`
	ImagePath  string `json:"image_path"`
	Type       string `json:"type"`
	SubType    string `json:"sub_type"`
}

type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"home_team"`
}
