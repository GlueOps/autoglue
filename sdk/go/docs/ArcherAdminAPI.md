# \ArcherAdminAPI

All URIs are relative to */api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AdminCancelArcherJob**](ArcherAdminAPI.md#AdminCancelArcherJob) | **Post** /admin/archer/jobs/{id}/cancel | Cancel an Archer job (admin)
[**AdminEnqueueArcherJob**](ArcherAdminAPI.md#AdminEnqueueArcherJob) | **Post** /admin/archer/jobs | Enqueue a new Archer job (admin)
[**AdminListArcherJobs**](ArcherAdminAPI.md#AdminListArcherJobs) | **Get** /admin/archer/jobs | List Archer jobs (admin)
[**AdminListArcherQueues**](ArcherAdminAPI.md#AdminListArcherQueues) | **Get** /admin/archer/queues | List Archer queues (admin)
[**AdminRetryArcherJob**](ArcherAdminAPI.md#AdminRetryArcherJob) | **Post** /admin/archer/jobs/{id}/retry | Retry a failed/canceled Archer job (admin)



## AdminCancelArcherJob

> DtoJob AdminCancelArcherJob(ctx, id).Execute()

Cancel an Archer job (admin)



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
	id := "id_example" // string | Job ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArcherAdminAPI.AdminCancelArcherJob(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArcherAdminAPI.AdminCancelArcherJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminCancelArcherJob`: DtoJob
	fmt.Fprintf(os.Stdout, "Response from `ArcherAdminAPI.AdminCancelArcherJob`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Job ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminCancelArcherJobRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DtoJob**](DtoJob.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminEnqueueArcherJob

> DtoJob AdminEnqueueArcherJob(ctx).Body(body).Execute()

Enqueue a new Archer job (admin)



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
	body := map[string]interface{}{ ... } // map[string]interface{} | Job parameters

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArcherAdminAPI.AdminEnqueueArcherJob(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArcherAdminAPI.AdminEnqueueArcherJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminEnqueueArcherJob`: DtoJob
	fmt.Fprintf(os.Stdout, "Response from `ArcherAdminAPI.AdminEnqueueArcherJob`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAdminEnqueueArcherJobRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | **map[string]interface{}** | Job parameters | 

### Return type

[**DtoJob**](DtoJob.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminListArcherJobs

> DtoPageJob AdminListArcherJobs(ctx).Status(status).Queue(queue).Q(q).Page(page).PageSize(pageSize).Execute()

List Archer jobs (admin)



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
	status := "status_example" // string | Filter by status (optional)
	queue := "queue_example" // string | Filter by queue name / worker name (optional)
	q := "q_example" // string | Free-text search (optional)
	page := int32(56) // int32 | Page number (optional) (default to 1)
	pageSize := int32(56) // int32 | Items per page (optional) (default to 25)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArcherAdminAPI.AdminListArcherJobs(context.Background()).Status(status).Queue(queue).Q(q).Page(page).PageSize(pageSize).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArcherAdminAPI.AdminListArcherJobs``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminListArcherJobs`: DtoPageJob
	fmt.Fprintf(os.Stdout, "Response from `ArcherAdminAPI.AdminListArcherJobs`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAdminListArcherJobsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **status** | **string** | Filter by status | 
 **queue** | **string** | Filter by queue name / worker name | 
 **q** | **string** | Free-text search | 
 **page** | **int32** | Page number | [default to 1]
 **pageSize** | **int32** | Items per page | [default to 25]

### Return type

[**DtoPageJob**](DtoPageJob.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminListArcherQueues

> []DtoQueueInfo AdminListArcherQueues(ctx).Execute()

List Archer queues (admin)



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
	resp, r, err := apiClient.ArcherAdminAPI.AdminListArcherQueues(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArcherAdminAPI.AdminListArcherQueues``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminListArcherQueues`: []DtoQueueInfo
	fmt.Fprintf(os.Stdout, "Response from `ArcherAdminAPI.AdminListArcherQueues`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiAdminListArcherQueuesRequest struct via the builder pattern


### Return type

[**[]DtoQueueInfo**](DtoQueueInfo.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AdminRetryArcherJob

> DtoJob AdminRetryArcherJob(ctx, id).Execute()

Retry a failed/canceled Archer job (admin)



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
	id := "id_example" // string | Job ID

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.ArcherAdminAPI.AdminRetryArcherJob(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `ArcherAdminAPI.AdminRetryArcherJob``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AdminRetryArcherJob`: DtoJob
	fmt.Fprintf(os.Stdout, "Response from `ArcherAdminAPI.AdminRetryArcherJob`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Job ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiAdminRetryArcherJobRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**DtoJob**](DtoJob.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

