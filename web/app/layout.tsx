import "./globals.css"
import type { Metadata } from "next"
import { Inter } from "next/font/google"
import { YouTubeLayout } from "@/components/youtube-layout"
import SessionProviderWrapper from "@/components/session-provider"
import AccountProviderWrapper from "@/components/account-provider"

const inter = Inter({ subsets: ["latin"] })

export const metadata: Metadata = {
  title: "Sweet Reel - Video Sharing Platform",
  description:
    "A modern video sharing platform built with Next.js and NextAuth.js",
}

export default function RootLayout({ children }: React.PropsWithChildren) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <SessionProviderWrapper>
          <AccountProviderWrapper>
            <YouTubeLayout>{children}</YouTubeLayout>
          </AccountProviderWrapper>
        </SessionProviderWrapper>
      </body>
    </html>
  )
}
