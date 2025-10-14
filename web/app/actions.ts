"use server"

import { signIn, signOut } from "auth"

export async function handleSignIn(provider?: string, opts?: FormData | ({
    redirectTo?: string | undefined;
    redirect?: true | undefined;
} & Record<string, any>) | undefined) {
  await signIn(provider, opts)
}

export async function handleSignOut(opts?: {
    redirectTo?: string | undefined;
    redirect?: true | undefined;
} | undefined) {
  await signOut(opts)
}
