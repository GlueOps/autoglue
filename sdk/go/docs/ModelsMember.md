# ModelsMember

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Organization** | Pointer to [**ModelsOrganization**](ModelsOrganization.md) |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**Role** | Pointer to [**ModelsMemberRole**](ModelsMemberRole.md) | e.g. admin, member | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 
**User** | Pointer to [**ModelsUser**](ModelsUser.md) |  | [optional] 
**UserId** | Pointer to **string** |  | [optional] 

## Methods

### NewModelsMember

`func NewModelsMember() *ModelsMember`

NewModelsMember instantiates a new ModelsMember object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewModelsMemberWithDefaults

`func NewModelsMemberWithDefaults() *ModelsMember`

NewModelsMemberWithDefaults instantiates a new ModelsMember object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *ModelsMember) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ModelsMember) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ModelsMember) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *ModelsMember) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetId

`func (o *ModelsMember) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ModelsMember) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ModelsMember) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ModelsMember) HasId() bool`

HasId returns a boolean if a field has been set.

### GetOrganization

`func (o *ModelsMember) GetOrganization() ModelsOrganization`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *ModelsMember) GetOrganizationOk() (*ModelsOrganization, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *ModelsMember) SetOrganization(v ModelsOrganization)`

SetOrganization sets Organization field to given value.

### HasOrganization

`func (o *ModelsMember) HasOrganization() bool`

HasOrganization returns a boolean if a field has been set.

### GetOrganizationId

`func (o *ModelsMember) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *ModelsMember) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *ModelsMember) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *ModelsMember) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetRole

`func (o *ModelsMember) GetRole() ModelsMemberRole`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *ModelsMember) GetRoleOk() (*ModelsMemberRole, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *ModelsMember) SetRole(v ModelsMemberRole)`

SetRole sets Role field to given value.

### HasRole

`func (o *ModelsMember) HasRole() bool`

HasRole returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *ModelsMember) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ModelsMember) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ModelsMember) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *ModelsMember) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetUser

`func (o *ModelsMember) GetUser() ModelsUser`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *ModelsMember) GetUserOk() (*ModelsUser, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *ModelsMember) SetUser(v ModelsUser)`

SetUser sets User field to given value.

### HasUser

`func (o *ModelsMember) HasUser() bool`

HasUser returns a boolean if a field has been set.

### GetUserId

`func (o *ModelsMember) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *ModelsMember) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *ModelsMember) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *ModelsMember) HasUserId() bool`

HasUserId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


