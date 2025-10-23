# TaintsNodePoolResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Servers** | Pointer to [**[]TaintsServerBrief**](TaintsServerBrief.md) |  | [optional] 

## Methods

### NewTaintsNodePoolResponse

`func NewTaintsNodePoolResponse() *TaintsNodePoolResponse`

NewTaintsNodePoolResponse instantiates a new TaintsNodePoolResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTaintsNodePoolResponseWithDefaults

`func NewTaintsNodePoolResponseWithDefaults() *TaintsNodePoolResponse`

NewTaintsNodePoolResponseWithDefaults instantiates a new TaintsNodePoolResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *TaintsNodePoolResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *TaintsNodePoolResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *TaintsNodePoolResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *TaintsNodePoolResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *TaintsNodePoolResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *TaintsNodePoolResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *TaintsNodePoolResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *TaintsNodePoolResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetServers

`func (o *TaintsNodePoolResponse) GetServers() []TaintsServerBrief`

GetServers returns the Servers field if non-nil, zero value otherwise.

### GetServersOk

`func (o *TaintsNodePoolResponse) GetServersOk() (*[]TaintsServerBrief, bool)`

GetServersOk returns a tuple with the Servers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServers

`func (o *TaintsNodePoolResponse) SetServers(v []TaintsServerBrief)`

SetServers sets Servers field to given value.

### HasServers

`func (o *TaintsNodePoolResponse) HasServers() bool`

HasServers returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


