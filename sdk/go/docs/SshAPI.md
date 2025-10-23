# \SshAPI

All URIs are relative to *https://autoglue.apps.nonprod.earth.onglueops.rocks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateSSHKey**](SshAPI.md#CreateSSHKey) | **Post** /api/v1/ssh | Create ssh keypair (org scoped)
[**DeleteSSHKey**](SshAPI.md#DeleteSSHKey) | **Delete** /api/v1/ssh/{id} | Delete ssh keypair (org scoped)
[**DownloadSSHKey**](SshAPI.md#DownloadSSHKey) | **Get** /api/v1/ssh/{id}/download | Download ssh key files by ID (org scoped)
[**GetSSHKey**](SshAPI.md#GetSSHKey) | **Get** /api/v1/ssh/{id} | Get ssh key by ID (org scoped)
[**ListPublicKeys**](SshAPI.md#ListPublicKeys) | **Get** /api/v1/ssh | List ssh keys (org scoped)



## CreateSSHKey

> SshSshResponse CreateSSHKey(ctx).XOrgID(xOrgID).Body(body).Execute()

Create ssh keypair (org scoped)



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
	body := *openapiclient.NewSshCreateSSHRequest() // SshCreateSSHRequest | Key generation options

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SshAPI.CreateSSHKey(context.Background()).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SshAPI.CreateSSHKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateSSHKey`: SshSshResponse
	fmt.Fprintf(os.Stdout, "Response from `SshAPI.CreateSSHKey`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateSSHKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **body** | [**SshCreateSSHRequest**](SshCreateSSHRequest.md) | Key generation options | 

### Return type

[**SshSshResponse**](SshSshResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteSSHKey

> string DeleteSSHKey(ctx, id).XOrgID(xOrgID).Execute()

Delete ssh keypair (org scoped)



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
	id := "id_example" // string | SSH Key ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SshAPI.DeleteSSHKey(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SshAPI.DeleteSSHKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteSSHKey`: string
	fmt.Fprintf(os.Stdout, "Response from `SshAPI.DeleteSSHKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | SSH Key ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteSSHKeyRequest struct via the builder pattern


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


## DownloadSSHKey

> string DownloadSSHKey(ctx, id).XOrgID(xOrgID).Part(part).Execute()

Download ssh key files by ID (org scoped)



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
	id := "id_example" // string | SSH Key ID (UUID)
	part := "part_example" // string | Which part to download

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SshAPI.DownloadSSHKey(context.Background(), id).XOrgID(xOrgID).Part(part).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SshAPI.DownloadSSHKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DownloadSSHKey`: string
	fmt.Fprintf(os.Stdout, "Response from `SshAPI.DownloadSSHKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | SSH Key ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDownloadSSHKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **part** | **string** | Which part to download | 

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


## GetSSHKey

> SshSshRevealResponse GetSSHKey(ctx, id).XOrgID(xOrgID).Reveal(reveal).Execute()

Get ssh key by ID (org scoped)



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
	id := "id_example" // string | SSH Key ID (UUID)
	reveal := true // bool | Reveal private key PEM (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SshAPI.GetSSHKey(context.Background(), id).XOrgID(xOrgID).Reveal(reveal).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SshAPI.GetSSHKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetSSHKey`: SshSshRevealResponse
	fmt.Fprintf(os.Stdout, "Response from `SshAPI.GetSSHKey`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | SSH Key ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetSSHKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **reveal** | **bool** | Reveal private key PEM | 

### Return type

[**SshSshRevealResponse**](SshSshRevealResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListPublicKeys

> []SshSshResponse ListPublicKeys(ctx).XOrgID(xOrgID).Execute()

List ssh keys (org scoped)



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

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.SshAPI.ListPublicKeys(context.Background()).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `SshAPI.ListPublicKeys``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListPublicKeys`: []SshSshResponse
	fmt.Fprintf(os.Stdout, "Response from `SshAPI.ListPublicKeys`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListPublicKeysRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**[]SshSshResponse**](SshSshResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

