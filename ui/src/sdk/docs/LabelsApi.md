# LabelsApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createLabel**](LabelsApi.md#createlabel) | **POST** /labels | Create label (org scoped) |
| [**deleteLabel**](LabelsApi.md#deletelabel) | **DELETE** /labels/{id} | Delete label (org scoped) |
| [**getLabel**](LabelsApi.md#getlabel) | **GET** /labels/{id} | Get label by ID (org scoped) |
| [**listLabels**](LabelsApi.md#listlabels) | **GET** /labels | List node labels (org scoped) |
| [**updateLabel**](LabelsApi.md#updatelabel) | **PATCH** /labels/{id} | Update label (org scoped) |



## createLabel

> DtoLabelResponse createLabel(body, xOrgID)

Create label (org scoped)

Creates a label.

### Example

```ts
import {
  Configuration,
  LabelsApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateLabelRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new LabelsApi(config);

  const body = {
    // DtoCreateLabelRequest | Label payload
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateLabelRequest;

  try {
    const data = await api.createLabel(body);
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
| **body** | [DtoCreateLabelRequest](DtoCreateLabelRequest.md) | Label payload | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoLabelResponse**](DtoLabelResponse.md)

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


## deleteLabel

> string deleteLabel(id, xOrgID)

Delete label (org scoped)

Permanently deletes the label.

### Example

```ts
import {
  Configuration,
  LabelsApi,
} from '@glueops/autoglue-sdk-go';
import type { DeleteLabelRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new LabelsApi(config);

  const body = {
    // string | Label ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteLabelRequest;

  try {
    const data = await api.deleteLabel(body);
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
| **id** | `string` | Label ID (UUID) | [Defaults to `undefined`] |
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


## getLabel

> DtoLabelResponse getLabel(id, xOrgID)

Get label by ID (org scoped)

Returns one label.

### Example

```ts
import {
  Configuration,
  LabelsApi,
} from '@glueops/autoglue-sdk-go';
import type { GetLabelRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new LabelsApi(config);

  const body = {
    // string | Label ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies GetLabelRequest;

  try {
    const data = await api.getLabel(body);
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
| **id** | `string` | Label ID (UUID) | [Defaults to `undefined`] |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoLabelResponse**](DtoLabelResponse.md)

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


## listLabels

> Array&lt;DtoLabelResponse&gt; listLabels(xOrgID, key, value, q)

List node labels (org scoped)

Returns node labels for the organization in X-Org-ID. Filters: &#x60;key&#x60;, &#x60;value&#x60;, and &#x60;q&#x60; (key contains). Add &#x60;include&#x3D;node_pools&#x60; to include linked node groups.

### Example

```ts
import {
  Configuration,
  LabelsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListLabelsRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new LabelsApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // string | Exact key (optional)
    key: key_example,
    // string | Exact value (optional)
    value: value_example,
    // string | Key contains (case-insensitive) (optional)
    q: q_example,
  } satisfies ListLabelsRequest;

  try {
    const data = await api.listLabels(body);
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
| **q** | `string` | Key contains (case-insensitive) | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoLabelResponse&gt;**](DtoLabelResponse.md)

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


## updateLabel

> DtoLabelResponse updateLabel(id, body, xOrgID)

Update label (org scoped)

Partially update label fields.

### Example

```ts
import {
  Configuration,
  LabelsApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateLabelRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new LabelsApi(config);

  const body = {
    // string | Label ID (UUID)
    id: id_example,
    // DtoUpdateLabelRequest | Fields to update
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateLabelRequest;

  try {
    const data = await api.updateLabel(body);
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
| **id** | `string` | Label ID (UUID) | [Defaults to `undefined`] |
| **body** | [DtoUpdateLabelRequest](DtoUpdateLabelRequest.md) | Fields to update | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoLabelResponse**](DtoLabelResponse.md)

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

