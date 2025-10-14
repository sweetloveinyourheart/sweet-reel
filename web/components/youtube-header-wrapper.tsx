"use client"

import { YouTubeHeader as YouTubeHeaderClient } from "./youtube-header"
import { useState } from "react"
import { useSession } from "next-auth/react"

export function YouTubeHeader() {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const { data: session, status } = useSession()

  if (status === "loading") {
    return (
      <header className="sticky top-0 z-50 flex h-14 items-center justify-between border-b bg-background px-4">
        <div className="h-6 w-32 animate-pulse rounded bg-muted" />
      </header>
    )
  }

  return (
    <YouTubeHeaderClient
      onMenuClick={() => setSidebarOpen(!sidebarOpen)}
      session={session}
    />
  )
}
