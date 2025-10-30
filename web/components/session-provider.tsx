import { auth } from "auth"
import { SessionProvider } from "next-auth/react"
import SessionErrorHandler from "@/components/session-error-handler"

export default async function SessionProviderWrapper({ children }: { children: React.ReactNode }) {
  const session = await auth()
  if (session?.user) {
    session.user = {
      name: session.user.name,
      email: session.user.email,
      image: session.user.image,
    }
  }

  return (
    <SessionProvider basePath={"/auth"} session={session}>
      <SessionErrorHandler error={session?.error}>
        {children}
      </SessionErrorHandler>
    </SessionProvider>
  )
}
