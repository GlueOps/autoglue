# NodepoolsNodePoolResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Servers** | Pointer to [**[]NodepoolsServerBrief**](NodepoolsServerBrief.md) |  | [optional] 

## Methods

### NewNodepoolsNodePoolResponse

`func NewNodepoolsNodePoolResponse() *NodepoolsNodePoolResponse`

NewNodepoolsNodePoolResponse instantiates a new NodepoolsNodePoolResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodepoolsNodePoolResponseWithDefaults

`func NewNodepoolsNodePoolResponseWithDefaults() *NodepoolsNodePoolResponse`

NewNodepoolsNodePoolResponseWithDefaults instantiates a new NodepoolsNodePoolResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *NodepoolsNodePoolResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NodepoolsNodePoolResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NodepoolsNodePoolResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *NodepoolsNodePoolResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *NodepoolsNodePoolResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NodepoolsNodePoolResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NodepoolsNodePoolResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *NodepoolsNodePoolResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetServers

`func (o *NodepoolsNodePoolResponse) GetServers() []NodepoolsServerBrief`

GetServers returns the Servers field if non-nil, zero value otherwise.

### GetServersOk

`func (o *NodepoolsNodePoolResponse) GetServersOk() (*[]NodepoolsServerBrief, bool)`

GetServersOk returns a tuple with the Servers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServers

`func (o *NodepoolsNodePoolResponse) SetServers(v []NodepoolsServerBrief)`

SetServers sets Servers field to given value.

### HasServers

`func (o *NodepoolsNodePoolResponse) HasServers() bool`

HasServers returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


