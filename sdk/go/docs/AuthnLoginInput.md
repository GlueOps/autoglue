# AuthnLoginInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | Pointer to **string** |  | [optional] 
**Password** | Pointer to **string** |  | [optional] 

## Methods

### NewAuthnLoginInput

`func NewAuthnLoginInput() *AuthnLoginInput`

NewAuthnLoginInput instantiates a new AuthnLoginInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthnLoginInputWithDefaults

`func NewAuthnLoginInputWithDefaults() *AuthnLoginInput`

NewAuthnLoginInputWithDefaults instantiates a new AuthnLoginInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *AuthnLoginInput) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *AuthnLoginInput) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *AuthnLoginInput) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *AuthnLoginInput) HasEmail() bool`

HasEmail returns a boolean if a field has been set.

### GetPassword

`func (o *AuthnLoginInput) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *AuthnLoginInput) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *AuthnLoginInput) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *AuthnLoginInput) HasPassword() bool`

HasPassword returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


