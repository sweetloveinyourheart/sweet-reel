import * as z from "zod"; 
 
export const PresignedUrlRequest = z.object({
  channel_id: z.uuidv7(),
  title: z.string().min(3),
  description: z.string().min(3),
  file_name: z.string().nonempty()
});

export type PresignedUrlRequest = z.infer<typeof PresignedUrlRequest>;

export type PresignedUrlResponse = {
    video_id: string
    presigned_url: string
    expires_in: string
}

export type ChannelVideos = {
  video_id: string
  title: string
  thumbnail_url: string
  total_duration: number
  processed_at: number
}

export type GetChannelVideosResponse = {
  videos: ChannelVideos[]
}

export type UserMetadata = {
  email: string
  name: string
  picture: string
}

export type ChannelMetadata = {
  name: string
  handle: string
  onwer_metadata: UserMetadata
}

export type GetVideoMetadataResponse = {
  video_title?: string
  video_description?: string
  total_view?: number
  available_qualities?: string[]
  processed_at?: number
  channel_metadata: ChannelMetadata
}