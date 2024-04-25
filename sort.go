package entity

const (
	ASCENDING  = "asc"
	DESCENDING = "desc"
)

type Sort struct {
	Key       string `json:"key"`
	Direction string `json:"d"`
}
