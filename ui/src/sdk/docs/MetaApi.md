# MetaApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**versionOperationId**](MetaApi.md#versionoperationid) | **GET** /version | Service version information |



## versionOperationId

> HandlersVersionResponse versionOperationId()

Service version information

Returns build/runtime metadata for the running service.

### Example

```ts
import {
  Configuration,
  MetaApi,
} from '@glueops/autoglue-sdk-go';
import type { VersionOperationIdRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const api = new MetaApi();

  try {
    const data = await api.versionOperationId();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**HandlersVersionResponse**](HandlersVersionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

