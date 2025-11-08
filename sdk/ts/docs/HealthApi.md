# HealthApi

All URIs are relative to _/api/v1_

| Method                                                            | HTTP request     | Description        |
| ----------------------------------------------------------------- | ---------------- | ------------------ |
| [**healthCheckOperationId**](HealthApi.md#healthcheckoperationid) | **GET** /healthz | Basic health check |

## healthCheckOperationId

> HandlersHealthStatus healthCheckOperationId()

Basic health check

Returns 200 OK when the service is up

### Example

```ts
import { Configuration, HealthApi } from "@glueops/autoglue-sdk-go";
import type { HealthCheckOperationIdRequest } from "@glueops/autoglue-sdk-go";

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const api = new HealthApi();

  try {
    const data = await api.healthCheckOperationId();
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

[**HandlersHealthStatus**](HandlersHealthStatus.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **200**     | OK          | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
