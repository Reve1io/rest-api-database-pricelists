package dto

type SearchRequest struct {
	Mapping map[string]string `json:"mapping"`
	Data    [][]interface{}   `json:"data"`
	Mode    string            `json:"mode"`
}
