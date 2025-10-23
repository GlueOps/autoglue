# \ServersAPI

All URIs are relative to *https://autoglue.apps.nonprod.earth.onglueops.rocks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateServer**](ServersAPI.md#CreateServer) | **Post** /api/v1/servers | Create server (org scoped)
[**DeleteServer**](ServersAPI.md#DeleteServer) | **Delete** /api/v1/servers/{id} | Delete server (org scoped)
[**GetServer**](ServersAPI.md#GetServer) | **Get** /api/v1/servers/{id} | Get server by ID (org scoped)
[**ListServers**](ServersAPI.md#ListServers) | **Get** /api/v1/servers | List servers (org scoped)
[**UpdateServer**](ServersAPI.md#UpdateServer) | **Patch** /api/v1/servers/{id} | Update server (org scoped)



## CreateServer

> ServersServerResponse CreateServer(ctx).XOrgID(xOrgID).Body(body).Execute()

Create server (org scoped)



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
	xOrgID := "xOrgID_example" // string | Organization UUID
	body := *openapiclient.NewServersCreateServerRequest() // ServersCreateServerRequest | Server payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServersAPI.CreateServer(context.Background()).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServersAPI.CreateServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateServer`: ServersServerResponse
	fmt.Fprintf(os.Stdout, "Response from `ServersAPI.CreateServer`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateServerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **body** | [**ServersCreateServerRequest**](ServersCreateServerRequest.md) | Server payload | 

### Return type

[**ServersServerResponse**](ServersServerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteServer

> string DeleteServer(ctx, id).XOrgID(xOrgID).Execute()

Delete server (org scoped)



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
	xOrgID := "xOrgID_example" // string | Organization UUID
	id := "id_example" // string | Server ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServersAPI.DeleteServer(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServersAPI.DeleteServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteServer`: string
	fmt.Fprintf(os.Stdout, "Response from `ServersAPI.DeleteServer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Server ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteServerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


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


## GetServer

> ServersServerResponse GetServer(ctx, id).XOrgID(xOrgID).Execute()

Get server by ID (org scoped)



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
	xOrgID := "xOrgID_example" // string | Organization UUID
	id := "id_example" // string | Server ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServersAPI.GetServer(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServersAPI.GetServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetServer`: ServersServerResponse
	fmt.Fprintf(os.Stdout, "Response from `ServersAPI.GetServer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Server ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetServerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


### Return type

[**ServersServerResponse**](ServersServerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListServers

> []ServersServerResponse ListServers(ctx).XOrgID(xOrgID).Status(status).Role(role).Execute()

List servers (org scoped)



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
	xOrgID := "xOrgID_example" // string | Organization UUID
	status := "status_example" // string | Filter by status (pending|provisioning|ready|failed) (optional)
	role := "role_example" // string | Filter by role (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServersAPI.ListServers(context.Background()).XOrgID(xOrgID).Status(status).Role(role).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServersAPI.ListServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListServers`: []ServersServerResponse
	fmt.Fprintf(os.Stdout, "Response from `ServersAPI.ListServers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **status** | **string** | Filter by status (pending|provisioning|ready|failed) | 
 **role** | **string** | Filter by role | 

### Return type

[**[]ServersServerResponse**](ServersServerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateServer

> ServersServerResponse UpdateServer(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Update server (org scoped)



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
	xOrgID := "xOrgID_example" // string | Organization UUID
	id := "id_example" // string | Server ID (UUID)
	body := *openapiclient.NewServersUpdateServerRequest() // ServersUpdateServerRequest | Fields to update

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ServersAPI.UpdateServer(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ServersAPI.UpdateServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateServer`: ServersServerResponse
	fmt.Fprintf(os.Stdout, "Response from `ServersAPI.UpdateServer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Server ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateServerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**ServersUpdateServerRequest**](ServersUpdateServerRequest.md) | Fields to update | 

### Return type

[**ServersServerResponse**](ServersServerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

