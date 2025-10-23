"use client"

import { Search, Menu, Video, Bell, Upload } from "lucide-react"
import { Button } from "./ui/button"
import { Input } from "./ui/input"
import { UserButtonClient } from "./user-button-client"
import { useState } from "react"
import type { Session } from "next-auth"
import Link from "next/link"

export function YouTubeHeader({
  onMenuClick,
  session
}: {
  onMenuClick?: () => void
  session: Session | null
}) {
  const [searchQuery, setSearchQuery] = useState("")

  return (
    <header className="sticky top-0 z-50 flex h-14 items-center justify-between border-b bg-background px-4">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          className="lg:hidden"
          onClick={onMenuClick}
        >
          <Menu className="h-6 w-6" />
        </Button>
        <Button variant="ghost" size="icon" className="hidden lg:flex">
          <Menu className="h-6 w-6" />
        </Button>
        <Link href={"/"}>
          <div className="flex items-center gap-1">
            <Video className="h-7 w-7 text-red-600" />
            <span className="text-xl font-bold">Sweet Reel</span>
          </div>
        </Link>
      </div>

      <div className="flex flex-1 max-w-2xl mx-4">
        <div className="flex w-full">
          <Input
            type="text"
            placeholder="Search"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="rounded-r-none border-r-0 focus-visible:ring-0 focus-visible:ring-offset-0"
          />
          <Button
            variant="outline"
            className="rounded-l-none border-l px-6"
            size="default"
          >
            <Search className="h-5 w-5" />
          </Button>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <Link href="/upload">
          <Button variant="ghost" size="icon" className="hidden sm:flex">
            <Upload className="h-6 w-6" />
          </Button>
        </Link>
        <Button variant="ghost" size="icon" className="hidden sm:flex">
          <Bell className="h-6 w-6" />
        </Button>
        <UserButtonClient session={session} />
      </div>
    </header>
  )
}
