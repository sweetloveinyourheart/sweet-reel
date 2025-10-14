import { ApiResponse } from "../types"
import { API_CONFIG } from "./config"

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: any
  ) {
    super(message)
    this.name = "ApiError"
  }
}

interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>
  token?: string
  timeout?: number
  retry?: number
  skipInterceptors?: boolean
}

type RequestInterceptor = (config: RequestInit, url: string) => RequestInit | Promise<RequestInit>
type ResponseInterceptor = (response: Response) => Response | Promise<Response>
type ErrorInterceptor = (error: ApiError) => ApiError | Promise<ApiError>

class ApiClient {
  private baseUrl: string
  private defaultTimeout: number = API_CONFIG.timeout
  private requestInterceptors: RequestInterceptor[] = []
  private responseInterceptors: ResponseInterceptor[] = []
  private errorInterceptors: ErrorInterceptor[] = []
  private abortControllers: Map<string, AbortController> = new Map()

  constructor(baseUrl: string = API_CONFIG.baseUrl) {
    this.baseUrl = baseUrl
  }

  // Interceptor management
  addRequestInterceptor(interceptor: RequestInterceptor) {
    this.requestInterceptors.push(interceptor)
  }

  addResponseInterceptor(interceptor: ResponseInterceptor) {
    this.responseInterceptors.push(interceptor)
  }

  addErrorInterceptor(interceptor: ErrorInterceptor) {
    this.errorInterceptors.push(interceptor)
  }

  // Token management - Use this to integrate with NextAuth or other auth systems
  setAuthToken(token: string | null) {
    if (typeof window !== "undefined") {
      if (token) {
        localStorage.setItem("auth_token", token)
      } else {
        localStorage.removeItem("auth_token")
      }
    }
  }

  getAuthToken(): string | null {
    if (typeof window !== "undefined") {
      return localStorage.getItem("auth_token")
    }
    return null
  }

  // Request cancellation
  cancelRequest(key: string) {
    const controller = this.abortControllers.get(key)
    if (controller) {
      controller.abort()
      this.abortControllers.delete(key)
    }
  }

  cancelAllRequests() {
    this.abortControllers.forEach((controller) => controller.abort())
    this.abortControllers.clear()
  }

  private async applyRequestInterceptors(config: RequestInit, url: string): Promise<RequestInit> {
    let modifiedConfig = config
    for (const interceptor of this.requestInterceptors) {
      modifiedConfig = await interceptor(modifiedConfig, url)
    }
    return modifiedConfig
  }

  private async applyResponseInterceptors(response: Response): Promise<Response> {
    let modifiedResponse = response
    for (const interceptor of this.responseInterceptors) {
      modifiedResponse = await interceptor(modifiedResponse)
    }
    return modifiedResponse
  }

  private async applyErrorInterceptors(error: ApiError): Promise<ApiError> {
    let modifiedError = error
    for (const interceptor of this.errorInterceptors) {
      modifiedError = await interceptor(modifiedError)
    }
    return modifiedError
  }

  private async requestWithTimeout(
    url: string,
    config: RequestInit,
    timeout: number
  ): Promise<Response> {
    const controller = new AbortController()
    const id = Math.random().toString(36)
    this.abortControllers.set(id, controller)

    const timeoutId = setTimeout(() => {
      controller.abort()
      this.abortControllers.delete(id)
    }, timeout)

    try {
      const response = await fetch(url, {
        ...config,
        signal: controller.signal,
      })
      clearTimeout(timeoutId)
      this.abortControllers.delete(id)
      return response
    } catch (error) {
      clearTimeout(timeoutId)
      this.abortControllers.delete(id)
      if (error instanceof Error && error.name === "AbortError") {
        throw new ApiError("Request timeout", 408)
      }
      throw error
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestOptions = {},
    attempt: number = 0
  ): Promise<ApiResponse<T>> {
    const { params, token, timeout, retry = 0, skipInterceptors = false, ...fetchOptions } = options

    // Build URL with query parameters
    let url = `${this.baseUrl}${endpoint}`
    if (params) {
      const searchParams = new URLSearchParams()
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          searchParams.append(key, String(value))
        }
      })
      const queryString = searchParams.toString()
      if (queryString) {
        url += `?${queryString}`
      }
    }

    const defaultHeaders: HeadersInit = {
      "Content-Type": "application/json",
    }

    // Add authorization token if available
    const authToken = token || this.getAuthToken()
    if (authToken) {
      defaultHeaders["Authorization"] = `Bearer ${authToken}`
    }

    let config: RequestInit = {
      ...fetchOptions,
      headers: {
        ...defaultHeaders,
        ...fetchOptions.headers,
      },
    }

    // Apply request interceptors
    if (!skipInterceptors) {
      config = await this.applyRequestInterceptors(config, url)
    }

    try {
      let response = await this.requestWithTimeout(
        url,
        config,
        timeout || this.defaultTimeout
      )

      // Apply response interceptors
      if (!skipInterceptors) {
        response = await this.applyResponseInterceptors(response)
      }

      // Handle non-JSON responses
      const contentType = response.headers.get("content-type")
      if (!contentType || !contentType.includes("application/json")) {
        if (!response.ok) {
          throw new ApiError(
            `HTTP error! status: ${response.status}`,
            response.status
          )
        }
        return { data: undefined as T }
      }

      const data = await response.json()

      if (!response.ok) {
        throw new ApiError(
          data.message || data.error || `HTTP error! status: ${response.status}`,
          response.status,
          data
        )
      }

      return data
    } catch (error) {
      let apiError: ApiError

      if (error instanceof ApiError) {
        apiError = error
      } else {
        apiError = new ApiError(
          error instanceof Error ? error.message : "An unknown error occurred",
          500
        )
      }

      // Retry logic for network errors or 5xx errors
      if (
        attempt < retry &&
        (apiError.status >= 500 || apiError.status === 408 || apiError.message.includes("network"))
      ) {
        // Exponential backoff
        const delay = Math.min(1000 * Math.pow(2, attempt), 10000)
        await new Promise((resolve) => setTimeout(resolve, delay))
        return this.request<T>(endpoint, options, attempt + 1)
      }

      // Apply error interceptors
      if (!skipInterceptors) {
        apiError = await this.applyErrorInterceptors(apiError)
      }

      throw apiError
    }
  }

  async get<T>(endpoint: string, options?: RequestOptions): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { ...options, method: "GET" })
  }

  async post<T>(
    endpoint: string,
    body?: any,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      ...options,
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async put<T>(
    endpoint: string,
    body?: any,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      ...options,
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async patch<T>(
    endpoint: string,
    body?: any,
    options?: RequestOptions
  ): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      ...options,
      method: "PATCH",
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  async delete<T>(endpoint: string, options?: RequestOptions): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { ...options, method: "DELETE" })
  }

  async uploadFile(
    endpoint: string,
    file: File,
    onProgress?: (progress: number) => void,
    additionalData?: Record<string, string>
  ): Promise<ApiResponse<any>> {
    return new Promise((resolve, reject) => {
      const xhr = new XMLHttpRequest()

      xhr.upload.addEventListener("progress", (e) => {
        if (e.lengthComputable && onProgress) {
          const percentage = (e.loaded / e.total) * 100
          onProgress(percentage)
        }
      })

      xhr.addEventListener("load", () => {
        if (xhr.status >= 200 && xhr.status < 300) {
          try {
            const data = JSON.parse(xhr.responseText)
            resolve(data)
          } catch (error) {
            resolve({ data: xhr.responseText })
          }
        } else {
          try {
            const errorData = JSON.parse(xhr.responseText)
            reject(
              new ApiError(
                errorData.message || `Upload failed with status ${xhr.status}`,
                xhr.status,
                errorData
              )
            )
          } catch {
            reject(
              new ApiError(
                `Upload failed with status ${xhr.status}`,
                xhr.status
              )
            )
          }
        }
      })

      xhr.addEventListener("error", () => {
        reject(new ApiError("Upload failed - Network error", 0))
      })

      xhr.addEventListener("abort", () => {
        reject(new ApiError("Upload cancelled", 0))
      })

      xhr.addEventListener("timeout", () => {
        reject(new ApiError("Upload timeout", 408))
      })

      xhr.open("POST", `${this.baseUrl}${endpoint}`)

      // Add auth token if available
      const authToken = this.getAuthToken()
      if (authToken) {
        xhr.setRequestHeader("Authorization", `Bearer ${authToken}`)
      }

      // Set timeout
      xhr.timeout = this.defaultTimeout

      const formData = new FormData()
      formData.append("file", file)

      // Add additional form data if provided
      if (additionalData) {
        Object.entries(additionalData).forEach(([key, value]) => {
          formData.append(key, value)
        })
      }

      xhr.send(formData)
    })
  }
}

export const apiClient = new ApiClient()

// Example: Add a request interceptor to log all requests (optional)
if (process.env.NODE_ENV === "development") {
  apiClient.addRequestInterceptor((config, url) => {
    console.log(`[API] ${config.method || "GET"} ${url}`)
    return config
  })
}

// Example: Add an error interceptor for global error handling
apiClient.addErrorInterceptor((error) => {
  // You can add global error handling here (e.g., toast notifications)
  if (error.status === 401) {
    // Handle unauthorized - redirect to login
    console.warn("[API] Unauthorized request")
  }
  return error
})
