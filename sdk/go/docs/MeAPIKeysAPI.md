# \MeAPIKeysAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateUserAPIKey**](MeAPIKeysAPI.md#CreateUserAPIKey) | **Post** /me/api-keys | Create a new user API key
[**DeleteUserAPIKey**](MeAPIKeysAPI.md#DeleteUserAPIKey) | **Delete** /me/api-keys/{id} | Delete a user API key
[**ListUserAPIKeys**](MeAPIKeysAPI.md#ListUserAPIKeys) | **Get** /me/api-keys | List my API keys



## CreateUserAPIKey

> HandlersUserAPIKeyOut CreateUserAPIKey(ctx).Body(body).Execute()

Create a new user API key



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {
	body := *openapiclient.NewHandlersCreateUserKeyRequest() // HandlersCreateUserKeyRequest | Key options

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MeAPIKeysAPI.CreateUserAPIKey(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MeAPIKeysAPI.CreateUserAPIKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateUserAPIKey`: HandlersUserAPIKeyOut
	fmt.Fprintf(os.Stdout, "Response from `MeAPIKeysAPI.CreateUserAPIKey`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateUserAPIKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**HandlersCreateUserKeyRequest**](HandlersCreateUserKeyRequest.md) | Key options | 

### Return type

[**HandlersUserAPIKeyOut**](HandlersUserAPIKeyOut.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteUserAPIKey

> DeleteUserAPIKey(ctx, id).Execute()

Delete a user API key

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {
	id := "id_example" // string | Key ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MeAPIKeysAPI.DeleteUserAPIKey(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MeAPIKeysAPI.DeleteUserAPIKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Key ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteUserAPIKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListUserAPIKeys

> []HandlersUserAPIKeyOut ListUserAPIKeys(ctx).Execute()

List my API keys

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MeAPIKeysAPI.ListUserAPIKeys(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MeAPIKeysAPI.ListUserAPIKeys``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListUserAPIKeys`: []HandlersUserAPIKeyOut
	fmt.Fprintf(os.Stdout, "Response from `MeAPIKeysAPI.ListUserAPIKeys`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListUserAPIKeysRequest struct via the builder pattern


### Return type

[**[]HandlersUserAPIKeyOut**](HandlersUserAPIKeyOut.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

