import Image from "next/image"
import Link from "next/link"
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"

export interface VideoCardProps {
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
}

export function VideoCard({
  id,
  thumbnail,
  title,
  channel,
  views,
  timestamp,
  duration = "10:24",
}: VideoCardProps) {
  return (
    <Link href={`/video/${id}`} className="group flex flex-col gap-2">
      <div className="relative aspect-video w-full overflow-hidden rounded-xl bg-muted">
        <Image
          src={thumbnail}
          alt={title}
          fill
          className="object-cover transition-transform group-hover:scale-105"
        />
        {duration && (
          <div className="absolute bottom-2 right-2 rounded bg-black/80 px-1.5 py-0.5 text-xs font-semibold text-white">
            {duration}
          </div>
        )}
      </div>
      <div className="flex gap-3">
        <Avatar className="h-9 w-9 flex-shrink-0">
          <AvatarImage src={channel.avatar} alt={channel.name} />
          <AvatarFallback>{channel.name[0]}</AvatarFallback>
        </Avatar>
        <div className="flex flex-col gap-1">
          <h3 className="line-clamp-2 text-sm font-semibold leading-tight group-hover:text-primary">
            {title}
          </h3>
          <div className="flex flex-col text-xs text-muted-foreground">
            <span className="hover:text-foreground">{channel.name}</span>
            <span>
              {views} views • {timestamp}
            </span>
          </div>
        </div>
      </div>
    </Link>
  )
}
