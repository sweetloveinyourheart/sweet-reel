import { VideoCard } from "@/components/video-card"

// Mock video data
const mockVideos = [
  {
    id: "1",
    thumbnail: "https://images.unsplash.com/photo-1611162617474-5b21e879e113?w=500&h=281&fit=crop",
    title: "Building Modern Web Applications with Next.js 14 - Complete Tutorial",
    channel: {
      name: "Tech Masters",
      avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=TechMasters",
    },
    views: "124K",
    timestamp: "2 days ago",
    duration: "15:42",
  },
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
  {
    id: "4",
    thumbnail: "https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=500&h=281&fit=crop",
    title: "CSS Grid vs Flexbox: When to Use Each One",
    channel: {
      name: "Design Weekly",
      avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=DesignWeekly",
    },
    views: "43K",
    timestamp: "3 days ago",
    duration: "18:27",
  },
  {
    id: "5",
    thumbnail: "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=500&h=281&fit=crop",
    title: "Node.js Crash Course - Building a REST API from Scratch",
    channel: {
      name: "Backend Pro",
      avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=BackendPro",
    },
    views: "178K",
    timestamp: "4 days ago",
    duration: "45:33",
  },
  {
    id: "6",
    thumbnail: "https://images.unsplash.com/photo-1551650975-87deedd944c3?w=500&h=281&fit=crop",
    title: "Docker Tutorial for Absolute Beginners",
    channel: {
      name: "DevOps Guide",
      avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=DevOpsGuide",
    },
    views: "312K",
    timestamp: "2 weeks ago",
    duration: "28:14",
  },
  {
    id: "7",
    thumbnail: "https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=500&h=281&fit=crop",
    title: "JavaScript ES6+ Features You Must Know in 2024",
    channel: {
      name: "JS Mastery",
      avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=JSMastery",
    },
    views: "201K",
    timestamp: "1 week ago",
    duration: "24:56",
  },
  {
    id: "8",
    thumbnail: "https://images.unsplash.com/photo-1605379399642-870262d3d051?w=500&h=281&fit=crop",
    title: "Tailwind CSS Full Course - Build Beautiful Websites Fast",
    channel: {
      name: "CSS Pro",
      avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=CSSPro",
    },
    views: "167K",
    timestamp: "5 days ago",
    duration: "38:42",
  },
]

export default function Home() {
  return (
    <div className="w-full">
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
        {mockVideos.map((video) => (
          <VideoCard key={video.id} {...video} />
        ))}
      </div>
    </div>
  )
}
