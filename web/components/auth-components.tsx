import { handleSignIn, handleSignOut } from "@/app/actions"
import { Button } from "./ui/button"

export function SignIn({
  provider,
  ...props
}: { provider?: string } & React.ComponentPropsWithRef<typeof Button>) {
  return (
    <form action={handleSignIn.bind(null, provider)}>
      <Button {...props}>Sign In</Button>
    </form>
  )
}

export function SignOut(props: React.ComponentPropsWithRef<typeof Button>) {
  return (
    <form action={handleSignOut.bind(null, {})} className="w-full">
      <Button variant="ghost" className="w-full p-0" {...props}>
        Sign Out
      </Button>
    </form>
  )
}
