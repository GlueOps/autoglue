# ClustersCreateClusterRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BastionServerId** | Pointer to **string** |  | [optional] 
**ClusterLoadBalancer** | Pointer to **string** |  | [optional] 
**ControlLoadBalancer** | Pointer to **string** |  | [optional] 
**Kubeconfig** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**NodePoolIds** | Pointer to **[]string** |  | [optional] 
**Provider** | Pointer to **string** |  | [optional] 
**Region** | Pointer to **string** |  | [optional] 

## Methods

### NewClustersCreateClusterRequest

`func NewClustersCreateClusterRequest() *ClustersCreateClusterRequest`

NewClustersCreateClusterRequest instantiates a new ClustersCreateClusterRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClustersCreateClusterRequestWithDefaults

`func NewClustersCreateClusterRequestWithDefaults() *ClustersCreateClusterRequest`

NewClustersCreateClusterRequestWithDefaults instantiates a new ClustersCreateClusterRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBastionServerId

`func (o *ClustersCreateClusterRequest) GetBastionServerId() string`

GetBastionServerId returns the BastionServerId field if non-nil, zero value otherwise.

### GetBastionServerIdOk

`func (o *ClustersCreateClusterRequest) GetBastionServerIdOk() (*string, bool)`

GetBastionServerIdOk returns a tuple with the BastionServerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBastionServerId

`func (o *ClustersCreateClusterRequest) SetBastionServerId(v string)`

SetBastionServerId sets BastionServerId field to given value.

### HasBastionServerId

`func (o *ClustersCreateClusterRequest) HasBastionServerId() bool`

HasBastionServerId returns a boolean if a field has been set.

### GetClusterLoadBalancer

`func (o *ClustersCreateClusterRequest) GetClusterLoadBalancer() string`

GetClusterLoadBalancer returns the ClusterLoadBalancer field if non-nil, zero value otherwise.

### GetClusterLoadBalancerOk

`func (o *ClustersCreateClusterRequest) GetClusterLoadBalancerOk() (*string, bool)`

GetClusterLoadBalancerOk returns a tuple with the ClusterLoadBalancer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterLoadBalancer

`func (o *ClustersCreateClusterRequest) SetClusterLoadBalancer(v string)`

SetClusterLoadBalancer sets ClusterLoadBalancer field to given value.

### HasClusterLoadBalancer

`func (o *ClustersCreateClusterRequest) HasClusterLoadBalancer() bool`

HasClusterLoadBalancer returns a boolean if a field has been set.

### GetControlLoadBalancer

`func (o *ClustersCreateClusterRequest) GetControlLoadBalancer() string`

GetControlLoadBalancer returns the ControlLoadBalancer field if non-nil, zero value otherwise.

### GetControlLoadBalancerOk

`func (o *ClustersCreateClusterRequest) GetControlLoadBalancerOk() (*string, bool)`

GetControlLoadBalancerOk returns a tuple with the ControlLoadBalancer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControlLoadBalancer

`func (o *ClustersCreateClusterRequest) SetControlLoadBalancer(v string)`

SetControlLoadBalancer sets ControlLoadBalancer field to given value.

### HasControlLoadBalancer

`func (o *ClustersCreateClusterRequest) HasControlLoadBalancer() bool`

HasControlLoadBalancer returns a boolean if a field has been set.

### GetKubeconfig

`func (o *ClustersCreateClusterRequest) GetKubeconfig() string`

GetKubeconfig returns the Kubeconfig field if non-nil, zero value otherwise.

### GetKubeconfigOk

`func (o *ClustersCreateClusterRequest) GetKubeconfigOk() (*string, bool)`

GetKubeconfigOk returns a tuple with the Kubeconfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKubeconfig

`func (o *ClustersCreateClusterRequest) SetKubeconfig(v string)`

SetKubeconfig sets Kubeconfig field to given value.

### HasKubeconfig

`func (o *ClustersCreateClusterRequest) HasKubeconfig() bool`

HasKubeconfig returns a boolean if a field has been set.

### GetName

`func (o *ClustersCreateClusterRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ClustersCreateClusterRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ClustersCreateClusterRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ClustersCreateClusterRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetNodePoolIds

`func (o *ClustersCreateClusterRequest) GetNodePoolIds() []string`

GetNodePoolIds returns the NodePoolIds field if non-nil, zero value otherwise.

### GetNodePoolIdsOk

`func (o *ClustersCreateClusterRequest) GetNodePoolIdsOk() (*[]string, bool)`

GetNodePoolIdsOk returns a tuple with the NodePoolIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodePoolIds

`func (o *ClustersCreateClusterRequest) SetNodePoolIds(v []string)`

SetNodePoolIds sets NodePoolIds field to given value.

### HasNodePoolIds

`func (o *ClustersCreateClusterRequest) HasNodePoolIds() bool`

HasNodePoolIds returns a boolean if a field has been set.

### GetProvider

`func (o *ClustersCreateClusterRequest) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *ClustersCreateClusterRequest) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *ClustersCreateClusterRequest) SetProvider(v string)`

SetProvider sets Provider field to given value.

### HasProvider

`func (o *ClustersCreateClusterRequest) HasProvider() bool`

HasProvider returns a boolean if a field has been set.

### GetRegion

`func (o *ClustersCreateClusterRequest) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *ClustersCreateClusterRequest) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *ClustersCreateClusterRequest) SetRegion(v string)`

SetRegion sets Region field to given value.

### HasRegion

`func (o *ClustersCreateClusterRequest) HasRegion() bool`

HasRegion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


