"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { cn } from "@/lib/utils"
import {
  Home,
  TrendingUp,
  Users,
  Library,
  History,
  PlaySquare,
  Clock,
  ThumbsUp,
  Flame,
  Music,
  Gamepad2,
  Newspaper,
  Trophy,
} from "lucide-react"

interface NavItem {
  title: string
  href: string
  icon: React.ComponentType<{ className?: string }>
}

const mainNav: NavItem[] = [
  { title: "Home", href: "/", icon: Home },
  { title: "Trending", href: "/trending", icon: TrendingUp },
  { title: "Subscriptions", href: "/subscriptions", icon: Users },
]

const libraryNav: NavItem[] = [
  { title: "Library", href: "/library", icon: Library },
  { title: "History", href: "/history", icon: History },
  { title: "Your videos", href: "/your-videos", icon: PlaySquare },
  { title: "Watch later", href: "/watch-later", icon: Clock },
  { title: "Liked videos", href: "/liked", icon: ThumbsUp },
]

const exploreNav: NavItem[] = [
  { title: "Popular", href: "/popular", icon: Flame },
  { title: "Music", href: "/music", icon: Music },
  { title: "Gaming", href: "/gaming", icon: Gamepad2 },
  { title: "News", href: "/news", icon: Newspaper },
  { title: "Sports", href: "/sports", icon: Trophy },
]

export function Sidebar() {
  const pathname = usePathname()

  return (
    <aside className="hidden lg:flex w-60 flex-col gap-2 border-r bg-background px-3 py-4">
      <nav className="flex flex-col gap-1">
        {mainNav.map((item) => (
          <SidebarLink
            key={item.href}
            item={item}
            isActive={pathname === item.href}
          />
        ))}
      </nav>

      <div className="my-2 h-px bg-border" />

      <nav className="flex flex-col gap-1">
        {libraryNav.map((item) => (
          <SidebarLink
            key={item.href}
            item={item}
            isActive={pathname === item.href}
          />
        ))}
      </nav>

      <div className="my-2 h-px bg-border" />

      <div className="px-3 pb-2 text-sm font-semibold text-muted-foreground">
        Explore
      </div>
      <nav className="flex flex-col gap-1">
        {exploreNav.map((item) => (
          <SidebarLink
            key={item.href}
            item={item}
            isActive={pathname === item.href}
          />
        ))}
      </nav>
    </aside>
  )
}

function SidebarLink({
  item,
  isActive,
}: {
  item: NavItem
  isActive: boolean
}) {
  const Icon = item.icon
  return (
    <Link
      href={item.href}
      className={cn(
        "flex items-center gap-4 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
        isActive
          ? "bg-accent text-accent-foreground"
          : "text-muted-foreground hover:bg-accent hover:text-accent-foreground"
      )}
    >
      <Icon className="h-5 w-5" />
      {item.title}
    </Link>
  )
}
