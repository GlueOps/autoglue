# ModelsUserEmail

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** | example: 3fa85f64-5717-4562-b3fc-2c963f66afa6 | [optional] 
**IsPrimary** | Pointer to **bool** |  | [optional] 
**IsVerified** | Pointer to **bool** |  | [optional] 
**UpdatedAt** | Pointer to **time.Time** |  | [optional] 
**User** | Pointer to [**ModelsUser**](ModelsUser.md) |  | [optional] 
**UserId** | Pointer to **string** |  | [optional] 

## Methods

### NewModelsUserEmail

`func NewModelsUserEmail() *ModelsUserEmail`

NewModelsUserEmail instantiates a new ModelsUserEmail object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewModelsUserEmailWithDefaults

`func NewModelsUserEmailWithDefaults() *ModelsUserEmail`

NewModelsUserEmailWithDefaults instantiates a new ModelsUserEmail object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCreatedAt

`func (o *ModelsUserEmail) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *ModelsUserEmail) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *ModelsUserEmail) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *ModelsUserEmail) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetEmail

`func (o *ModelsUserEmail) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *ModelsUserEmail) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *ModelsUserEmail) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *ModelsUserEmail) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetId

`func (o *ModelsUserEmail) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ModelsUserEmail) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ModelsUserEmail) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ModelsUserEmail) HasId() bool`

HasId returns a boolean if a field has been set.

### GetIsPrimary

`func (o *ModelsUserEmail) GetIsPrimary() bool`

GetIsPrimary returns the IsPrimary field if non-nil, zero value otherwise.

### GetIsPrimaryOk

`func (o *ModelsUserEmail) GetIsPrimaryOk() (*bool, bool)`

GetIsPrimaryOk returns a tuple with the IsPrimary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPrimary

`func (o *ModelsUserEmail) SetIsPrimary(v bool)`

SetIsPrimary sets IsPrimary field to given value.

### HasIsPrimary

`func (o *ModelsUserEmail) HasIsPrimary() bool`

HasIsPrimary returns a boolean if a field has been set.

### GetIsVerified

`func (o *ModelsUserEmail) GetIsVerified() bool`

GetIsVerified returns the IsVerified field if non-nil, zero value otherwise.

### GetIsVerifiedOk

`func (o *ModelsUserEmail) GetIsVerifiedOk() (*bool, bool)`

GetIsVerifiedOk returns a tuple with the IsVerified field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsVerified

`func (o *ModelsUserEmail) SetIsVerified(v bool)`

SetIsVerified sets IsVerified field to given value.

### HasIsVerified

`func (o *ModelsUserEmail) HasIsVerified() bool`

HasIsVerified returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *ModelsUserEmail) GetUpdatedAt() time.Time`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *ModelsUserEmail) GetUpdatedAtOk() (*time.Time, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *ModelsUserEmail) SetUpdatedAt(v time.Time)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *ModelsUserEmail) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetUser

`func (o *ModelsUserEmail) GetUser() ModelsUser`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *ModelsUserEmail) GetUserOk() (*ModelsUser, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *ModelsUserEmail) SetUser(v ModelsUser)`

SetUser sets User field to given value.

### HasUser

`func (o *ModelsUserEmail) HasUser() bool`

HasUser returns a boolean if a field has been set.

### GetUserId

`func (o *ModelsUserEmail) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *ModelsUserEmail) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *ModelsUserEmail) SetUserId(v string)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *ModelsUserEmail) HasUserId() bool`

HasUserId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


