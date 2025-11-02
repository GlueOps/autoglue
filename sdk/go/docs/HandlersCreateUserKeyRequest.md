# HandlersCreateUserKeyRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExpiresInHours** | Pointer to **int32** | optional TTL | [optional] 
**Name** | Pointer to **string** |  | [optional] 

## Methods

### NewHandlersCreateUserKeyRequest

`func NewHandlersCreateUserKeyRequest() *HandlersCreateUserKeyRequest`

NewHandlersCreateUserKeyRequest instantiates a new HandlersCreateUserKeyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHandlersCreateUserKeyRequestWithDefaults

`func NewHandlersCreateUserKeyRequestWithDefaults() *HandlersCreateUserKeyRequest`

NewHandlersCreateUserKeyRequestWithDefaults instantiates a new HandlersCreateUserKeyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExpiresInHours

`func (o *HandlersCreateUserKeyRequest) GetExpiresInHours() int32`

GetExpiresInHours returns the ExpiresInHours field if non-nil, zero value otherwise.

### GetExpiresInHoursOk

`func (o *HandlersCreateUserKeyRequest) GetExpiresInHoursOk() (*int32, bool)`

GetExpiresInHoursOk returns a tuple with the ExpiresInHours field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresInHours

`func (o *HandlersCreateUserKeyRequest) SetExpiresInHours(v int32)`

SetExpiresInHours sets ExpiresInHours field to given value.

### HasExpiresInHours

`func (o *HandlersCreateUserKeyRequest) HasExpiresInHours() bool`

HasExpiresInHours returns a boolean if a field has been set.

### GetName

`func (o *HandlersCreateUserKeyRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *HandlersCreateUserKeyRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *HandlersCreateUserKeyRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *HandlersCreateUserKeyRequest) HasName() bool`

HasName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


