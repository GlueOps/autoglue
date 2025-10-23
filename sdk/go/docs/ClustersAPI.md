# \ClustersAPI

All URIs are relative to *https://autoglue.apps.nonprod.earth.onglueops.rocks*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AttachNodePools**](ClustersAPI.md#AttachNodePools) | **Post** /api/v1/clusters/{id}/node_pools | Attach node pools to cluster (org scoped)
[**CreateCluster**](ClustersAPI.md#CreateCluster) | **Post** /api/v1/clusters | Create cluster (org scoped)
[**DeleteBastion**](ClustersAPI.md#DeleteBastion) | **Delete** /api/v1/clusters/{id}/bastion | Clear cluster bastion (org scoped)
[**DeleteCluster**](ClustersAPI.md#DeleteCluster) | **Delete** /api/v1/clusters/{id} | Delete cluster (org scoped)
[**DetachNodePool**](ClustersAPI.md#DetachNodePool) | **Delete** /api/v1/clusters/{id}/node_pools/{poolId} | Detach one node pool from a cluster (org scoped)
[**GetBastion**](ClustersAPI.md#GetBastion) | **Get** /api/v1/clusters/{id}/bastion | Get cluster bastion (org scoped)
[**GetCluster**](ClustersAPI.md#GetCluster) | **Get** /api/v1/clusters/{id} | Get cluster by ID (org scoped)
[**ListClusterNodePools**](ClustersAPI.md#ListClusterNodePools) | **Get** /api/v1/clusters/{id}/node_pools | List node pools attached to a cluster (org scoped)
[**ListClusters**](ClustersAPI.md#ListClusters) | **Get** /api/v1/clusters | List clusters (org scoped)
[**PutBastion**](ClustersAPI.md#PutBastion) | **Post** /api/v1/clusters/{id}/bastion | Set/replace cluster bastion (org scoped)
[**UpdateCluster**](ClustersAPI.md#UpdateCluster) | **Patch** /api/v1/clusters/{id} | Update cluster (org scoped). If &#x60;kubeconfig&#x60; is provided and non-empty, it will be encrypted per-organization and stored (never returned). Sending an empty string for &#x60;kubeconfig&#x60; is ignored (no change).



## AttachNodePools

> string AttachNodePools(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Attach node pools to cluster (org scoped)

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
	id := "id_example" // string | Cluster ID (UUID)
	body := *openapiclient.NewClustersAttachNodePoolsRequest() // ClustersAttachNodePoolsRequest | node_pool_ids

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.AttachNodePools(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.AttachNodePools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AttachNodePools`: string
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.AttachNodePools`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiAttachNodePoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**ClustersAttachNodePoolsRequest**](ClustersAttachNodePoolsRequest.md) | node_pool_ids | 

### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateCluster

> ClustersClusterResponse CreateCluster(ctx).XOrgID(xOrgID).Body(body).Execute()

Create cluster (org scoped)



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
	body := *openapiclient.NewClustersCreateClusterRequest() // ClustersCreateClusterRequest | payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.CreateCluster(context.Background()).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.CreateCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateCluster`: ClustersClusterResponse
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.CreateCluster`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **body** | [**ClustersCreateClusterRequest**](ClustersCreateClusterRequest.md) | payload | 

### Return type

[**ClustersClusterResponse**](ClustersClusterResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteBastion

> string DeleteBastion(ctx, id).XOrgID(xOrgID).Execute()

Clear cluster bastion (org scoped)

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
	id := "id_example" // string | Cluster ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteBastion(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteBastion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteBastion`: string
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteBastion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteBastionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteCluster

> string DeleteCluster(ctx, id).XOrgID(xOrgID).Execute()

Delete cluster (org scoped)

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
	id := "id_example" // string | Cluster ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DeleteCluster(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DeleteCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteCluster`: string
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DeleteCluster`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DetachNodePool

> string DetachNodePool(ctx, id, poolId).XOrgID(xOrgID).Execute()

Detach one node pool from a cluster (org scoped)

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
	id := "id_example" // string | Cluster ID (UUID)
	poolId := "poolId_example" // string | Node Pool ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.DetachNodePool(context.Background(), id, poolId).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.DetachNodePool``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DetachNodePool`: string
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.DetachNodePool`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 
**poolId** | **string** | Node Pool ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiDetachNodePoolRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 



### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetBastion

> ClustersServerBrief GetBastion(ctx, id).XOrgID(xOrgID).Execute()

Get cluster bastion (org scoped)

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
	id := "id_example" // string | Cluster ID (UUID)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetBastion(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetBastion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetBastion`: ClustersServerBrief
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetBastion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetBastionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 


### Return type

[**ClustersServerBrief**](ClustersServerBrief.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCluster

> ClustersClusterResponse GetCluster(ctx, id).XOrgID(xOrgID).Include(include).Execute()

Get cluster by ID (org scoped)



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
	id := "id_example" // string | Cluster ID (UUID)
	include := "include_example" // string | Optional: node_pools,bastion (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.GetCluster(context.Background(), id).XOrgID(xOrgID).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.GetCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetCluster`: ClustersClusterResponse
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.GetCluster`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **include** | **string** | Optional: node_pools,bastion | 

### Return type

[**ClustersClusterResponse**](ClustersClusterResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusterNodePools

> []ClustersNodePoolBrief ListClusterNodePools(ctx, id).XOrgID(xOrgID).Q(q).Execute()

List node pools attached to a cluster (org scoped)

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
	id := "id_example" // string | Cluster ID (UUID)
	q := "q_example" // string | Name contains (case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusterNodePools(context.Background(), id).XOrgID(xOrgID).Q(q).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusterNodePools``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusterNodePools`: []ClustersNodePoolBrief
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusterNodePools`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiListClusterNodePoolsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **q** | **string** | Name contains (case-insensitive) | 

### Return type

[**[]ClustersNodePoolBrief**](ClustersNodePoolBrief.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListClusters

> []ClustersClusterResponse ListClusters(ctx).XOrgID(xOrgID).Q(q).Include(include).Execute()

List clusters (org scoped)



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
	include := "include_example" // string | Optional: node_pools,bastion (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.ListClusters(context.Background()).XOrgID(xOrgID).Q(q).Include(include).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.ListClusters``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListClusters`: []ClustersClusterResponse
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.ListClusters`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListClustersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 
 **q** | **string** | Name contains (case-insensitive) | 
 **include** | **string** | Optional: node_pools,bastion | 

### Return type

[**[]ClustersClusterResponse**](ClustersClusterResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PutBastion

> string PutBastion(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Set/replace cluster bastion (org scoped)

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
	id := "id_example" // string | Cluster ID (UUID)
	body := *openapiclient.NewClustersSetBastionRequest() // ClustersSetBastionRequest | server_id with role=bastion

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.PutBastion(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.PutBastion``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PutBastion`: string
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.PutBastion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiPutBastionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**ClustersSetBastionRequest**](ClustersSetBastionRequest.md) | server_id with role&#x3D;bastion | 

### Return type

**string**

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateCluster

> ClustersClusterResponse UpdateCluster(ctx, id).XOrgID(xOrgID).Body(body).Execute()

Update cluster (org scoped). If `kubeconfig` is provided and non-empty, it will be encrypted per-organization and stored (never returned). Sending an empty string for `kubeconfig` is ignored (no change).

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
	id := "id_example" // string | Cluster ID (UUID)
	body := *openapiclient.NewClustersUpdateClusterRequest() // ClustersUpdateClusterRequest | payload

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ClustersAPI.UpdateCluster(context.Background(), id).XOrgID(xOrgID).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ClustersAPI.UpdateCluster``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateCluster`: ClustersClusterResponse
	fmt.Fprintf(os.Stdout, "Response from `ClustersAPI.UpdateCluster`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Cluster ID (UUID) | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateClusterRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xOrgID** | **string** | Organization UUID | 

 **body** | [**ClustersUpdateClusterRequest**](ClustersUpdateClusterRequest.md) | payload | 

### Return type

[**ClustersClusterResponse**](ClustersClusterResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

