# ServersUpdateServerRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Hostname** | Pointer to **string** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**Role** | Pointer to **string** |  | [optional] 
**SshKeyId** | Pointer to **string** |  | [optional] 
**SshUser** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** | enum: pending,provisioning,ready,failed | [optional] 

## Methods

### NewServersUpdateServerRequest

`func NewServersUpdateServerRequest() *ServersUpdateServerRequest`

NewServersUpdateServerRequest instantiates a new ServersUpdateServerRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServersUpdateServerRequestWithDefaults

`func NewServersUpdateServerRequestWithDefaults() *ServersUpdateServerRequest`

NewServersUpdateServerRequestWithDefaults instantiates a new ServersUpdateServerRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHostname

`func (o *ServersUpdateServerRequest) GetHostname() string`

GetHostname returns the Hostname field if non-nil, zero value otherwise.

### GetHostnameOk

`func (o *ServersUpdateServerRequest) GetHostnameOk() (*string, bool)`

GetHostnameOk returns a tuple with the Hostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostname

`func (o *ServersUpdateServerRequest) SetHostname(v string)`

SetHostname sets Hostname field to given value.

### HasHostname

`func (o *ServersUpdateServerRequest) HasHostname() bool`

HasHostname returns a boolean if a field has been set.

### GetIpAddress

`func (o *ServersUpdateServerRequest) GetIpAddress() string`

GetIpAddress returns the IpAddress field if non-nil, zero value otherwise.

### GetIpAddressOk

`func (o *ServersUpdateServerRequest) GetIpAddressOk() (*string, bool)`

GetIpAddressOk returns a tuple with the IpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpAddress

`func (o *ServersUpdateServerRequest) SetIpAddress(v string)`

SetIpAddress sets IpAddress field to given value.

### HasIpAddress

`func (o *ServersUpdateServerRequest) HasIpAddress() bool`

HasIpAddress returns a boolean if a field has been set.

### GetRole

`func (o *ServersUpdateServerRequest) GetRole() string`

GetRole returns the Role field if non-nil, zero value otherwise.

### GetRoleOk

`func (o *ServersUpdateServerRequest) GetRoleOk() (*string, bool)`

GetRoleOk returns a tuple with the Role field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRole

`func (o *ServersUpdateServerRequest) SetRole(v string)`

SetRole sets Role field to given value.

### HasRole

`func (o *ServersUpdateServerRequest) HasRole() bool`

HasRole returns a boolean if a field has been set.

### GetSshKeyId

`func (o *ServersUpdateServerRequest) GetSshKeyId() string`

GetSshKeyId returns the SshKeyId field if non-nil, zero value otherwise.

### GetSshKeyIdOk

`func (o *ServersUpdateServerRequest) GetSshKeyIdOk() (*string, bool)`

GetSshKeyIdOk returns a tuple with the SshKeyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyId

`func (o *ServersUpdateServerRequest) SetSshKeyId(v string)`

SetSshKeyId sets SshKeyId field to given value.

### HasSshKeyId

`func (o *ServersUpdateServerRequest) HasSshKeyId() bool`

HasSshKeyId returns a boolean if a field has been set.

### GetSshUser

`func (o *ServersUpdateServerRequest) GetSshUser() string`

GetSshUser returns the SshUser field if non-nil, zero value otherwise.

### GetSshUserOk

`func (o *ServersUpdateServerRequest) GetSshUserOk() (*string, bool)`

GetSshUserOk returns a tuple with the SshUser field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshUser

`func (o *ServersUpdateServerRequest) SetSshUser(v string)`

SetSshUser sets SshUser field to given value.

### HasSshUser

`func (o *ServersUpdateServerRequest) HasSshUser() bool`

HasSshUser returns a boolean if a field has been set.

### GetStatus

`func (o *ServersUpdateServerRequest) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ServersUpdateServerRequest) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ServersUpdateServerRequest) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ServersUpdateServerRequest) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


