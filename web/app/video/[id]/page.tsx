import { VideoCard } from "@/components/video-card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { ThumbsUp, ThumbsDown, Share2, Download, Flag } from "lucide-react"
import { getServerApiClient } from "@/lib/api/server"
import type { GetVideoMetadataResponse } from "@/types"
import { redirect } from "next/navigation"
import moment from "moment"

// Mock related videos (TODO: Replace with actual API call)
const relatedVideos = [
  {
    id: "2",
    thumbnail: "https://images.unsplash.com/photo-1498050108023-c5249f4df085?w=500&h=281&fit=crop",
    title: "React Hooks Explained: useState, useEffect, and Custom Hooks",
    channel: {
      name: "Code Academy",
      avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=CodeAcademy",
    },
    views: "89K",
    timestamp: "5 days ago",
    duration: "22:15",
  },
  {
    id: "3",
    thumbnail: "https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=500&h=281&fit=crop",
    title: "TypeScript for Beginners - Learn in 30 Minutes",
    channel: {
      name: "Dev Channel",
      avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=DevChannel",
    },
    views: "256K",
    timestamp: "1 week ago",
    duration: "31:08",
  },
]

const comments = [
  {
    id: "1",
    author: "John Developer",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=John",
    timestamp: "2 days ago",
    content: "This is exactly what I needed! Clear explanations and great examples. Thank you!",
    likes: 45,
  },
  {
    id: "2",
    author: "Sarah Code",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Sarah",
    timestamp: "1 day ago",
    content: "The section on Server Components really helped me understand the new paradigm. Subscribed!",
    likes: 23,
  },
  {
    id: "3",
    author: "Mike Frontend",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Mike",
    timestamp: "5 hours ago",
    content: "Great tutorial! Could you make one about authentication with NextAuth?",
    likes: 12,
  },
]

export default async function VideoPage({
  params
}: {
  params: Promise<{ id: string }>
}) {
  const { id } = await params
  
  // Get authenticated API client
  const api = await getServerApiClient()
  
  if (!api) {
    redirect("/signin")
  }

  // Fetch video metadata
  let metadata
  let error = null

  try {
    metadata = await api.get<GetVideoMetadataResponse>(`/videos/${id}/metadata`)
  } catch (err) {
    error = err instanceof Error ? err.message : "Failed to load video"
  }

  // If error occurred, show error page
  if (error || !metadata) {
    return (
      <div className="flex items-center justify-center min-h-[50vh]">
        <div className="text-center max-w-md">
          <h1 className="text-2xl font-bold mb-2">Video Not Found</h1>
          <p className="text-muted-foreground mb-4">
            {error || "The video you're looking for doesn't exist or has been removed."}
          </p>
          <Button asChild>
            <a href="/">Go Back Home</a>
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6 lg:flex-row">
      {/* Main content */}
      <div className="flex-1 space-y-4">
        {/* Video player placeholder */}
        <div className="aspect-video w-full rounded-xl bg-black flex items-center justify-center">
          <div className="text-white text-center">
            <div className="text-6xl mb-4">▶️</div>
            <p className="text-lg">Video Player Placeholder</p>
            <p className="text-sm text-gray-400">Video ID: {id}</p>
          </div>
        </div>

        {/* Video info */}
        <div className="space-y-4">
          <h1 className="text-xl font-bold">{metadata.video_title}</h1>

          {/* Channel info and actions */}
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <Avatar className="h-10 w-10">
                <AvatarImage src={metadata.channel_metadata.onwer_metadata.picture} />
                <AvatarFallback>{metadata.channel_metadata.name[0]}</AvatarFallback>
              </Avatar>
              <div className="flex flex-col">
                <span className="font-semibold">{metadata.channel_metadata.name}</span>
                <span className="text-sm text-muted-foreground">
                  @{metadata.channel_metadata.handle}
                </span>
              </div>
              <Button className="ml-4" size="sm">
                Subscribe
              </Button>
            </div>

            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" className="gap-2">
                <ThumbsUp className="h-4 w-4" />
                Like
              </Button>
              <Button variant="outline" size="sm">
                <ThumbsDown className="h-4 w-4" />
              </Button>
              <Button variant="outline" size="sm" className="gap-2">
                <Share2 className="h-4 w-4" />
                Share
              </Button>
              <Button variant="outline" size="sm" className="gap-2">
                <Download className="h-4 w-4" />
              </Button>
            </div>
          </div>

          {/* Description */}
          <div className="rounded-xl bg-muted p-4">
            <div className="mb-2 flex gap-4 text-sm font-semibold">
              <span>{metadata.total_view?.toLocaleString() || 0} views</span>
              {metadata.processed_at && (
                <span>{moment(metadata.processed_at * 1000).format("MMM D, YYYY")}</span>
              )}
            </div>
            <p className="whitespace-pre-line text-sm">{metadata.video_description}</p>
            
            {/* Available Qualities */}
            {metadata.available_qualities && metadata.available_qualities.length > 0 && (
              <div className="mt-4 pt-4 border-t">
                <p className="text-sm font-semibold mb-2">Available Qualities:</p>
                <div className="flex flex-wrap gap-2">
                  {metadata.available_qualities.map((quality: string) => (
                    <span
                      key={quality}
                      className="px-2 py-1 text-xs font-medium bg-primary/10 text-primary rounded"
                    >
                      {quality}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Comments section */}
          <div className="space-y-6">
            <div className="flex items-center gap-4">
              <h2 className="text-xl font-semibold">{comments.length} Comments</h2>
            </div>

            {/* Comment input */}
            <div className="flex gap-3">
              <Avatar className="h-10 w-10">
                <AvatarImage src="https://api.dicebear.com/7.x/avataaars/svg?seed=CurrentUser" />
                <AvatarFallback>U</AvatarFallback>
              </Avatar>
              <div className="flex-1">
                <input
                  type="text"
                  placeholder="Add a comment..."
                  className="w-full border-b bg-transparent px-2 py-2 text-sm outline-none focus:border-foreground"
                />
              </div>
            </div>

            {/* Comments list */}
            <div className="space-y-6">
              {comments.map((comment) => (
                <div key={comment.id} className="flex gap-3">
                  <Avatar className="h-10 w-10">
                    <AvatarImage src={comment.avatar} />
                    <AvatarFallback>{comment.author[0]}</AvatarFallback>
                  </Avatar>
                  <div className="flex-1 space-y-1">
                    <div className="flex items-center gap-2 text-sm">
                      <span className="font-semibold">{comment.author}</span>
                      <span className="text-muted-foreground">
                        {comment.timestamp}
                      </span>
                    </div>
                    <p className="text-sm">{comment.content}</p>
                    <div className="flex items-center gap-2">
                      <Button variant="ghost" size="sm" className="h-8 gap-2 px-2">
                        <ThumbsUp className="h-4 w-4" />
                        <span className="text-xs">{comment.likes}</span>
                      </Button>
                      <Button variant="ghost" size="sm" className="h-8 px-2">
                        <ThumbsDown className="h-4 w-4" />
                      </Button>
                      <Button variant="ghost" size="sm" className="h-8 text-xs">
                        Reply
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Sidebar - Related videos */}
      <div className="w-full space-y-2 lg:w-96">
        <h2 className="text-lg font-semibold">Related Videos</h2>
        <div className="space-y-2">
          {relatedVideos.map((video) => (
            <VideoCard key={video.id} {...video} />
          ))}
        </div>
      </div>
    </div>
  )
}
