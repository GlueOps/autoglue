# ClustersNodePoolBrief

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Annotations** | Pointer to [**[]ClustersAnnotationBrief**](ClustersAnnotationBrief.md) |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to [**[]ClustersLabelBrief**](ClustersLabelBrief.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Servers** | Pointer to [**[]ClustersServerBrief**](ClustersServerBrief.md) |  | [optional] 
**Taints** | Pointer to [**[]ClustersTaintBrief**](ClustersTaintBrief.md) |  | [optional] 

## Methods

### NewClustersNodePoolBrief

`func NewClustersNodePoolBrief() *ClustersNodePoolBrief`

NewClustersNodePoolBrief instantiates a new ClustersNodePoolBrief object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewClustersNodePoolBriefWithDefaults

`func NewClustersNodePoolBriefWithDefaults() *ClustersNodePoolBrief`

NewClustersNodePoolBriefWithDefaults instantiates a new ClustersNodePoolBrief object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAnnotations

`func (o *ClustersNodePoolBrief) GetAnnotations() []ClustersAnnotationBrief`

GetAnnotations returns the Annotations field if non-nil, zero value otherwise.

### GetAnnotationsOk

`func (o *ClustersNodePoolBrief) GetAnnotationsOk() (*[]ClustersAnnotationBrief, bool)`

GetAnnotationsOk returns a tuple with the Annotations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnnotations

`func (o *ClustersNodePoolBrief) SetAnnotations(v []ClustersAnnotationBrief)`

SetAnnotations sets Annotations field to given value.

### HasAnnotations

`func (o *ClustersNodePoolBrief) HasAnnotations() bool`

HasAnnotations returns a boolean if a field has been set.

### GetId

`func (o *ClustersNodePoolBrief) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ClustersNodePoolBrief) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ClustersNodePoolBrief) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ClustersNodePoolBrief) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLabels

`func (o *ClustersNodePoolBrief) GetLabels() []ClustersLabelBrief`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *ClustersNodePoolBrief) GetLabelsOk() (*[]ClustersLabelBrief, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *ClustersNodePoolBrief) SetLabels(v []ClustersLabelBrief)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *ClustersNodePoolBrief) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetName

`func (o *ClustersNodePoolBrief) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ClustersNodePoolBrief) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ClustersNodePoolBrief) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ClustersNodePoolBrief) HasName() bool`

HasName returns a boolean if a field has been set.

### GetServers

`func (o *ClustersNodePoolBrief) GetServers() []ClustersServerBrief`

GetServers returns the Servers field if non-nil, zero value otherwise.

### GetServersOk

`func (o *ClustersNodePoolBrief) GetServersOk() (*[]ClustersServerBrief, bool)`

GetServersOk returns a tuple with the Servers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServers

`func (o *ClustersNodePoolBrief) SetServers(v []ClustersServerBrief)`

SetServers sets Servers field to given value.

### HasServers

`func (o *ClustersNodePoolBrief) HasServers() bool`

HasServers returns a boolean if a field has been set.

### GetTaints

`func (o *ClustersNodePoolBrief) GetTaints() []ClustersTaintBrief`

GetTaints returns the Taints field if non-nil, zero value otherwise.

### GetTaintsOk

`func (o *ClustersNodePoolBrief) GetTaintsOk() (*[]ClustersTaintBrief, bool)`

GetTaintsOk returns a tuple with the Taints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaints

`func (o *ClustersNodePoolBrief) SetTaints(v []ClustersTaintBrief)`

SetTaints sets Taints field to given value.

### HasTaints

`func (o *ClustersNodePoolBrief) HasTaints() bool`

HasTaints returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


