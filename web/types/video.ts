export type PresignedUrlRequest = {
    title: string
    description: string
    file_name: string
}

export type PresignedUrlResponse = {
    video_id: string
    presigned_url: string
    expires_in: string
}
