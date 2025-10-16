/**
 * API Client Library
 * 
 * Organized structure:
 * - core/    - Core API client functionality
 * - client.ts - Client-side hook (useApiClient) for React components
 * - server.ts - Server-side function (getServerApiClient) for server components/actions
 * 
 * @example Client-side usage (React components):
 * ```tsx
 * "use client"
 * import { useApiClient } from "@/lib/api"
 * 
 * function MyComponent() {
 *   const api = useApiClient()
 *   // Use api.get(), api.post(), etc.
 * }
 * ```
 * 
 * @example Server-side usage (Server components, actions, route handlers):
 * ```tsx
 * import { getServerApiClient } from "@/lib/api/server"
 * 
 * async function MyServerComponent() {
 *   const api = await getServerApiClient()
 *   // Use api.get(), api.post(), etc.
 * }
 * ```
 */

// Core exports
export * from "./core"

// Client-side hook
export * from "./client"

// Note: Server-side function must be imported directly from "./server"
// to avoid "use client" directive conflicts
