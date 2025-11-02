# \TaintsAPI

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateTaint**](TaintsAPI.md#CreateTaint) | **Post** /taints | Create node taint (org scoped)
[**DeleteTaint**](TaintsAPI.md#DeleteTaint) | **Delete** /taints/{id} | Delete taint (org scoped)
[**GetTaint**](TaintsAPI.md#GetTaint) | **Get** /taints/{id} | Get node taint by ID (org scoped)
[**ListTaints**](TaintsAPI.md#ListTaints) | **Get** /taints | List node pool taints (org scoped)
[**UpdateTaint**](TaintsAPI.md#UpdateTaint) | **Patch** /taints/{id} | Update node taint (org scoped)



## CreateTaint

> DtoTaintResponse CreateTaint(ctx).Body(body).XOrgID(xOrgID).Execute()

Create node taint (org scoped)



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
	body := *openapiclient.NewDtoCreateTaintRequest() // DtoCreateTaintRequest | Taint payload
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.CreateTaint(context.Background()).Body(body).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.CreateTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateTaint`: DtoTaintResponse
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.CreateTaint`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateTaintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**DtoCreateTaintRequest**](DtoCreateTaintRequest.md) | Taint payload | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoTaintResponse**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteTaint

> string DeleteTaint(ctx, id).XOrgID(xOrgID).Execute()

Delete taint (org scoped)



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
	id := "id_example" // string | Node Taint ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.DeleteTaint(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.DeleteTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteTaint`: string
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.DeleteTaint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Taint ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteTaintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xOrgID** | **string** | Organization UUID | 

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTaint

> DtoTaintResponse GetTaint(ctx, id).XOrgID(xOrgID).Execute()

Get node taint by ID (org scoped)

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
	id := "id_example" // string | Node Taint ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.GetTaint(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.GetTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTaint`: DtoTaintResponse
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.GetTaint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Taint ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetTaintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoTaintResponse**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTaints

> []DtoTaintResponse ListTaints(ctx).XOrgID(xOrgID).Key(key).Value(value).Q(q).Execute()

List node pool taints (org scoped)



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
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)
	key := "key_example" // string | Exact key (optional)
	value := "value_example" // string | Exact value (optional)
	q := "q_example" // string | key contains (case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.ListTaints(context.Background()).XOrgID(xOrgID).Key(key).Value(value).Q(q).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.ListTaints``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTaints`: []DtoTaintResponse
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.ListTaints`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListTaintsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **key** | **string** | Exact key | 
 **value** | **string** | Exact value | 
 **q** | **string** | key contains (case-insensitive) | 

### Return type

[**[]DtoTaintResponse**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateTaint

> DtoTaintResponse UpdateTaint(ctx, id).Body(body).XOrgID(xOrgID).Execute()

Update node taint (org scoped)



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
	id := "id_example" // string | Node Taint ID (UUID)
	body := *openapiclient.NewDtoUpdateTaintRequest() // DtoUpdateTaintRequest | Fields to update
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.UpdateTaint(context.Background(), id).Body(body).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.UpdateTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateTaint`: DtoTaintResponse
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.UpdateTaint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Taint ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateTaintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**DtoUpdateTaintRequest**](DtoUpdateTaintRequest.md) | Fields to update | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoTaintResponse**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

