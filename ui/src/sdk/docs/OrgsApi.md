# OrgsApi

All URIs are relative to *http://localhost:8080/api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addOrUpdateMember**](OrgsApi.md#addorupdatemember) | **POST** /orgs/{id}/members | Add or update a member (owner/admin) |
| [**createOrg**](OrgsApi.md#createorg) | **POST** /orgs | Create organization |
| [**createOrgKey**](OrgsApi.md#createorgkey) | **POST** /orgs/{id}/api-keys | Create org key/secret pair (owner/admin) |
| [**deleteOrg**](OrgsApi.md#deleteorg) | **DELETE** /orgs/{id} | Delete organization (owner) |
| [**deleteOrgKey**](OrgsApi.md#deleteorgkey) | **DELETE** /orgs/{id}/api-keys/{key_id} | Delete org key (owner/admin) |
| [**getOrg**](OrgsApi.md#getorg) | **GET** /orgs/{id} | Get organization |
| [**listMembers**](OrgsApi.md#listmembers) | **GET** /orgs/{id}/members | List members in org |
| [**listMyOrgs**](OrgsApi.md#listmyorgs) | **GET** /orgs | List organizations I belong to |
| [**listOrgKeys**](OrgsApi.md#listorgkeys) | **GET** /orgs/{id}/api-keys | List org-scoped API keys (no secrets) |
| [**removeMember**](OrgsApi.md#removemember) | **DELETE** /orgs/{id}/members/{user_id} | Remove a member (owner/admin) |
| [**updateOrg**](OrgsApi.md#updateorg) | **PATCH** /orgs/{id} | Update organization (owner/admin) |



## addOrUpdateMember

> HandlersMemberOut addOrUpdateMember(id, body)

Add or update a member (owner/admin)

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { AddOrUpdateMemberRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
    // HandlersMemberUpsertReq | User & role
    body: ...,
  } satisfies AddOrUpdateMemberRequest;

  try {
    const data = await api.addOrUpdateMember(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |
| **body** | [HandlersMemberUpsertReq](HandlersMemberUpsertReq.md) | User &amp; role | |

### Return type

[**HandlersMemberOut**](HandlersMemberOut.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createOrg

> ModelsOrganization createOrg(body)

Create organization

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateOrgRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // HandlersOrgCreateReq | Org payload
    body: ...,
  } satisfies CreateOrgRequest;

  try {
    const data = await api.createOrg(body);
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
| **body** | [HandlersOrgCreateReq](HandlersOrgCreateReq.md) | Org payload | |

### Return type

[**ModelsOrganization**](ModelsOrganization.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **400** | Bad Request |  -  |
| **401** | Unauthorized |  -  |
| **409** | Conflict |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createOrgKey

> HandlersOrgKeyCreateResp createOrgKey(id, body)

Create org key/secret pair (owner/admin)

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateOrgKeyRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
    // HandlersOrgKeyCreateReq | Key name + optional expiry
    body: ...,
  } satisfies CreateOrgKeyRequest;

  try {
    const data = await api.createOrgKey(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |
| **body** | [HandlersOrgKeyCreateReq](HandlersOrgKeyCreateReq.md) | Key name + optional expiry | |

### Return type

[**HandlersOrgKeyCreateResp**](HandlersOrgKeyCreateResp.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteOrg

> deleteOrg(id)

Delete organization (owner)

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { DeleteOrgRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
  } satisfies DeleteOrgRequest;

  try {
    const data = await api.deleteOrg(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Deleted |  -  |
| **401** | Unauthorized |  -  |
| **404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteOrgKey

> deleteOrgKey(id, keyId)

Delete org key (owner/admin)

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { DeleteOrgKeyRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
    // string | Key ID (UUID)
    keyId: keyId_example,
  } satisfies DeleteOrgKeyRequest;

  try {
    const data = await api.deleteOrgKey(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |
| **keyId** | `string` | Key ID (UUID) | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Deleted |  -  |
| **401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getOrg

> ModelsOrganization getOrg(id)

Get organization

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { GetOrgRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
  } satisfies GetOrgRequest;

  try {
    const data = await api.getOrg(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |

### Return type

[**ModelsOrganization**](ModelsOrganization.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |
| **404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listMembers

> Array&lt;HandlersMemberOut&gt; listMembers(id)

List members in org

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListMembersRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
  } satisfies ListMembersRequest;

  try {
    const data = await api.listMembers(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |

### Return type

[**Array&lt;HandlersMemberOut&gt;**](HandlersMemberOut.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listMyOrgs

> Array&lt;ModelsOrganization&gt; listMyOrgs()

List organizations I belong to

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListMyOrgsRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  try {
    const data = await api.listMyOrgs();
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

[**Array&lt;ModelsOrganization&gt;**](ModelsOrganization.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listOrgKeys

> Array&lt;ModelsAPIKey&gt; listOrgKeys(id)

List org-scoped API keys (no secrets)

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { ListOrgKeysRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
  } satisfies ListOrgKeysRequest;

  try {
    const data = await api.listOrgKeys(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsAPIKey&gt;**](ModelsAPIKey.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## removeMember

> removeMember(id, userId)

Remove a member (owner/admin)

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { RemoveMemberRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
    // string | User ID (UUID)
    userId: userId_example,
  } satisfies RemoveMemberRequest;

  try {
    const data = await api.removeMember(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |
| **userId** | `string` | User ID (UUID) | [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Removed |  -  |
| **401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateOrg

> ModelsOrganization updateOrg(id, body)

Update organization (owner/admin)

### Example

```ts
import {
  Configuration,
  OrgsApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateOrgRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new OrgsApi(config);

  const body = {
    // string | Org ID (UUID)
    id: id_example,
    // HandlersOrgUpdateReq | Update payload
    body: ...,
  } satisfies UpdateOrgRequest;

  try {
    const data = await api.updateOrg(body);
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
| **id** | `string` | Org ID (UUID) | [Defaults to `undefined`] |
| **body** | [HandlersOrgUpdateReq](HandlersOrgUpdateReq.md) | Update payload | |

### Return type

[**ModelsOrganization**](ModelsOrganization.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **401** | Unauthorized |  -  |
| **404** | Not Found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

