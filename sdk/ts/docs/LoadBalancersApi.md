# LoadBalancersApi

All URIs are relative to *https://autoglue.glueopshosted.com/api/v1*

| Method                                                           | HTTP request                    | Description                         |
| ---------------------------------------------------------------- | ------------------------------- | ----------------------------------- |
| [**createLoadBalancer**](LoadBalancersApi.md#createloadbalancer) | **POST** /load-balancers        | Create a load balancer              |
| [**deleteLoadBalancer**](LoadBalancersApi.md#deleteloadbalancer) | **DELETE** /load-balancers/{id} | Delete a load balancer              |
| [**getLoadBalancers**](LoadBalancersApi.md#getloadbalancers)     | **GET** /load-balancers/{id}    | Get a load balancer (org scoped)    |
| [**listLoadBalancers**](LoadBalancersApi.md#listloadbalancers)   | **GET** /load-balancers         | List load balancers (org scoped)    |
| [**updateLoadBalancer**](LoadBalancersApi.md#updateloadbalancer) | **PATCH** /load-balancers/{id}  | Update a load balancer (org scoped) |

## createLoadBalancer

> DtoLoadBalancerResponse createLoadBalancer(dtoCreateLoadBalancerRequest, xOrgID)

Create a load balancer

### Example

```ts
import {
  Configuration,
  LoadBalancersApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateLoadBalancerRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new LoadBalancersApi(config);

  const body = {
    // DtoCreateLoadBalancerRequest | Record set payload
    dtoCreateLoadBalancerRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateLoadBalancerRequest;

  try {
    const data = await api.createLoadBalancer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                             | Type                                                            | Description        | Notes                                |
| -------------------------------- | --------------------------------------------------------------- | ------------------ | ------------------------------------ |
| **dtoCreateLoadBalancerRequest** | [DtoCreateLoadBalancerRequest](DtoCreateLoadBalancerRequest.md) | Record set payload |                                      |
| **xOrgID**                       | `string`                                                        | Organization UUID  | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoLoadBalancerResponse**](DtoLoadBalancerResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **201**     | Created               | -                |
| **400**     | validation error      | -                |
| **403**     | organization required | -                |
| **404**     | domain not found      | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## deleteLoadBalancer

> deleteLoadBalancer(id, xOrgID)

Delete a load balancer

### Example

```ts
import { Configuration, LoadBalancersApi } from "@glueops/autoglue-sdk-go";
import type { DeleteLoadBalancerRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new LoadBalancersApi(config);

  const body = {
    // string | Load Balancer ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteLoadBalancerRequest;

  try {
    const data = await api.deleteLoadBalancer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type     | Description             | Notes                                |
| ---------- | -------- | ----------------------- | ------------------------------------ |
| **id**     | `string` | Load Balancer ID (UUID) | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID       | [Optional] [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **204**     | No Content            | -                |
| **403**     | organization required | -                |
| **404**     | not found             | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## getLoadBalancers

> Array&lt;DtoLoadBalancerResponse&gt; getLoadBalancers(id, xOrgID)

Get a load balancer (org scoped)

Returns load balancer for the organization in X-Org-ID.

### Example

```ts
import { Configuration, LoadBalancersApi } from "@glueops/autoglue-sdk-go";
import type { GetLoadBalancersRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new LoadBalancersApi(config);

  const body = {
    // string | LoadBalancer ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies GetLoadBalancersRequest;

  try {
    const data = await api.getLoadBalancers(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type     | Description            | Notes                                |
| ---------- | -------- | ---------------------- | ------------------------------------ |
| **id**     | `string` | LoadBalancer ID (UUID) | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID      | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoLoadBalancerResponse&gt;**](DtoLoadBalancerResponse.md)

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

## listLoadBalancers

> Array&lt;DtoLoadBalancerResponse&gt; listLoadBalancers(xOrgID)

List load balancers (org scoped)

Returns load balancers for the organization in X-Org-ID.

### Example

```ts
import { Configuration, LoadBalancersApi } from "@glueops/autoglue-sdk-go";
import type { ListLoadBalancersRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new LoadBalancersApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies ListLoadBalancersRequest;

  try {
    const data = await api.listLoadBalancers(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type     | Description       | Notes                                |
| ---------- | -------- | ----------------- | ------------------------------------ |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoLoadBalancerResponse&gt;**](DtoLoadBalancerResponse.md)

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

## updateLoadBalancer

> DtoLoadBalancerResponse updateLoadBalancer(id, dtoUpdateLoadBalancerRequest, xOrgID)

Update a load balancer (org scoped)

### Example

```ts
import {
  Configuration,
  LoadBalancersApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateLoadBalancerRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new LoadBalancersApi(config);

  const body = {
    // string | Load Balancer ID (UUID)
    id: id_example,
    // DtoUpdateLoadBalancerRequest | Fields to update
    dtoUpdateLoadBalancerRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateLoadBalancerRequest;

  try {
    const data = await api.updateLoadBalancer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                             | Type                                                            | Description             | Notes                                |
| -------------------------------- | --------------------------------------------------------------- | ----------------------- | ------------------------------------ |
| **id**                           | `string`                                                        | Load Balancer ID (UUID) | [Defaults to `undefined`]            |
| **dtoUpdateLoadBalancerRequest** | [DtoUpdateLoadBalancerRequest](DtoUpdateLoadBalancerRequest.md) | Fields to update        |                                      |
| **xOrgID**                       | `string`                                                        | Organization UUID       | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoLoadBalancerResponse**](DtoLoadBalancerResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | validation error      | -                |
| **403**     | organization required | -                |
| **404**     | not found             | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
