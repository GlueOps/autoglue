# \AdminAPI

All URIs are relative to *https://autoglue.apps.nonprod.earth.onglueops.rocks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AdminCreateUser**](AdminAPI.md#AdminCreateUser) | **Post** /api/v1/admin/users | Admin: create user
[**AdminDeleteUser**](AdminAPI.md#AdminDeleteUser) | **Delete** /api/v1/admin/users/{userId} | Admin: delete user
[**AdminListUsers**](AdminAPI.md#AdminListUsers) | **Get** /api/v1/admin/users | Admin: list all users
[**AdminUpdateUser**](AdminAPI.md#AdminUpdateUser) | **Patch** /api/v1/admin/users/{userId} | Admin: update user



## AdminCreateUser

> AuthnUserOut AdminCreateUser(ctx).Body(body).Execute()

Admin: create user

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/autoglue-sdk"
)

func main() {
	body := *openapiclient.NewAuthnAdminCreateUserRequest() // AuthnAdminCreateUserRequest | payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.AdminCreateUser(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.AdminCreateUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminCreateUser`: AuthnUserOut
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.AdminCreateUser`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAdminCreateUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AuthnAdminCreateUserRequest**](AuthnAdminCreateUserRequest.md) | payload | 

### Return type

[**AuthnUserOut**](AuthnUserOut.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminDeleteUser

> string AdminDeleteUser(ctx, userId).Execute()

Admin: delete user

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/autoglue-sdk"
)

func main() {
	userId := "userId_example" // string | User ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.AdminDeleteUser(context.Background(), userId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.AdminDeleteUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminDeleteUser`: string
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.AdminDeleteUser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | User ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminDeleteUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminListUsers

> AuthnListUsersOut AdminListUsers(ctx).Page(page).PageSize(pageSize).Execute()

Admin: list all users



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/autoglue-sdk"
)

func main() {
	page := int32(56) // int32 | Page number (1-based) (optional)
	pageSize := int32(56) // int32 | Page size (max 200) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.AdminListUsers(context.Background()).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.AdminListUsers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminListUsers`: AuthnListUsersOut
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.AdminListUsers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAdminListUsersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** | Page number (1-based) | 
 **pageSize** | **int32** | Page size (max 200) | 

### Return type

[**AuthnListUsersOut**](AuthnListUsersOut.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminUpdateUser

> AuthnUserOut AdminUpdateUser(ctx, userId).Body(body).Execute()

Admin: update user

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/autoglue-sdk"
)

func main() {
	userId := "userId_example" // string | User ID
	body := *openapiclient.NewAuthnAdminUpdateUserRequest() // AuthnAdminUpdateUserRequest | payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AdminAPI.AdminUpdateUser(context.Background(), userId).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AdminAPI.AdminUpdateUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminUpdateUser`: AuthnUserOut
	fmt.Fprintf(os.Stdout, "Response from `AdminAPI.AdminUpdateUser`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**userId** | **string** | User ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminUpdateUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**AuthnAdminUpdateUserRequest**](AuthnAdminUpdateUserRequest.md) | payload | 

### Return type

[**AuthnUserOut**](AuthnUserOut.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

