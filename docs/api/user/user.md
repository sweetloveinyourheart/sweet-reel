# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [user.proto](#user-proto)
    - [Channel](#com-sweetloveinyourheart-srl-user-Channel)
    - [GetChannelByHandleRequest](#com-sweetloveinyourheart-srl-user-GetChannelByHandleRequest)
    - [GetChannelByHandleResponse](#com-sweetloveinyourheart-srl-user-GetChannelByHandleResponse)
    - [GetChannelByIDRequest](#com-sweetloveinyourheart-srl-user-GetChannelByIDRequest)
    - [GetChannelByIDResponse](#com-sweetloveinyourheart-srl-user-GetChannelByIDResponse)
    - [GetChannelByUserRequest](#com-sweetloveinyourheart-srl-user-GetChannelByUserRequest)
    - [GetChannelByUserResponse](#com-sweetloveinyourheart-srl-user-GetChannelByUserResponse)
    - [GetUserByIDRequest](#com-sweetloveinyourheart-srl-user-GetUserByIDRequest)
    - [GetUserByIDResponse](#com-sweetloveinyourheart-srl-user-GetUserByIDResponse)
    - [UpsertOAuthUserRequest](#com-sweetloveinyourheart-srl-user-UpsertOAuthUserRequest)
    - [UpsertOAuthUserResponse](#com-sweetloveinyourheart-srl-user-UpsertOAuthUserResponse)
    - [User](#com-sweetloveinyourheart-srl-user-User)
  
    - [UserService](#com-sweetloveinyourheart-srl-user-UserService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="user-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## user.proto



<a name="com-sweetloveinyourheart-srl-user-Channel"></a>

### Channel



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| owner_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| handle | [string](#string) |  |  |
| description | [string](#string) |  |  |
| banner_url | [string](#string) |  |  |
| subscriber_count | [int32](#int32) |  |  |
| total_views | [int64](#int64) |  |  |
| total_videos | [int32](#int32) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-user-GetChannelByHandleRequest"></a>

### GetChannelByHandleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| handle | [string](#string) |  | e.g., &#34;@johndoe&#34; |






<a name="com-sweetloveinyourheart-srl-user-GetChannelByHandleResponse"></a>

### GetChannelByHandleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel | [Channel](#com-sweetloveinyourheart-srl-user-Channel) |  |  |
| owner | [User](#com-sweetloveinyourheart-srl-user-User) |  |  |






<a name="com-sweetloveinyourheart-srl-user-GetChannelByIDRequest"></a>

### GetChannelByIDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| video_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-user-GetChannelByIDResponse"></a>

### GetChannelByIDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel | [Channel](#com-sweetloveinyourheart-srl-user-Channel) |  |  |
| owner | [User](#com-sweetloveinyourheart-srl-user-User) |  |  |






<a name="com-sweetloveinyourheart-srl-user-GetChannelByUserRequest"></a>

### GetChannelByUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-user-GetChannelByUserResponse"></a>

### GetChannelByUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel | [Channel](#com-sweetloveinyourheart-srl-user-Channel) |  |  |
| owner | [User](#com-sweetloveinyourheart-srl-user-User) |  |  |






<a name="com-sweetloveinyourheart-srl-user-GetUserByIDRequest"></a>

### GetUserByIDRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-user-GetUserByIDResponse"></a>

### GetUserByIDResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-srl-user-User) |  |  |






<a name="com-sweetloveinyourheart-srl-user-UpsertOAuthUserRequest"></a>

### UpsertOAuthUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| provider | [string](#string) |  | e.g. &#34;google&#34;, &#34;github&#34; |
| provider_user_id | [string](#string) |  | ID from provider (e.g. Google sub) |
| email | [string](#string) |  |  |
| name | [string](#string) |  |  |
| picture | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-user-UpsertOAuthUserResponse"></a>

### UpsertOAuthUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-srl-user-User) |  |  |
| is_new_user | [bool](#bool) |  |  |






<a name="com-sweetloveinyourheart-srl-user-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| email | [string](#string) |  |  |
| name | [string](#string) |  |  |
| picture | [string](#string) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |





 

 

 


<a name="com-sweetloveinyourheart-srl-user-UserService"></a>

### UserService
The UserService manages user profiles and linked OAuth identities.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| UpsertOAuthUser | [UpsertOAuthUserRequest](#com-sweetloveinyourheart-srl-user-UpsertOAuthUserRequest) | [UpsertOAuthUserResponse](#com-sweetloveinyourheart-srl-user-UpsertOAuthUserResponse) | Called by AuthService after verifying an OAuth provider token. |
| GetUserByID | [GetUserByIDRequest](#com-sweetloveinyourheart-srl-user-GetUserByIDRequest) | [GetUserByIDResponse](#com-sweetloveinyourheart-srl-user-GetUserByIDResponse) | Fetch user info by ID (used internally by other services). |
| GetChannelByID | [GetChannelByIDRequest](#com-sweetloveinyourheart-srl-user-GetChannelByIDRequest) | [GetChannelByIDResponse](#com-sweetloveinyourheart-srl-user-GetChannelByIDResponse) | Fetch channel info by ID. |
| GetChannelByUser | [GetChannelByUserRequest](#com-sweetloveinyourheart-srl-user-GetChannelByUserRequest) | [GetChannelByUserResponse](#com-sweetloveinyourheart-srl-user-GetChannelByUserResponse) | Fetch channel info by user. |
| GetChannelByHandle | [GetChannelByHandleRequest](#com-sweetloveinyourheart-srl-user-GetChannelByHandleRequest) | [GetChannelByHandleResponse](#com-sweetloveinyourheart-srl-user-GetChannelByHandleResponse) | Fetch channel info by handle (e.g., @username). |

 



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

