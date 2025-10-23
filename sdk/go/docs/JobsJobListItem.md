# JobsJobListItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ElapsedMs** | Pointer to **int32** |  | [optional] 
**Failed** | Pointer to **int32** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**LastError** | Pointer to **string** |  | [optional] 
**MaxRetry** | Pointer to **int32** |  | [optional] 
**Processed** | Pointer to **int32** |  | [optional] 
**QueueName** | Pointer to **string** |  | [optional] 
**Ready** | Pointer to **int32** |  | [optional] 
**ResultStatus** | Pointer to **string** |  | [optional] 
**RetryCount** | Pointer to **int32** |  | [optional] 
**ScheduledAt** | Pointer to **string** |  | [optional] 
**StartedAt** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 

## Methods

### NewJobsJobListItem

`func NewJobsJobListItem() *JobsJobListItem`

NewJobsJobListItem instantiates a new JobsJobListItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJobsJobListItemWithDefaults

`func NewJobsJobListItemWithDefaults() *JobsJobListItem`

NewJobsJobListItemWithDefaults instantiates a new JobsJobListItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetElapsedMs

`func (o *JobsJobListItem) GetElapsedMs() int32`

GetElapsedMs returns the ElapsedMs field if non-nil, zero value otherwise.

### GetElapsedMsOk

`func (o *JobsJobListItem) GetElapsedMsOk() (*int32, bool)`

GetElapsedMsOk returns a tuple with the ElapsedMs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetElapsedMs

`func (o *JobsJobListItem) SetElapsedMs(v int32)`

SetElapsedMs sets ElapsedMs field to given value.

### HasElapsedMs

`func (o *JobsJobListItem) HasElapsedMs() bool`

HasElapsedMs returns a boolean if a field has been set.

### GetFailed

`func (o *JobsJobListItem) GetFailed() int32`

GetFailed returns the Failed field if non-nil, zero value otherwise.

### GetFailedOk

`func (o *JobsJobListItem) GetFailedOk() (*int32, bool)`

GetFailedOk returns a tuple with the Failed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailed

`func (o *JobsJobListItem) SetFailed(v int32)`

SetFailed sets Failed field to given value.

### HasFailed

`func (o *JobsJobListItem) HasFailed() bool`

HasFailed returns a boolean if a field has been set.

### GetId

`func (o *JobsJobListItem) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *JobsJobListItem) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *JobsJobListItem) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *JobsJobListItem) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLastError

`func (o *JobsJobListItem) GetLastError() string`

GetLastError returns the LastError field if non-nil, zero value otherwise.

### GetLastErrorOk

`func (o *JobsJobListItem) GetLastErrorOk() (*string, bool)`

GetLastErrorOk returns a tuple with the LastError field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastError

`func (o *JobsJobListItem) SetLastError(v string)`

SetLastError sets LastError field to given value.

### HasLastError

`func (o *JobsJobListItem) HasLastError() bool`

HasLastError returns a boolean if a field has been set.

### GetMaxRetry

`func (o *JobsJobListItem) GetMaxRetry() int32`

GetMaxRetry returns the MaxRetry field if non-nil, zero value otherwise.

### GetMaxRetryOk

`func (o *JobsJobListItem) GetMaxRetryOk() (*int32, bool)`

GetMaxRetryOk returns a tuple with the MaxRetry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxRetry

`func (o *JobsJobListItem) SetMaxRetry(v int32)`

SetMaxRetry sets MaxRetry field to given value.

### HasMaxRetry

`func (o *JobsJobListItem) HasMaxRetry() bool`

HasMaxRetry returns a boolean if a field has been set.

### GetProcessed

`func (o *JobsJobListItem) GetProcessed() int32`

GetProcessed returns the Processed field if non-nil, zero value otherwise.

### GetProcessedOk

`func (o *JobsJobListItem) GetProcessedOk() (*int32, bool)`

GetProcessedOk returns a tuple with the Processed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProcessed

`func (o *JobsJobListItem) SetProcessed(v int32)`

SetProcessed sets Processed field to given value.

### HasProcessed

`func (o *JobsJobListItem) HasProcessed() bool`

HasProcessed returns a boolean if a field has been set.

### GetQueueName

`func (o *JobsJobListItem) GetQueueName() string`

GetQueueName returns the QueueName field if non-nil, zero value otherwise.

### GetQueueNameOk

`func (o *JobsJobListItem) GetQueueNameOk() (*string, bool)`

GetQueueNameOk returns a tuple with the QueueName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueueName

`func (o *JobsJobListItem) SetQueueName(v string)`

SetQueueName sets QueueName field to given value.

### HasQueueName

`func (o *JobsJobListItem) HasQueueName() bool`

HasQueueName returns a boolean if a field has been set.

### GetReady

`func (o *JobsJobListItem) GetReady() int32`

GetReady returns the Ready field if non-nil, zero value otherwise.

### GetReadyOk

`func (o *JobsJobListItem) GetReadyOk() (*int32, bool)`

GetReadyOk returns a tuple with the Ready field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReady

`func (o *JobsJobListItem) SetReady(v int32)`

SetReady sets Ready field to given value.

### HasReady

`func (o *JobsJobListItem) HasReady() bool`

HasReady returns a boolean if a field has been set.

### GetResultStatus

`func (o *JobsJobListItem) GetResultStatus() string`

GetResultStatus returns the ResultStatus field if non-nil, zero value otherwise.

### GetResultStatusOk

`func (o *JobsJobListItem) GetResultStatusOk() (*string, bool)`

GetResultStatusOk returns a tuple with the ResultStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResultStatus

`func (o *JobsJobListItem) SetResultStatus(v string)`

SetResultStatus sets ResultStatus field to given value.

### HasResultStatus

`func (o *JobsJobListItem) HasResultStatus() bool`

HasResultStatus returns a boolean if a field has been set.

### GetRetryCount

`func (o *JobsJobListItem) GetRetryCount() int32`

GetRetryCount returns the RetryCount field if non-nil, zero value otherwise.

### GetRetryCountOk

`func (o *JobsJobListItem) GetRetryCountOk() (*int32, bool)`

GetRetryCountOk returns a tuple with the RetryCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetryCount

`func (o *JobsJobListItem) SetRetryCount(v int32)`

SetRetryCount sets RetryCount field to given value.

### HasRetryCount

`func (o *JobsJobListItem) HasRetryCount() bool`

HasRetryCount returns a boolean if a field has been set.

### GetScheduledAt

`func (o *JobsJobListItem) GetScheduledAt() string`

GetScheduledAt returns the ScheduledAt field if non-nil, zero value otherwise.

### GetScheduledAtOk

`func (o *JobsJobListItem) GetScheduledAtOk() (*string, bool)`

GetScheduledAtOk returns a tuple with the ScheduledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheduledAt

`func (o *JobsJobListItem) SetScheduledAt(v string)`

SetScheduledAt sets ScheduledAt field to given value.

### HasScheduledAt

`func (o *JobsJobListItem) HasScheduledAt() bool`

HasScheduledAt returns a boolean if a field has been set.

### GetStartedAt

`func (o *JobsJobListItem) GetStartedAt() string`

GetStartedAt returns the StartedAt field if non-nil, zero value otherwise.

### GetStartedAtOk

`func (o *JobsJobListItem) GetStartedAtOk() (*string, bool)`

GetStartedAtOk returns a tuple with the StartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartedAt

`func (o *JobsJobListItem) SetStartedAt(v string)`

SetStartedAt sets StartedAt field to given value.

### HasStartedAt

`func (o *JobsJobListItem) HasStartedAt() bool`

HasStartedAt returns a boolean if a field has been set.

### GetStatus

`func (o *JobsJobListItem) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *JobsJobListItem) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *JobsJobListItem) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *JobsJobListItem) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *JobsJobListItem) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *JobsJobListItem) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *JobsJobListItem) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *JobsJobListItem) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


