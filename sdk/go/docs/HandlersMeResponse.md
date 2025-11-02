# HandlersMeResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvatarUrl** | Pointer to **string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**DisplayName** | Pointer to **string** |  | [optional] 
**Emails** | Pointer to [**[]ModelsUserEmail**](ModelsUserEmail.md) |  | [optional] 
**Id** | Pointer to **string** | example: 3fa85f64-5717-4562-b3fc-2c963f66afa6 | [optional] 
**IsDisabled** | Pointer to **bool** |  | [optional] 
**Organizations** | Pointer to [**[]ModelsOrganization**](ModelsOrganization.md) |  | [optional] 
**PrimaryEmail** | Pointer to **string** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewHandlersMeResponse

`func NewHandlersMeResponse() *HandlersMeResponse`

NewHandlersMeResponse instantiates a new HandlersMeResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHandlersMeResponseWithDefaults

`func NewHandlersMeResponseWithDefaults() *HandlersMeResponse`

NewHandlersMeResponseWithDefaults instantiates a new HandlersMeResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvatarUrl

`func (o *HandlersMeResponse) GetAvatarUrl() string`

GetAvatarUrl returns the AvatarUrl field if non-nil, zero value otherwise.

### GetAvatarUrlOk

`func (o *HandlersMeResponse) GetAvatarUrlOk() (*string, bool)`

GetAvatarUrlOk returns a tuple with the AvatarUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvatarUrl

`func (o *HandlersMeResponse) SetAvatarUrl(v string)`

SetAvatarUrl sets AvatarUrl field to given value.

### HasAvatarUrl

`func (o *HandlersMeResponse) HasAvatarUrl() bool`

HasAvatarUrl returns a boolean if a field has been set.

### GetCreatedAt

`func (o *HandlersMeResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *HandlersMeResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *HandlersMeResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *HandlersMeResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetDisplayName

`func (o *HandlersMeResponse) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *HandlersMeResponse) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *HandlersMeResponse) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.

### HasDisplayName

`func (o *HandlersMeResponse) HasDisplayName() bool`

HasDisplayName returns a boolean if a field has been set.

### GetEmails

`func (o *HandlersMeResponse) GetEmails() []ModelsUserEmail`

GetEmails returns the Emails field if non-nil, zero value otherwise.

### GetEmailsOk

`func (o *HandlersMeResponse) GetEmailsOk() (*[]ModelsUserEmail, bool)`

GetEmailsOk returns a tuple with the Emails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmails

`func (o *HandlersMeResponse) SetEmails(v []ModelsUserEmail)`

SetEmails sets Emails field to given value.

### HasEmails

`func (o *HandlersMeResponse) HasEmails() bool`

HasEmails returns a boolean if a field has been set.

### GetId

`func (o *HandlersMeResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *HandlersMeResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *HandlersMeResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *HandlersMeResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetIsDisabled

`func (o *HandlersMeResponse) GetIsDisabled() bool`

GetIsDisabled returns the IsDisabled field if non-nil, zero value otherwise.

### GetIsDisabledOk

`func (o *HandlersMeResponse) GetIsDisabledOk() (*bool, bool)`

GetIsDisabledOk returns a tuple with the IsDisabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsDisabled

`func (o *HandlersMeResponse) SetIsDisabled(v bool)`

SetIsDisabled sets IsDisabled field to given value.

### HasIsDisabled

`func (o *HandlersMeResponse) HasIsDisabled() bool`

HasIsDisabled returns a boolean if a field has been set.

### GetOrganizations

`func (o *HandlersMeResponse) GetOrganizations() []ModelsOrganization`

GetOrganizations returns the Organizations field if non-nil, zero value otherwise.

### GetOrganizationsOk

`func (o *HandlersMeResponse) GetOrganizationsOk() (*[]ModelsOrganization, bool)`

GetOrganizationsOk returns a tuple with the Organizations field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizations

`func (o *HandlersMeResponse) SetOrganizations(v []ModelsOrganization)`

SetOrganizations sets Organizations field to given value.

### HasOrganizations

`func (o *HandlersMeResponse) HasOrganizations() bool`

HasOrganizations returns a boolean if a field has been set.

### GetPrimaryEmail

`func (o *HandlersMeResponse) GetPrimaryEmail() string`

GetPrimaryEmail returns the PrimaryEmail field if non-nil, zero value otherwise.

### GetPrimaryEmailOk

`func (o *HandlersMeResponse) GetPrimaryEmailOk() (*string, bool)`

GetPrimaryEmailOk returns a tuple with the PrimaryEmail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrimaryEmail

`func (o *HandlersMeResponse) SetPrimaryEmail(v string)`

SetPrimaryEmail sets PrimaryEmail field to given value.

### HasPrimaryEmail

`func (o *HandlersMeResponse) HasPrimaryEmail() bool`

HasPrimaryEmail returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *HandlersMeResponse) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *HandlersMeResponse) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *HandlersMeResponse) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *HandlersMeResponse) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


