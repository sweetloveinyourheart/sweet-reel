"use client"

import { useCallback } from "react"
import { useApiClient } from "../lib/api"
import { ApiError } from "../lib/api/core/errors"
import type { GetVideoMetadataResponse, PresignedUrlRequest, PresignedUrlResponse } from "../types"

/**
 * Hook for video-related API operations
 * 
 * @example
 * ```tsx
 * "use client"
 * 
 * function VideoMetadataComponent({ videoId }: { videoId: string }) {
 *   const { getVideoMetadata } = useVideoAPI()
 *   const [metadata, setMetadata] = useState(null)
 *   
 *   useEffect(() => {
 *     getVideoMetadata(videoId).then(setMetadata)
 *   }, [videoId])
 *   
 *   return <div>{metadata?.video_title}</div>
 * }
 * ```
 */
export function useVideoAPI() {
  const api = useApiClient()

  const getVideoMetadata = useCallback(async (videoId: string): Promise<GetVideoMetadataResponse> => {
    try {
      const metadata = await api.get<GetVideoMetadataResponse>(`/videos/${videoId}/metadata`)
      return metadata
    } catch (err) {
      throw err instanceof ApiError ? err : new Error("Failed to fetch video metadata")
    }
  }, [api])

  const generatePresignedUrl = useCallback(async (data: PresignedUrlRequest): Promise<PresignedUrlResponse> => {
    try {
      const response = await api.post<PresignedUrlResponse>("/videos/presigned-url", data)
      return response
    } catch (err) {
      throw err instanceof ApiError ? err : new Error("Failed to generate presigned URL")
    }
  }, [api])

  return {
    getVideoMetadata,
    generatePresignedUrl,
  }
}
