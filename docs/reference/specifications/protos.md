<!-- WARNING: THIS DOC IS AUTO-GENERATED. DO NOT EDIT! -->
# Protocol Documentation
<a name="top"></a>


## schema/v1/schema.proto
Flag evaluation API

This proto forms the basis of a flag-evaluation API.
It supports single and bulk evaluation RPCs, and flags of various types, as well as establishing a stream for getting notifications about changes in a flag definition.
It supports the inclusion of a &#34;context&#34; with each evaluation, which may contain arbitrary attributes relevant to flag evaluation.


<a name="schema-v1-AnyFlag"></a>

### AnyFlag
A variant type flag response.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reason | [string](#string) |  | The reason for the given return value, see https://openfeature.dev/docs/specification/types#resolution-details |
| variant | [string](#string) |  | The variant name of the returned flag value. |
| bool_value | [bool](#bool) |  |  |
| string_value | [string](#string) |  |  |
| double_value | [double](#double) |  |  |
| object_value | [google.protobuf.Struct](#google-protobuf-Struct) |  |  |






<a name="schema-v1-EventStreamRequest"></a>

### EventStreamRequest
Empty stream request body






<a name="schema-v1-EventStreamResponse"></a>

### EventStreamResponse
Response body for the EventStream stream response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [string](#string) |  | String key indicating the type of event that is being received, for example, provider_ready or configuration_change |
| data | [google.protobuf.Struct](#google-protobuf-Struct) |  | Object structure for use when sending relevant metadata to provide context to the event. Can be left unset when it is not required. |






<a name="schema-v1-ResolveAllRequest"></a>

### ResolveAllRequest
Request body for bulk flag evaluation, used by the ResolveAll rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| context | [google.protobuf.Struct](#google-protobuf-Struct) |  | Object structure describing the EvaluationContext used in the flag evaluation, see https://openfeature.dev/docs/reference/concepts/evaluation-context |






<a name="schema-v1-ResolveAllResponse"></a>

### ResolveAllResponse
Response body for bulk flag evaluation, used by the ResolveAll rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flags | [ResolveAllResponse.FlagsEntry](#schema-v1-ResolveAllResponse-FlagsEntry) | repeated | Object structure describing the evaluated flags for the provided context. |






<a name="schema-v1-ResolveAllResponse-FlagsEntry"></a>

### ResolveAllResponse.FlagsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [AnyFlag](#schema-v1-AnyFlag) |  |  |






<a name="schema-v1-ResolveBooleanRequest"></a>

### ResolveBooleanRequest
Request body for boolean flag evaluation, used by the ResolveBoolean rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flag_key | [string](#string) |  | Flag key of the requested flag. |
| context | [google.protobuf.Struct](#google-protobuf-Struct) |  | Object structure describing the EvaluationContext used in the flag evaluation, see https://openfeature.dev/docs/reference/concepts/evaluation-context |






<a name="schema-v1-ResolveBooleanResponse"></a>

### ResolveBooleanResponse
Response body for boolean flag evaluation. used by the ResolveBoolean rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [bool](#bool) |  | The response value of the boolean flag evaluation, will be unset in the case of error. |
| reason | [string](#string) |  | The reason for the given return value, see https://openfeature.dev/docs/specification/types#resolution-details |
| variant | [string](#string) |  | The variant name of the returned flag value. |
| metadata | [google.protobuf.Struct](#google-protobuf-Struct) |  | Metadata for this evaluation |






<a name="schema-v1-ResolveFloatRequest"></a>

### ResolveFloatRequest
Request body for float flag evaluation, used by the ResolveFloat rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flag_key | [string](#string) |  | Flag key of the requested flag. |
| context | [google.protobuf.Struct](#google-protobuf-Struct) |  | Object structure describing the EvaluationContext used in the flag evaluation, see https://openfeature.dev/docs/reference/concepts/evaluation-context |






<a name="schema-v1-ResolveFloatResponse"></a>

### ResolveFloatResponse
Response body for float flag evaluation. used by the ResolveFloat rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [double](#double) |  | The response value of the float flag evaluation, will be empty in the case of error. |
| reason | [string](#string) |  | The reason for the given return value, see https://openfeature.dev/docs/specification/types#resolution-details |
| variant | [string](#string) |  | The variant name of the returned flag value. |
| metadata | [google.protobuf.Struct](#google-protobuf-Struct) |  | Metadata for this evaluation |






<a name="schema-v1-ResolveIntRequest"></a>

### ResolveIntRequest
Request body for int flag evaluation, used by the ResolveInt rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flag_key | [string](#string) |  | Flag key of the requested flag. |
| context | [google.protobuf.Struct](#google-protobuf-Struct) |  | Object structure describing the EvaluationContext used in the flag evaluation, see https://openfeature.dev/docs/reference/concepts/evaluation-context |






<a name="schema-v1-ResolveIntResponse"></a>

### ResolveIntResponse
Response body for int flag evaluation. used by the ResolveInt rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [int64](#int64) |  | The response value of the int flag evaluation, will be unset in the case of error. |
| reason | [string](#string) |  | The reason for the given return value, see https://openfeature.dev/docs/specification/types#resolution-details |
| variant | [string](#string) |  | The variant name of the returned flag value. |
| metadata | [google.protobuf.Struct](#google-protobuf-Struct) |  | Metadata for this evaluation |






<a name="schema-v1-ResolveObjectRequest"></a>

### ResolveObjectRequest
Request body for object flag evaluation, used by the ResolveObject rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flag_key | [string](#string) |  | Flag key of the requested flag. |
| context | [google.protobuf.Struct](#google-protobuf-Struct) |  | Object structure describing the EvaluationContext used in the flag evaluation, see https://openfeature.dev/docs/reference/concepts/evaluation-context |






<a name="schema-v1-ResolveObjectResponse"></a>

### ResolveObjectResponse
Response body for object flag evaluation. used by the ResolveObject rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [google.protobuf.Struct](#google-protobuf-Struct) |  | The response value of the object flag evaluation, will be unset in the case of error.

NOTE: This structure will need to be decoded from google/protobuf/struct.proto before it is returned to the SDK |
| reason | [string](#string) |  | The reason for the given return value, see https://openfeature.dev/docs/specification/types#resolution-details |
| variant | [string](#string) |  | The variant name of the returned flag value. |
| metadata | [google.protobuf.Struct](#google-protobuf-Struct) |  | Metadata for this evaluation |






<a name="schema-v1-ResolveStringRequest"></a>

### ResolveStringRequest
Request body for string flag evaluation, used by the ResolveString rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flag_key | [string](#string) |  | Flag key of the requested flag. |
| context | [google.protobuf.Struct](#google-protobuf-Struct) |  | Object structure describing the EvaluationContext used in the flag evaluation, see https://openfeature.dev/docs/reference/concepts/evaluation-context |






<a name="schema-v1-ResolveStringResponse"></a>

### ResolveStringResponse
Response body for string flag evaluation. used by the ResolveString rpc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  | The response value of the string flag evaluation, will be unset in the case of error. |
| reason | [string](#string) |  | The reason for the given return value, see https://openfeature.dev/docs/specification/types#resolution-details |
| variant | [string](#string) |  | The variant name of the returned flag value. |
| metadata | [google.protobuf.Struct](#google-protobuf-Struct) |  | Metadata for this evaluation |





 

 

 


<a name="schema-v1-Service"></a>

### Service
Service defines the exposed rpcs of flagd

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ResolveAll | [ResolveAllRequest](#schema-v1-ResolveAllRequest) | [ResolveAllResponse](#schema-v1-ResolveAllResponse) |  |
| ResolveBoolean | [ResolveBooleanRequest](#schema-v1-ResolveBooleanRequest) | [ResolveBooleanResponse](#schema-v1-ResolveBooleanResponse) |  |
| ResolveString | [ResolveStringRequest](#schema-v1-ResolveStringRequest) | [ResolveStringResponse](#schema-v1-ResolveStringResponse) |  |
| ResolveFloat | [ResolveFloatRequest](#schema-v1-ResolveFloatRequest) | [ResolveFloatResponse](#schema-v1-ResolveFloatResponse) |  |
| ResolveInt | [ResolveIntRequest](#schema-v1-ResolveIntRequest) | [ResolveIntResponse](#schema-v1-ResolveIntResponse) |  |
| ResolveObject | [ResolveObjectRequest](#schema-v1-ResolveObjectRequest) | [ResolveObjectResponse](#schema-v1-ResolveObjectResponse) |  |
| EventStream | [EventStreamRequest](#schema-v1-EventStreamRequest) | [EventStreamResponse](#schema-v1-EventStreamResponse) stream |  |

 



<a name="sync_v1_sync_service-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## sync/v1/sync_service.proto
Flag definition sync API

This proto defines a simple API to synchronize a feature flag definition.
It supports establishing a stream for getting notifications about changes in a flag definition.


<a name="sync-v1-FetchAllFlagsRequest"></a>

### FetchAllFlagsRequest
FetchAllFlagsRequest is the request to fetch all flags. Flagd sends this request as the client in order to resync its internal state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| provider_id | [string](#string) |  | Optional: A unique identifier for flagd(grpc client) initiating the request. The server implementations may utilize this identifier to uniquely identify, validate(ex:- enforce authentication/authorization) and filter flag configurations that it can expose to this request. This field is intended to be optional. However server implementations may enforce it. ex:- provider_id: flagd-weatherapp-sidecar |
| selector | [string](#string) |  | Optional: A selector for the flag configuration request. The server implementation may utilize this to select flag configurations from a collection, select the source of the flag or combine this to any desired underlying filtering mechanism. ex:- selector: &#39;source=database,app=weatherapp&#39; |






<a name="sync-v1-FetchAllFlagsResponse"></a>

### FetchAllFlagsResponse
FetchAllFlagsResponse is the server response containing feature flag configurations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flag_configuration | [string](#string) |  | flagd feature flag configuration. Must be validated to schema - https://raw.githubusercontent.com/open-feature/schemas/main/json/flagd-definitions.json |






<a name="sync-v1-SyncFlagsRequest"></a>

### SyncFlagsRequest
SyncFlagsRequest is the request initiating the sever-streaming rpc. Flagd sends this request, acting as the client


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| provider_id | [string](#string) |  | Optional: A unique identifier for flagd(grpc client) initiating the request. The server implementations may utilize this identifier to uniquely identify, validate(ex:- enforce authentication/authorization) and filter flag configurations that it can expose to this request. This field is intended to be optional. However server implementations may enforce it. ex:- provider_id: flagd-weatherapp-sidecar |
| selector | [string](#string) |  | Optional: A selector for the flag configuration request. The server implementation may utilize this to select flag configurations from a collection, select the source of the flag or combine this to any desired underlying filtering mechanism. ex:- selector: &#39;source=database,app=weatherapp&#39; |






<a name="sync-v1-SyncFlagsResponse"></a>

### SyncFlagsResponse
SyncFlagsResponse is the server response containing feature flag configurations and the state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| flag_configuration | [string](#string) |  | flagd feature flag configuration. Must be validated to schema - https://raw.githubusercontent.com/open-feature/schemas/main/json/flagd-definitions.json |
| state | [SyncState](#sync-v1-SyncState) |  | State conveying the operation to be performed by flagd. See the descriptions of SyncState for an explanation of supported values |





 


<a name="sync-v1-SyncState"></a>

### SyncState
SyncState conveys the state of the payload. These states are related to flagd isync.go type definitions but
contains extras to optimize grpc use case. Refer - https://github.com/open-feature/flagd/blob/main/pkg/sync/isync.go

| Name | Number | Description |
| ---- | ------ | ----------- |
| SYNC_STATE_UNSPECIFIED | 0 | Value is ignored by the listening flagd |
| SYNC_STATE_ALL | 1 | All the flags matching the request. This is the default response and other states can be ignored by the implementation. Flagd internally replaces all existing flags for this response state. |
| SYNC_STATE_ADD | 2 | Convey an addition of a flag. Flagd internally handles this by combining new flags with existing ones |
| SYNC_STATE_UPDATE | 3 | Convey an update of a flag. Flagd internally attempts to update if the updated flag already exist OR if it does not, it will get added |
| SYNC_STATE_DELETE | 4 | Convey a deletion of a flag. Flagd internally removes the flag |
| SYNC_STATE_PING | 5 | Optional server ping to check client connectivity. Handling is ignored by flagd and is to merely support live check |


 

 


<a name="sync-v1-FlagSyncService"></a>

### FlagSyncService
FlagService implements a server streaming to provide realtime flag configurations

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SyncFlags | [SyncFlagsRequest](#sync-v1-SyncFlagsRequest) | [SyncFlagsResponse](#sync-v1-SyncFlagsResponse) stream |  |
| FetchAllFlags | [FetchAllFlagsRequest](#sync-v1-FetchAllFlagsRequest) | [FetchAllFlagsResponse](#sync-v1-FetchAllFlagsResponse) |  |

 



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

