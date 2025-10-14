import { apiClient } from "../api-client"
import type { ApiResponse, GoogleOAuthRequest, GoogleOAuthResponse, RefreshTokenResponse } from "../../types"

export class UserService {
  /**
   * Login by Google OAuth
   */
  static async googleOAuth(data: GoogleOAuthRequest): Promise<ApiResponse<GoogleOAuthResponse>> {
    return apiClient.post<GoogleOAuthResponse>("/oauth/google", data)
  }

  /**
   * Refresh token
   */
  static async refreshToken(): Promise<ApiResponse<RefreshTokenResponse>> {
    return apiClient.get<RefreshTokenResponse>("/auth/refresh-token")
  }
}
