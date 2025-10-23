# ClustersClusterResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BastionServer** | Pointer to [**ClustersServerBrief**](ClustersServerBrief.md) |  | [optional] 
**ClusterLoadBalancer** | Pointer to **string** |  | [optional] 
**ControlLoadBalancer** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**NodePools** | Pointer to [**[]ClustersNodePoolBrief**](ClustersNodePoolBrief.md) |  | [optional] 
**Provider** | Pointer to **string** |  | [optional] 
**Region** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Methods

### NewClustersClusterResponse

`func NewClustersClusterResponse() *ClustersClusterResponse`

NewClustersClusterResponse instantiates a new ClustersClusterResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClustersClusterResponseWithDefaults

`func NewClustersClusterResponseWithDefaults() *ClustersClusterResponse`

NewClustersClusterResponseWithDefaults instantiates a new ClustersClusterResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBastionServer

`func (o *ClustersClusterResponse) GetBastionServer() ClustersServerBrief`

GetBastionServer returns the BastionServer field if non-nil, zero value otherwise.

### GetBastionServerOk

`func (o *ClustersClusterResponse) GetBastionServerOk() (*ClustersServerBrief, bool)`

GetBastionServerOk returns a tuple with the BastionServer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBastionServer

`func (o *ClustersClusterResponse) SetBastionServer(v ClustersServerBrief)`

SetBastionServer sets BastionServer field to given value.

### HasBastionServer

`func (o *ClustersClusterResponse) HasBastionServer() bool`

HasBastionServer returns a boolean if a field has been set.

### GetClusterLoadBalancer

`func (o *ClustersClusterResponse) GetClusterLoadBalancer() string`

GetClusterLoadBalancer returns the ClusterLoadBalancer field if non-nil, zero value otherwise.

### GetClusterLoadBalancerOk

`func (o *ClustersClusterResponse) GetClusterLoadBalancerOk() (*string, bool)`

GetClusterLoadBalancerOk returns a tuple with the ClusterLoadBalancer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterLoadBalancer

`func (o *ClustersClusterResponse) SetClusterLoadBalancer(v string)`

SetClusterLoadBalancer sets ClusterLoadBalancer field to given value.

### HasClusterLoadBalancer

`func (o *ClustersClusterResponse) HasClusterLoadBalancer() bool`

HasClusterLoadBalancer returns a boolean if a field has been set.

### GetControlLoadBalancer

`func (o *ClustersClusterResponse) GetControlLoadBalancer() string`

GetControlLoadBalancer returns the ControlLoadBalancer field if non-nil, zero value otherwise.

### GetControlLoadBalancerOk

`func (o *ClustersClusterResponse) GetControlLoadBalancerOk() (*string, bool)`

GetControlLoadBalancerOk returns a tuple with the ControlLoadBalancer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControlLoadBalancer

`func (o *ClustersClusterResponse) SetControlLoadBalancer(v string)`

SetControlLoadBalancer sets ControlLoadBalancer field to given value.

### HasControlLoadBalancer

`func (o *ClustersClusterResponse) HasControlLoadBalancer() bool`

HasControlLoadBalancer returns a boolean if a field has been set.

### GetId

`func (o *ClustersClusterResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ClustersClusterResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ClustersClusterResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ClustersClusterResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetName

`func (o *ClustersClusterResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ClustersClusterResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ClustersClusterResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ClustersClusterResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetNodePools

`func (o *ClustersClusterResponse) GetNodePools() []ClustersNodePoolBrief`

GetNodePools returns the NodePools field if non-nil, zero value otherwise.

### GetNodePoolsOk

`func (o *ClustersClusterResponse) GetNodePoolsOk() (*[]ClustersNodePoolBrief, bool)`

GetNodePoolsOk returns a tuple with the NodePools field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodePools

`func (o *ClustersClusterResponse) SetNodePools(v []ClustersNodePoolBrief)`

SetNodePools sets NodePools field to given value.

### HasNodePools

`func (o *ClustersClusterResponse) HasNodePools() bool`

HasNodePools returns a boolean if a field has been set.

### GetProvider

`func (o *ClustersClusterResponse) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *ClustersClusterResponse) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *ClustersClusterResponse) SetProvider(v string)`

SetProvider sets Provider field to given value.

### HasProvider

`func (o *ClustersClusterResponse) HasProvider() bool`

HasProvider returns a boolean if a field has been set.

### GetRegion

`func (o *ClustersClusterResponse) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *ClustersClusterResponse) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *ClustersClusterResponse) SetRegion(v string)`

SetRegion sets Region field to given value.

### HasRegion

`func (o *ClustersClusterResponse) HasRegion() bool`

HasRegion returns a boolean if a field has been set.

### GetStatus

`func (o *ClustersClusterResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ClustersClusterResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ClustersClusterResponse) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ClustersClusterResponse) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


