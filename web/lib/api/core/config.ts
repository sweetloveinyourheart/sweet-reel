/**
 * API Configuration
 */
export const API_CONFIG = {
  baseUrl: (typeof window === 'undefined' 
    ? process.env.API_URL 
    : process.env.NEXT_PUBLIC_API_URL) || "http://localhost:8080/api/v1",
  timeout: 30000,
} as const
