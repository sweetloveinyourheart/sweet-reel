import NextAuth from "next-auth"
import "next-auth/jwt"

import Google from "next-auth/providers/google"
import { createStorage } from "unstorage"
import { UnstorageAdapter } from "@auth/unstorage-adapter"
import { ApiClient } from "./lib/api/core/client"
import { API_CONFIG } from "./lib/api/core/config"
import moment from "moment"
import type { GoogleOAuthResponse, RefreshTokenResponse } from "./types"

const storage = createStorage()
const authApiClient = new ApiClient(API_CONFIG.baseUrl)

export const { handlers, auth, signIn, signOut } = NextAuth({
  debug: !!process.env.AUTH_DEBUG,
  adapter: UnstorageAdapter(storage),
  providers: [
    Google,
  ],
  basePath: "/auth",
  session: { strategy: "jwt" },
  pages: {
    signIn: "/signin",
  },
  callbacks: {
    async jwt({ token, trigger, session, account }) {
      if (trigger === "update") token.name = session.user.name

      if (account?.provider === "google" && account.access_token) {
        const oauthResponse = await authApiClient.post<GoogleOAuthResponse>(
          "/oauth/google", 
          { access_token: account.access_token }
        )

        return {
          ...token,
          accessToken: oauthResponse.jwt_token,
          accessTokenExp: moment().add(1, 'hour').valueOf(),
          refreshToken: oauthResponse.jwt_refresh_token,
        }
      }

      if (token.refreshToken && token.accessTokenExp && Date.now() > token.accessTokenExp) {
        try {
          const newTokens = await authApiClient.get<RefreshTokenResponse>(
            "/auth/refresh-token", 
            { params: { "token": token.refreshToken } }
          )

          return {
            ...token,
            accessToken: newTokens.jwt_token,
            accessTokenExp: moment().add(1, 'hour').valueOf(),
          }
        } catch (err) {
          return { ...token, error: RefreshTokenError}
        }
      }

      return token
    },
    async session({ session, token }) {
      if (token?.accessToken) session.accessToken = token.accessToken

      return session
    },
  },
  experimental: { enableWebAuthn: true },
})

declare module "next-auth" {
  interface Session {
    accessToken?: string
    error?: SessionAuthError
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    accessToken?: string
    accessTokenExp?: number
    refreshToken?: string
  }
}

export type SessionAuthError = string
export const RefreshTokenError = "RefreshTokenError" as SessionAuthError
