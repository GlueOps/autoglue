# AuthnListUsersOut

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Page** | Pointer to **int32** |  | [optional] 
**PageSize** | Pointer to **int32** |  | [optional] 
**Total** | Pointer to **int32** |  | [optional] 
**Users** | Pointer to [**[]AuthnUserListItem**](AuthnUserListItem.md) |  | [optional] 

## Methods

### NewAuthnListUsersOut

`func NewAuthnListUsersOut() *AuthnListUsersOut`

NewAuthnListUsersOut instantiates a new AuthnListUsersOut object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthnListUsersOutWithDefaults

`func NewAuthnListUsersOutWithDefaults() *AuthnListUsersOut`

NewAuthnListUsersOutWithDefaults instantiates a new AuthnListUsersOut object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPage

`func (o *AuthnListUsersOut) GetPage() int32`

GetPage returns the Page field if non-nil, zero value otherwise.

### GetPageOk

`func (o *AuthnListUsersOut) GetPageOk() (*int32, bool)`

GetPageOk returns a tuple with the Page field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPage

`func (o *AuthnListUsersOut) SetPage(v int32)`

SetPage sets Page field to given value.

### HasPage

`func (o *AuthnListUsersOut) HasPage() bool`

HasPage returns a boolean if a field has been set.

### GetPageSize

`func (o *AuthnListUsersOut) GetPageSize() int32`

GetPageSize returns the PageSize field if non-nil, zero value otherwise.

### GetPageSizeOk

`func (o *AuthnListUsersOut) GetPageSizeOk() (*int32, bool)`

GetPageSizeOk returns a tuple with the PageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageSize

`func (o *AuthnListUsersOut) SetPageSize(v int32)`

SetPageSize sets PageSize field to given value.

### HasPageSize

`func (o *AuthnListUsersOut) HasPageSize() bool`

HasPageSize returns a boolean if a field has been set.

### GetTotal

`func (o *AuthnListUsersOut) GetTotal() int32`

GetTotal returns the Total field if non-nil, zero value otherwise.

### GetTotalOk

`func (o *AuthnListUsersOut) GetTotalOk() (*int32, bool)`

GetTotalOk returns a tuple with the Total field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotal

`func (o *AuthnListUsersOut) SetTotal(v int32)`

SetTotal sets Total field to given value.

### HasTotal

`func (o *AuthnListUsersOut) HasTotal() bool`

HasTotal returns a boolean if a field has been set.

### GetUsers

`func (o *AuthnListUsersOut) GetUsers() []AuthnUserListItem`

GetUsers returns the Users field if non-nil, zero value otherwise.

### GetUsersOk

`func (o *AuthnListUsersOut) GetUsersOk() (*[]AuthnUserListItem, bool)`

GetUsersOk returns a tuple with the Users field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsers

`func (o *AuthnListUsersOut) SetUsers(v []AuthnUserListItem)`

SetUsers sets Users field to given value.

### HasUsers

`func (o *AuthnListUsersOut) HasUsers() bool`

HasUsers returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


