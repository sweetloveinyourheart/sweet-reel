# FFmpeg Go Wrapper

A comprehensive Go wrapper for FFmpeg that provides transcoding, segmentation, and video processing capabilities.

## Features

- **Video Transcoding**: Convert videos between different formats, codecs, and quality settings
- **HLS Segmentation**: Create HLS segments for adaptive streaming with multiple quality levels
- **DASH Support**: Generate DASH (Dynamic Adaptive Streaming over HTTP) segments
- **Video Utilities**: Extract audio, create thumbnails, trim videos, concatenate files
- **Progress Monitoring**: Real-time progress tracking for long-running operations
- **File Probing**: Extract detailed media file information using ffprobe
- **Hardware Acceleration**: Support for hardware-accelerated encoding (CUDA, VideoToolbox, etc.)

## Prerequisites

- FFmpeg must be installed and available in PATH
- Both `ffmpeg` and `ffprobe` binaries are required

### Installing FFmpeg

**macOS (with Homebrew):**
```bash
brew install ffmpeg
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install ffmpeg
```

**CentOS/RHEL:**
```bash
sudo yum install epel-release
sudo yum install ffmpeg
```

## Usage

### Basic Setup

```go
package main

import (
    "context"
    "log"
    
    "github.com/sweetloveinyourheart/sweet-reel/pkg/ffmpeg"
    "go.uber.org/zap"
)

func main() {
    logger, _ := zap.NewDevelopment()
    ff := ffmpeg.New(logger)
    
    ctx := context.Background()
    
    // Check if FFmpeg is available
    if err := ff.IsAvailable(ctx); err != nil {
        log.Fatal("FFmpeg not available:", err)
    }
    
    // Get version
    version, _ := ff.GetVersion(ctx)
    logger.Info("FFmpeg version", zap.String("version", version))
}
```

### Video Transcoding

```go
// Basic transcoding
options := ffmpeg.TranscodeOptions{
    VideoCodec:   "libx264",
    VideoQuality: "23",
    AudioCodec:   "aac",
    AudioBitrate: "128k",
    Resolution:   "1280x720",
    Format:       "mp4",
}

progressCallback := func(progress ffmpeg.ProgressInfo) {
    fmt.Printf("Progress: %.2f%% (Speed: %s)\n", 
        progress.Percentage, progress.Speed)
}

err := ff.Transcode(ctx, "input.mp4", "output.mp4", options, progressCallback)
```

### HLS Segmentation

```go
// Single quality HLS
options := ffmpeg.SegmentationOptions{
    SegmentDuration: "10",
    PlaylistType:    "vod",
    PlaylistName:    "playlist.m3u8",
    VideoCodec:      "libx264",
    VideoQuality:    "23",
    AudioCodec:      "aac",
    AudioBitrate:    "128k",
}

err := ff.SegmentVideo(ctx, "input.mp4", "output/hls", options, progressCallback)
```

### Multi-Quality Adaptive Streaming

```go
// Multiple quality levels
qualities := []ffmpeg.SegmentationOptions{
    {
        SegmentDuration: "10",
        PlaylistType:    "vod",
        VideoCodec:      "libx264",
        VideoBitrate:    "800k",
        AudioBitrate:    "96k",
        Resolution:      "640x360", // 360p
    },
    {
        SegmentDuration: "10",
        PlaylistType:    "vod",
        VideoCodec:      "libx264",
        VideoBitrate:    "1400k",
        AudioBitrate:    "128k",
        Resolution:      "1280x720", // 720p
    },
    {
        SegmentDuration: "10",
        PlaylistType:    "vod",
        VideoCodec:      "libx264",
        VideoBitrate:    "2800k",
        AudioBitrate:    "192k",
        Resolution:      "1920x1080", // 1080p
    },
}

err := ff.SegmentVideoMultiQuality(ctx, "input.mp4", "output/adaptive", qualities, progressCallback)
```

### File Probing

```go
// Get detailed file information
probeInfo, err := ff.ProbeFile(ctx, "input.mp4")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Format: %s\n", probeInfo.Format.FormatName)
fmt.Printf("Duration: %s\n", probeInfo.Format.Duration)
fmt.Printf("Streams: %d\n", len(probeInfo.Streams))

// Get video-specific information
videoInfo, err := ff.GetVideoInfo(ctx, "input.mp4")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Codec: %s\n", videoInfo.CodecName)
fmt.Printf("Resolution: %dx%d\n", videoInfo.Width, videoInfo.Height)
```

### Utility Functions

```go
// Create thumbnail
err := ff.CreateThumbnail(ctx, "input.mp4", "thumbnail.jpg", "00:00:10", 320, 240)

// Extract audio
audioOptions := ffmpeg.TranscodeOptions{
    AudioCodec:   "mp3",
    AudioBitrate: "192k",
}
err := ff.ExtractAudio(ctx, "input.mp4", "audio.mp3", audioOptions, nil)

// Trim video
err := ff.TrimVideo(ctx, "input.mp4", "trimmed.mp4", "00:01:00", "00:02:00", nil)

// Concatenate videos
inputFiles := []string{"part1.mp4", "part2.mp4", "part3.mp4"}
err := ff.ConcatenateVideos(ctx, inputFiles, "combined.mp4", progressCallback)
```

## Configuration Options

### TranscodeOptions

| Field | Description | Example |
|-------|-------------|---------|
| VideoCodec | Video codec | "libx264", "libx265", "libvpx" |
| VideoBitrate | Video bitrate | "1000k", "2M" |
| VideoQuality | CRF value (0-51) | "23", "18" |
| Resolution | Output resolution | "1920x1080", "1280x720" |
| FrameRate | Frame rate | "30", "24", "60" |
| AudioCodec | Audio codec | "aac", "mp3", "opus" |
| AudioBitrate | Audio bitrate | "128k", "192k" |
| AudioChannels | Audio channels | "2", "1" |
| AudioSampleRate | Sample rate | "44100", "48000" |
| Format | Container format | "mp4", "mkv", "webm" |
| StartTime | Start time for trimming | "00:01:30" |
| Duration | Duration for trimming | "00:02:00" |
| HWAccel | Hardware acceleration | "cuda", "videotoolbox" |

### SegmentationOptions

| Field | Description | Example |
|-------|-------------|---------|
| SegmentDuration | Segment length in seconds | "10", "6" |
| PlaylistType | HLS playlist type | "vod", "live" |
| PlaylistName | Playlist filename | "playlist.m3u8" |
| SegmentPrefix | Segment file prefix | "segment" |
| SegmentFormat | Segment file format | "ts", "mp4" |
| EnableEncryption | Enable AES-128 encryption | true, false |
| KeyInfoFile | Path to key info file | "/path/to/key.info" |

## Progress Monitoring

The wrapper provides real-time progress information through callback functions:

```go
progressCallback := func(progress ffmpeg.ProgressInfo) {
    fmt.Printf("Progress: %.2f%%\n", progress.Percentage)
    fmt.Printf("Current: %v / Total: %v\n", progress.Current, progress.Duration)
    fmt.Printf("Speed: %s\n", progress.Speed)
    fmt.Printf("Bitrate: %s\n", progress.Bitrate)
}
```

## Error Handling

All functions return detailed errors that can be checked and handled:

```go
if err := ff.Transcode(ctx, input, output, options, nil); err != nil {
    if strings.Contains(err.Error(), "No such file") {
        log.Printf("Input file not found: %v", err)
    } else if strings.Contains(err.Error(), "Invalid data") {
        log.Printf("Invalid input format: %v", err)
    } else {
        log.Printf("Transcoding failed: %v", err)
    }
}
```

## Hardware Acceleration

Enable hardware acceleration for faster processing:

```go
options := ffmpeg.TranscodeOptions{
    HWAccel:      "cuda",        // NVIDIA GPU acceleration
    VideoCodec:   "h264_nvenc",  // Hardware encoder
    VideoBitrate: "2000k",
    // ... other options
}
```

Supported hardware acceleration:
- **CUDA** (NVIDIA): `"cuda"`
- **VideoToolbox** (macOS): `"videotoolbox"`
- **VAAPI** (Intel): `"vaapi"`
- **QSV** (Intel Quick Sync): `"qsv"`

## Best Practices

1. **Always check FFmpeg availability** before using the wrapper
2. **Use progress callbacks** for long-running operations
3. **Validate input files** before processing
4. **Handle context cancellation** for graceful shutdowns
5. **Use appropriate quality settings** based on your use case
6. **Consider hardware acceleration** for production workloads
7. **Monitor resource usage** during batch processing

## Performance Tips

- Use CRF (Constant Rate Factor) for quality-based encoding
- Enable hardware acceleration when available
- Use multi-pass encoding for optimal quality/size ratio
- Segment videos for adaptive streaming in production
- Batch process multiple files efficiently

## Integration Example

For integration with the sweet-reel video processing service:

```go
func processVideoSplitter(ctx context.Context, msg *models.VideoSplitterMessage, storageClient storage.Storage) error {
    logger := logger.Global()
    ff := ffmpeg.New(logger)
    
    // Download video from storage
    videoData, err := storageClient.Download(msg.Metadata.Key, msg.Metadata.Bucket)
    if err != nil {
        return err
    }
    
    // Save to temporary file
    tmpFile := "/tmp/" + msg.VideoID.String() + ".mp4"
    if err := os.WriteFile(tmpFile, videoData, 0644); err != nil {
        return err
    }
    defer os.Remove(tmpFile)
    
    // Create HLS segments
    outputDir := "/tmp/hls_" + msg.VideoID.String()
    defer os.RemoveAll(outputDir)
    
    qualities := []ffmpeg.SegmentationOptions{
        // Define quality levels...
    }
    
    if err := ff.SegmentVideoMultiQuality(ctx, tmpFile, outputDir, qualities, nil); err != nil {
        return err
    }
    
    // Upload segments back to storage
    // ... upload logic
    
    return nil
}
```