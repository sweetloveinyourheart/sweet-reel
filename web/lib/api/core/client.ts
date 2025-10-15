import { API_CONFIG } from "./config"
import { ApiError } from "./errors"
import type {
  RequestOptions,
  RequestInterceptor,
  ResponseInterceptor,
  ErrorInterceptor,
  TokenProvider,
} from "./types"

/**
 * Core API Client class
 * Handles HTTP requests with interceptors, retries, and authentication
 */
export class ApiClient {
  private baseUrl: string
  private defaultTimeout: number = API_CONFIG.timeout
  private requestInterceptors: RequestInterceptor[] = []
  private responseInterceptors: ResponseInterceptor[] = []
  private errorInterceptors: ErrorInterceptor[] = []
  private abortControllers: Map<string, AbortController> = new Map()
  private tokenProvider: TokenProvider | null = null

  constructor(baseUrl: string = API_CONFIG.baseUrl) {
    this.baseUrl = baseUrl
  }

  /**
   * Set a custom token provider (e.g., from useSession hook in client components)
   */
  setTokenProvider(provider: TokenProvider) {
    this.tokenProvider = provider
  }

  /**
   * Add a request interceptor
   */
  addRequestInterceptor(interceptor: RequestInterceptor) {
    this.requestInterceptors.push(interceptor)
  }

  /**
   * Add a response interceptor
   */
  addResponseInterceptor(interceptor: ResponseInterceptor) {
    this.responseInterceptors.push(interceptor)
  }

  /**
   * Add an error interceptor
   */
  addErrorInterceptor(interceptor: ErrorInterceptor) {
    this.errorInterceptors.push(interceptor)
  }

  /**
   * Get auth token - uses custom provider if set, otherwise returns null
   */
  async getAuthToken(): Promise<string | null> {
    if (this.tokenProvider) {
      return await this.tokenProvider()
    }
    return null
  }

  /**
   * Cancel a specific request by key
   */
  cancelRequest(key: string) {
    const controller = this.abortControllers.get(key)
    if (controller) {
      controller.abort()
      this.abortControllers.delete(key)
    }
  }

  /**
   * Cancel all pending requests
   */
  cancelAllRequests() {
    this.abortControllers.forEach((controller) => controller.abort())
    this.abortControllers.clear()
  }

  private async applyRequestInterceptors(
    config: RequestInit,
    url: string
  ): Promise<RequestInit> {
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
  ): Promise<T> {
    const {
      params,
      token,
      timeout,
      retry = 0,
      skipInterceptors = false,
      ...fetchOptions
    } = options

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
    const authToken = token || (await this.getAuthToken())
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
          throw new ApiError(`HTTP error! status: ${response.status}`, response.status)
        }
        return undefined as T
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
        (apiError.status >= 500 ||
          apiError.status === 408 ||
          apiError.message.includes("network"))
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

  /**
   * Perform a GET request
   */
  async get<T>(endpoint: string, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: "GET" })
  }

  /**
   * Perform a POST request
   */
  async post<T>(endpoint: string, body?: any, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  /**
   * Perform a PUT request
   */
  async put<T>(endpoint: string, body?: any, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  /**
   * Perform a PATCH request
   */
  async patch<T>(endpoint: string, body?: any, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: "PATCH",
      body: body ? JSON.stringify(body) : undefined,
    })
  }

  /**
   * Perform a DELETE request
   */
  async delete<T>(endpoint: string, options?: RequestOptions): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: "DELETE" })
  }

  /**
   * Upload a file with progress tracking
   */
  async uploadFile(
    endpoint: string,
    file: File,
    onProgress?: (progress: number) => void,
    additionalData?: Record<string, string>
  ): Promise<any> {
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
            resolve(xhr.responseText)
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
            reject(new ApiError(`Upload failed with status ${xhr.status}`, xhr.status))
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
      this.getAuthToken().then((authToken) => {
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
    })
  }
}

/**
 * Default API client instance
 */
export const apiClient = new ApiClient()

// Development logging
if (process.env.NODE_ENV === "development") {
  apiClient.addRequestInterceptor((config, url) => {
    console.log(`[API] ${config.method || "GET"} ${url}`)
    return config
  })
}

// Global error handling
apiClient.addErrorInterceptor((error) => {
  if (error.status === 401) {
    console.warn("[API] Unauthorized request")
  }
  return error
})
