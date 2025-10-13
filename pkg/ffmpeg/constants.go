package ffmpeg

// Video codecs
const (
	CodecLibX264 = "libx264"
	CodecLibX265 = "libx265"
	CodecH264    = "h264"
	CodecH265    = "h265"
	CodecVP8     = "vp8"
	CodecVP9     = "vp9"
	CodecAV1     = "av1"
)

// Audio codecs
const (
	CodecAAC    = "aac"
	CodecMP3    = "mp3"
	CodecOpus   = "opus"
	CodecVorbis = "vorbis"
)

// Container formats
const (
	FormatMP4  = "mp4"
	FormatMKV  = "mkv"
	FormatWebM = "webm"
	FormatAVI  = "avi"
	FormatMOV  = "mov"
	FormatFLV  = "flv"
	FormatTS   = "ts"
	FormatHLS  = "hls"
	FormatDASH = "dash"
)

// File extensions
const (
	ExtMP4  = ".mp4"
	ExtMKV  = ".mkv"
	ExtWebM = ".webm"
	ExtAVI  = ".avi"
	ExtMOV  = ".mov"
	ExtFLV  = ".flv"
	ExtTS   = ".ts"
	ExtM3U8 = ".m3u8"
	ExtJPG  = ".jpg"
	ExtJPEG = ".jpeg"
	ExtPNG  = ".png"
)

// MIME types
const (
	MimeTypeMP4         = "video/mp4"
	MimeTypeMKV         = "video/x-matroska"
	MimeTypeWebM        = "video/webm"
	MimeTypeAVI         = "video/x-msvideo"
	MimeTypeMOV         = "video/quicktime"
	MimeTypeFLV         = "video/x-flv"
	MimeTypeTS          = "video/mp2t"
	MimeTypeM3U8        = "application/vnd.apple.mpegurl"
	MimeTypeJPEG        = "image/jpeg"
	MimeTypePNG         = "image/png"
	MimeTypeOctetStream = "application/octet-stream"
)

// HLS playlist settings
const (
	PlaylistTypeVOD   = "vod"
	PlaylistTypeLive  = "live"
	PlaylistTypeEvent = "event"

	DefaultPlaylistName    = "playlist.m3u8"
	MasterPlaylistName     = "master.m3u8"
	DefaultSegmentPrefix   = "segment"
	DefaultSegmentFormat   = "ts"
	DefaultSegmentDuration = "10"
)

// Quality presets
const (
	PresetUltraFast = "ultrafast"
	PresetSuperFast = "superfast"
	PresetVeryFast  = "veryfast"
	PresetFaster    = "faster"
	PresetFast      = "fast"
	PresetMedium    = "medium"
	PresetSlow      = "slow"
	PresetSlower    = "slower"
	PresetVerySlow  = "veryslow"
)

// CRF (Constant Rate Factor) quality values
const (
	CRFLossless   = "0"
	CRFExcellent  = "18"
	CRFGood       = "23"
	CRFAcceptable = "28"
	CRFPoor       = "35"
)

// Standard resolutions
const (
	Resolution240p  = "426x240"
	Resolution360p  = "640x360"
	Resolution480p  = "854x480"
	Resolution720p  = "1280x720"
	Resolution1080p = "1920x1080"
	Resolution1440p = "2560x1440"
	Resolution2160p = "3840x2160" // 4K
)

// Standard bitrates for video
const (
	VideoBitrate240p  = "300k"
	VideoBitrate360p  = "500k"
	VideoBitrate480p  = "800k"
	VideoBitrate720p  = "1400k"
	VideoBitrate1080p = "2800k"
	VideoBitrate1440p = "5000k"
	VideoBitrate2160p = "8000k"
)

// Standard bitrates for audio
const (
	AudioBitrate64k  = "64k"
	AudioBitrate96k  = "96k"
	AudioBitrate128k = "128k"
	AudioBitrate192k = "192k"
	AudioBitrate256k = "256k"
	AudioBitrate320k = "320k"
)

// Audio channels
const (
	AudioChannelsMono   = "1"
	AudioChannelsStereo = "2"
	AudioChannels5_1    = "6"
	AudioChannels7_1    = "8"
)

// Audio sample rates
const (
	AudioSampleRate8kHz  = "8000"
	AudioSampleRate16kHz = "16000"
	AudioSampleRate22kHz = "22050"
	AudioSampleRate44kHz = "44100"
	AudioSampleRate48kHz = "48000"
	AudioSampleRate96kHz = "96000"
)

// Pixel formats
const (
	PixFmtYUV420P = "yuv420p"
	PixFmtYUV422P = "yuv422p"
	PixFmtYUV444P = "yuv444p"
	PixFmtRGB24   = "rgb24"
	PixFmtRGBA    = "rgba"
)

// Hardware acceleration
const (
	HWAccelCUDA         = "cuda"
	HWAccelVideoToolbox = "videotoolbox"
	HWAccelQSV          = "qsv"
	HWAccelVAAPI        = "vaapi"
	HWAccelDXVA2        = "dxva2"
	HWAccelD3D11VA      = "d3d11va"
)

// DASH settings
const (
	DASHManifestName        = "manifest.mpd"
	DASHInitSegmentPattern  = "init_$RepresentationID$.m4s"
	DASHMediaSegmentPattern = "chunk_$RepresentationID$_$Number$.m4s"
)

// Watermark positions
const (
	WatermarkTopLeft     = "top-left"
	WatermarkTopRight    = "top-right"
	WatermarkBottomLeft  = "bottom-left"
	WatermarkBottomRight = "bottom-right"
	WatermarkCenter      = "center"
)

// FFmpeg options
const (
	OptionOverwrite       = "-y"
	OptionInput           = "-i"
	OptionVideoCodec      = "-c:v"
	OptionAudioCodec      = "-c:a"
	OptionVideoBitrate    = "-b:v"
	OptionAudioBitrate    = "-b:a"
	OptionFormat          = "-f"
	OptionCRF             = "-crf"
	OptionPreset          = "-preset"
	OptionScale           = "-s"
	OptionFrameRate       = "-r"
	OptionAudioChannels   = "-ac"
	OptionAudioSampleRate = "-ar"
	OptionStartTime       = "-ss"
	OptionDuration        = "-t"
	OptionVideoFilter     = "-vf"
	OptionCopy            = "copy"
	OptionDisableVideo    = "-vn"
	OptionDisableAudio    = "-an"
	OptionFrameCount      = "-vframes"
	OptionHWAccel         = "-hwaccel"
	OptionMovFlags        = "-movflags"
	OptionFastStart       = "+faststart"
)
