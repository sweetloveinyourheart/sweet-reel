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