import { apiClient } from "../api-client"
import type { GoogleOAuthRequest, GoogleOAuthResponse, RefreshTokenResponse } from "../../types"

export class AuthService {
  /**
   * Login by Google OAuth
   */
  static async googleOAuth(data: GoogleOAuthRequest): Promise<GoogleOAuthResponse> {
    return apiClient.post<GoogleOAuthResponse>("/oauth/google", data)
  }

  /**
   * Refresh token
   */
  static async refreshToken(): Promise<RefreshTokenResponse> {
    return apiClient.get<RefreshTokenResponse>("/auth/refresh-token")
  }
}
