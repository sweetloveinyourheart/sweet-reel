"use client"

import { useState } from "react"
import { VideoCard } from "@/components/video-card"
import { Button } from "@/components/ui/button"
import { Video, Loader2 } from "lucide-react"
import { GetChannelVideosResponse, ChannelVideos } from "@/types"
import { useApiClient } from "@/lib/api/client"

interface UserVideoListProps {
  initialVideos: Array<{
    id: string
    thumbnail: string
    title: string
    channel: {
      name: string
      avatar: string
    }
    views: string
    timestamp: string
    duration?: string
  }>
  userName: string
  userAvatar: string
}

export function UserVideoList({ initialVideos, userName, userAvatar }: UserVideoListProps) {
  const api = useApiClient()

  const [videos, setVideos] = useState(initialVideos)
  const [offset, setOffset] = useState(25)
  const [isLoading, setIsLoading] = useState(false)
  const [hasMore, setHasMore] = useState(initialVideos.length === 25)

  const loadMore = async () => {
    setIsLoading(true)
    try {
      const data = await api.get<GetChannelVideosResponse>("/videos/user", {
        params: { limit: 25, offset }
      })

      const newVideos = data.videos.map((video: ChannelVideos) => ({
        id: video.video_id,
        thumbnail: video.thumbnail_url || "https://images.unsplash.com/photo-1611162617474-5b21e879e113?w=500&h=281&fit=crop",
        title: video.title,
        channel: {
          name: userName,
          avatar: userAvatar,
        },
        views: "0", // Placeholder - not implemented in API yet
        timestamp: "Recently", // Placeholder - not implemented in API yet
        duration: "0:00", // Placeholder - not implemented in API yet
      }))

      setVideos([...videos, ...newVideos])
      setOffset(offset + 25)
      setHasMore(newVideos.length === 25)
    } catch (error) {
      console.error("Failed to load more videos:", error)
    } finally {
      setIsLoading(false)
    }
  }

  if (videos.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <div className="rounded-full bg-muted p-6 mb-4">
          <Video className="h-12 w-12 text-muted-foreground" />
        </div>
        <h3 className="text-lg font-semibold mb-2">No videos yet</h3>
        <p className="text-muted-foreground mb-4">
          Upload your first video to get started
        </p>
        <Button>Upload Video</Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {videos.map((video) => (
          <VideoCard key={video.id} {...video} />
        ))}
      </div>

      {hasMore && (
        <div className="flex justify-center">
          <Button
            onClick={loadMore}
            disabled={isLoading}
            variant="outline"
            size="lg"
          >
            {isLoading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Loading...
              </>
            ) : (
              "Load More"
            )}
          </Button>
        </div>
      )}
    </div>
  )
}
