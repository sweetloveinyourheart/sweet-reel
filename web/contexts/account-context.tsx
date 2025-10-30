"use client"

import React, { createContext, useContext } from "react"
import type { ChannelResponse } from "@/types/channel"

type AccountContextValue = {
  channel: ChannelResponse | null
}

const AccountContext = createContext<AccountContextValue | undefined>(undefined)

export function AccountProvider({
  channel,
  children,
}: {
  channel: ChannelResponse | null
  children: React.ReactNode
}) {
  return (
    <AccountContext.Provider value={{ channel }}>
      {children}
    </AccountContext.Provider>
  )
}

export function useAccount() {
  const context = useContext(AccountContext)
  if (!context) {
    throw new Error("useAccount must be used within an AccountProvider")
  }
  return context
}
