package entity

const (
	ASCENDING  = "ascending"
	DESCENDING = "descending"

	ASC  = "asc"
	DESC = "desc"
)

type Sort struct {
	Key       string `json:"key" url:"key"`
	Direction string `json:"d" url:"d"`
}
