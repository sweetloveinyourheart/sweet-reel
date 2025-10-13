package ffmpeg

import "time"

// ProgressCallback is a function type for progress monitoring
type ProgressCallback func(ProgressInfo)

// ProgressInfo contains information about the current processing progress
type ProgressInfo struct {
	Percentage float64       // Progress percentage (0-100)
	Duration   time.Duration // Total duration of the media
	Current    time.Duration // Current processing time
	Speed      string        // Processing speed (e.g., "2.34x")
	Bitrate    string        // Current bitrate (e.g., "1234.5kbits/s")
}

// TranscodeOptions contains options for video transcoding
type TranscodeOptions struct {
	// Video options
	VideoCodec   string // Video codec (e.g., "libx264", "libx265")
	VideoBitrate string // Video bitrate (e.g., "1000k", "2M")
	VideoQuality string // Video quality/CRF (e.g., "23")
	Resolution   string // Output resolution (e.g., "1920x1080", "1280x720")
	FrameRate    string // Frame rate (e.g., "30", "24")

	// Audio options
	AudioCodec      string // Audio codec (e.g., "aac", "mp3")
	AudioBitrate    string // Audio bitrate (e.g., "128k", "192k")
	AudioChannels   string // Number of audio channels (e.g., "2", "1")
	AudioSampleRate string // Audio sample rate (e.g., "44100", "48000")

	// Container options
	Format string // Output format/container (e.g., "mp4", "mkv", "webm")

	// Processing options
	StartTime string // Start time for trimming (e.g., "00:01:30")
	Duration  string // Duration for trimming (e.g., "00:02:00")

	// Hardware acceleration
	HWAccel string // Hardware acceleration (e.g., "cuda", "videotoolbox")

	// Custom options
	CustomArgs []string // Additional custom FFmpeg arguments
}

// SegmentationOptions contains options for video segmentation
type SegmentationOptions struct {
	QualityName string

	// Segment duration
	SegmentDuration string // Duration of each segment (e.g., "10", "30")

	// Playlist options
	PlaylistType string // HLS playlist type ("vod" or "live")
	PlaylistName string // Name of the playlist file (default: "playlist.m3u8")

	// Segment naming
	SegmentPrefix string // Prefix for segment files (default: "segment")
	SegmentFormat string // Segment file format (default: "ts")

	// Video options (inherited from transcoding)
	VideoCodec   string
	VideoBitrate string
	VideoQuality string
	Resolution   string
	FrameRate    string

	// Audio options (inherited from transcoding)
	AudioCodec      string
	AudioBitrate    string
	AudioChannels   string
	AudioSampleRate string

	// Encryption options
	EnableEncryption bool   // Enable AES-128 encryption
	KeyInfoFile      string // Path to key info file for encryption

	// Custom options
	CustomArgs []string // Additional custom FFmpeg arguments
}

// ProbeInfo contains media file information
type ProbeInfo struct {
	Format  FormatInfo   `json:"format"`
	Streams []StreamInfo `json:"streams"`
}

// FormatInfo contains format information
type FormatInfo struct {
	Filename       string `json:"filename"`
	NBStreams      int    `json:"nb_streams"`
	NBPrograms     int    `json:"nb_programs"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime      string `json:"start_time"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
}

// StreamInfo contains stream information
type StreamInfo struct {
	Index              int               `json:"index"`
	CodecName          string            `json:"codec_name"`
	CodecLongName      string            `json:"codec_long_name"`
	Profile            string            `json:"profile"`
	CodecType          string            `json:"codec_type"`
	CodecTimeBase      string            `json:"codec_time_base"`
	CodecTagString     string            `json:"codec_tag_string"`
	CodecTag           string            `json:"codec_tag"`
	Width              int               `json:"width,omitempty"`
	Height             int               `json:"height,omitempty"`
	CodedWidth         int               `json:"coded_width,omitempty"`
	CodedHeight        int               `json:"coded_height,omitempty"`
	HasBFrames         int               `json:"has_b_frames,omitempty"`
	SampleAspectRatio  string            `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio string            `json:"display_aspect_ratio,omitempty"`
	PixFmt             string            `json:"pix_fmt,omitempty"`
	Level              int               `json:"level,omitempty"`
	ChromaLocation     string            `json:"chroma_location,omitempty"`
	Refs               int               `json:"refs,omitempty"`
	RFrameRate         string            `json:"r_frame_rate,omitempty"`
	AvgFrameRate       string            `json:"avg_frame_rate,omitempty"`
	TimeBase           string            `json:"time_base"`
	StartPts           int64             `json:"start_pts"`
	StartTime          string            `json:"start_time"`
	Duration           string            `json:"duration,omitempty"`
	DurationTs         int64             `json:"duration_ts,omitempty"`
	BitRate            string            `json:"bit_rate,omitempty"`
	BitsPerRawSample   string            `json:"bits_per_raw_sample,omitempty"`
	NBFrames           string            `json:"nb_frames,omitempty"`
	SampleFmt          string            `json:"sample_fmt,omitempty"`
	SampleRate         string            `json:"sample_rate,omitempty"`
	Channels           int               `json:"channels,omitempty"`
	ChannelLayout      string            `json:"channel_layout,omitempty"`
	BitsPerSample      int               `json:"bits_per_sample,omitempty"`
	Tags               map[string]string `json:"tags,omitempty"`
}
