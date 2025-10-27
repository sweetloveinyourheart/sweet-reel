import { auth, RefreshTokenError } from "auth"
import { SessionProvider, signOut } from "next-auth/react"

export default async function SessionProviderWrapper({ children }: { children: React.ReactNode }) {
  const session = await auth()
  if (session?.user) {
    session.user = {
      name: session.user.name,
      email: session.user.email,
      image: session.user.image,
    }
  }

  if (session?.error == RefreshTokenError) {
    signOut()
  }

  return (
    <SessionProvider basePath={"/auth"} session={session}>
      {children}
    </SessionProvider>
  )
}
