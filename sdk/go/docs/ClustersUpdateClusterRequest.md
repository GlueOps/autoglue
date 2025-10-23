# ClustersUpdateClusterRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BastionServerId** | Pointer to **string** |  | [optional] 
**ClusterLoadBalancer** | Pointer to **string** |  | [optional] 
**ControlLoadBalancer** | Pointer to **string** |  | [optional] 
**Kubeconfig** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Provider** | Pointer to **string** |  | [optional] 
**Region** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Methods

### NewClustersUpdateClusterRequest

`func NewClustersUpdateClusterRequest() *ClustersUpdateClusterRequest`

NewClustersUpdateClusterRequest instantiates a new ClustersUpdateClusterRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClustersUpdateClusterRequestWithDefaults

`func NewClustersUpdateClusterRequestWithDefaults() *ClustersUpdateClusterRequest`

NewClustersUpdateClusterRequestWithDefaults instantiates a new ClustersUpdateClusterRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBastionServerId

`func (o *ClustersUpdateClusterRequest) GetBastionServerId() string`

GetBastionServerId returns the BastionServerId field if non-nil, zero value otherwise.

### GetBastionServerIdOk

`func (o *ClustersUpdateClusterRequest) GetBastionServerIdOk() (*string, bool)`

GetBastionServerIdOk returns a tuple with the BastionServerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBastionServerId

`func (o *ClustersUpdateClusterRequest) SetBastionServerId(v string)`

SetBastionServerId sets BastionServerId field to given value.

### HasBastionServerId

`func (o *ClustersUpdateClusterRequest) HasBastionServerId() bool`

HasBastionServerId returns a boolean if a field has been set.

### GetClusterLoadBalancer

`func (o *ClustersUpdateClusterRequest) GetClusterLoadBalancer() string`

GetClusterLoadBalancer returns the ClusterLoadBalancer field if non-nil, zero value otherwise.

### GetClusterLoadBalancerOk

`func (o *ClustersUpdateClusterRequest) GetClusterLoadBalancerOk() (*string, bool)`

GetClusterLoadBalancerOk returns a tuple with the ClusterLoadBalancer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterLoadBalancer

`func (o *ClustersUpdateClusterRequest) SetClusterLoadBalancer(v string)`

SetClusterLoadBalancer sets ClusterLoadBalancer field to given value.

### HasClusterLoadBalancer

`func (o *ClustersUpdateClusterRequest) HasClusterLoadBalancer() bool`

HasClusterLoadBalancer returns a boolean if a field has been set.

### GetControlLoadBalancer

`func (o *ClustersUpdateClusterRequest) GetControlLoadBalancer() string`

GetControlLoadBalancer returns the ControlLoadBalancer field if non-nil, zero value otherwise.

### GetControlLoadBalancerOk

`func (o *ClustersUpdateClusterRequest) GetControlLoadBalancerOk() (*string, bool)`

GetControlLoadBalancerOk returns a tuple with the ControlLoadBalancer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControlLoadBalancer

`func (o *ClustersUpdateClusterRequest) SetControlLoadBalancer(v string)`

SetControlLoadBalancer sets ControlLoadBalancer field to given value.

### HasControlLoadBalancer

`func (o *ClustersUpdateClusterRequest) HasControlLoadBalancer() bool`

HasControlLoadBalancer returns a boolean if a field has been set.

### GetKubeconfig

`func (o *ClustersUpdateClusterRequest) GetKubeconfig() string`

GetKubeconfig returns the Kubeconfig field if non-nil, zero value otherwise.

### GetKubeconfigOk

`func (o *ClustersUpdateClusterRequest) GetKubeconfigOk() (*string, bool)`

GetKubeconfigOk returns a tuple with the Kubeconfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKubeconfig

`func (o *ClustersUpdateClusterRequest) SetKubeconfig(v string)`

SetKubeconfig sets Kubeconfig field to given value.

### HasKubeconfig

`func (o *ClustersUpdateClusterRequest) HasKubeconfig() bool`

HasKubeconfig returns a boolean if a field has been set.

### GetName

`func (o *ClustersUpdateClusterRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ClustersUpdateClusterRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ClustersUpdateClusterRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ClustersUpdateClusterRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetProvider

`func (o *ClustersUpdateClusterRequest) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *ClustersUpdateClusterRequest) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *ClustersUpdateClusterRequest) SetProvider(v string)`

SetProvider sets Provider field to given value.

### HasProvider

`func (o *ClustersUpdateClusterRequest) HasProvider() bool`

HasProvider returns a boolean if a field has been set.

### GetRegion

`func (o *ClustersUpdateClusterRequest) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *ClustersUpdateClusterRequest) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *ClustersUpdateClusterRequest) SetRegion(v string)`

SetRegion sets Region field to given value.

### HasRegion

`func (o *ClustersUpdateClusterRequest) HasRegion() bool`

HasRegion returns a boolean if a field has been set.

### GetStatus

`func (o *ClustersUpdateClusterRequest) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ClustersUpdateClusterRequest) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ClustersUpdateClusterRequest) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ClustersUpdateClusterRequest) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


