# NodepoolsCreateNodePoolRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**ServerIds** | Pointer to **[]string** |  | [optional] 

## Methods

### NewNodepoolsCreateNodePoolRequest

`func NewNodepoolsCreateNodePoolRequest() *NodepoolsCreateNodePoolRequest`

NewNodepoolsCreateNodePoolRequest instantiates a new NodepoolsCreateNodePoolRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodepoolsCreateNodePoolRequestWithDefaults

`func NewNodepoolsCreateNodePoolRequestWithDefaults() *NodepoolsCreateNodePoolRequest`

NewNodepoolsCreateNodePoolRequestWithDefaults instantiates a new NodepoolsCreateNodePoolRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *NodepoolsCreateNodePoolRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NodepoolsCreateNodePoolRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NodepoolsCreateNodePoolRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *NodepoolsCreateNodePoolRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetServerIds

`func (o *NodepoolsCreateNodePoolRequest) GetServerIds() []string`

GetServerIds returns the ServerIds field if non-nil, zero value otherwise.

### GetServerIdsOk

`func (o *NodepoolsCreateNodePoolRequest) GetServerIdsOk() (*[]string, bool)`

GetServerIdsOk returns a tuple with the ServerIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerIds

`func (o *NodepoolsCreateNodePoolRequest) SetServerIds(v []string)`

SetServerIds sets ServerIds field to given value.

### HasServerIds

`func (o *NodepoolsCreateNodePoolRequest) HasServerIds() bool`

HasServerIds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


