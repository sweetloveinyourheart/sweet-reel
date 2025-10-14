import { apiClient } from "../api-client"
import type { ApiResponse, PresignedUrlResponse, } from "../../types"

export class VideoService {
  /**
   * Generate the presigned-url for video uploading
   */
  static async generatePresignedURL(): Promise<ApiResponse<PresignedUrlResponse>> {
    return apiClient.post<PresignedUrlResponse>("/videos/presigned-url")
  }
}
