"use client"

import { signOut } from "next-auth/react"
import { useEffect } from "react"
import { RefreshTokenError } from "auth"

export default function SessionErrorHandler({ children, error }: { children: React.ReactNode, error?: string }) {
  useEffect(() => {
    if (error === RefreshTokenError) {
      signOut({ callbackUrl: "/signin" })
    }
  }, [error])

  if (error === RefreshTokenError) {
    return null
  }

  return <>{children}</>
}
