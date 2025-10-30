# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [video_management.proto](#video_management-proto)
    - [ChannelVideo](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-ChannelVideo)
    - [GetChannelVideosRequest](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetChannelVideosRequest)
    - [GetChannelVideosResponse](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetChannelVideosResponse)
    - [GetVideoMetadataByIdRequest](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetVideoMetadataByIdRequest)
    - [GetVideoMetadataByIdResponse](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetVideoMetadataByIdResponse)
    - [PresignedUrlRequest](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-PresignedUrlRequest)
    - [PresignedUrlResponse](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-PresignedUrlResponse)
    - [ServePlaylistRequest](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistRequest)
    - [ServePlaylistResponse](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistResponse)
    - [ServePlaylistVariant](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistVariant)
  
    - [VideoManagement](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-VideoManagement)
  
- [Scalar Value Types](#scalar-value-types)



<a name="video_management-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## video_management.proto



<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-ChannelVideo"></a>

### ChannelVideo



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| video_id | [string](#string) |  |  |
| video_title | [string](#string) |  |  |
| thumbnail_url | [string](#string) |  |  |
| total_duration | [int32](#int32) |  |  |
| total_view | [int64](#int64) |  |  |
| processed_at | [int64](#int64) |  |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetChannelVideosRequest"></a>

### GetChannelVideosRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetChannelVideosResponse"></a>

### GetChannelVideosResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| videos | [ChannelVideo](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-ChannelVideo) | repeated |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetVideoMetadataByIdRequest"></a>

### GetVideoMetadataByIdRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| video_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetVideoMetadataByIdResponse"></a>

### GetVideoMetadataByIdResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| video_id | [string](#string) |  |  |
| channel_id | [string](#string) |  |  |
| video_title | [string](#string) |  |  |
| video_description | [string](#string) |  |  |
| total_view | [int64](#int64) |  |  |
| available_qualities | [string](#string) | repeated |  |
| processed_at | [int64](#int64) |  |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-PresignedUrlRequest"></a>

### PresignedUrlRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [string](#string) |  |  |
| description | [string](#string) |  |  |
| file_name | [string](#string) |  |  |
| uploader_id | [string](#string) |  |  |
| channel_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-PresignedUrlResponse"></a>

### PresignedUrlResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| video_id | [string](#string) |  |  |
| presigned_url | [string](#string) |  |  |
| expires_in | [int32](#int32) |  |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistRequest"></a>

### ServePlaylistRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| video_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistResponse"></a>

### ServePlaylistResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| playlist_url | [string](#string) |  |  |
| variants | [ServePlaylistVariant](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistVariant) | repeated |  |






<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistVariant"></a>

### ServePlaylistVariant



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| quality | [string](#string) |  |  |
| playlist_url | [string](#string) |  |  |
| segment_urls | [string](#string) | repeated |  |





 

 

 


<a name="com-sweetloveinyourheart-srl-videomanagement-dataproviders-VideoManagement"></a>

### VideoManagement


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| PresignedUrl | [PresignedUrlRequest](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-PresignedUrlRequest) | [PresignedUrlResponse](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-PresignedUrlResponse) |  |
| GetChannelVideos | [GetChannelVideosRequest](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetChannelVideosRequest) | [GetChannelVideosResponse](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetChannelVideosResponse) |  |
| GetVideoMetadataById | [GetVideoMetadataByIdRequest](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetVideoMetadataByIdRequest) | [GetVideoMetadataByIdResponse](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-GetVideoMetadataByIdResponse) |  |
| ServePlaylist | [ServePlaylistRequest](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistRequest) | [ServePlaylistResponse](#com-sweetloveinyourheart-srl-videomanagement-dataproviders-ServePlaylistResponse) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

