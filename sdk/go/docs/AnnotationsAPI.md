# \AnnotationsAPI

All URIs are relative to *https://autoglue.apps.nonprod.earth.onglueops.rocks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddAnnotationToNodePools**](AnnotationsAPI.md#AddAnnotationToNodePools) | **Post** /api/v1/annotations/{id}/node_pools | Attach annotation to node pools (org scoped)
[**CreateAnnotation**](AnnotationsAPI.md#CreateAnnotation) | **Post** /api/v1/annotations | Create annotation (org scoped)
[**DeleteAnnotation**](AnnotationsAPI.md#DeleteAnnotation) | **Delete** /api/v1/annotations/{id} | Delete annotation (org scoped)
[**GetAnnotation**](AnnotationsAPI.md#GetAnnotation) | **Get** /api/v1/annotations/{id} | Get annotation by ID (org scoped)
[**ListAnnotations**](AnnotationsAPI.md#ListAnnotations) | **Get** /api/v1/annotations | List annotations (org scoped)
[**ListNodePoolsWithAnnotation**](AnnotationsAPI.md#ListNodePoolsWithAnnotation) | **Get** /api/v1/annotations/{id}/node_pools | List node pools linked to an annotation (org scoped)
[**RemoveAnnotationFromNodePool**](AnnotationsAPI.md#RemoveAnnotationFromNodePool) | **Delete** /api/v1/annotations/{id}/node_pools/{poolId} | Detach annotation from a node pool (org scoped)
[**UpdateAnnotation**](AnnotationsAPI.md#UpdateAnnotation) | **Patch** /api/v1/annotations/{id} | Update annotation (org scoped)



## AddAnnotationToNodePools

> AnnotationsAnnotationResponse AddAnnotationToNodePools(ctx, id).XOrgID(xOrgID).Body(body).Include(include).Execute()

Attach annotation to node pools (org scoped)



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
	id := "id_example" // string | Annotation ID (UUID)
	body := *openapiclient.NewAnnotationsAddAnnotationToNodePool() // AnnotationsAddAnnotationToNodePool | IDs to attach
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.AddAnnotationToNodePools(context.Background(), id).XOrgID(xOrgID).Body(body).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.AddAnnotationToNodePools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddAnnotationToNodePools`: AnnotationsAnnotationResponse
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.AddAnnotationToNodePools`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Annotation ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddAnnotationToNodePoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**AnnotationsAddAnnotationToNodePool**](AnnotationsAddAnnotationToNodePool.md) | IDs to attach | 
 **include** | **string** | Optional: node_pools | 

### Return type

[**AnnotationsAnnotationResponse**](AnnotationsAnnotationResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateAnnotation

> AnnotationsAnnotationResponse CreateAnnotation(ctx).XOrgID(xOrgID).Body(body).Execute()

Create annotation (org scoped)



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
	body := *openapiclient.NewAnnotationsCreateAnnotationRequest() // AnnotationsCreateAnnotationRequest | Annotation payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.CreateAnnotation(context.Background()).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.CreateAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAnnotation`: AnnotationsAnnotationResponse
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.CreateAnnotation`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAnnotationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **body** | [**AnnotationsCreateAnnotationRequest**](AnnotationsCreateAnnotationRequest.md) | Annotation payload | 

### Return type

[**AnnotationsAnnotationResponse**](AnnotationsAnnotationResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteAnnotation

> string DeleteAnnotation(ctx, id).XOrgID(xOrgID).Execute()

Delete annotation (org scoped)



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
	id := "id_example" // string | Annotation ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.DeleteAnnotation(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.DeleteAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteAnnotation`: string
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.DeleteAnnotation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Annotation ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteAnnotationRequest struct via the builder pattern


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


## GetAnnotation

> AnnotationsAnnotationResponse GetAnnotation(ctx, id).XOrgID(xOrgID).Include(include).Execute()

Get annotation by ID (org scoped)



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
	id := "id_example" // string | Annotation ID (UUID)
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.GetAnnotation(context.Background(), id).XOrgID(xOrgID).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.GetAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAnnotation`: AnnotationsAnnotationResponse
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.GetAnnotation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Annotation ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAnnotationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **include** | **string** | Optional: node_pools | 

### Return type

[**AnnotationsAnnotationResponse**](AnnotationsAnnotationResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAnnotations

> []AnnotationsAnnotationResponse ListAnnotations(ctx).XOrgID(xOrgID).Name(name).Value(value).Q(q).Include(include).Execute()

List annotations (org scoped)



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
	name := "name_example" // string | Exact name (optional)
	value := "value_example" // string | Exact value (optional)
	q := "q_example" // string | name contains (case-insensitive) (optional)
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.ListAnnotations(context.Background()).XOrgID(xOrgID).Name(name).Value(value).Q(q).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.ListAnnotations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAnnotations`: []AnnotationsAnnotationResponse
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.ListAnnotations`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListAnnotationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **name** | **string** | Exact name | 
 **value** | **string** | Exact value | 
 **q** | **string** | name contains (case-insensitive) | 
 **include** | **string** | Optional: node_pools | 

### Return type

[**[]AnnotationsAnnotationResponse**](AnnotationsAnnotationResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolsWithAnnotation

> []AnnotationsNodePoolBrief ListNodePoolsWithAnnotation(ctx, id).XOrgID(xOrgID).Q(q).Execute()

List node pools linked to an annotation (org scoped)



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
	id := "id_example" // string | Annotation ID (UUID)
	q := "q_example" // string | Name contains (case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.ListNodePoolsWithAnnotation(context.Background(), id).XOrgID(xOrgID).Q(q).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.ListNodePoolsWithAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolsWithAnnotation`: []AnnotationsNodePoolBrief
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.ListNodePoolsWithAnnotation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Annotation ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolsWithAnnotationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **q** | **string** | Name contains (case-insensitive) | 

### Return type

[**[]AnnotationsNodePoolBrief**](AnnotationsNodePoolBrief.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveAnnotationFromNodePool

> string RemoveAnnotationFromNodePool(ctx, id, poolId).XOrgID(xOrgID).Execute()

Detach annotation from a node pool (org scoped)



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
	id := "id_example" // string | Annotation ID (UUID)
	poolId := "poolId_example" // string | Node Pool ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.RemoveAnnotationFromNodePool(context.Background(), id, poolId).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.RemoveAnnotationFromNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveAnnotationFromNodePool`: string
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.RemoveAnnotationFromNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Annotation ID (UUID) | 
**poolId** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveAnnotationFromNodePoolRequest struct via the builder pattern


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


## UpdateAnnotation

> AnnotationsAnnotationResponse UpdateAnnotation(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Update annotation (org scoped)



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
	id := "id_example" // string | Annotation ID (UUID)
	body := *openapiclient.NewAnnotationsUpdateAnnotationRequest() // AnnotationsUpdateAnnotationRequest | Fields to update

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.UpdateAnnotation(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.UpdateAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateAnnotation`: AnnotationsAnnotationResponse
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.UpdateAnnotation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Annotation ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateAnnotationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**AnnotationsUpdateAnnotationRequest**](AnnotationsUpdateAnnotationRequest.md) | Fields to update | 

### Return type

[**AnnotationsAnnotationResponse**](AnnotationsAnnotationResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

