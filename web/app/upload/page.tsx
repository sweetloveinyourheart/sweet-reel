"use client"

import { useState, useRef, DragEvent, ChangeEvent } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Upload, Video, X, CheckCircle2, AlertCircle } from "lucide-react"
import { cn } from "@/lib/utils"
import { validateFileSize, validateFileType } from "@/lib/s3"
import { useVideoUpload } from "@/hooks"
import { PresignedUrlRequest } from "@/types"

interface UploadedFile {
  file: File
  preview?: string
  progress: number
  status: "pending" | "uploading" | "success" | "error"
  error?: string
}

const maxSize = 500 * 1024 * 1024 // 500MB
const allowedTypes = ["video/mp4", "video/webm", "video/ogg", "video/quicktime"]

export default function VideoUploadPage() {
  const { uploading, progress, error, uploadVideo, reset } = useVideoUpload()

  const [uploadedFile, setUploadedFile] = useState<UploadedFile | null>(null)
  const [isDragging, setIsDragging] = useState(false)
  const [title, setTitle] = useState("")
  const [description, setDescription] = useState("")
  const fileInputRef = useRef<HTMLInputElement>(null)

  const validateFile = (file: File): string | null => {
    if (!validateFileType(file, allowedTypes)) {
      return "Invalid file type. Please upload MP4, WebM, OGG, or MOV files."
    }

    if (!validateFileSize(file, maxSize)) {
      return "File size exceeds 500MB limit."
    }

    return null
  }

  const handleFiles = (files: FileList | null) => {
    if (!files || files.length === 0) return

    const file = files[0]
    const error = validateFile(file)

    // Clean up previous preview if exists
    if (uploadedFile?.preview) {
      URL.revokeObjectURL(uploadedFile.preview)
    }

    setUploadedFile({
      file,
      preview: file.type.startsWith("video/") ? URL.createObjectURL(file) : undefined,
      progress: 0,
      status: error ? "error" : "pending",
      error: error || undefined,
    })
  }

  const handleDragEnter = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }

  const handleDragLeave = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)
  }

  const handleDragOver = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
  }

  const handleDrop = (e: DragEvent<HTMLDivElement>) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)

    const files = e.dataTransfer.files
    handleFiles(files)
  }

  const handleFileInputChange = (e: ChangeEvent<HTMLInputElement>) => {
    handleFiles(e.target.files)
  }

  const removeFile = () => {
    if (uploadedFile?.preview) {
      URL.revokeObjectURL(uploadedFile.preview)
    }
    setUploadedFile(null)
    setTitle("")
    setDescription("")
  }

  const simulateUpload = () => {
    if (!uploadedFile) return

    setUploadedFile((prev) => prev ? { ...prev, status: "uploading" } : null)

    const interval = setInterval(() => {
      setUploadedFile((prev) => {
        if (!prev) return null
        if (prev.progress < 100) {
          return { ...prev, progress: prev.progress + 10 }
        } else {
          clearInterval(interval)
          return { ...prev, status: "success" }
        }
      })
    }, 200)
  }

  const handleUpload = () => {
    try {
      if (uploadedFile && uploadedFile.status === "pending") {
        const fileMetadata = PresignedUrlRequest.parse({ 
          title, 
          description, 
          file_name: uploadedFile.file.name 
        })
      }

      // TODO: handle upload video with the useVideoUpload() hook
    } catch (error) {
      // TODO: handle errors
    }
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return "0 Bytes"
    const k = 1024
    const sizes = ["Bytes", "KB", "MB", "GB"]
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + " " + sizes[i]
  }

  return (
    <div className="w-full max-w-5xl mx-auto p-6">
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Upload Video</h1>
        <p className="text-muted-foreground">
          Share your videos with the world. Upload MP4, WebM, OGG, or MOV files (max 500MB).
        </p>
      </div>

      <div className="space-y-6">
        {/* Drag and Drop Zone */}
        <div
          onDragEnter={handleDragEnter}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          className={cn(
            "border-2 border-dashed rounded-lg p-12 text-center transition-colors",
            isDragging
              ? "border-primary bg-primary/5"
              : "border-muted-foreground/25 hover:border-muted-foreground/50"
          )}
        >
          <input
            ref={fileInputRef}
            type="file"
            accept="video/*"
            onChange={handleFileInputChange}
            className="hidden"
          />
          <div className="flex flex-col items-center gap-4">
            <div className="rounded-full bg-primary/10 p-6">
              <Upload className="h-12 w-12 text-primary" />
            </div>
            <div>
              <p className="text-lg font-medium mb-1">
                Drag and drop a video file here
              </p>
              <p className="text-sm text-muted-foreground mb-4">
                or click the button below to browse
              </p>
            </div>
            <Button
              onClick={() => fileInputRef.current?.click()}
              size="lg"
            >
              Select File
            </Button>
          </div>
        </div>

        {/* Video Details Form */}
        {uploadedFile && (
          <div className="space-y-4 border rounded-lg p-6">
            <h2 className="text-xl font-semibold mb-4">Video Details</h2>
            <div>
              <label htmlFor="title" className="block text-sm font-medium mb-2">
                Title
              </label>
              <Input
                id="title"
                placeholder="Enter video title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="description" className="block text-sm font-medium mb-2">
                Description
              </label>
              <textarea
                id="description"
                placeholder="Tell viewers about your video"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              />
            </div>
          </div>
        )}

        {/* Uploaded File Preview */}
        {uploadedFile && (
          <div className="space-y-4">
            <h2 className="text-xl font-semibold">File to Upload</h2>
            <div className="border rounded-lg p-4 flex items-start gap-4">
              <div className="flex-shrink-0">
                {uploadedFile.preview ? (
                  <video
                    src={uploadedFile.preview}
                    className="w-24 h-16 rounded object-cover bg-muted"
                  />
                ) : (
                  <div className="w-24 h-16 rounded bg-muted flex items-center justify-center">
                    <Video className="h-8 w-8 text-muted-foreground" />
                  </div>
                )}
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-start justify-between gap-2 mb-2">
                  <div className="flex-1 min-w-0">
                    <p className="font-medium truncate">
                      {uploadedFile.file.name}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      {formatFileSize(uploadedFile.file.size)}
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    {uploadedFile.status === "success" && (
                      <CheckCircle2 className="h-5 w-5 text-green-600" />
                    )}
                    {uploadedFile.status === "error" && (
                      <AlertCircle className="h-5 w-5 text-destructive" />
                    )}
                    <button
                      onClick={removeFile}
                      className="text-muted-foreground hover:text-foreground"
                    >
                      <X className="h-5 w-5" />
                    </button>
                  </div>
                </div>
                {uploadedFile.status === "error" && uploadedFile.error && (
                  <p className="text-sm text-destructive mb-2">
                    {uploadedFile.error}
                  </p>
                )}
                {(uploadedFile.status === "uploading" ||
                  uploadedFile.status === "success") && (
                    <div className="space-y-1">
                      <div className="w-full bg-muted rounded-full h-2 overflow-hidden">
                        <div
                          className="bg-primary h-full transition-all duration-300"
                          style={{ width: `${uploadedFile.progress}%` }}
                        />
                      </div>
                      <p className="text-xs text-muted-foreground">
                        {uploadedFile.progress}% uploaded
                      </p>
                    </div>
                  )}
              </div>
            </div>

            <div className="flex gap-3 pt-4">
              <Button
                onClick={handleUpload}
                disabled={
                  uploadedFile.status === "uploading" ||
                  uploadedFile.status === "success" ||
                  uploadedFile.status === "error"
                }
                className="flex-1"
                size="lg"
              >
                Upload Video
              </Button>
              <Button
                onClick={removeFile}
                variant="outline"
                size="lg"
              >
                Clear
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
