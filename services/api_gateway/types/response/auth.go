package response

type GoogleOAuthUser struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type GoogleOAuthResponse struct {
	JwtToken string          `json:"jwt_token"`
	User     GoogleOAuthUser `json:"user"`
	IsNew    bool            `json:"is_new"`
}

type RefreshTokenResponse struct {
	JwtToken string `json:"jwt_token"`
}
