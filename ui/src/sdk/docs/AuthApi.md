# AuthApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**authCallback**](AuthApi.md#authcallback) | **GET** /auth/{provider}/callback | Handle social login callback |
| [**authStart**](AuthApi.md#authstart) | **POST** /auth/{provider}/start | Begin social login |
| [**getJWKS**](AuthApi.md#getjwks) | **GET** /.well-known/jwks.json | Get JWKS |
| [**logout**](AuthApi.md#logout) | **POST** /auth/logout | Revoke refresh token family (logout everywhere) |
| [**refresh**](AuthApi.md#refresh) | **POST** /auth/refresh | Rotate refresh token |



## authCallback

> DtoTokenPair authCallback(provider)

Handle social login callback

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/autoglue-sdk-go';
import type { AuthCallbackRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const api = new AuthApi();

  const body = {
    // string | google|github
    provider: provider_example,
  } satisfies AuthCallbackRequest;

  try {
    const data = await api.authCallback(body);
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
| **provider** | `string` | google|github | [Defaults to `undefined`] |

### Return type

[**DtoTokenPair**](DtoTokenPair.md)

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


## authStart

> DtoAuthStartResponse authStart(provider)

Begin social login

Returns provider authorization URL for the frontend to redirect

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/autoglue-sdk-go';
import type { AuthStartRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const api = new AuthApi();

  const body = {
    // string | google|github
    provider: provider_example,
  } satisfies AuthStartRequest;

  try {
    const data = await api.authStart(body);
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
| **provider** | `string` | google|github | [Defaults to `undefined`] |

### Return type

[**DtoAuthStartResponse**](DtoAuthStartResponse.md)

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


## getJWKS

> DtoJWKS getJWKS()

Get JWKS

Returns the JSON Web Key Set for token verification

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/autoglue-sdk-go';
import type { GetJWKSRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const api = new AuthApi();

  try {
    const data = await api.getJWKS();
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

[**DtoJWKS**](DtoJWKS.md)

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


## logout

> logout(body)

Revoke refresh token family (logout everywhere)

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/autoglue-sdk-go';
import type { LogoutRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const api = new AuthApi();

  const body = {
    // DtoLogoutRequest | Refresh token
    body: ...,
  } satisfies LogoutRequest;

  try {
    const data = await api.logout(body);
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
| **body** | [DtoLogoutRequest](DtoLogoutRequest.md) | Refresh token | |

### Return type

`void` (Empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: Not defined


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## refresh

> DtoTokenPair refresh(body)

Rotate refresh token

### Example

```ts
import {
  Configuration,
  AuthApi,
} from '@glueops/autoglue-sdk-go';
import type { RefreshRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const api = new AuthApi();

  const body = {
    // DtoRefreshRequest | Refresh token
    body: ...,
  } satisfies RefreshRequest;

  try {
    const data = await api.refresh(body);
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
| **body** | [DtoRefreshRequest](DtoRefreshRequest.md) | Refresh token | |

### Return type

[**DtoTokenPair**](DtoTokenPair.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

