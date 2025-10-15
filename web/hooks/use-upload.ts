"use client"

import { useState, useCallback } from "react"
import { useApiClient } from "../lib/api"
import { ApiError } from "../lib/api/core/errors"
import { PresignedUrlRequest } from "../types"
import type { PresignedUrlResponse } from "../types"

export function useVideoUpload() {
  const api = useApiClient()

  const [uploading, setUploading] = useState(false)
  const [progress, setProgress] = useState(0)
  const [error, setError] = useState<string | null>(null)

  const uploadVideo = useCallback(async (file: File, metadata: PresignedUrlRequest) => {
    try {
      setUploading(true)
      setError(null)
      setProgress(0)

      // Get the presigned url
      const uploadResponse = await api.post<PresignedUrlResponse>("/videos/presigned-url")
      if (!uploadResponse) {
        throw new Error("Failed to generate presigned url")
      }

      setProgress(100)
      return true
    } catch (err) {
      const errorMessage = err instanceof ApiError ? err.message : "Failed to upload video"
      setError(errorMessage)

      return false
    } finally {
      setUploading(false)
    }
  }, [api])

  const reset = useCallback(() => {
    setProgress(0)
    setError(null)
    setUploading(false)
  }, [])

  return {
    uploading,
    progress,
    error,
    uploadVideo,
    reset,
  }
}
