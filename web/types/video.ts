import * as z from "zod"; 
 
export const PresignedUrlRequest = z.object({ 
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
