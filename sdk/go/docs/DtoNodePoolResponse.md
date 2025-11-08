# DtoNodePoolResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Annotations** | Pointer to [**[]DtoAnnotationResponse**](DtoAnnotationResponse.md) |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to [**[]DtoLabelResponse**](DtoLabelResponse.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**Role** | Pointer to **string** |  | [optional] 
**Servers** | Pointer to [**[]DtoServerResponse**](DtoServerResponse.md) |  | [optional] 
**Taints** | Pointer to [**[]DtoTaintResponse**](DtoTaintResponse.md) |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 

## Methods

### NewDtoNodePoolResponse

`func NewDtoNodePoolResponse() *DtoNodePoolResponse`

NewDtoNodePoolResponse instantiates a new DtoNodePoolResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDtoNodePoolResponseWithDefaults

`func NewDtoNodePoolResponseWithDefaults() *DtoNodePoolResponse`

NewDtoNodePoolResponseWithDefaults instantiates a new DtoNodePoolResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAnnotations

`func (o *DtoNodePoolResponse) GetAnnotations() []DtoAnnotationResponse`

GetAnnotations returns the Annotations field if non-nil, zero value otherwise.

### GetAnnotationsOk

`func (o *DtoNodePoolResponse) GetAnnotationsOk() (*[]DtoAnnotationResponse, bool)`

GetAnnotationsOk returns a tuple with the Annotations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnnotations

`func (o *DtoNodePoolResponse) SetAnnotations(v []DtoAnnotationResponse)`

SetAnnotations sets Annotations field to given value.

### HasAnnotations

`func (o *DtoNodePoolResponse) HasAnnotations() bool`

HasAnnotations returns a boolean if a field has been set.

### GetCreatedAt

`func (o *DtoNodePoolResponse) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *DtoNodePoolResponse) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *DtoNodePoolResponse) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *DtoNodePoolResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetId

`func (o *DtoNodePoolResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DtoNodePoolResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DtoNodePoolResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *DtoNodePoolResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLabels

`func (o *DtoNodePoolResponse) GetLabels() []DtoLabelResponse`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *DtoNodePoolResponse) GetLabelsOk() (*[]DtoLabelResponse, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *DtoNodePoolResponse) SetLabels(v []DtoLabelResponse)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *DtoNodePoolResponse) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetName

`func (o *DtoNodePoolResponse) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DtoNodePoolResponse) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DtoNodePoolResponse) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DtoNodePoolResponse) HasName() bool`

HasName returns a boolean if a field has been set.

### GetOrganizationId

`func (o *DtoNodePoolResponse) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *DtoNodePoolResponse) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *DtoNodePoolResponse) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *DtoNodePoolResponse) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetRole

`func (o *DtoNodePoolResponse) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *DtoNodePoolResponse) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *DtoNodePoolResponse) SetRole(v string)`

SetRole sets Role field to given value.

### HasRole

`func (o *DtoNodePoolResponse) HasRole() bool`

HasRole returns a boolean if a field has been set.

### GetServers

`func (o *DtoNodePoolResponse) GetServers() []DtoServerResponse`

GetServers returns the Servers field if non-nil, zero value otherwise.

### GetServersOk

`func (o *DtoNodePoolResponse) GetServersOk() (*[]DtoServerResponse, bool)`

GetServersOk returns a tuple with the Servers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServers

`func (o *DtoNodePoolResponse) SetServers(v []DtoServerResponse)`

SetServers sets Servers field to given value.

### HasServers

`func (o *DtoNodePoolResponse) HasServers() bool`

HasServers returns a boolean if a field has been set.

### GetTaints

`func (o *DtoNodePoolResponse) GetTaints() []DtoTaintResponse`

GetTaints returns the Taints field if non-nil, zero value otherwise.

### GetTaintsOk

`func (o *DtoNodePoolResponse) GetTaintsOk() (*[]DtoTaintResponse, bool)`

GetTaintsOk returns a tuple with the Taints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaints

`func (o *DtoNodePoolResponse) SetTaints(v []DtoTaintResponse)`

SetTaints sets Taints field to given value.

### HasTaints

`func (o *DtoNodePoolResponse) HasTaints() bool`

HasTaints returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *DtoNodePoolResponse) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *DtoNodePoolResponse) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *DtoNodePoolResponse) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *DtoNodePoolResponse) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


