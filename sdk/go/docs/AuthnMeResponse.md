# AuthnMeResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Claims** | Pointer to [**AuthnAuthClaimsDTO**](AuthnAuthClaimsDTO.md) |  | [optional] 
**OrgRole** | Pointer to **string** |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**UserId** | Pointer to [**AuthnUserDTO**](AuthnUserDTO.md) |  | [optional] 

## Methods

### NewAuthnMeResponse

`func NewAuthnMeResponse() *AuthnMeResponse`

NewAuthnMeResponse instantiates a new AuthnMeResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthnMeResponseWithDefaults

`func NewAuthnMeResponseWithDefaults() *AuthnMeResponse`

NewAuthnMeResponseWithDefaults instantiates a new AuthnMeResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClaims

`func (o *AuthnMeResponse) GetClaims() AuthnAuthClaimsDTO`

GetClaims returns the Claims field if non-nil, zero value otherwise.

### GetClaimsOk

`func (o *AuthnMeResponse) GetClaimsOk() (*AuthnAuthClaimsDTO, bool)`

GetClaimsOk returns a tuple with the Claims field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClaims

`func (o *AuthnMeResponse) SetClaims(v AuthnAuthClaimsDTO)`

SetClaims sets Claims field to given value.

### HasClaims

`func (o *AuthnMeResponse) HasClaims() bool`

HasClaims returns a boolean if a field has been set.

### GetOrgRole

`func (o *AuthnMeResponse) GetOrgRole() string`

GetOrgRole returns the OrgRole field if non-nil, zero value otherwise.

### GetOrgRoleOk

`func (o *AuthnMeResponse) GetOrgRoleOk() (*string, bool)`

GetOrgRoleOk returns a tuple with the OrgRole field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrgRole

`func (o *AuthnMeResponse) SetOrgRole(v string)`

SetOrgRole sets OrgRole field to given value.

### HasOrgRole

`func (o *AuthnMeResponse) HasOrgRole() bool`

HasOrgRole returns a boolean if a field has been set.

### GetOrganizationId

`func (o *AuthnMeResponse) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *AuthnMeResponse) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *AuthnMeResponse) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *AuthnMeResponse) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetUserId

`func (o *AuthnMeResponse) GetUserId() AuthnUserDTO`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *AuthnMeResponse) GetUserIdOk() (*AuthnUserDTO, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *AuthnMeResponse) SetUserId(v AuthnUserDTO)`

SetUserId sets UserId field to given value.

### HasUserId

`func (o *AuthnMeResponse) HasUserId() bool`

HasUserId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


