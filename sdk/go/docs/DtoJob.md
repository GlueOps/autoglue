# DtoJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Attempts** | Pointer to **int32** |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**LastError** | Pointer to **string** |  | [optional] 
**MaxAttempts** | Pointer to **int32** |  | [optional] 
**Payload** | Pointer to **map[string]interface{}** |  | [optional] 
**Queue** | Pointer to **string** |  | [optional] 
**RunAt** | Pointer to **string** |  | [optional] 
**Status** | Pointer to [**DtoJobStatus**](DtoJobStatus.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 

## Methods

### NewDtoJob

`func NewDtoJob() *DtoJob`

NewDtoJob instantiates a new DtoJob object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDtoJobWithDefaults

`func NewDtoJobWithDefaults() *DtoJob`

NewDtoJobWithDefaults instantiates a new DtoJob object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAttempts

`func (o *DtoJob) GetAttempts() int32`

GetAttempts returns the Attempts field if non-nil, zero value otherwise.

### GetAttemptsOk

`func (o *DtoJob) GetAttemptsOk() (*int32, bool)`

GetAttemptsOk returns a tuple with the Attempts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttempts

`func (o *DtoJob) SetAttempts(v int32)`

SetAttempts sets Attempts field to given value.

### HasAttempts

`func (o *DtoJob) HasAttempts() bool`

HasAttempts returns a boolean if a field has been set.

### GetCreatedAt

`func (o *DtoJob) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *DtoJob) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *DtoJob) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *DtoJob) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetId

`func (o *DtoJob) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DtoJob) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DtoJob) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *DtoJob) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLastError

`func (o *DtoJob) GetLastError() string`

GetLastError returns the LastError field if non-nil, zero value otherwise.

### GetLastErrorOk

`func (o *DtoJob) GetLastErrorOk() (*string, bool)`

GetLastErrorOk returns a tuple with the LastError field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastError

`func (o *DtoJob) SetLastError(v string)`

SetLastError sets LastError field to given value.

### HasLastError

`func (o *DtoJob) HasLastError() bool`

HasLastError returns a boolean if a field has been set.

### GetMaxAttempts

`func (o *DtoJob) GetMaxAttempts() int32`

GetMaxAttempts returns the MaxAttempts field if non-nil, zero value otherwise.

### GetMaxAttemptsOk

`func (o *DtoJob) GetMaxAttemptsOk() (*int32, bool)`

GetMaxAttemptsOk returns a tuple with the MaxAttempts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxAttempts

`func (o *DtoJob) SetMaxAttempts(v int32)`

SetMaxAttempts sets MaxAttempts field to given value.

### HasMaxAttempts

`func (o *DtoJob) HasMaxAttempts() bool`

HasMaxAttempts returns a boolean if a field has been set.

### GetPayload

`func (o *DtoJob) GetPayload() map[string]interface{}`

GetPayload returns the Payload field if non-nil, zero value otherwise.

### GetPayloadOk

`func (o *DtoJob) GetPayloadOk() (*map[string]interface{}, bool)`

GetPayloadOk returns a tuple with the Payload field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayload

`func (o *DtoJob) SetPayload(v map[string]interface{})`

SetPayload sets Payload field to given value.

### HasPayload

`func (o *DtoJob) HasPayload() bool`

HasPayload returns a boolean if a field has been set.

### GetQueue

`func (o *DtoJob) GetQueue() string`

GetQueue returns the Queue field if non-nil, zero value otherwise.

### GetQueueOk

`func (o *DtoJob) GetQueueOk() (*string, bool)`

GetQueueOk returns a tuple with the Queue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueue

`func (o *DtoJob) SetQueue(v string)`

SetQueue sets Queue field to given value.

### HasQueue

`func (o *DtoJob) HasQueue() bool`

HasQueue returns a boolean if a field has been set.

### GetRunAt

`func (o *DtoJob) GetRunAt() string`

GetRunAt returns the RunAt field if non-nil, zero value otherwise.

### GetRunAtOk

`func (o *DtoJob) GetRunAtOk() (*string, bool)`

GetRunAtOk returns a tuple with the RunAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunAt

`func (o *DtoJob) SetRunAt(v string)`

SetRunAt sets RunAt field to given value.

### HasRunAt

`func (o *DtoJob) HasRunAt() bool`

HasRunAt returns a boolean if a field has been set.

### GetStatus

`func (o *DtoJob) GetStatus() DtoJobStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DtoJob) GetStatusOk() (*DtoJobStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DtoJob) SetStatus(v DtoJobStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *DtoJob) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetType

`func (o *DtoJob) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DtoJob) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DtoJob) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *DtoJob) HasType() bool`

HasType returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *DtoJob) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *DtoJob) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *DtoJob) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *DtoJob) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


