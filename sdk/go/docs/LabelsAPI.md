# \LabelsAPI

All URIs are relative to *https://autoglue.apps.nonprod.earth.onglueops.rocks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddLabelToNodePool**](LabelsAPI.md#AddLabelToNodePool) | **Post** /api/v1/labels/{id}/node_pools | Attach label to node pools (org scoped)
[**CreateLabel**](LabelsAPI.md#CreateLabel) | **Post** /api/v1/labels | Create label (org scoped)
[**DeleteLabel**](LabelsAPI.md#DeleteLabel) | **Delete** /api/v1/labels/{id} | Delete label (org scoped)
[**GetLabel**](LabelsAPI.md#GetLabel) | **Get** /api/v1/labels/{id} | Get label by ID (org scoped)
[**ListLabels**](LabelsAPI.md#ListLabels) | **Get** /api/v1/labels | List node labels (org scoped)
[**ListNodePoolsWithLabel**](LabelsAPI.md#ListNodePoolsWithLabel) | **Get** /api/v1/labels/{id}/node_pools | List node pools linked to a label (org scoped)
[**RemoveLabelFromNodePool**](LabelsAPI.md#RemoveLabelFromNodePool) | **Delete** /api/v1/labels/{id}/node_pools/{poolId} | Detach label from a node pool (org scoped)
[**UpdateLabel**](LabelsAPI.md#UpdateLabel) | **Patch** /api/v1/labels/{id} | Update label (org scoped)



## AddLabelToNodePool

> LabelsLabelResponse AddLabelToNodePool(ctx, id).XOrgID(xOrgID).Body(body).Include(include).Execute()

Attach label to node pools (org scoped)



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
	id := "id_example" // string | Label ID (UUID)
	body := *openapiclient.NewLabelsAddLabelToPoolRequest() // LabelsAddLabelToPoolRequest | IDs to attach
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.AddLabelToNodePool(context.Background(), id).XOrgID(xOrgID).Body(body).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.AddLabelToNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AddLabelToNodePool`: LabelsLabelResponse
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.AddLabelToNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Label ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAddLabelToNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**LabelsAddLabelToPoolRequest**](LabelsAddLabelToPoolRequest.md) | IDs to attach | 
 **include** | **string** | Optional: node_pools | 

### Return type

[**LabelsLabelResponse**](LabelsLabelResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateLabel

> LabelsLabelResponse CreateLabel(ctx).XOrgID(xOrgID).Body(body).Execute()

Create label (org scoped)



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
	body := *openapiclient.NewLabelsCreateLabelRequest() // LabelsCreateLabelRequest | Label payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.CreateLabel(context.Background()).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.CreateLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateLabel`: LabelsLabelResponse
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.CreateLabel`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateLabelRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **body** | [**LabelsCreateLabelRequest**](LabelsCreateLabelRequest.md) | Label payload | 

### Return type

[**LabelsLabelResponse**](LabelsLabelResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteLabel

> string DeleteLabel(ctx, id).XOrgID(xOrgID).Execute()

Delete label (org scoped)



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
	id := "id_example" // string | Label ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.DeleteLabel(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.DeleteLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteLabel`: string
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.DeleteLabel`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Label ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteLabelRequest struct via the builder pattern


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


## GetLabel

> LabelsLabelResponse GetLabel(ctx, id).XOrgID(xOrgID).Include(include).Execute()

Get label by ID (org scoped)



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
	id := "id_example" // string | Label ID (UUID)
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.GetLabel(context.Background(), id).XOrgID(xOrgID).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.GetLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLabel`: LabelsLabelResponse
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.GetLabel`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Label ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetLabelRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **include** | **string** | Optional: node_pools | 

### Return type

[**LabelsLabelResponse**](LabelsLabelResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLabels

> []LabelsLabelResponse ListLabels(ctx).XOrgID(xOrgID).Key(key).Value(value).Q(q).Include(include).Execute()

List node labels (org scoped)



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
	q := "q_example" // string | Key contains (case-insensitive) (optional)
	include := "include_example" // string | Optional: node_pools (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.ListLabels(context.Background()).XOrgID(xOrgID).Key(key).Value(value).Q(q).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.ListLabels``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLabels`: []LabelsLabelResponse
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.ListLabels`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListLabelsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **key** | **string** | Exact key | 
 **value** | **string** | Exact value | 
 **q** | **string** | Key contains (case-insensitive) | 
 **include** | **string** | Optional: node_pools | 

### Return type

[**[]LabelsLabelResponse**](LabelsLabelResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolsWithLabel

> []LabelsNodePoolBrief ListNodePoolsWithLabel(ctx, id).XOrgID(xOrgID).Q(q).Execute()

List node pools linked to a label (org scoped)



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
	id := "id_example" // string | Label ID (UUID)
	q := "q_example" // string | Name contains (case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.ListNodePoolsWithLabel(context.Background(), id).XOrgID(xOrgID).Q(q).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.ListNodePoolsWithLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolsWithLabel`: []LabelsNodePoolBrief
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.ListNodePoolsWithLabel`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Label ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolsWithLabelRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **q** | **string** | Name contains (case-insensitive) | 

### Return type

[**[]LabelsNodePoolBrief**](LabelsNodePoolBrief.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RemoveLabelFromNodePool

> string RemoveLabelFromNodePool(ctx, id, poolId).XOrgID(xOrgID).Execute()

Detach label from a node pool (org scoped)



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
	id := "id_example" // string | Label ID (UUID)
	poolId := "poolId_example" // string | Node Pool ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.RemoveLabelFromNodePool(context.Background(), id, poolId).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.RemoveLabelFromNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `RemoveLabelFromNodePool`: string
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.RemoveLabelFromNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Label ID (UUID) | 
**poolId** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiRemoveLabelFromNodePoolRequest struct via the builder pattern


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


## UpdateLabel

> LabelsLabelResponse UpdateLabel(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Update label (org scoped)



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
	id := "id_example" // string | Label ID (UUID)
	body := *openapiclient.NewLabelsUpdateLabelRequest() // LabelsUpdateLabelRequest | Fields to update

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.UpdateLabel(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.UpdateLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLabel`: LabelsLabelResponse
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.UpdateLabel`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Label ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateLabelRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**LabelsUpdateLabelRequest**](LabelsUpdateLabelRequest.md) | Fields to update | 

### Return type

[**LabelsLabelResponse**](LabelsLabelResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

