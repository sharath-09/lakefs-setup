package main


// -----------------------------
// Data Types (schemas)
// -----------------------------

// AuthRequestBody represents the request body for POST /auth
type AuthRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponseBody is the response for a successful POST /auth
type AuthResponseBody struct {
	ExternalUserIdentifier string `json:"external_user_identifier"`
}

// ErrorResponse is used for sending error messages
type ErrorResponse struct {
	Message string `json:"message"`
}

// UserDetail is returned by GET /auth/users/{userid}
type UserDetail struct {
	Username          string `json:"username"`
	CreationDate      int64  `json:"creation_date"`
	FriendlyName      string `json:"friendly_name"`
	Email             string `json:"email"`
	Source            string `json:"source"`
	EncryptedPassword string `json:"encryptedPassword"`
	ExternalID        string `json:"external_id"`
}
