# LabelsCreateLabelRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | Pointer to **string** |  | [optional] 
**NodePoolIds** | Pointer to **[]string** |  | [optional] 
**Value** | Pointer to **string** |  | [optional] 

## Methods

### NewLabelsCreateLabelRequest

`func NewLabelsCreateLabelRequest() *LabelsCreateLabelRequest`

NewLabelsCreateLabelRequest instantiates a new LabelsCreateLabelRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLabelsCreateLabelRequestWithDefaults

`func NewLabelsCreateLabelRequestWithDefaults() *LabelsCreateLabelRequest`

NewLabelsCreateLabelRequestWithDefaults instantiates a new LabelsCreateLabelRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *LabelsCreateLabelRequest) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *LabelsCreateLabelRequest) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *LabelsCreateLabelRequest) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *LabelsCreateLabelRequest) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetNodePoolIds

`func (o *LabelsCreateLabelRequest) GetNodePoolIds() []string`

GetNodePoolIds returns the NodePoolIds field if non-nil, zero value otherwise.

### GetNodePoolIdsOk

`func (o *LabelsCreateLabelRequest) GetNodePoolIdsOk() (*[]string, bool)`

GetNodePoolIdsOk returns a tuple with the NodePoolIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodePoolIds

`func (o *LabelsCreateLabelRequest) SetNodePoolIds(v []string)`

SetNodePoolIds sets NodePoolIds field to given value.

### HasNodePoolIds

`func (o *LabelsCreateLabelRequest) HasNodePoolIds() bool`

HasNodePoolIds returns a boolean if a field has been set.

### GetValue

`func (o *LabelsCreateLabelRequest) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *LabelsCreateLabelRequest) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *LabelsCreateLabelRequest) SetValue(v string)`

SetValue sets Value field to given value.

### HasValue

`func (o *LabelsCreateLabelRequest) HasValue() bool`

HasValue returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


