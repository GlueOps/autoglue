# MeAPIKeysApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUserAPIKey**](MeAPIKeysApi.md#createuserapikey) | **POST** /me/api-keys | Create a new user API key |
| [**deleteUserAPIKey**](MeAPIKeysApi.md#deleteuserapikey) | **DELETE** /me/api-keys/{id} | Delete a user API key |
| [**listUserAPIKeys**](MeAPIKeysApi.md#listuserapikeys) | **GET** /me/api-keys | List my API keys |



## createUserAPIKey

> HandlersUserAPIKeyOut createUserAPIKey(body)

Create a new user API key

Returns the plaintext key once. Store it securely on the client side.

### Example

```ts
import {
  Configuration,
  MeAPIKeysApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateUserAPIKeyRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: ApiKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new MeAPIKeysApi(config);

  const body = {
    // HandlersCreateUserKeyRequest | Key options
    body: ...,
  } satisfies CreateUserAPIKeyRequest;

  try {
    const data = await api.createUserAPIKey(body);
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
| **body** | [HandlersCreateUserKeyRequest](HandlersCreateUserKeyRequest.md) | Key options | |

### Return type

[**HandlersUserAPIKeyOut**](HandlersUserAPIKeyOut.md)

### Authorization

[ApiKeyAuth](../README.md#ApiKeyAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteUserAPIKey

> deleteUserAPIKey(id)

Delete a user API key

### Example

```ts
import {
  Configuration,
  MeAPIKeysApi,
} from '@glueops/autoglue-sdk-go';
import type { DeleteUserAPIKeyRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new MeAPIKeysApi(config);

  const body = {
    // string | Key ID (UUID)
    id: id_example,
  } satisfies DeleteUserAPIKeyRequest;

  try {
    const data = await api.deleteUserAPIKey(body);
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
| **id** | `string` | Key ID (UUID) | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listUserAPIKeys

> Array&lt;HandlersUserAPIKeyOut&gt; listUserAPIKeys()

List my API keys

### Example

```ts
import {
  Configuration,
  MeAPIKeysApi,
} from '@glueops/autoglue-sdk-go';
import type { ListUserAPIKeysRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: ApiKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new MeAPIKeysApi(config);

  try {
    const data = await api.listUserAPIKeys();
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

[**Array&lt;HandlersUserAPIKeyOut&gt;**](HandlersUserAPIKeyOut.md)

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

