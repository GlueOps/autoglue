# \TaintsAPI

All URIs are relative to *https://autoglue.apps.nonprod.earth.onglueops.rocks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddTaintToNodePool**](TaintsAPI.md#AddTaintToNodePool) | **Post** /api/v1/taints/{id}/node_pools | Attach taint to node pools (org scoped)
[**CreateTaint**](TaintsAPI.md#CreateTaint) | **Post** /api/v1/taints | Create node taint (org scoped)
[**DeleteTaint**](TaintsAPI.md#DeleteTaint) | **Delete** /api/v1/taints/{id} | Delete taint (org scoped)
[**GetTaint**](TaintsAPI.md#GetTaint) | **Get** /api/v1/taints/{id} | Get node taint by ID (org scoped)
[**ListNodePoolsWithTaint**](TaintsAPI.md#ListNodePoolsWithTaint) | **Get** /api/v1/taints/{id}/node_pools | List node pools linked to a taint (org scoped)
[**ListTaints**](TaintsAPI.md#ListTaints) | **Get** /api/v1/taints | List node taints (org scoped)
[**RemoveTaintFromNodePool**](TaintsAPI.md#RemoveTaintFromNodePool) | **Delete** /api/v1/taints/{id}/node_pools/{poolId} | Detach taint from a node pool (org scoped)
[**UpdateTaint**](TaintsAPI.md#UpdateTaint) | **Patch** /api/v1/taints/{id} | Update node taint (org scoped)



## AddTaintToNodePool

> TaintsTaintResponse AddTaintToNodePool(ctx, id).XOrgID(xOrgID).Body(body).Include(include).Execute()

Attach taint to node pools (org scoped)



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
	id := "id_example" // string | Taint ID (UUID)
	body := *openapiclient.NewTaintsAddTaintToPoolRequest() // TaintsAddTaintToPoolRequest | IDs to attach
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.AddTaintToNodePool(context.Background(), id).XOrgID(xOrgID).Body(body).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.AddTaintToNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddTaintToNodePool`: TaintsTaintResponse
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.AddTaintToNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Taint ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddTaintToNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**TaintsAddTaintToPoolRequest**](TaintsAddTaintToPoolRequest.md) | IDs to attach | 
 **include** | **string** | Optional: node_pools | 

### Return type

[**TaintsTaintResponse**](TaintsTaintResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateTaint

> TaintsTaintResponse CreateTaint(ctx).XOrgID(xOrgID).Body(body).Execute()

Create node taint (org scoped)



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
	body := *openapiclient.NewTaintsCreateTaintRequest() // TaintsCreateTaintRequest | Taint payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.CreateTaint(context.Background()).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.CreateTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateTaint`: TaintsTaintResponse
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.CreateTaint`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateTaintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **body** | [**TaintsCreateTaintRequest**](TaintsCreateTaintRequest.md) | Taint payload | 

### Return type

[**TaintsTaintResponse**](TaintsTaintResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

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
	openapiclient "github.com/glueops/autoglue-sdk"
)

func main() {
	xOrgID := "xOrgID_example" // string | Organization UUID
	id := "id_example" // string | Node Taint ID (UUID)

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

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetTaint

> TaintsTaintResponse GetTaint(ctx, id).XOrgID(xOrgID).Include(include).Execute()

Get node taint by ID (org scoped)



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
	id := "id_example" // string | Node Taint ID (UUID)
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.GetTaint(context.Background(), id).XOrgID(xOrgID).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.GetTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetTaint`: TaintsTaintResponse
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

 **include** | **string** | Optional: node_pools | 

### Return type

[**TaintsTaintResponse**](TaintsTaintResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolsWithTaint

> []TaintsNodePoolResponse ListNodePoolsWithTaint(ctx, id).XOrgID(xOrgID).Q(q).Execute()

List node pools linked to a taint (org scoped)



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
	id := "id_example" // string | Taint ID (UUID)
	q := "q_example" // string | Name contains (case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.ListNodePoolsWithTaint(context.Background(), id).XOrgID(xOrgID).Q(q).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.ListNodePoolsWithTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolsWithTaint`: []TaintsNodePoolResponse
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.ListNodePoolsWithTaint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Taint ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolsWithTaintRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **q** | **string** | Name contains (case-insensitive) | 

### Return type

[**[]TaintsNodePoolResponse**](TaintsNodePoolResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListTaints

> []TaintsTaintResponse ListTaints(ctx).XOrgID(xOrgID).Key(key).Value(value).Q(q).Include(include).Execute()

List node taints (org scoped)



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
	key := "key_example" // string | Exact key (optional)
	value := "value_example" // string | Exact value (optional)
	q := "q_example" // string | key contains (case-insensitive) (optional)
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.ListTaints(context.Background()).XOrgID(xOrgID).Key(key).Value(value).Q(q).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.ListTaints``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListTaints`: []TaintsTaintResponse
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
 **include** | **string** | Optional: node_pools | 

### Return type

[**[]TaintsTaintResponse**](TaintsTaintResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveTaintFromNodePool

> string RemoveTaintFromNodePool(ctx, id, poolId).XOrgID(xOrgID).Execute()

Detach taint from a node pool (org scoped)



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
	id := "id_example" // string | Taint ID (UUID)
	poolId := "poolId_example" // string | Node Pool ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.RemoveTaintFromNodePool(context.Background(), id, poolId).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.RemoveTaintFromNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveTaintFromNodePool`: string
	fmt.Fprintf(os.Stdout, "Response from `TaintsAPI.RemoveTaintFromNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Taint ID (UUID) | 
**poolId** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveTaintFromNodePoolRequest struct via the builder pattern


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


## UpdateTaint

> TaintsTaintResponse UpdateTaint(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Update node taint (org scoped)



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
	id := "id_example" // string | Node Taint ID (UUID)
	body := *openapiclient.NewTaintsUpdateTaintRequest() // TaintsUpdateTaintRequest | Fields to update

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TaintsAPI.UpdateTaint(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TaintsAPI.UpdateTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateTaint`: TaintsTaintResponse
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
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**TaintsUpdateTaintRequest**](TaintsUpdateTaintRequest.md) | Fields to update | 

### Return type

[**TaintsTaintResponse**](TaintsTaintResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

