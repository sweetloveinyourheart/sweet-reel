import { apiClient } from "../api-client"
import type { PresignedUrlResponse } from "../../types"

export class VideoService {
  /**
   * Generate the presigned-url for video uploading
   */
  static async generatePresignedURL(): Promise<PresignedUrlResponse> {
    return apiClient.post<PresignedUrlResponse>("/videos/presigned-url")
  }
}
