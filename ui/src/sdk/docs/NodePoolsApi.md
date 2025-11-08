# NodePoolsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**attachNodePoolAnnotations**](NodePoolsApi.md#attachnodepoolannotations) | **POST** /node-pools/{id}/annotations | Attach annotation to a node pool (org scoped) |
| [**attachNodePoolLabels**](NodePoolsApi.md#attachnodepoollabels) | **POST** /node-pools/{id}/labels | Attach labels to a node pool (org scoped) |
| [**attachNodePoolServers**](NodePoolsApi.md#attachnodepoolservers) | **POST** /node-pools/{id}/servers | Attach servers to a node pool (org scoped) |
| [**attachNodePoolTaints**](NodePoolsApi.md#attachnodepooltaints) | **POST** /node-pools/{id}/taints | Attach taints to a node pool (org scoped) |
| [**createNodePool**](NodePoolsApi.md#createnodepool) | **POST** /node-pools | Create node pool (org scoped) |
| [**deleteNodePool**](NodePoolsApi.md#deletenodepool) | **DELETE** /node-pools/{id} | Delete node pool (org scoped) |
| [**detachNodePoolAnnotation**](NodePoolsApi.md#detachnodepoolannotation) | **DELETE** /node-pools/{id}/annotations/{annotationId} | Detach one annotation from a node pool (org scoped) |
| [**detachNodePoolLabel**](NodePoolsApi.md#detachnodepoollabel) | **DELETE** /node-pools/{id}/labels/{labelId} | Detach one label from a node pool (org scoped) |
| [**detachNodePoolServer**](NodePoolsApi.md#detachnodepoolserver) | **DELETE** /node-pools/{id}/servers/{serverId} | Detach one server from a node pool (org scoped) |
| [**detachNodePoolTaint**](NodePoolsApi.md#detachnodepooltaint) | **DELETE** /node-pools/{id}/taints/{taintId} | Detach one taint from a node pool (org scoped) |
| [**getNodePool**](NodePoolsApi.md#getnodepool) | **GET** /node-pools/{id} | Get node pool by ID (org scoped) |
| [**listNodePoolAnnotations**](NodePoolsApi.md#listnodepoolannotations) | **GET** /node-pools/{id}/annotations | List annotations attached to a node pool (org scoped) |
| [**listNodePoolLabels**](NodePoolsApi.md#listnodepoollabels) | **GET** /node-pools/{id}/labels | List labels attached to a node pool (org scoped) |
| [**listNodePoolServers**](NodePoolsApi.md#listnodepoolservers) | **GET** /node-pools/{id}/servers | List servers attached to a node pool (org scoped) |
| [**listNodePoolTaints**](NodePoolsApi.md#listnodepooltaints) | **GET** /node-pools/{id}/taints | List taints attached to a node pool (org scoped) |
| [**listNodePools**](NodePoolsApi.md#listnodepools) | **GET** /node-pools | List node pools (org scoped) |
| [**updateNodePool**](NodePoolsApi.md#updatenodepool) | **PATCH** /node-pools/{id} | Update node pool (org scoped) |



## attachNodePoolAnnotations

> string attachNodePoolAnnotations(id, body, xOrgID)

Attach annotation to a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachNodePoolAnnotationsRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Group ID (UUID)
    id: id_example,
    // DtoAttachAnnotationsRequest | Annotation IDs to attach
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachNodePoolAnnotationsRequest;

  try {
    const data = await api.attachNodePoolAnnotations(body);
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
| **id** | `string` | Node Group ID (UUID) | [Defaults to `undefined`] |
| **body** | [DtoAttachAnnotationsRequest](DtoAttachAnnotationsRequest.md) | Annotation IDs to attach | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **400** | invalid id / invalid server_ids |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **404** | not found |  -  |
| **500** | attach failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## attachNodePoolLabels

> string attachNodePoolLabels(id, body, xOrgID)

Attach labels to a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachNodePoolLabelsRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // DtoAttachLabelsRequest | Label IDs to attach
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachNodePoolLabelsRequest;

  try {
    const data = await api.attachNodePoolLabels(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **body** | [DtoAttachLabelsRequest](DtoAttachLabelsRequest.md) | Label IDs to attach | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **400** | invalid id / invalid server_ids |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **404** | not found |  -  |
| **500** | attach failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## attachNodePoolServers

> string attachNodePoolServers(id, body, xOrgID)

Attach servers to a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachNodePoolServersRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // DtoAttachServersRequest | Server IDs to attach
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachNodePoolServersRequest;

  try {
    const data = await api.attachNodePoolServers(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **body** | [DtoAttachServersRequest](DtoAttachServersRequest.md) | Server IDs to attach | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **400** | invalid id / invalid server_ids |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **404** | not found |  -  |
| **500** | attach failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## attachNodePoolTaints

> string attachNodePoolTaints(id, body, xOrgID)

Attach taints to a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachNodePoolTaintsRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // DtoAttachTaintsRequest | Taint IDs to attach
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachNodePoolTaintsRequest;

  try {
    const data = await api.attachNodePoolTaints(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **body** | [DtoAttachTaintsRequest](DtoAttachTaintsRequest.md) | Taint IDs to attach | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | No Content |  -  |
| **400** | invalid id / invalid taint_ids |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **404** | not found |  -  |
| **500** | attach failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createNodePool

> DtoNodePoolResponse createNodePool(body, xOrgID)

Create node pool (org scoped)

Creates a node pool. Optionally attach initial servers.

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateNodePoolRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // DtoCreateNodePoolRequest | NodePool payload
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateNodePoolRequest;

  try {
    const data = await api.createNodePool(body);
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
| **body** | [DtoCreateNodePoolRequest](DtoCreateNodePoolRequest.md) | NodePool payload | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoNodePoolResponse**](DtoNodePoolResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | invalid json / missing fields / invalid server_ids |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **500** | create failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteNodePool

> string deleteNodePool(id, xOrgID)

Delete node pool (org scoped)

Permanently deletes the node pool.

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { DeleteNodePoolRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteNodePoolRequest;

  try {
    const data = await api.deleteNodePool(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
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


## detachNodePoolAnnotation

> string detachNodePoolAnnotation(id, annotationId, xOrgID)

Detach one annotation from a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { DetachNodePoolAnnotationRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Annotation ID (UUID)
    annotationId: annotationId_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachNodePoolAnnotationRequest;

  try {
    const data = await api.detachNodePoolAnnotation(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **annotationId** | `string` | Annotation ID (UUID) | [Defaults to `undefined`] |
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
| **404** | not found |  -  |
| **500** | detach failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## detachNodePoolLabel

> string detachNodePoolLabel(id, labelId, xOrgID)

Detach one label from a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { DetachNodePoolLabelRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Label ID (UUID)
    labelId: labelId_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachNodePoolLabelRequest;

  try {
    const data = await api.detachNodePoolLabel(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **labelId** | `string` | Label ID (UUID) | [Defaults to `undefined`] |
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
| **404** | not found |  -  |
| **500** | detach failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## detachNodePoolServer

> string detachNodePoolServer(id, serverId, xOrgID)

Detach one server from a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { DetachNodePoolServerRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Server ID (UUID)
    serverId: serverId_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachNodePoolServerRequest;

  try {
    const data = await api.detachNodePoolServer(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **serverId** | `string` | Server ID (UUID) | [Defaults to `undefined`] |
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
| **404** | not found |  -  |
| **500** | detach failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## detachNodePoolTaint

> string detachNodePoolTaint(id, taintId, xOrgID)

Detach one taint from a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { DetachNodePoolTaintRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Taint ID (UUID)
    taintId: taintId_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachNodePoolTaintRequest;

  try {
    const data = await api.detachNodePoolTaint(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **taintId** | `string` | Taint ID (UUID) | [Defaults to `undefined`] |
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
| **404** | not found |  -  |
| **500** | detach failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getNodePool

> DtoNodePoolResponse getNodePool(id, xOrgID)

Get node pool by ID (org scoped)

Returns one node pool. Add &#x60;include&#x3D;servers&#x60; to include servers.

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { GetNodePoolRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies GetNodePoolRequest;

  try {
    const data = await api.getNodePool(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoNodePoolResponse**](DtoNodePoolResponse.md)

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


## listNodePoolAnnotations

> Array&lt;DtoAnnotationResponse&gt; listNodePoolAnnotations(id, xOrgID)

List annotations attached to a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListNodePoolAnnotationsRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies ListNodePoolAnnotationsRequest;

  try {
    const data = await api.listNodePoolAnnotations(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoAnnotationResponse&gt;**](DtoAnnotationResponse.md)

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


## listNodePoolLabels

> Array&lt;DtoLabelResponse&gt; listNodePoolLabels(id, xOrgID)

List labels attached to a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListNodePoolLabelsRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Label Pool ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies ListNodePoolLabelsRequest;

  try {
    const data = await api.listNodePoolLabels(body);
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
| **id** | `string` | Label Pool ID (UUID) | [Defaults to `undefined`] |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

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
| **400** | invalid id |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **404** | not found |  -  |
| **500** | fetch failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listNodePoolServers

> Array&lt;DtoServerResponse&gt; listNodePoolServers(id, xOrgID)

List servers attached to a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListNodePoolServersRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies ListNodePoolServersRequest;

  try {
    const data = await api.listNodePoolServers(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoServerResponse&gt;**](DtoServerResponse.md)

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


## listNodePoolTaints

> Array&lt;DtoTaintResponse&gt; listNodePoolTaints(id, xOrgID)

List taints attached to a node pool (org scoped)

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListNodePoolTaintsRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies ListNodePoolTaintsRequest;

  try {
    const data = await api.listNodePoolTaints(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

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
| **400** | invalid id |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **404** | not found |  -  |
| **500** | fetch failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listNodePools

> Array&lt;DtoNodePoolResponse&gt; listNodePools(xOrgID, q)

List node pools (org scoped)

Returns node pools for the organization in X-Org-ID.

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListNodePoolsRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // string | Name contains (case-insensitive) (optional)
    q: q_example,
  } satisfies ListNodePoolsRequest;

  try {
    const data = await api.listNodePools(body);
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
| **q** | `string` | Name contains (case-insensitive) | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoNodePoolResponse&gt;**](DtoNodePoolResponse.md)

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
| **500** | failed to list node pools |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateNodePool

> DtoNodePoolResponse updateNodePool(id, body, xOrgID)

Update node pool (org scoped)

Partially update node pool fields.

### Example

```ts
import {
  Configuration,
  NodePoolsApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateNodePoolRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new NodePoolsApi(config);

  const body = {
    // string | Node Pool ID (UUID)
    id: id_example,
    // DtoUpdateNodePoolRequest | Fields to update
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateNodePoolRequest;

  try {
    const data = await api.updateNodePool(body);
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
| **id** | `string` | Node Pool ID (UUID) | [Defaults to `undefined`] |
| **body** | [DtoUpdateNodePoolRequest](DtoUpdateNodePoolRequest.md) | Fields to update | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoNodePoolResponse**](DtoNodePoolResponse.md)

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

