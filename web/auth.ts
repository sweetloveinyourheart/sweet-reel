import NextAuth from "next-auth"
import "next-auth/jwt"

import Google from "next-auth/providers/google"
import { createStorage } from "unstorage"
import { UnstorageAdapter } from "@auth/unstorage-adapter"
import { AuthService } from "./lib/api/services"
import moment from "moment"

const RefreshTokenError = "RefreshTokenError"
const storage = createStorage()

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
        const oauthResponse = await AuthService.googleOAuth({ access_token: account.access_token })

        return {
          ...token,
          accessToken: oauthResponse.jwt_token,
          accessTokenExp: moment().add(1, 'hour').valueOf(),
          refreshToken: oauthResponse.jwt_refresh_token,
        }
      }

      if (token.refreshToken && token.accessTokenExp && Date.now() > token.accessTokenExp) {
        try {
          const newTokens = await AuthService.refreshToken(token.refreshToken)

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
  }
}

declare module "next-auth/jwt" {
  interface JWT {
    accessToken?: string
    accessTokenExp?: number
    refreshToken?: string
  }
}
