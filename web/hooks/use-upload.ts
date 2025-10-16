"use client"

import { useState, useCallback } from "react"
import { useApiClient } from "../lib/api"
import { ApiError } from "../lib/api/core/errors"
import { PresignedUrlRequest } from "../types"
import type { PresignedUrlResponse } from "../types"
import { ALLOWED_VIDEO_TYPES, formatFileSize, MAX_VIDEO_SIZE, validateFileSize, validateFileType, getS3Client } from "../lib/s3"

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

      // Validate file type
      if (!validateFileType(file, ALLOWED_VIDEO_TYPES)) {
        throw new Error('Invalid file type. Only video files are allowed.');
      }

      // Validate file size
      if (!validateFileSize(file, MAX_VIDEO_SIZE)) {
        throw new Error(
          `File is too large. Maximum size is ${formatFileSize(MAX_VIDEO_SIZE)}.`
        );
      }

      // Get the presigned url
      const uploadResponse = await api.post<PresignedUrlResponse>("/videos/presigned-url", metadata)
      if (!uploadResponse) {
        throw new Error("Failed to generate presigned url")
      }
      const presignedUrl = uploadResponse.presigned_url

      // Get S3 client instance
      const s3Client = getS3Client();

      // Upload with progress tracking
      await s3Client.upload(presignedUrl, file, {
        contentType: file.type,
        onProgress: (progress) => {
          setProgress(progress)
        },
      });

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
