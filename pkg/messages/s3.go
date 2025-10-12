package messages

import "time"

// S3EventMessage represents the top-level structure of MinIO S3 event notifications
type S3EventMessage struct {
	EventName string     `json:"EventName"`
	Key       string     `json:"Key"`
	Records   []S3Record `json:"Records"`
}

type S3Record struct {
	EventVersion      string              `json:"eventVersion"`
	EventSource       string              `json:"eventSource"`
	AWSRegion         string              `json:"awsRegion"`
	EventTime         time.Time           `json:"eventTime"`
	EventName         string              `json:"eventName"`
	UserIdentity      S3UserIdentity      `json:"userIdentity"`
	RequestParameters S3RequestParameters `json:"requestParameters"`
	ResponseElements  S3ResponseElements  `json:"responseElements"`
	S3                S3Details           `json:"s3"`
	Source            S3Source            `json:"source"`
}

type S3UserIdentity struct {
	PrincipalID string `json:"principalId"`
}

type S3RequestParameters struct {
	PrincipalID     string `json:"principalId"`
	Region          string `json:"region"`
	SourceIPAddress string `json:"sourceIPAddress"`
}

type S3ResponseElements struct {
	XAmzID2              string `json:"x-amz-id-2"`
	XAmzRequestID        string `json:"x-amz-request-id"`
	XMinioDeploymentID   string `json:"x-minio-deployment-id"`
	XMinioOriginEndpoint string `json:"x-minio-origin-endpoint"`
}

type S3Details struct {
	S3SchemaVersion string   `json:"s3SchemaVersion"`
	ConfigurationID string   `json:"configurationId"`
	Bucket          S3Bucket `json:"bucket"`
	Object          S3Object `json:"object"`
}

type S3Bucket struct {
	Name          string         `json:"name"`
	OwnerIdentity S3UserIdentity `json:"ownerIdentity"`
	ARN           string         `json:"arn"`
}

type S3Object struct {
	Key          string            `json:"key"`
	Size         int64             `json:"size"`
	ETag         string            `json:"eTag"`
	ContentType  string            `json:"contentType"`
	UserMetadata map[string]string `json:"userMetadata"`
	Sequencer    string            `json:"sequencer"`
}

type S3Source struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	UserAgent string `json:"userAgent"`
}
