# ClustersApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createCluster**](ClustersApi.md#createcluster) | **POST** /clusters | Create cluster (org scoped) |
| [**listClusters**](ClustersApi.md#listclusters) | **GET** /clusters | List clusters (org scoped) |



## createCluster

> DtoClusterResponse createCluster(body, xOrgID)

Create cluster (org scoped)

Creates a cluster. If &#x60;kubeconfig&#x60; is provided, it will be encrypted per-organization and stored securely (never returned).

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateClusterRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new ClustersApi(config);

  const body = {
    // DtoCreateClusterRequest | payload
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateClusterRequest;

  try {
    const data = await api.createCluster(body);
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
| **body** | [DtoCreateClusterRequest](DtoCreateClusterRequest.md) | payload | |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | invalid json |  -  |
| **401** | Unauthorized |  -  |
| **403** | organization required |  -  |
| **500** | create failed |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listClusters

> Array&lt;DtoClusterResponse&gt; listClusters(xOrgID, q)

List clusters (org scoped)

Returns clusters for the organization in X-Org-ID. Filter by &#x60;q&#x60; (name contains).

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { ListClustersRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new ClustersApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // string | Name contains (case-insensitive) (optional)
    q: q_example,
  } satisfies ListClustersRequest;

  try {
    const data = await api.listClusters(body);
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

[**Array&lt;DtoClusterResponse&gt;**](DtoClusterResponse.md)

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
| **500** | failed to list clusters |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

