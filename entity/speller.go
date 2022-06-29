package entity

type Result struct {
	Status   bool `json:"status"`
	Response `json:"response"`
}

type Response struct {
	Errors []Error `json:"errors"`
}
type Error struct {
	Id      string   `json:"id"`
	Offset  int      `json:"offset"`
	Length  int      `json:"length"`
	Bad     string   `json:"bad"`
	Better  []string `json:"better"`
	TypeRes string   `json:"type"`
}
