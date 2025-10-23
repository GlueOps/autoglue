# AuthnAdminCreateUserRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Password** | Pointer to **string** |  | [optional] 
**Role** | Pointer to **string** | Role allowed values: \&quot;user\&quot; or \&quot;admin\&quot; | [optional] 

## Methods

### NewAuthnAdminCreateUserRequest

`func NewAuthnAdminCreateUserRequest() *AuthnAdminCreateUserRequest`

NewAuthnAdminCreateUserRequest instantiates a new AuthnAdminCreateUserRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthnAdminCreateUserRequestWithDefaults

`func NewAuthnAdminCreateUserRequestWithDefaults() *AuthnAdminCreateUserRequest`

NewAuthnAdminCreateUserRequestWithDefaults instantiates a new AuthnAdminCreateUserRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *AuthnAdminCreateUserRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AuthnAdminCreateUserRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AuthnAdminCreateUserRequest) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *AuthnAdminCreateUserRequest) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetName

`func (o *AuthnAdminCreateUserRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AuthnAdminCreateUserRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AuthnAdminCreateUserRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AuthnAdminCreateUserRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetPassword

`func (o *AuthnAdminCreateUserRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *AuthnAdminCreateUserRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *AuthnAdminCreateUserRequest) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *AuthnAdminCreateUserRequest) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetRole

`func (o *AuthnAdminCreateUserRequest) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *AuthnAdminCreateUserRequest) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *AuthnAdminCreateUserRequest) SetRole(v string)`

SetRole sets Role field to given value.

### HasRole

`func (o *AuthnAdminCreateUserRequest) HasRole() bool`

HasRole returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


