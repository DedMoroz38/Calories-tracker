package dto

// BaseResponse is the uniform envelope returned by every endpoint.
//
//	success: { "data": ... }
//	failure: { "message": ... }
//
// Both fields are omitempty so a response carries only the relevant half.
type BaseResponse struct {
	Data    any `json:"data,omitempty"`
	Message any `json:"message,omitempty"`
}
