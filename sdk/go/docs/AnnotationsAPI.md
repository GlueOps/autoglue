# \AnnotationsAPI

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateAnnotation**](AnnotationsAPI.md#CreateAnnotation) | **Post** /annotations | Create annotation (org scoped)
[**DeleteAnnotation**](AnnotationsAPI.md#DeleteAnnotation) | **Delete** /annotations/{id} | Delete annotation (org scoped)
[**GetAnnotation**](AnnotationsAPI.md#GetAnnotation) | **Get** /annotations/{id} | Get annotation by ID (org scoped)
[**ListAnnotations**](AnnotationsAPI.md#ListAnnotations) | **Get** /annotations | List annotations (org scoped)
[**UpdateAnnotation**](AnnotationsAPI.md#UpdateAnnotation) | **Patch** /annotations/{id} | Update annotation (org scoped)



## CreateAnnotation

> DtoAnnotationResponse CreateAnnotation(ctx).Body(body).XOrgID(xOrgID).Execute()

Create annotation (org scoped)



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
	body := *openapiclient.NewDtoCreateAnnotationRequest() // DtoCreateAnnotationRequest | Annotation payload
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.CreateAnnotation(context.Background()).Body(body).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.CreateAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateAnnotation`: DtoAnnotationResponse
	fmt.Fprintf(os.Stdout, "Response from `AnnotationsAPI.CreateAnnotation`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateAnnotationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**DtoCreateAnnotationRequest**](DtoCreateAnnotationRequest.md) | Annotation payload | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoAnnotationResponse**](DtoAnnotationResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

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
	openapiclient "github.com/glueops/autoglue-sdk-go"
)

func main() {
	id := "id_example" // string | Annotation ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

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

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAnnotation

> DtoAnnotationResponse GetAnnotation(ctx, id).XOrgID(xOrgID).Execute()

Get annotation by ID (org scoped)



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
	id := "id_example" // string | Annotation ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.GetAnnotation(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.GetAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetAnnotation`: DtoAnnotationResponse
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

### Return type

[**DtoAnnotationResponse**](DtoAnnotationResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListAnnotations

> []DtoAnnotationResponse ListAnnotations(ctx).XOrgID(xOrgID).Name(name).Value(value).Q(q).Execute()

List annotations (org scoped)



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
	name := "name_example" // string | Exact name (optional)
	value := "value_example" // string | Exact value (optional)
	q := "q_example" // string | name contains (case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.ListAnnotations(context.Background()).XOrgID(xOrgID).Name(name).Value(value).Q(q).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.ListAnnotations``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListAnnotations`: []DtoAnnotationResponse
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


## UpdateAnnotation

> DtoAnnotationResponse UpdateAnnotation(ctx, id).Body(body).XOrgID(xOrgID).Execute()

Update annotation (org scoped)



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
	id := "id_example" // string | Annotation ID (UUID)
	body := *openapiclient.NewDtoUpdateAnnotationRequest() // DtoUpdateAnnotationRequest | Fields to update
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AnnotationsAPI.UpdateAnnotation(context.Background(), id).Body(body).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AnnotationsAPI.UpdateAnnotation``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateAnnotation`: DtoAnnotationResponse
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

 **body** | [**DtoUpdateAnnotationRequest**](DtoUpdateAnnotationRequest.md) | Fields to update | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoAnnotationResponse**](DtoAnnotationResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

