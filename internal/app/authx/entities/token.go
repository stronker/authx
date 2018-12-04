package entities

// TokenData is the information that the system stores.
type TokenData struct {
	Username       string `cql:"user_name"`
	TokenID        string `cql:"token_id"`
	RefreshToken   []byte `cql:"refresh_token"`
	ExpirationDate int64 `cql:"expiration_date"`
}

// NewTokenData creates an instance of the structure
func NewTokenData(username string, tokenID string, refreshToken []byte,
	expirationDate int64) *TokenData {

	return &TokenData{Username: username,
		TokenID:        tokenID,
		RefreshToken:   refreshToken,
		ExpirationDate: expirationDate}
}
