"use client"

import { useSession } from "next-auth/react"
import { useEffect, useMemo } from "react"
import { ApiClient } from "./core/client"
import { API_CONFIG } from "./core/config"

/**
 * Hook to get API client configured with NextAuth session token (client-side only)
 * 
 * This hook automatically provides the access token from the user's session
 * to the API client for authenticated requests.
 * 
 * @example
 * ```tsx
 * "use client"
 * 
 * function MyComponent() {
 *   const api = useApiClient()
 *   
 *   const handleFetch = async () => {
 *     const data = await api.get("/videos")
 *     console.log(data)
 *   }
 *   
 *   return <button onClick={handleFetch}>Fetch Videos</button>
 * }
 * ```
 */
export function useApiClient() {
  const { data: session } = useSession()
  
  const client = useMemo(() => new ApiClient(API_CONFIG.baseUrl), [])
  
  useEffect(() => {
    // Configure API client with session token
    client.setTokenProvider(() => {
      return session?.accessToken || null
    })
  }, [client, session?.accessToken])
  
  return client
}
