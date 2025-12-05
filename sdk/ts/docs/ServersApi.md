# ServersApi

All URIs are relative to *https://autoglue.glueopshosted.com/api/v1*

| Method                                                     | HTTP request                         | Description                     |
| ---------------------------------------------------------- | ------------------------------------ | ------------------------------- |
| [**createServer**](ServersApi.md#createserver)             | **POST** /servers                    | Create server (org scoped)      |
| [**deleteServer**](ServersApi.md#deleteserver)             | **DELETE** /servers/{id}             | Delete server (org scoped)      |
| [**getServer**](ServersApi.md#getserver)                   | **GET** /servers/{id}                | Get server by ID (org scoped)   |
| [**listServers**](ServersApi.md#listservers)               | **GET** /servers                     | List servers (org scoped)       |
| [**resetServerHostKey**](ServersApi.md#resetserverhostkey) | **POST** /servers/{id}/reset-hostkey | Reset SSH host key (org scoped) |
| [**updateServer**](ServersApi.md#updateserver)             | **PATCH** /servers/{id}              | Update server (org scoped)      |

## createServer

> DtoServerResponse createServer(dtoCreateServerRequest, xOrgID)

Create server (org scoped)

Creates a server bound to the org in X-Org-ID. Validates that ssh_key_id belongs to the org.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateServerRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new ServersApi(config);

  const body = {
    // DtoCreateServerRequest | Server payload
    dtoCreateServerRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateServerRequest;

  try {
    const data = await api.createServer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                       | Type                                                | Description       | Notes                                |
| -------------------------- | --------------------------------------------------- | ----------------- | ------------------------------------ |
| **dtoCreateServerRequest** | [DtoCreateServerRequest](DtoCreateServerRequest.md) | Server payload    |                                      |
| **xOrgID**                 | `string`                                            | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoServerResponse**](DtoServerResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                                                         | Response headers |
| ----------- | ------------------------------------------------------------------- | ---------------- |
| **201**     | Created                                                             | -                |
| **400**     | invalid json / missing fields / invalid status / invalid ssh_key_id | -                |
| **401**     | Unauthorized                                                        | -                |
| **403**     | organization required                                               | -                |
| **500**     | create failed                                                       | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## deleteServer

> deleteServer(id, xOrgID)

Delete server (org scoped)

Permanently deletes the server.

### Example

```ts
import { Configuration, ServersApi } from "@glueops/autoglue-sdk-go";
import type { DeleteServerRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new ServersApi(config);

  const body = {
    // string | Server ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteServerRequest;

  try {
    const data = await api.deleteServer(body);
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
| **id**     | `string` | Server ID (UUID)  | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

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
| **400**     | invalid id            | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **500**     | delete failed         | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## getServer

> DtoServerResponse getServer(id, xOrgID)

Get server by ID (org scoped)

Returns one server in the given organization.

### Example

```ts
import { Configuration, ServersApi } from "@glueops/autoglue-sdk-go";
import type { GetServerRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new ServersApi(config);

  const body = {
    // string | Server ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies GetServerRequest;

  try {
    const data = await api.getServer(body);
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
| **id**     | `string` | Server ID (UUID)  | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoServerResponse**](DtoServerResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | invalid id            | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | not found             | -                |
| **500**     | fetch failed          | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## listServers

> Array&lt;DtoServerResponse&gt; listServers(xOrgID, status, role)

List servers (org scoped)

Returns servers for the organization in X-Org-ID. Optional filters: status, role.

### Example

```ts
import { Configuration, ServersApi } from "@glueops/autoglue-sdk-go";
import type { ListServersRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new ServersApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // string | Filter by status (pending|provisioning|ready|failed) (optional)
    status: status_example,
    // string | Filter by role (optional)
    role: role_example,
  } satisfies ListServersRequest;

  try {
    const data = await api.listServers(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type     | Description               | Notes                                |
| ---------- | -------- | ------------------------- | ------------------------------------ | ----- | ------- | ------------------------------------ |
| **xOrgID** | `string` | Organization UUID         | [Optional] [Defaults to `undefined`] |
| **status** | `string` | Filter by status (pending | provisioning                         | ready | failed) | [Optional] [Defaults to `undefined`] |
| **role**   | `string` | Filter by role            | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoServerResponse&gt;**](DtoServerResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description            | Response headers |
| ----------- | ---------------------- | ---------------- |
| **200**     | OK                     | -                |
| **401**     | Unauthorized           | -                |
| **403**     | organization required  | -                |
| **500**     | failed to list servers | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## resetServerHostKey

> DtoServerResponse resetServerHostKey(id, xOrgID, body)

Reset SSH host key (org scoped)

Clears the stored SSH host key for this server. The next SSH connection will re-learn the host key (trust-on-first-use).

### Example

```ts
import { Configuration, ServersApi } from "@glueops/autoglue-sdk-go";
import type { ResetServerHostKeyRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new ServersApi(config);

  const body = {
    // string | Server ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // object (optional)
    body: Object,
  } satisfies ResetServerHostKeyRequest;

  try {
    const data = await api.resetServerHostKey(body);
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
| **id**     | `string` | Server ID (UUID)  | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |
| **body**   | `object` |                   | [Optional]                           |

### Return type

[**DtoServerResponse**](DtoServerResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **400**     | invalid id            | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | not found             | -                |
| **500**     | reset failed          | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## updateServer

> DtoServerResponse updateServer(id, dtoUpdateServerRequest, xOrgID)

Update server (org scoped)

Partially update fields; changing ssh_key_id validates ownership.

### Example

```ts
import {
  Configuration,
  ServersApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateServerRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new ServersApi(config);

  const body = {
    // string | Server ID (UUID)
    id: id_example,
    // DtoUpdateServerRequest | Fields to update
    dtoUpdateServerRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateServerRequest;

  try {
    const data = await api.updateServer(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                       | Type                                                | Description       | Notes                                |
| -------------------------- | --------------------------------------------------- | ----------------- | ------------------------------------ |
| **id**                     | `string`                                            | Server ID (UUID)  | [Defaults to `undefined`]            |
| **dtoUpdateServerRequest** | [DtoUpdateServerRequest](DtoUpdateServerRequest.md) | Fields to update  |                                      |
| **xOrgID**                 | `string`                                            | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoServerResponse**](DtoServerResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                                                     | Response headers |
| ----------- | --------------------------------------------------------------- | ---------------- |
| **200**     | OK                                                              | -                |
| **400**     | invalid id / invalid json / invalid status / invalid ssh_key_id | -                |
| **401**     | Unauthorized                                                    | -                |
| **403**     | organization required                                           | -                |
| **404**     | not found                                                       | -                |
| **500**     | update failed                                                   | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
