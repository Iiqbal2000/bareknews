package restapi

// RespBody represents the common response body for JSON type.
type RespBody struct {
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
}

// ErrRespBody represents an error response body for JSON type.
type ErrRespBody struct {
	Err map[string]interface{} `json:"error" swaggertype:"object"`
}