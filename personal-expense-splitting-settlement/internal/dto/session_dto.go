package dto

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" bindings:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token" bindings:"required"`
	RefreshToken string `json:"refresh_token" bindings:"required"`
	ExpiresIn    int    `json:"expires_in"`
}
