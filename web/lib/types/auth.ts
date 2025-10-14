export type User = {
  id: string
  email: string
  name: string
  picture: string
  created_at: string
  updated_at: string
}

export type GoogleOAuthResponse = {
  jwt_token: string
  user: User
  is_new: boolean
}

export type GoogleOAuthRequest = {
  access_token: string
}

export type RefreshTokenResponse = {
	jwt_token: string
}
