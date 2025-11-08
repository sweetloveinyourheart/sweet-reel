"use client"

import { useEffect, useState } from "react"
import { useVideoAPI } from "@/hooks/use-video"
import { GetVideoMetadataResponse } from "@/types/video"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar"
import moment from "moment"

interface VideoMetadataCardProps {
  videoId: string
}

/**
 * VideoMetadataCard Component
 * 
 * Displays video metadata fetched from the API endpoint:
 * GET /api/v1/videos/{video_id}/metadata
 * 
 * @example
 * ```tsx
 * <VideoMetadataCard videoId="123e4567-e89b-12d3-a456-426614174000" />
 * ```
 */
export function VideoMetadataCard({ videoId }: VideoMetadataCardProps) {
  const { getVideoMetadata } = useVideoAPI()
  const [metadata, setMetadata] = useState<GetVideoMetadataResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchMetadata = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await getVideoMetadata(videoId)
        setMetadata(data)
      } catch (err) {
        console.error("Failed to fetch video metadata:", err)
        setError(err instanceof Error ? err.message : "Failed to load video metadata")
      } finally {
        setLoading(false)
      }
    }

    if (videoId) {
      fetchMetadata()
    }
  }, [videoId, getVideoMetadata])

  if (loading) {
    return (
      <div className="w-full border rounded-lg p-6">
        <div className="animate-pulse space-y-4">
          <div className="h-8 bg-gray-200 rounded w-3/4"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2"></div>
          <div className="space-y-2">
            <div className="h-4 bg-gray-200 rounded"></div>
            <div className="h-4 bg-gray-200 rounded"></div>
            <div className="h-4 bg-gray-200 rounded w-2/3"></div>
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="w-full border border-red-300 rounded-lg p-6 bg-red-50">
        <h3 className="text-lg font-semibold text-red-800 mb-2">Error</h3>
        <p className="text-red-600">{error}</p>
      </div>
    )
  }

  if (!metadata) {
    return null
  }

  return (
    <div className="w-full border rounded-lg overflow-hidden">
      {/* Video Header */}
      <div className="p-6 border-b">
        <h2 className="text-2xl font-bold mb-2">{metadata.video_title}</h2>
        <p className="text-gray-600">{metadata.video_description}</p>
      </div>

      {/* Channel Info */}
      <div className="p-6 border-b bg-gray-50">
        <div className="flex items-center gap-3">
          <Avatar className="h-12 w-12">
            <AvatarImage
              src={metadata.channel_metadata.onwer_metadata.picture}
              alt={metadata.channel_metadata.name}
            />
            <AvatarFallback>
              {metadata.channel_metadata.name[0]?.toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div>
            <p className="font-semibold text-lg">{metadata.channel_metadata.name}</p>
            <p className="text-sm text-gray-600">@{metadata.channel_metadata.handle}</p>
          </div>
        </div>
      </div>

      {/* Stats & Info */}
      <div className="p-6 space-y-4">
        {/* Stats Row */}
        <div className="flex flex-wrap gap-6">
          {metadata.total_view !== undefined && (
            <div className="flex items-center gap-2">
              <svg
                className="h-5 w-5 text-gray-500"
                fill="none"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <span className="text-sm font-medium">
                {metadata.total_view.toLocaleString()} views
              </span>
            </div>
          )}

          {metadata.processed_at && (
            <div className="flex items-center gap-2">
              <svg
                className="h-5 w-5 text-gray-500"
                fill="none"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              <span className="text-sm font-medium">
                {moment(metadata.processed_at * 1000).format("MMM D, YYYY")}
              </span>
            </div>
          )}
        </div>

        {/* Available Qualities */}
        {metadata.available_qualities && metadata.available_qualities.length > 0 && (
          <div>
            <div className="flex items-center gap-2 mb-3">
              <svg
                className="h-5 w-5 text-gray-500"
                fill="none"
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span className="text-sm font-semibold text-gray-700">Available Qualities</span>
            </div>
            <div className="flex flex-wrap gap-2">
              {metadata.available_qualities.map((quality) => (
                <span
                  key={quality}
                  className="px-3 py-1 text-sm font-medium bg-blue-100 text-blue-800 rounded-full"
                >
                  {quality}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
