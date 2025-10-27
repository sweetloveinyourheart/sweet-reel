import { auth } from "@/auth"
import { redirect } from "next/navigation"
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { UserVideoList } from "@/components/user-video-list"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Video, Users, Eye, Calendar } from "lucide-react"
import { getServerApiClient } from "@/lib/api/server"
import { GetChannelVideosResponse } from "@/types"
import moment from "moment"
import { ChannelResponse } from "@/types/channel"

export default async function ChannelPage() {
  const session = await auth()
  const api = await getServerApiClient()

  if (!session?.user || !api) {
    redirect("/signin")
  }

  const user = session.user

  // Fetch channel and channel videos from API
  const videosResponse = await api.get<GetChannelVideosResponse>("/channels/videos", {
    params: { limit: 25, offset: 0 }
  })

  const channelResponse = await api.get<ChannelResponse>("/channels")
  const channel = channelResponse

  // Transform API data to match VideoCard component props
  const userVideos = videosResponse.videos.map((video) => ({
    id: video.video_id,
    thumbnail: video.thumbnail_url || "https://images.unsplash.com/photo-1611162617474-5b21e879e113?w=500&h=281&fit=crop",
    title: video.title,
    channel: {
      name: channel.name,
      avatar: channel.owner.picture,
    },
    timestamp: moment(video.processed_at * 1000).fromNow(),
    duration: moment.utc(video.total_duration * 1000).format('mm:ss'),
    views: "0", // Placeholder - not implemented in API yet
  }))

  return (
    <div className="w-full">
      {/* Channel Header */}
      <div className="relative">
        {/* Channel Info */}
        <div className="px-4 sm:px-6 lg:px-8 pb-4">
          <div className="flex flex-col sm:flex-row gap-4 sm:gap-6">
            {/* Avatar */}
            <Avatar className="h-24 w-24 sm:h-32 sm:w-32 border-4 border-background">
              <AvatarImage
                src={channel.owner.picture}
                alt={channel.name}
              />
              <AvatarFallback className="text-2xl sm:text-4xl">
                {channel.name?.[0]?.toUpperCase() ?? "U"}
              </AvatarFallback>
            </Avatar>

            {/* Channel Info */}
            <div className="flex-1 flex flex-col justify-end gap-2">
              <div>
                <h1 className="text-2xl sm:text-3xl font-bold">{channel.name}</h1>
                <div className="flex flex-wrap gap-2 text-sm text-muted-foreground mt-1">
                  <span>{channel.handle}</span>
                  <span>•</span>
                  <span>{channel.total_videos} videos</span>
                  <span>•</span>
                  <span>{channel.total_views.toLocaleString()} views</span>
                </div>
              </div>
              <p className="text-sm text-muted-foreground max-w-2xl">
                {channel.description || "Welcome to my channel! I create content about web development, programming tutorials, and tech reviews."}
              </p>
            </div>

            {/* Action Buttons */}
            <div className="flex gap-2 sm:self-end">
              <Button variant="outline">Edit Channel</Button>
            </div>
          </div>

          {/* Stats Cards */}
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4 mt-6">
            <div className="border rounded-lg p-4 flex items-center gap-3">
              <div className="rounded-full bg-primary/10 p-2">
                <Video className="h-5 w-5 text-primary" />
              </div>
              <div>
                <div className="text-2xl font-bold">{channel.total_videos}</div>
                <div className="text-xs text-muted-foreground">Videos</div>
              </div>
            </div>
            <div className="border rounded-lg p-4 flex items-center gap-3">
              <div className="rounded-full bg-blue-500/10 p-2">
                <Eye className="h-5 w-5 text-blue-500" />
              </div>
              <div>
                <div className="text-2xl font-bold">{channel.total_views.toLocaleString()}</div>
                <div className="text-xs text-muted-foreground">Total Views</div>
              </div>
            </div>
            <div className="border rounded-lg p-4 flex items-center gap-3">
              <div className="rounded-full bg-green-500/10 p-2">
                <Users className="h-5 w-5 text-green-500" />
              </div>
              <div>
                <div className="text-2xl font-bold">{channel.subscriber_count.toLocaleString()}</div>
                <div className="text-xs text-muted-foreground">Subscribers</div>
              </div>
            </div>
            <div className="border rounded-lg p-4 flex items-center gap-3">
              <div className="rounded-full bg-purple-500/10 p-2">
                <Calendar className="h-5 w-5 text-purple-500" />
              </div>
              <div>
                <div className="text-2xl font-bold">{moment(channel.created_at).format('MMM YYYY')}</div>
                <div className="text-xs text-muted-foreground">Joined</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Content Tabs */}
      <div className="px-4 sm:px-6 lg:px-8 mt-6">
        <Tabs defaultValue="videos" className="w-full">
          <TabsList className="w-full justify-start border-b rounded-none h-auto p-0 bg-transparent">
            <TabsTrigger
              value="videos"
              className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent"
            >
              Videos
            </TabsTrigger>
            <TabsTrigger
              value="about"
              className="rounded-none border-b-2 border-transparent data-[state=active]:border-primary data-[state=active]:bg-transparent"
            >
              About
            </TabsTrigger>
          </TabsList>

          <TabsContent value="videos" className="mt-6">
            <UserVideoList
              initialVideos={userVideos}
              userName={channel.name}
              userAvatar={channel.owner.picture}
            />
          </TabsContent>

          <TabsContent value="about" className="mt-6">
            <div className="max-w-3xl space-y-6">
              <div>
                <h3 className="text-lg font-semibold mb-2">Description</h3>
                <p className="text-muted-foreground">
                  {channel.description || "Welcome to my channel! I create content about web development, programming tutorials, and tech reviews. Subscribe to stay updated with the latest videos."}
                </p>
              </div>

              <div>
                <h3 className="text-lg font-semibold mb-2">Details</h3>
                <div className="space-y-2 text-sm">
                  <div className="flex gap-2">
                    <span className="text-muted-foreground">Email:</span>
                    <span>{user.email}</span>
                  </div>
                  <div className="flex gap-2">
                    <span className="text-muted-foreground">Joined:</span>
                    <span>{moment(channel.created_at).format('MMMM D, YYYY')}</span>
                  </div>
                  <div className="flex gap-2">
                    <span className="text-muted-foreground">Total views:</span>
                    <span>{channel.total_views.toLocaleString()}</span>
                  </div>
                </div>
              </div>

              <div>
                <h3 className="text-lg font-semibold mb-2">Stats</h3>
                <div className="grid grid-cols-2 gap-4">
                  <div className="border rounded-lg p-4">
                    <div className="text-2xl font-bold">{channel.total_videos}</div>
                    <div className="text-sm text-muted-foreground">Total Videos</div>
                  </div>
                  <div className="border rounded-lg p-4">
                    <div className="text-2xl font-bold">{channel.subscriber_count.toLocaleString()}</div>
                    <div className="text-sm text-muted-foreground">Subscribers</div>
                  </div>
                </div>
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}
