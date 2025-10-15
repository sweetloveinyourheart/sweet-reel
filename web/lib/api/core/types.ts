/**
 * Request options interface
 */
export interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>
  token?: string
  timeout?: number
  retry?: number
  skipInterceptors?: boolean
}

/**
 * Interceptor types
 */
export type RequestInterceptor = (
  config: RequestInit,
  url: string
) => RequestInit | Promise<RequestInit>

export type ResponseInterceptor = (
  response: Response
) => Response | Promise<Response>

export type ErrorInterceptor = (
  error: import('./errors').ApiError
) => import('./errors').ApiError | Promise<import('./errors').ApiError>

export type TokenProvider = () => string | null | Promise<string | null>
