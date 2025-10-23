# JobsEnqueueReq

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Args** | Pointer to **map[string]interface{}** |  | [optional] 
**MaxRetries** | Pointer to **int32** |  | [optional] 
**Queue** | Pointer to **string** |  | [optional] 
**ScheduleAt** | Pointer to **string** |  | [optional] 

## Methods

### NewJobsEnqueueReq

`func NewJobsEnqueueReq() *JobsEnqueueReq`

NewJobsEnqueueReq instantiates a new JobsEnqueueReq object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewJobsEnqueueReqWithDefaults

`func NewJobsEnqueueReqWithDefaults() *JobsEnqueueReq`

NewJobsEnqueueReqWithDefaults instantiates a new JobsEnqueueReq object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetArgs

`func (o *JobsEnqueueReq) GetArgs() map[string]interface{}`

GetArgs returns the Args field if non-nil, zero value otherwise.

### GetArgsOk

`func (o *JobsEnqueueReq) GetArgsOk() (*map[string]interface{}, bool)`

GetArgsOk returns a tuple with the Args field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArgs

`func (o *JobsEnqueueReq) SetArgs(v map[string]interface{})`

SetArgs sets Args field to given value.

### HasArgs

`func (o *JobsEnqueueReq) HasArgs() bool`

HasArgs returns a boolean if a field has been set.

### GetMaxRetries

`func (o *JobsEnqueueReq) GetMaxRetries() int32`

GetMaxRetries returns the MaxRetries field if non-nil, zero value otherwise.

### GetMaxRetriesOk

`func (o *JobsEnqueueReq) GetMaxRetriesOk() (*int32, bool)`

GetMaxRetriesOk returns a tuple with the MaxRetries field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxRetries

`func (o *JobsEnqueueReq) SetMaxRetries(v int32)`

SetMaxRetries sets MaxRetries field to given value.

### HasMaxRetries

`func (o *JobsEnqueueReq) HasMaxRetries() bool`

HasMaxRetries returns a boolean if a field has been set.

### GetQueue

`func (o *JobsEnqueueReq) GetQueue() string`

GetQueue returns the Queue field if non-nil, zero value otherwise.

### GetQueueOk

`func (o *JobsEnqueueReq) GetQueueOk() (*string, bool)`

GetQueueOk returns a tuple with the Queue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueue

`func (o *JobsEnqueueReq) SetQueue(v string)`

SetQueue sets Queue field to given value.

### HasQueue

`func (o *JobsEnqueueReq) HasQueue() bool`

HasQueue returns a boolean if a field has been set.

### GetScheduleAt

`func (o *JobsEnqueueReq) GetScheduleAt() string`

GetScheduleAt returns the ScheduleAt field if non-nil, zero value otherwise.

### GetScheduleAtOk

`func (o *JobsEnqueueReq) GetScheduleAtOk() (*string, bool)`

GetScheduleAtOk returns a tuple with the ScheduleAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheduleAt

`func (o *JobsEnqueueReq) SetScheduleAt(v string)`

SetScheduleAt sets ScheduleAt field to given value.

### HasScheduleAt

`func (o *JobsEnqueueReq) HasScheduleAt() bool`

HasScheduleAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


