# DtoCreateSSHRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Bits** | Pointer to **int32** | Only for RSA | [optional] 
**Comment** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** | \&quot;rsa\&quot; (default) or \&quot;ed25519\&quot; | [optional] 

## Methods

### NewDtoCreateSSHRequest

`func NewDtoCreateSSHRequest() *DtoCreateSSHRequest`

NewDtoCreateSSHRequest instantiates a new DtoCreateSSHRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDtoCreateSSHRequestWithDefaults

`func NewDtoCreateSSHRequestWithDefaults() *DtoCreateSSHRequest`

NewDtoCreateSSHRequestWithDefaults instantiates a new DtoCreateSSHRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBits

`func (o *DtoCreateSSHRequest) GetBits() int32`

GetBits returns the Bits field if non-nil, zero value otherwise.

### GetBitsOk

`func (o *DtoCreateSSHRequest) GetBitsOk() (*int32, bool)`

GetBitsOk returns a tuple with the Bits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBits

`func (o *DtoCreateSSHRequest) SetBits(v int32)`

SetBits sets Bits field to given value.

### HasBits

`func (o *DtoCreateSSHRequest) HasBits() bool`

HasBits returns a boolean if a field has been set.

### GetComment

`func (o *DtoCreateSSHRequest) GetComment() string`

GetComment returns the Comment field if non-nil, zero value otherwise.

### GetCommentOk

`func (o *DtoCreateSSHRequest) GetCommentOk() (*string, bool)`

GetCommentOk returns a tuple with the Comment field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComment

`func (o *DtoCreateSSHRequest) SetComment(v string)`

SetComment sets Comment field to given value.

### HasComment

`func (o *DtoCreateSSHRequest) HasComment() bool`

HasComment returns a boolean if a field has been set.

### GetName

`func (o *DtoCreateSSHRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DtoCreateSSHRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DtoCreateSSHRequest) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DtoCreateSSHRequest) HasName() bool`

HasName returns a boolean if a field has been set.

### GetType

`func (o *DtoCreateSSHRequest) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DtoCreateSSHRequest) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DtoCreateSSHRequest) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *DtoCreateSSHRequest) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


