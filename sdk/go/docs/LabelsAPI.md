# \LabelsAPI

All URIs are relative to *http://localhost:8080/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateLabel**](LabelsAPI.md#CreateLabel) | **Post** /labels | Create label (org scoped)
[**DeleteLabel**](LabelsAPI.md#DeleteLabel) | **Delete** /labels/{id} | Delete label (org scoped)
[**GetLabel**](LabelsAPI.md#GetLabel) | **Get** /labels/{id} | Get label by ID (org scoped)
[**ListLabels**](LabelsAPI.md#ListLabels) | **Get** /labels | List node labels (org scoped)
[**UpdateLabel**](LabelsAPI.md#UpdateLabel) | **Patch** /labels/{id} | Update label (org scoped)



## CreateLabel

> DtoLabelResponse CreateLabel(ctx).Body(body).XOrgID(xOrgID).Execute()

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
	body := *openapiclient.NewDtoCreateLabelRequest() // DtoCreateLabelRequest | Label payload
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.CreateLabel(context.Background()).Body(body).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.CreateLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateLabel`: DtoLabelResponse
	fmt.Fprintf(os.Stdout, "Response from `LabelsAPI.CreateLabel`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateLabelRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**DtoCreateLabelRequest**](DtoCreateLabelRequest.md) | Label payload | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoLabelResponse**](DtoLabelResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

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
	id := "id_example" // string | Label ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

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

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetLabel

> DtoLabelResponse GetLabel(ctx, id).XOrgID(xOrgID).Execute()

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
	id := "id_example" // string | Label ID (UUID)
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.GetLabel(context.Background(), id).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.GetLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GetLabel`: DtoLabelResponse
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

### Return type

[**DtoLabelResponse**](DtoLabelResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListLabels

> []DtoLabelResponse ListLabels(ctx).XOrgID(xOrgID).Key(key).Value(value).Q(q).Execute()

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
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)
	key := "key_example" // string | Exact key (optional)
	value := "value_example" // string | Exact value (optional)
	q := "q_example" // string | Key contains (case-insensitive) (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.ListLabels(context.Background()).XOrgID(xOrgID).Key(key).Value(value).Q(q).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.ListLabels``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ListLabels`: []DtoLabelResponse
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


## UpdateLabel

> DtoLabelResponse UpdateLabel(ctx, id).Body(body).XOrgID(xOrgID).Execute()

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
	id := "id_example" // string | Label ID (UUID)
	body := *openapiclient.NewDtoUpdateLabelRequest() // DtoUpdateLabelRequest | Fields to update
	xOrgID := "xOrgID_example" // string | Organization UUID (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.LabelsAPI.UpdateLabel(context.Background(), id).Body(body).XOrgID(xOrgID).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `LabelsAPI.UpdateLabel``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateLabel`: DtoLabelResponse
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

 **body** | [**DtoUpdateLabelRequest**](DtoUpdateLabelRequest.md) | Fields to update | 
 **xOrgID** | **string** | Organization UUID | 

### Return type

[**DtoLabelResponse**](DtoLabelResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

