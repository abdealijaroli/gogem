package model

type Link struct {
	ID          int    `json:"id"`
	Link        string `json:"link"`
	RawData     string `json:"raw_data"`
	RefinedData string `json:"refined_data"`
}
