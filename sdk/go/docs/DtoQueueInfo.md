# DtoQueueInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Failed** | Pointer to **int32** | example: 5 | [optional] 
**Name** | Pointer to **string** | example: default | [optional] 
**Pending** | Pointer to **int32** | example: 42 | [optional] 
**Running** | Pointer to **int32** | example: 3 | [optional] 
**Scheduled** | Pointer to **int32** | example: 7 | [optional] 

## Methods

### NewDtoQueueInfo

`func NewDtoQueueInfo() *DtoQueueInfo`

NewDtoQueueInfo instantiates a new DtoQueueInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDtoQueueInfoWithDefaults

`func NewDtoQueueInfoWithDefaults() *DtoQueueInfo`

NewDtoQueueInfoWithDefaults instantiates a new DtoQueueInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFailed

`func (o *DtoQueueInfo) GetFailed() int32`

GetFailed returns the Failed field if non-nil, zero value otherwise.

### GetFailedOk

`func (o *DtoQueueInfo) GetFailedOk() (*int32, bool)`

GetFailedOk returns a tuple with the Failed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFailed

`func (o *DtoQueueInfo) SetFailed(v int32)`

SetFailed sets Failed field to given value.

### HasFailed

`func (o *DtoQueueInfo) HasFailed() bool`

HasFailed returns a boolean if a field has been set.

### GetName

`func (o *DtoQueueInfo) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DtoQueueInfo) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DtoQueueInfo) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DtoQueueInfo) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPending

`func (o *DtoQueueInfo) GetPending() int32`

GetPending returns the Pending field if non-nil, zero value otherwise.

### GetPendingOk

`func (o *DtoQueueInfo) GetPendingOk() (*int32, bool)`

GetPendingOk returns a tuple with the Pending field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPending

`func (o *DtoQueueInfo) SetPending(v int32)`

SetPending sets Pending field to given value.

### HasPending

`func (o *DtoQueueInfo) HasPending() bool`

HasPending returns a boolean if a field has been set.

### GetRunning

`func (o *DtoQueueInfo) GetRunning() int32`

GetRunning returns the Running field if non-nil, zero value otherwise.

### GetRunningOk

`func (o *DtoQueueInfo) GetRunningOk() (*int32, bool)`

GetRunningOk returns a tuple with the Running field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunning

`func (o *DtoQueueInfo) SetRunning(v int32)`

SetRunning sets Running field to given value.

### HasRunning

`func (o *DtoQueueInfo) HasRunning() bool`

HasRunning returns a boolean if a field has been set.

### GetScheduled

`func (o *DtoQueueInfo) GetScheduled() int32`

GetScheduled returns the Scheduled field if non-nil, zero value otherwise.

### GetScheduledOk

`func (o *DtoQueueInfo) GetScheduledOk() (*int32, bool)`

GetScheduledOk returns a tuple with the Scheduled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheduled

`func (o *DtoQueueInfo) SetScheduled(v int32)`

SetScheduled sets Scheduled field to given value.

### HasScheduled

`func (o *DtoQueueInfo) HasScheduled() bool`

HasScheduled returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


