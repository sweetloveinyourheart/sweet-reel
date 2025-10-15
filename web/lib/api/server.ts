import { auth } from "@/auth"
import { ApiClient } from "./core/client"
import { API_CONFIG } from "./core/config"
import { cache } from "react"

/**
 * Get API client instance configured for server-side usage
 * 
 * This function must be called within server components, server actions,
 * or route handlers. It automatically retrieves the access token from
 * the NextAuth session and configures the API client.
 * 
 * Uses React's cache() to deduplicate requests within the same render pass.
 * 
 * @example
 * Server Component:
 * ```tsx
 * import { getServerApiClient } from "@/lib/api/server"
 * 
 * export default async function DashboardPage() {
 *   const api = await getServerApiClient()
 *   const videos = await api.get("/videos")
 *   
 *   return <div>{videos.map(...)}</div>
 * }
 * ```
 * 
 * @example
 * Server Action:
 * ```tsx
 * "use server"
 * 
 * import { getServerApiClient } from "@/lib/api/server"
 * 
 * export async function deleteVideo(id: string) {
 *   const api = await getServerApiClient()
 *   await api.delete(`/videos/${id}`)
 * }
 * ```
 * 
 * @example
 * Route Handler:
 * ```tsx
 * import { getServerApiClient } from "@/lib/api/server"
 * 
 * export async function GET() {
 *   const api = await getServerApiClient()
 *   const data = await api.get("/videos")
 *   return Response.json(data)
 * }
 * ```
 */
export const getServerApiClient = cache(async () => {
  const session = await auth()
  const client = new ApiClient(API_CONFIG.baseUrl)
  
  if (session?.accessToken) {
    client.setTokenProvider(() => session.accessToken || null)
  }
  
  return client
})
