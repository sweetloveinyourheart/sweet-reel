import { getServerApiClient } from "@/lib/api/server"
import { ChannelResponse } from "@/types/channel"
import { AccountProvider } from "@/contexts/account-context"
import { auth } from "@/auth"

export default async function AccountProviderWrapper({ children }: { children: React.ReactNode }) {
    const api = await getServerApiClient()
    const session = await auth()

    const channelInfo = session ? await api.get<ChannelResponse>("/channels") : null

    return (
        <AccountProvider channel={channelInfo}>
            {children}
        </AccountProvider>
    )
}
