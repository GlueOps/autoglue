# MeApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getMe**](MeApi.md#getme) | **GET** /me | Get current user profile |
| [**updateMe**](MeApi.md#updateme) | **PATCH** /me | Update current user profile |



## getMe

> HandlersMeResponse getMe()

Get current user profile

### Example

```ts
import {
  Configuration,
  MeApi,
} from '@glueops/autoglue-sdk';
import type { GetMeRequest } from '@glueops/autoglue-sdk';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: ApiKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new MeApi(config);

  try {
    const data = await api.getMe();
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

[**HandlersMeResponse**](HandlersMeResponse.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateMe

> ModelsUser updateMe(body)

Update current user profile

### Example

```ts
import {
  Configuration,
  MeApi,
} from '@glueops/autoglue-sdk';
import type { UpdateMeRequest } from '@glueops/autoglue-sdk';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: ApiKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new MeApi(config);

  const body = {
    // HandlersUpdateMeRequest | Patch profile
    body: ...,
  } satisfies UpdateMeRequest;

  try {
    const data = await api.updateMe(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **body** | [HandlersUpdateMeRequest](HandlersUpdateMeRequest.md) | Patch profile | |

### Return type

[**ModelsUser**](ModelsUser.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

