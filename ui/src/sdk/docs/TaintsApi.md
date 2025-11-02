# TaintsApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createTaint**](TaintsApi.md#createtaint) | **POST** /taints | Create node taint (org scoped) |
| [**deleteTaint**](TaintsApi.md#deletetaint) | **DELETE** /taints/{id} | Delete taint (org scoped) |
| [**getTaint**](TaintsApi.md#gettaint) | **GET** /taints/{id} | Get node taint by ID (org scoped) |
| [**listTaints**](TaintsApi.md#listtaints) | **GET** /taints | List node pool taints (org scoped) |
| [**updateTaint**](TaintsApi.md#updatetaint) | **PATCH** /taints/{id} | Update node taint (org scoped) |



## createTaint

> DtoTaintResponse createTaint(body, xOrgID)

Create node taint (org scoped)

Creates a taint.

### Example

```ts
import {
  Configuration,
  TaintsApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateTaintRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: OrgKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: OrgSecretAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new TaintsApi(config);

  const body = {
    // DtoCreateTaintRequest | Taint payload
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateTaintRequest;

  try {
    const data = await api.createTaint(body);
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
| **body** | [DtoCreateTaintRequest](DtoCreateTaintRequest.md) | Taint payload | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoTaintResponse**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | invalid json / missing fields / invalid node_pool_ids |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **500** | create failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteTaint

> string deleteTaint(id, xOrgID)

Delete taint (org scoped)

Permanently deletes the taint.

### Example

```ts
import {
  Configuration,
  TaintsApi,
} from '@glueops/autoglue-sdk-go';
import type { DeleteTaintRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: OrgKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: OrgSecretAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new TaintsApi(config);

  const body = {
    // string | Node Taint ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteTaintRequest;

  try {
    const data = await api.deleteTaint(body);
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
| **id** | `string` | Node Taint ID (UUID) | [Defaults to `undefined`] |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **400** | invalid id |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **500** | delete failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getTaint

> DtoTaintResponse getTaint(id, xOrgID)

Get node taint by ID (org scoped)

### Example

```ts
import {
  Configuration,
  TaintsApi,
} from '@glueops/autoglue-sdk-go';
import type { GetTaintRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: OrgKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: OrgSecretAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new TaintsApi(config);

  const body = {
    // string | Node Taint ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies GetTaintRequest;

  try {
    const data = await api.getTaint(body);
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
| **id** | `string` | Node Taint ID (UUID) | [Defaults to `undefined`] |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoTaintResponse**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | invalid id |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **404** | not found |  -  |
| **500** | fetch failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listTaints

> Array&lt;DtoTaintResponse&gt; listTaints(xOrgID, key, value, q)

List node pool taints (org scoped)

Returns node taints for the organization in X-Org-ID. Filters: &#x60;key&#x60;, &#x60;value&#x60;, and &#x60;q&#x60; (key contains). Add &#x60;include&#x3D;node_pools&#x60; to include linked node pools.

### Example

```ts
import {
  Configuration,
  TaintsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListTaintsRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: OrgKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: OrgSecretAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new TaintsApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // string | Exact key (optional)
    key: key_example,
    // string | Exact value (optional)
    value: value_example,
    // string | key contains (case-insensitive) (optional)
    q: q_example,
  } satisfies ListTaintsRequest;

  try {
    const data = await api.listTaints(body);
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
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |
| **key** | `string` | Exact key | [Optional] [Defaults to `undefined`] |
| **value** | `string` | Exact value | [Optional] [Defaults to `undefined`] |
| **q** | `string` | key contains (case-insensitive) | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoTaintResponse&gt;**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **500** | failed to list node taints |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateTaint

> DtoTaintResponse updateTaint(id, body, xOrgID)

Update node taint (org scoped)

Partially update taint fields.

### Example

```ts
import {
  Configuration,
  TaintsApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateTaintRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: OrgKeyAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: OrgSecretAuth
    apiKey: "YOUR API KEY",
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new TaintsApi(config);

  const body = {
    // string | Node Taint ID (UUID)
    id: id_example,
    // DtoUpdateTaintRequest | Fields to update
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateTaintRequest;

  try {
    const data = await api.updateTaint(body);
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
| **id** | `string` | Node Taint ID (UUID) | [Defaults to `undefined`] |
| **body** | [DtoUpdateTaintRequest](DtoUpdateTaintRequest.md) | Fields to update | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoTaintResponse**](DtoTaintResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | invalid id / invalid json |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **404** | not found |  -  |
| **500** | update failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

