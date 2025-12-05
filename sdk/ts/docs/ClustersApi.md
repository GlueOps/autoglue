# ClustersApi

All URIs are relative to *https://autoglue.glueopshosted.com/api/v1*

| Method                                                                        | HTTP request                                              | Description                                                   |
| ----------------------------------------------------------------------------- | --------------------------------------------------------- | ------------------------------------------------------------- |
| [**attachAppsLoadBalancer**](ClustersApi.md#attachappsloadbalancer)           | **POST** /clusters/{clusterID}/apps-load-balancer         | Attach an apps load balancer to a cluster                     |
| [**attachBastionServer**](ClustersApi.md#attachbastionserver)                 | **POST** /clusters/{clusterID}/bastion                    | Attach a bastion server to a cluster                          |
| [**attachCaptainDomain**](ClustersApi.md#attachcaptaindomain)                 | **POST** /clusters/{clusterID}/captain-domain             | Attach a captain domain to a cluster                          |
| [**attachControlPlaneRecordSet**](ClustersApi.md#attachcontrolplanerecordset) | **POST** /clusters/{clusterID}/control-plane-record-set   | Attach a control plane record set to a cluster                |
| [**attachGlueOpsLoadBalancer**](ClustersApi.md#attachglueopsloadbalancer)     | **POST** /clusters/{clusterID}/glueops-load-balancer      | Attach a GlueOps/control-plane load balancer to a cluster     |
| [**attachNodePool**](ClustersApi.md#attachnodepool)                           | **POST** /clusters/{clusterID}/node-pools                 | Attach a node pool to a cluster                               |
| [**clearClusterKubeconfig**](ClustersApi.md#clearclusterkubeconfig)           | **DELETE** /clusters/{clusterID}/kubeconfig               | Clear the kubeconfig for a cluster                            |
| [**createCluster**](ClustersApi.md#createcluster)                             | **POST** /clusters                                        | Create cluster (org scoped)                                   |
| [**deleteCluster**](ClustersApi.md#deletecluster)                             | **DELETE** /clusters/{clusterID}                          | Delete a cluster (org scoped)                                 |
| [**detachAppsLoadBalancer**](ClustersApi.md#detachappsloadbalancer)           | **DELETE** /clusters/{clusterID}/apps-load-balancer       | Detach the apps load balancer from a cluster                  |
| [**detachBastionServer**](ClustersApi.md#detachbastionserver)                 | **DELETE** /clusters/{clusterID}/bastion                  | Detach the bastion server from a cluster                      |
| [**detachCaptainDomain**](ClustersApi.md#detachcaptaindomain)                 | **DELETE** /clusters/{clusterID}/captain-domain           | Detach the captain domain from a cluster                      |
| [**detachControlPlaneRecordSet**](ClustersApi.md#detachcontrolplanerecordset) | **DELETE** /clusters/{clusterID}/control-plane-record-set | Detach the control plane record set from a cluster            |
| [**detachGlueOpsLoadBalancer**](ClustersApi.md#detachglueopsloadbalancer)     | **DELETE** /clusters/{clusterID}/glueops-load-balancer    | Detach the GlueOps/control-plane load balancer from a cluster |
| [**detachNodePool**](ClustersApi.md#detachnodepool)                           | **DELETE** /clusters/{clusterID}/node-pools/{nodePoolID}  | Detach a node pool from a cluster                             |
| [**getCluster**](ClustersApi.md#getcluster)                                   | **GET** /clusters/{clusterID}                             | Get a single cluster by ID (org scoped)                       |
| [**listClusters**](ClustersApi.md#listclusters)                               | **GET** /clusters                                         | List clusters (org scoped)                                    |
| [**setClusterKubeconfig**](ClustersApi.md#setclusterkubeconfig)               | **POST** /clusters/{clusterID}/kubeconfig                 | Set (or replace) the kubeconfig for a cluster                 |
| [**updateCluster**](ClustersApi.md#updatecluster)                             | **PATCH** /clusters/{clusterID}                           | Update basic cluster details (org scoped)                     |

## attachAppsLoadBalancer

> DtoClusterResponse attachAppsLoadBalancer(clusterID, dtoAttachLoadBalancerRequest, xOrgID)

Attach an apps load balancer to a cluster

Sets apps_load_balancer_id on the cluster.

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachAppsLoadBalancerRequest } from '@glueops/autoglue-sdk-go';

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // DtoAttachLoadBalancerRequest | payload
    dtoAttachLoadBalancerRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachAppsLoadBalancerRequest;

  try {
    const data = await api.attachAppsLoadBalancer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                             | Type                                                            | Description       | Notes                                |
| -------------------------------- | --------------------------------------------------------------- | ----------------- | ------------------------------------ |
| **clusterID**                    | `string`                                                        | Cluster ID        | [Defaults to `undefined`]            |
| **dtoAttachLoadBalancerRequest** | [DtoAttachLoadBalancerRequest](DtoAttachLoadBalancerRequest.md) | payload           |                                      |
| **xOrgID**                       | `string`                                                        | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                        | Response headers |
| ----------- | ---------------------------------- | ---------------- |
| **200**     | OK                                 | -                |
| **400**     | bad request                        | -                |
| **401**     | Unauthorized                       | -                |
| **403**     | organization required              | -                |
| **404**     | cluster or load balancer not found | -                |
| **500**     | db error                           | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## attachBastionServer

> DtoClusterResponse attachBastionServer(clusterID, dtoAttachBastionRequest, xOrgID)

Attach a bastion server to a cluster

Sets bastion_server_id on the cluster.

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachBastionServerRequest } from '@glueops/autoglue-sdk-go';

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // DtoAttachBastionRequest | payload
    dtoAttachBastionRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachBastionServerRequest;

  try {
    const data = await api.attachBastionServer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                        | Type                                                  | Description       | Notes                                |
| --------------------------- | ----------------------------------------------------- | ----------------- | ------------------------------------ |
| **clusterID**               | `string`                                              | Cluster ID        | [Defaults to `undefined`]            |
| **dtoAttachBastionRequest** | [DtoAttachBastionRequest](DtoAttachBastionRequest.md) | payload           |                                      |
| **xOrgID**                  | `string`                                              | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                 | Response headers |
| ----------- | --------------------------- | ---------------- |
| **200**     | OK                          | -                |
| **400**     | bad request                 | -                |
| **401**     | Unauthorized                | -                |
| **403**     | organization required       | -                |
| **404**     | cluster or server not found | -                |
| **500**     | db error                    | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## attachCaptainDomain

> DtoClusterResponse attachCaptainDomain(clusterID, dtoAttachCaptainDomainRequest, xOrgID)

Attach a captain domain to a cluster

Sets captain_domain_id on the cluster. Validation of shape happens asynchronously.

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachCaptainDomainRequest } from '@glueops/autoglue-sdk-go';

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // DtoAttachCaptainDomainRequest | payload
    dtoAttachCaptainDomainRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachCaptainDomainRequest;

  try {
    const data = await api.attachCaptainDomain(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                              | Type                                                              | Description       | Notes                                |
| --------------------------------- | ----------------------------------------------------------------- | ----------------- | ------------------------------------ |
| **clusterID**                     | `string`                                                          | Cluster ID        | [Defaults to `undefined`]            |
| **dtoAttachCaptainDomainRequest** | [DtoAttachCaptainDomainRequest](DtoAttachCaptainDomainRequest.md) | payload           |                                      |
| **xOrgID**                        | `string`                                                          | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                 | Response headers |
| ----------- | --------------------------- | ---------------- |
| **200**     | OK                          | -                |
| **400**     | bad request                 | -                |
| **401**     | Unauthorized                | -                |
| **403**     | organization required       | -                |
| **404**     | cluster or domain not found | -                |
| **500**     | db error                    | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## attachControlPlaneRecordSet

> DtoClusterResponse attachControlPlaneRecordSet(clusterID, dtoAttachRecordSetRequest, xOrgID)

Attach a control plane record set to a cluster

Sets control_plane_record_set_id on the cluster.

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachControlPlaneRecordSetRequest } from '@glueops/autoglue-sdk-go';

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // DtoAttachRecordSetRequest | payload
    dtoAttachRecordSetRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachControlPlaneRecordSetRequest;

  try {
    const data = await api.attachControlPlaneRecordSet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                          | Type                                                      | Description       | Notes                                |
| ----------------------------- | --------------------------------------------------------- | ----------------- | ------------------------------------ |
| **clusterID**                 | `string`                                                  | Cluster ID        | [Defaults to `undefined`]            |
| **dtoAttachRecordSetRequest** | [DtoAttachRecordSetRequest](DtoAttachRecordSetRequest.md) | payload           |                                      |
| **xOrgID**                    | `string`                                                  | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                     | Response headers |
| ----------- | ------------------------------- | ---------------- |
| **200**     | OK                              | -                |
| **400**     | bad request                     | -                |
| **401**     | Unauthorized                    | -                |
| **403**     | organization required           | -                |
| **404**     | cluster or record set not found | -                |
| **500**     | db error                        | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## attachGlueOpsLoadBalancer

> DtoClusterResponse attachGlueOpsLoadBalancer(clusterID, dtoAttachLoadBalancerRequest, xOrgID)

Attach a GlueOps/control-plane load balancer to a cluster

Sets glueops_load_balancer_id on the cluster.

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachGlueOpsLoadBalancerRequest } from '@glueops/autoglue-sdk-go';

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // DtoAttachLoadBalancerRequest | payload
    dtoAttachLoadBalancerRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachGlueOpsLoadBalancerRequest;

  try {
    const data = await api.attachGlueOpsLoadBalancer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                             | Type                                                            | Description       | Notes                                |
| -------------------------------- | --------------------------------------------------------------- | ----------------- | ------------------------------------ |
| **clusterID**                    | `string`                                                        | Cluster ID        | [Defaults to `undefined`]            |
| **dtoAttachLoadBalancerRequest** | [DtoAttachLoadBalancerRequest](DtoAttachLoadBalancerRequest.md) | payload           |                                      |
| **xOrgID**                       | `string`                                                        | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                        | Response headers |
| ----------- | ---------------------------------- | ---------------- |
| **200**     | OK                                 | -                |
| **400**     | bad request                        | -                |
| **401**     | Unauthorized                       | -                |
| **403**     | organization required              | -                |
| **404**     | cluster or load balancer not found | -                |
| **500**     | db error                           | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## attachNodePool

> DtoClusterResponse attachNodePool(clusterID, dtoAttachNodePoolRequest, xOrgID)

Attach a node pool to a cluster

Adds an entry in the cluster_node_pools join table.

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { AttachNodePoolRequest } from '@glueops/autoglue-sdk-go';

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // DtoAttachNodePoolRequest | payload
    dtoAttachNodePoolRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies AttachNodePoolRequest;

  try {
    const data = await api.attachNodePool(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                         | Type                                                    | Description       | Notes                                |
| ---------------------------- | ------------------------------------------------------- | ----------------- | ------------------------------------ |
| **clusterID**                | `string`                                                | Cluster ID        | [Defaults to `undefined`]            |
| **dtoAttachNodePoolRequest** | [DtoAttachNodePoolRequest](DtoAttachNodePoolRequest.md) | payload           |                                      |
| **xOrgID**                   | `string`                                                | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                    | Response headers |
| ----------- | ------------------------------ | ---------------- |
| **200**     | OK                             | -                |
| **400**     | bad request                    | -                |
| **401**     | Unauthorized                   | -                |
| **403**     | organization required          | -                |
| **404**     | cluster or node pool not found | -                |
| **500**     | db error                       | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## clearClusterKubeconfig

> DtoClusterResponse clearClusterKubeconfig(clusterID, xOrgID)

Clear the kubeconfig for a cluster

Removes the encrypted kubeconfig, IV, and tag from the cluster record.

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { ClearClusterKubeconfigRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies ClearClusterKubeconfigRequest;

  try {
    const data = await api.clearClusterKubeconfig(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name          | Type     | Description       | Notes                                |
| ------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID** | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **xOrgID**    | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## createCluster

> DtoClusterResponse createCluster(dtoCreateClusterRequest, xOrgID)

Create cluster (org scoped)

Creates a cluster. Status is managed by the system and starts as &#x60;pre_pending&#x60; for validation.

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
    dtoCreateClusterRequest: ...,
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

| Name                        | Type                                                  | Description       | Notes                                |
| --------------------------- | ----------------------------------------------------- | ----------------- | ------------------------------------ |
| **dtoCreateClusterRequest** | [DtoCreateClusterRequest](DtoCreateClusterRequest.md) | payload           |                                      |
| **xOrgID**                  | `string`                                              | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **201**     | Created               | -                |
| **400**     | invalid json          | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **500**     | create failed         | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## deleteCluster

> string deleteCluster(clusterID, xOrgID)

Delete a cluster (org scoped)

Deletes the cluster. Related resources are cleaned up via DB constraints (e.g. CASCADE).

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { DeleteClusterRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteClusterRequest;

  try {
    const data = await api.deleteCluster(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name          | Type     | Description       | Notes                                |
| ------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID** | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **xOrgID**    | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **204**     | deleted               | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## detachAppsLoadBalancer

> DtoClusterResponse detachAppsLoadBalancer(clusterID, xOrgID)

Detach the apps load balancer from a cluster

Clears apps_load_balancer_id on the cluster.

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { DetachAppsLoadBalancerRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachAppsLoadBalancerRequest;

  try {
    const data = await api.detachAppsLoadBalancer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name          | Type     | Description       | Notes                                |
| ------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID** | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **xOrgID**    | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## detachBastionServer

> DtoClusterResponse detachBastionServer(clusterID, xOrgID)

Detach the bastion server from a cluster

Clears bastion_server_id on the cluster.

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { DetachBastionServerRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachBastionServerRequest;

  try {
    const data = await api.detachBastionServer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name          | Type     | Description       | Notes                                |
| ------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID** | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **xOrgID**    | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## detachCaptainDomain

> DtoClusterResponse detachCaptainDomain(clusterID, xOrgID)

Detach the captain domain from a cluster

Clears captain_domain_id on the cluster. This will likely cause the cluster to become incomplete.

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { DetachCaptainDomainRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachCaptainDomainRequest;

  try {
    const data = await api.detachCaptainDomain(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name          | Type     | Description       | Notes                                |
| ------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID** | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **xOrgID**    | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## detachControlPlaneRecordSet

> DtoClusterResponse detachControlPlaneRecordSet(clusterID, xOrgID)

Detach the control plane record set from a cluster

Clears control_plane_record_set_id on the cluster.

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { DetachControlPlaneRecordSetRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachControlPlaneRecordSetRequest;

  try {
    const data = await api.detachControlPlaneRecordSet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name          | Type     | Description       | Notes                                |
| ------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID** | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **xOrgID**    | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## detachGlueOpsLoadBalancer

> DtoClusterResponse detachGlueOpsLoadBalancer(clusterID, xOrgID)

Detach the GlueOps/control-plane load balancer from a cluster

Clears glueops_load_balancer_id on the cluster.

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { DetachGlueOpsLoadBalancerRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachGlueOpsLoadBalancerRequest;

  try {
    const data = await api.detachGlueOpsLoadBalancer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name          | Type     | Description       | Notes                                |
| ------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID** | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **xOrgID**    | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## detachNodePool

> DtoClusterResponse detachNodePool(clusterID, nodePoolID, xOrgID)

Detach a node pool from a cluster

Removes an entry from the cluster_node_pools join table.

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { DetachNodePoolRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Node Pool ID
    nodePoolID: nodePoolID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DetachNodePoolRequest;

  try {
    const data = await api.detachNodePool(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name           | Type     | Description       | Notes                                |
| -------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID**  | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **nodePoolID** | `string` | Node Pool ID      | [Defaults to `undefined`]            |
| **xOrgID**     | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                    | Response headers |
| ----------- | ------------------------------ | ---------------- |
| **200**     | OK                             | -                |
| **400**     | bad request                    | -                |
| **401**     | Unauthorized                   | -                |
| **403**     | organization required          | -                |
| **404**     | cluster or node pool not found | -                |
| **500**     | db error                       | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## getCluster

> DtoClusterResponse getCluster(clusterID, xOrgID)

Get a single cluster by ID (org scoped)

Returns a cluster with all related resources (domain, record set, load balancers, bastion, node pools).

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { GetClusterRequest } from "@glueops/autoglue-sdk-go";

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies GetClusterRequest;

  try {
    const data = await api.getCluster(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name          | Type     | Description       | Notes                                |
| ------------- | -------- | ----------------- | ------------------------------------ |
| **clusterID** | `string` | Cluster ID        | [Defaults to `undefined`]            |
| **xOrgID**    | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## listClusters

> Array&lt;DtoClusterResponse&gt; listClusters(xOrgID, q)

List clusters (org scoped)

Returns clusters for the organization in X-Org-ID. Filter by &#x60;q&#x60; (name contains).

### Example

```ts
import { Configuration, ClustersApi } from "@glueops/autoglue-sdk-go";
import type { ListClustersRequest } from "@glueops/autoglue-sdk-go";

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

| Name       | Type     | Description                      | Notes                                |
| ---------- | -------- | -------------------------------- | ------------------------------------ |
| **xOrgID** | `string` | Organization UUID                | [Optional] [Defaults to `undefined`] |
| **q**      | `string` | Name contains (case-insensitive) | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoClusterResponse&gt;**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description             | Response headers |
| ----------- | ----------------------- | ---------------- |
| **200**     | OK                      | -                |
| **401**     | Unauthorized            | -                |
| **403**     | organization required   | -                |
| **500**     | failed to list clusters | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## setClusterKubeconfig

> DtoClusterResponse setClusterKubeconfig(clusterID, dtoSetKubeconfigRequest, xOrgID)

Set (or replace) the kubeconfig for a cluster

Stores the kubeconfig encrypted per organization. The kubeconfig is never returned in responses.

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { SetClusterKubeconfigRequest } from '@glueops/autoglue-sdk-go';

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // DtoSetKubeconfigRequest | payload
    dtoSetKubeconfigRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies SetClusterKubeconfigRequest;

  try {
    const data = await api.setClusterKubeconfig(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                        | Type                                                  | Description       | Notes                                |
| --------------------------- | ----------------------------------------------------- | ----------------- | ------------------------------------ |
| **clusterID**               | `string`                                              | Cluster ID        | [Defaults to `undefined`]            |
| **dtoSetKubeconfigRequest** | [DtoSetKubeconfigRequest](DtoSetKubeconfigRequest.md) | payload           |                                      |
| **xOrgID**                  | `string`                                              | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## updateCluster

> DtoClusterResponse updateCluster(clusterID, dtoUpdateClusterRequest, xOrgID)

Update basic cluster details (org scoped)

Updates the cluster name, provider, and/or region. Status is managed by the system.

### Example

```ts
import {
  Configuration,
  ClustersApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateClusterRequest } from '@glueops/autoglue-sdk-go';

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
    // string | Cluster ID
    clusterID: clusterID_example,
    // DtoUpdateClusterRequest | payload
    dtoUpdateClusterRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateClusterRequest;

  try {
    const data = await api.updateCluster(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                        | Type                                                  | Description       | Notes                                |
| --------------------------- | ----------------------------------------------------- | ----------------- | ------------------------------------ |
| **clusterID**               | `string`                                              | Cluster ID        | [Defaults to `undefined`]            |
| **dtoUpdateClusterRequest** | [DtoUpdateClusterRequest](DtoUpdateClusterRequest.md) | payload           |                                      |
| **xOrgID**                  | `string`                                              | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoClusterResponse**](DtoClusterResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | bad request           | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | cluster not found     | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
