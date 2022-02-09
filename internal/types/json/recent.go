package typesjson

type Recent struct {
	Data  []interface{} `json:"data"`
	Error string        `json:"error,omitempty"`
}
