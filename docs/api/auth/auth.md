# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [auth.proto](#auth-proto)
    - [OAuthLoginRequest](#com-sweetloveinyourheart-srl-auth-OAuthLoginRequest)
    - [OAuthLoginResponse](#com-sweetloveinyourheart-srl-auth-OAuthLoginResponse)
    - [RefreshTokenRequest](#com-sweetloveinyourheart-srl-auth-RefreshTokenRequest)
    - [RefreshTokenResponse](#com-sweetloveinyourheart-srl-auth-RefreshTokenResponse)
    - [User](#com-sweetloveinyourheart-srl-auth-User)
  
    - [AuthService](#com-sweetloveinyourheart-srl-auth-AuthService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="auth-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## auth.proto



<a name="com-sweetloveinyourheart-srl-auth-OAuthLoginRequest"></a>

### OAuthLoginRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| provider | [string](#string) |  | &#34;google&#34;, &#34;github&#34;, ... |
| access_token | [string](#string) |  | token from provider |






<a name="com-sweetloveinyourheart-srl-auth-OAuthLoginResponse"></a>

### OAuthLoginResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#com-sweetloveinyourheart-srl-auth-User) |  |  |
| jwt_token | [string](#string) |  |  |
| jwt_refresh_token | [string](#string) |  |  |
| is_new_user | [bool](#bool) |  |  |






<a name="com-sweetloveinyourheart-srl-auth-RefreshTokenRequest"></a>

### RefreshTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jwt_refresh_token | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-auth-RefreshTokenResponse"></a>

### RefreshTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| jwt_token | [string](#string) |  |  |






<a name="com-sweetloveinyourheart-srl-auth-User"></a>

### User



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| email | [string](#string) |  |  |
| name | [string](#string) |  |  |
| picture | [string](#string) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |





 

 

 


<a name="com-sweetloveinyourheart-srl-auth-AuthService"></a>

### AuthService
AuthService handles all authentication-related logic.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| OAuthLogin | [OAuthLoginRequest](#com-sweetloveinyourheart-srl-auth-OAuthLoginRequest) | [OAuthLoginResponse](#com-sweetloveinyourheart-srl-auth-OAuthLoginResponse) | Handles OAuth login with external providers (Google, GitHub, etc.) |
| RefreshToken | [RefreshTokenRequest](#com-sweetloveinyourheart-srl-auth-RefreshTokenRequest) | [RefreshTokenResponse](#com-sweetloveinyourheart-srl-auth-RefreshTokenResponse) | Handle refresh tokens. |

 



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

