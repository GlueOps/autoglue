# \NodePoolsAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AttachNodePoolAnnotations**](NodePoolsAPI.md#AttachNodePoolAnnotations) | **Post** /node-pools/{id}/annotations | Attach annotation to a node pool (org scoped)
[**AttachNodePoolLabels**](NodePoolsAPI.md#AttachNodePoolLabels) | **Post** /node-pools/{id}/labels | Attach labels to a node pool (org scoped)
[**AttachNodePoolServers**](NodePoolsAPI.md#AttachNodePoolServers) | **Post** /node-pools/{id}/servers | Attach servers to a node pool (org scoped)
[**AttachNodePoolTaints**](NodePoolsAPI.md#AttachNodePoolTaints) | **Post** /node-pools/{id}/taints | Attach taints to a node pool (org scoped)
[**CreateNodePool**](NodePoolsAPI.md#CreateNodePool) | **Post** /node-pools | Create node pool (org scoped)
[**DeleteNodePool**](NodePoolsAPI.md#DeleteNodePool) | **Delete** /node-pools/{id} | Delete node pool (org scoped)
[**DetachNodePoolAnnotation**](NodePoolsAPI.md#DetachNodePoolAnnotation) | **Delete** /node-pools/{id}/annotations/{annotationId} | Detach one annotation from a node pool (org scoped)
[**DetachNodePoolLabel**](NodePoolsAPI.md#DetachNodePoolLabel) | **Delete** /node-pools/{id}/labels/{labelId} | Detach one label from a node pool (org scoped)
[**DetachNodePoolServer**](NodePoolsAPI.md#DetachNodePoolServer) | **Delete** /node-pools/{id}/servers/{serverId} | Detach one server from a node pool (org scoped)
[**DetachNodePoolTaint**](NodePoolsAPI.md#DetachNodePoolTaint) | **Delete** /node-pools/{id}/taints/{taintId} | Detach one taint from a node pool (org scoped)
[**GetNodePool**](NodePoolsAPI.md#GetNodePool) | **Get** /node-pools/{id} | Get node pool by ID (org scoped)
[**ListNodePoolAnnotations**](NodePoolsAPI.md#ListNodePoolAnnotations) | **Get** /node-pools/{id}/annotations | List annotations attached to a node pool (org scoped)
[**ListNodePoolLabels**](NodePoolsAPI.md#ListNodePoolLabels) | **Get** /node-pools/{id}/labels | List labels attached to a node pool (org scoped)
[**ListNodePoolServers**](NodePoolsAPI.md#ListNodePoolServers) | **Get** /node-pools/{id}/servers | List servers attached to a node pool (org scoped)
[**ListNodePoolTaints**](NodePoolsAPI.md#ListNodePoolTaints) | **Get** /node-pools/{id}/taints | List taints attached to a node pool (org scoped)
[**ListNodePools**](NodePoolsAPI.md#ListNodePools) | **Get** /node-pools | List node pools (org scoped)
[**UpdateNodePool**](NodePoolsAPI.md#UpdateNodePool) | **Patch** /node-pools/{id} | Update node pool (org scoped)



## AttachNodePoolAnnotations

> string AttachNodePoolAnnotations(ctx, id).Body(body).XOrgID(xOrgID).Execute()

Attach annotation to a node pool (org scoped)

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
	id := "id_example" // string | Node Group ID (UUID)
	body := *openapiclient.NewDtoAttachAnnotationsRequest() // DtoAttachAnnotationsRequest | Annotation IDs to attach
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.AttachNodePoolAnnotations(context.Background(), id).Body(body).XOrgID(xOrgID).Execute()
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
**id** | **string** | Node Group ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAttachNodePoolAnnotationsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**DtoAttachAnnotationsRequest**](DtoAttachAnnotationsRequest.md) | Annotation IDs to attach | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AttachNodePoolLabels

> string AttachNodePoolLabels(ctx, id).Body(body).XOrgID(xOrgID).Execute()

Attach labels to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	body := *openapiclient.NewDtoAttachLabelsRequest() // DtoAttachLabelsRequest | Label IDs to attach
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.AttachNodePoolLabels(context.Background(), id).Body(body).XOrgID(xOrgID).Execute()
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

 **body** | [**DtoAttachLabelsRequest**](DtoAttachLabelsRequest.md) | Label IDs to attach | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AttachNodePoolServers

> string AttachNodePoolServers(ctx, id).Body(body).XOrgID(xOrgID).Execute()

Attach servers to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	body := *openapiclient.NewDtoAttachServersRequest() // DtoAttachServersRequest | Server IDs to attach
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.AttachNodePoolServers(context.Background(), id).Body(body).XOrgID(xOrgID).Execute()
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
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAttachNodePoolServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**DtoAttachServersRequest**](DtoAttachServersRequest.md) | Server IDs to attach | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AttachNodePoolTaints

> string AttachNodePoolTaints(ctx, id).Body(body).XOrgID(xOrgID).Execute()

Attach taints to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	body := *openapiclient.NewDtoAttachTaintsRequest() // DtoAttachTaintsRequest | Taint IDs to attach
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.AttachNodePoolTaints(context.Background(), id).Body(body).XOrgID(xOrgID).Execute()
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

 **body** | [**DtoAttachTaintsRequest**](DtoAttachTaintsRequest.md) | Taint IDs to attach | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateNodePool

> DtoNodePoolResponse CreateNodePool(ctx).Body(body).XOrgID(xOrgID).Execute()

Create node pool (org scoped)



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
	body := *openapiclient.NewDtoCreateNodePoolRequest() // DtoCreateNodePoolRequest | NodePool payload
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.CreateNodePool(context.Background()).Body(body).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.CreateNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateNodePool`: DtoNodePoolResponse
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.CreateNodePool`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**DtoCreateNodePoolRequest**](DtoCreateNodePoolRequest.md) | NodePool payload | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoNodePoolResponse**](DtoNodePoolResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

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
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {
	id := "id_example" // string | Node Pool ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

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
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteNodePoolRequest struct via the builder pattern


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
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {
	id := "id_example" // string | Node Pool ID (UUID)
	annotationId := "annotationId_example" // string | Annotation ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

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

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

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
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {
	id := "id_example" // string | Node Pool ID (UUID)
	labelId := "labelId_example" // string | Label ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

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

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

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
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {
	id := "id_example" // string | Node Pool ID (UUID)
	serverId := "serverId_example" // string | Server ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

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

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

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
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {
	id := "id_example" // string | Node Pool ID (UUID)
	taintId := "taintId_example" // string | Taint ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

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

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetNodePool

> DtoNodePoolResponse GetNodePool(ctx, id).XOrgID(xOrgID).Execute()

Get node pool by ID (org scoped)



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
	id := "id_example" // string | Node Pool ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.GetNodePool(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.GetNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetNodePool`: DtoNodePoolResponse
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.GetNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoNodePoolResponse**](DtoNodePoolResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolAnnotations

> []DtoAnnotationResponse ListNodePoolAnnotations(ctx, id).XOrgID(xOrgID).Execute()

List annotations attached to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePoolAnnotations(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePoolAnnotations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolAnnotations`: []DtoAnnotationResponse
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

[**[]DtoAnnotationResponse**](DtoAnnotationResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolLabels

> []DtoLabelResponse ListNodePoolLabels(ctx, id).XOrgID(xOrgID).Execute()

List labels attached to a node pool (org scoped)

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
	id := "id_example" // string | Label Pool ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePoolLabels(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePoolLabels``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolLabels`: []DtoLabelResponse
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.ListNodePoolLabels`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Label Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolLabelsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xOrgID** | **string** | Organization UUID | 

### Return type

[**[]DtoLabelResponse**](DtoLabelResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolServers

> []DtoServerResponse ListNodePoolServers(ctx, id).XOrgID(xOrgID).Execute()

List servers attached to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePoolServers(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePoolServers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolServers`: []DtoServerResponse
	fmt.Fprintf(os.Stdout, "Response from `NodePoolsAPI.ListNodePoolServers`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListNodePoolServersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xOrgID** | **string** | Organization UUID | 

### Return type

[**[]DtoServerResponse**](DtoServerResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePoolTaints

> []DtoTaintResponse ListNodePoolTaints(ctx, id).XOrgID(xOrgID).Execute()

List taints attached to a node pool (org scoped)

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
	id := "id_example" // string | Node Pool ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePoolTaints(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePoolTaints``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePoolTaints`: []DtoTaintResponse
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

[**[]DtoTaintResponse**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListNodePools

> []DtoNodePoolResponse ListNodePools(ctx).XOrgID(xOrgID).Q(q).Execute()

List node pools (org scoped)



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
	q := "q_example" // string | Name contains (case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.ListNodePools(context.Background()).XOrgID(xOrgID).Q(q).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.ListNodePools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListNodePools`: []DtoNodePoolResponse
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

### Return type

[**[]DtoNodePoolResponse**](DtoNodePoolResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateNodePool

> DtoNodePoolResponse UpdateNodePool(ctx, id).Body(body).XOrgID(xOrgID).Execute()

Update node pool (org scoped)



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
	id := "id_example" // string | Node Pool ID (UUID)
	body := *openapiclient.NewDtoUpdateNodePoolRequest() // DtoUpdateNodePoolRequest | Fields to update
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.NodePoolsAPI.UpdateNodePool(context.Background(), id).Body(body).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `NodePoolsAPI.UpdateNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateNodePool`: DtoNodePoolResponse
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

 **body** | [**DtoUpdateNodePoolRequest**](DtoUpdateNodePoolRequest.md) | Fields to update | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoNodePoolResponse**](DtoNodePoolResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

