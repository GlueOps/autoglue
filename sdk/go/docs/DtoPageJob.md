# DtoPageJob

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | Pointer to [**[]DtoJob**](DtoJob.md) |  | [optional] 
**Page** | Pointer to **int32** | example: 1 | [optional] 
**PageSize** | Pointer to **int32** | example: 25 | [optional] 
**Total** | Pointer to **int32** | example: 120 | [optional] 

## Methods

### NewDtoPageJob

`func NewDtoPageJob() *DtoPageJob`

NewDtoPageJob instantiates a new DtoPageJob object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDtoPageJobWithDefaults

`func NewDtoPageJobWithDefaults() *DtoPageJob`

NewDtoPageJobWithDefaults instantiates a new DtoPageJob object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *DtoPageJob) GetItems() []DtoJob`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *DtoPageJob) GetItemsOk() (*[]DtoJob, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *DtoPageJob) SetItems(v []DtoJob)`

SetItems sets Items field to given value.

### HasItems

`func (o *DtoPageJob) HasItems() bool`

HasItems returns a boolean if a field has been set.

### GetPage

`func (o *DtoPageJob) GetPage() int32`

GetPage returns the Page field if non-nil, zero value otherwise.

### GetPageOk

`func (o *DtoPageJob) GetPageOk() (*int32, bool)`

GetPageOk returns a tuple with the Page field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPage

`func (o *DtoPageJob) SetPage(v int32)`

SetPage sets Page field to given value.

### HasPage

`func (o *DtoPageJob) HasPage() bool`

HasPage returns a boolean if a field has been set.

### GetPageSize

`func (o *DtoPageJob) GetPageSize() int32`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *DtoPageJob) GetPageSizeOk() (*int32, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *DtoPageJob) SetPageSize(v int32)`

SetPageSize sets PageSize field to given value.

### HasPageSize

`func (o *DtoPageJob) HasPageSize() bool`

HasPageSize returns a boolean if a field has been set.

### GetTotal

`func (o *DtoPageJob) GetTotal() int32`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *DtoPageJob) GetTotalOk() (*int32, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *DtoPageJob) SetTotal(v int32)`

SetTotal sets Total field to given value.

### HasTotal

`func (o *DtoPageJob) HasTotal() bool`

HasTotal returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


