package entity

const (
	ASCENDING  = "asc"
	DESCENDING = "desc"
)

type Sort struct {
	Key       string `json:"key" url:"key"`
	Direction string `json:"d" url:"d"`
}
