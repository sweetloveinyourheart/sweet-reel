import { VideoCard } from "@/components/video-card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { ThumbsUp, ThumbsDown, Share2, Download, Flag } from "lucide-react"

// Mock video data
const currentVideo = {
  id: "1",
  title: "Building Modern Web Applications with Next.js 14 - Complete Tutorial",
  channel: {
    name: "Tech Masters",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=TechMasters",
    subscribers: "1.2M",
  },
  views: "124,567",
  uploadDate: "Mar 10, 2024",
  likes: "12K",
  description: `In this comprehensive tutorial, we'll build a modern web application using Next.js 14.

We'll cover:
- App Router and Server Components
- Server Actions and Mutations
- Data Fetching and Caching
- Streaming and Suspense
- Metadata and SEO optimization

Perfect for developers who want to master Next.js and build production-ready applications!

🔗 Resources:
- Source code: github.com/example
- Documentation: nextjs.org

⏱️ Timestamps:
0:00 Introduction
2:30 Setting up Next.js 14
8:15 App Router basics
15:42 Server Components explained`,
}

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
          <h1 className="text-xl font-bold">{currentVideo.title}</h1>

          {/* Channel info and actions */}
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <Avatar className="h-10 w-10">
                <AvatarImage src={currentVideo.channel.avatar} />
                <AvatarFallback>{currentVideo.channel.name[0]}</AvatarFallback>
              </Avatar>
              <div className="flex flex-col">
                <span className="font-semibold">{currentVideo.channel.name}</span>
                <span className="text-sm text-muted-foreground">
                  {currentVideo.channel.subscribers} subscribers
                </span>
              </div>
              <Button className="ml-4" size="sm">
                Subscribe
              </Button>
            </div>

            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" className="gap-2">
                <ThumbsUp className="h-4 w-4" />
                {currentVideo.likes}
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
              <span>{currentVideo.views} views</span>
              <span>{currentVideo.uploadDate}</span>
            </div>
            <p className="whitespace-pre-line text-sm">{currentVideo.description}</p>
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
