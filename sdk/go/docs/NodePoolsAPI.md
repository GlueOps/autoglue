# \NodePoolsAPI

All URIs are relative to *https://autoglue.apps.nonprod.earth.onglueops.rocks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AttachNodePoolAnnotations**](NodePoolsAPI.md#AttachNodePoolAnnotations) | **Post** /api/v1/node-pools/{id}/annotations | Attach annotations to a node pool (org scoped)
[**AttachNodePoolLabels**](NodePoolsAPI.md#AttachNodePoolLabels) | **Post** /api/v1/node-pools/{id}/labels | Attach labels to a node pool (org scoped)
[**AttachNodePoolServers**](NodePoolsAPI.md#AttachNodePoolServers) | **Post** /api/v1/node-pools/{id}/servers | Attach servers to a node pool (org scoped)
[**AttachNodePoolTaints**](NodePoolsAPI.md#AttachNodePoolTaints) | **Post** /api/v1/node-pools/{id}/taints | Attach taints to a node pool (org scoped)
[**CreateNodePool**](NodePoolsAPI.md#CreateNodePool) | **Post** /api/v1/node-pools | Create node group (org scoped)
[**DeleteNodePool**](NodePoolsAPI.md#DeleteNodePool) | **Delete** /api/v1/node-pools/{id} | Delete node pool (org scoped)
[**DetachNodePoolAnnotation**](NodePoolsAPI.md#DetachNodePoolAnnotation) | **Delete** /api/v1/node-pools/{id}/annotations/{annotationId} | Detach one annotation from a node pool (org scoped)
[**DetachNodePoolLabel**](NodePoolsAPI.md#DetachNodePoolLabel) | **Delete** /api/v1/node-pools/{id}/labels/{labelId} | Detach one label from a node pool (org scoped)
[**DetachNodePoolServer**](NodePoolsAPI.md#DetachNodePoolServer) | **Delete** /api/v1/node-pools/{id}/servers/{serverId} | Detach one server from a node pool (org scoped)
[**DetachNodePoolTaint**](NodePoolsAPI.md#DetachNodePoolTaint) | **Delete** /api/v1/node-pools/{id}/taints/{taintId} | Detach one taint from a node pool (org scoped)
[**GetNodePool**](NodePoolsAPI.md#GetNodePool) | **Get** /api/v1/node-pools/{id} | Get node group by ID (org scoped)
[**ListNodePoolAnnotations**](NodePoolsAPI.md#ListNodePoolAnnotations) | **Get** /api/v1/node-pools/{id}/annotations | List annotations attached to a node pool (org scoped)
[**ListNodePoolLabels**](NodePoolsAPI.md#ListNodePoolLabels) | **Get** /api/v1/node-pools/{id}/labels | List labels attached to a node pool (org scoped)
[**ListNodePoolServers**](NodePoolsAPI.md#ListNodePoolServers) | **Get** /api/v1/node-pools/{id}/servers | List servers attached to a node pool (org scoped)
[**ListNodePoolTaints**](NodePoolsAPI.md#ListNodePoolTaints) | **Get** /api/v1/node-pools/{id}/taints | List taints attached to a node pool (org scoped)
[**ListNodePools**](NodePoolsAPI.md#ListNodePools) | **Get** /api/v1/node-pools | List node pools (org scoped)
[**UpdateNodePool**](NodePoolsAPI.md#UpdateNodePool) | **Patch** /api/v1/node-pools/{id} | Update node pool (org scoped)



## AttachNodePoolAnnotations

> string AttachNodePoolAnnotations(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Attach annotations to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	body := *openapiclient.NewNodepoolsAttachAnnotationsRequest() // NodepoolsAttachAnnotationsRequest | Annotation IDs to attach

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.AttachNodePoolAnnotations(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.AttachNodePoolAnnotations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AttachNodePoolAnnotations`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.AttachNodePoolAnnotations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAttachNodePoolAnnotationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**NodepoolsAttachAnnotationsRequest**](NodepoolsAttachAnnotationsRequest.md) | Annotation IDs to attach | 

### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AttachNodePoolLabels

> string AttachNodePoolLabels(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Attach labels to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	body := *openapiclient.NewNodepoolsAttachLabelsRequest() // NodepoolsAttachLabelsRequest | Label IDs to attach

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.AttachNodePoolLabels(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.AttachNodePoolLabels``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AttachNodePoolLabels`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.AttachNodePoolLabels`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAttachNodePoolLabelsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**NodepoolsAttachLabelsRequest**](NodepoolsAttachLabelsRequest.md) | Label IDs to attach | 

### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AttachNodePoolServers

> string AttachNodePoolServers(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Attach servers to a node pool (org scoped)

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
	id := "id_example" // string | Node Group ID (UUID)
	body := *openapiclient.NewNodepoolsAttachServersRequest() // NodepoolsAttachServersRequest | Server IDs to attach

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.AttachNodePoolServers(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.AttachNodePoolServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AttachNodePoolServers`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.AttachNodePoolServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Group ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAttachNodePoolServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**NodepoolsAttachServersRequest**](NodepoolsAttachServersRequest.md) | Server IDs to attach | 

### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AttachNodePoolTaints

> string AttachNodePoolTaints(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Attach taints to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	body := *openapiclient.NewNodepoolsAttachTaintsRequest() // NodepoolsAttachTaintsRequest | Taint IDs to attach

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.AttachNodePoolTaints(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.AttachNodePoolTaints``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AttachNodePoolTaints`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.AttachNodePoolTaints`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAttachNodePoolTaintsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**NodepoolsAttachTaintsRequest**](NodepoolsAttachTaintsRequest.md) | Taint IDs to attach | 

### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateNodePool

> NodepoolsNodePoolResponse CreateNodePool(ctx).XOrgID(xOrgID).Body(body).Execute()

Create node group (org scoped)



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
	body := *openapiclient.NewNodepoolsCreateNodePoolRequest() // NodepoolsCreateNodePoolRequest | NodeGroup payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.CreateNodePool(context.Background()).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.CreateNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateNodePool`: NodepoolsNodePoolResponse
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.CreateNodePool`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **body** | [**NodepoolsCreateNodePoolRequest**](NodepoolsCreateNodePoolRequest.md) | NodeGroup payload | 

### Return type

[**NodepoolsNodePoolResponse**](NodepoolsNodePoolResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteNodePool

> string DeleteNodePool(ctx, id).XOrgID(xOrgID).Execute()

Delete node pool (org scoped)



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
	id := "id_example" // string | Node Group ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.DeleteNodePool(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.DeleteNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteNodePool`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.DeleteNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Group ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteNodePoolRequest struct via the builder pattern


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


## DetachNodePoolAnnotation

> string DetachNodePoolAnnotation(ctx, id, annotationId).XOrgID(xOrgID).Execute()

Detach one annotation from a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	annotationId := "annotationId_example" // string | Annotation ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.DetachNodePoolAnnotation(context.Background(), id, annotationId).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.DetachNodePoolAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DetachNodePoolAnnotation`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.DetachNodePoolAnnotation`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 
**annotationId** | **string** | Annotation ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDetachNodePoolAnnotationRequest struct via the builder pattern


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


## DetachNodePoolLabel

> string DetachNodePoolLabel(ctx, id, labelId).XOrgID(xOrgID).Execute()

Detach one label from a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	labelId := "labelId_example" // string | Label ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.DetachNodePoolLabel(context.Background(), id, labelId).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.DetachNodePoolLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DetachNodePoolLabel`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.DetachNodePoolLabel`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 
**labelId** | **string** | Label ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDetachNodePoolLabelRequest struct via the builder pattern


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


## DetachNodePoolServer

> string DetachNodePoolServer(ctx, id, serverId).XOrgID(xOrgID).Execute()

Detach one server from a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	serverId := "serverId_example" // string | Server ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.DetachNodePoolServer(context.Background(), id, serverId).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.DetachNodePoolServer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DetachNodePoolServer`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.DetachNodePoolServer`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 
**serverId** | **string** | Server ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDetachNodePoolServerRequest struct via the builder pattern


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


## DetachNodePoolTaint

> string DetachNodePoolTaint(ctx, id, taintId).XOrgID(xOrgID).Execute()

Detach one taint from a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	taintId := "taintId_example" // string | Taint ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.DetachNodePoolTaint(context.Background(), id, taintId).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.DetachNodePoolTaint``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DetachNodePoolTaint`: string
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.DetachNodePoolTaint`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 
**taintId** | **string** | Taint ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDetachNodePoolTaintRequest struct via the builder pattern


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


## GetNodePool

> NodepoolsNodePoolResponse GetNodePool(ctx, id).XOrgID(xOrgID).Include(include).Execute()

Get node group by ID (org scoped)



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
	id := "id_example" // string | Node Group ID (UUID)
	include := "include_example" // string | Optional: servers (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.GetNodePool(context.Background(), id).XOrgID(xOrgID).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.GetNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodePool`: NodepoolsNodePoolResponse
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.GetNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Group ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **include** | **string** | Optional: servers | 

### Return type

[**NodepoolsNodePoolResponse**](NodepoolsNodePoolResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolAnnotations

> []NodepoolsAnnotationBrief ListNodePoolAnnotations(ctx, id).XOrgID(xOrgID).Execute()

List annotations attached to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePoolAnnotations(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePoolAnnotations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolAnnotations`: []NodepoolsAnnotationBrief
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.ListNodePoolAnnotations`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolAnnotationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


### Return type

[**[]NodepoolsAnnotationBrief**](NodepoolsAnnotationBrief.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolLabels

> []NodepoolsLabelBrief ListNodePoolLabels(ctx, id).XOrgID(xOrgID).Execute()

List labels attached to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePoolLabels(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePoolLabels``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolLabels`: []NodepoolsLabelBrief
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.ListNodePoolLabels`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolLabelsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


### Return type

[**[]NodepoolsLabelBrief**](NodepoolsLabelBrief.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolServers

> []NodepoolsServerBrief ListNodePoolServers(ctx, id).XOrgID(xOrgID).Execute()

List servers attached to a node pool (org scoped)

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
	id := "id_example" // string | Node Group ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePoolServers(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePoolServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolServers`: []NodepoolsServerBrief
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.ListNodePoolServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Group ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


### Return type

[**[]NodepoolsServerBrief**](NodepoolsServerBrief.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolTaints

> []NodepoolsTaintBrief ListNodePoolTaints(ctx, id).XOrgID(xOrgID).Execute()

List taints attached to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePoolTaints(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePoolTaints``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolTaints`: []NodepoolsTaintBrief
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.ListNodePoolTaints`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolTaintsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


### Return type

[**[]NodepoolsTaintBrief**](NodepoolsTaintBrief.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePools

> []NodepoolsNodePoolResponse ListNodePools(ctx).XOrgID(xOrgID).Q(q).Include(include).Execute()

List node pools (org scoped)



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
	q := "q_example" // string | Name contains (case-insensitive) (optional)
	include := "include_example" // string | Optional: servers (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePools(context.Background()).XOrgID(xOrgID).Q(q).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePools`: []NodepoolsNodePoolResponse
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.ListNodePools`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **q** | **string** | Name contains (case-insensitive) | 
 **include** | **string** | Optional: servers | 

### Return type

[**[]NodepoolsNodePoolResponse**](NodepoolsNodePoolResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNodePool

> NodepoolsNodePoolResponse UpdateNodePool(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Update node pool (org scoped)



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
	id := "id_example" // string | Node Pool ID (UUID)
	body := *openapiclient.NewNodepoolsUpdateNodePoolRequest() // NodepoolsUpdateNodePoolRequest | Fields to update

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.UpdateNodePool(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.UpdateNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateNodePool`: NodepoolsNodePoolResponse
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.UpdateNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**NodepoolsUpdateNodePoolRequest**](NodepoolsUpdateNodePoolRequest.md) | Fields to update | 

### Return type

[**NodepoolsNodePoolResponse**](NodepoolsNodePoolResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

