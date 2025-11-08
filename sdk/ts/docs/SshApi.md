# SshApi

All URIs are relative to _/api/v1_

| Method                                               | HTTP request               | Description                               |
| ---------------------------------------------------- | -------------------------- | ----------------------------------------- |
| [**createSSHKey**](SshApi.md#createsshkey)           | **POST** /ssh              | Create ssh keypair (org scoped)           |
| [**deleteSSHKey**](SshApi.md#deletesshkey)           | **DELETE** /ssh/{id}       | Delete ssh keypair (org scoped)           |
| [**downloadSSHKey**](SshApi.md#downloadsshkey)       | **GET** /ssh/{id}/download | Download ssh key files by ID (org scoped) |
| [**getSSHKey**](SshApi.md#getsshkey)                 | **GET** /ssh/{id}          | Get ssh key by ID (org scoped)            |
| [**listPublicSshKeys**](SshApi.md#listpublicsshkeys) | **GET** /ssh               | List ssh keys (org scoped)                |

## createSSHKey

> DtoSshResponse createSSHKey(body, xOrgID)

Create ssh keypair (org scoped)

Generates an RSA or ED25519 keypair, saves it, and returns metadata. For RSA you may set bits (2048/3072/4096). Default is 4096. ED25519 ignores bits.

### Example

```ts
import {
  Configuration,
  SshApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateSSHKeyRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new SshApi(config);

  const body = {
    // DtoCreateSSHRequest | Key generation options
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateSSHKeyRequest;

  try {
    const data = await api.createSSHKey(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type                                          | Description            | Notes                                |
| ---------- | --------------------------------------------- | ---------------------- | ------------------------------------ |
| **body**   | [DtoCreateSSHRequest](DtoCreateSSHRequest.md) | Key generation options |                                      |
| **xOrgID** | `string`                                      | Organization UUID      | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoSshResponse**](DtoSshResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                 | Response headers |
| ----------- | --------------------------- | ---------------- |
| **201**     | Created                     | -                |
| **400**     | invalid json / invalid bits | -                |
| **401**     | Unauthorized                | -                |
| **403**     | organization required       | -                |
| **500**     | generation/create failed    | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## deleteSSHKey

> string deleteSSHKey(id, xOrgID)

Delete ssh keypair (org scoped)

Permanently deletes a keypair.

### Example

```ts
import { Configuration, SshApi } from "@glueops/autoglue-sdk-go";
import type { DeleteSSHKeyRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new SshApi(config);

  const body = {
    // string | SSH Key ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteSSHKeyRequest;

  try {
    const data = await api.deleteSSHKey(body);
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
| **id**     | `string` | SSH Key ID (UUID) | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

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
| **204**     | No Content            | -                |
| **400**     | invalid id            | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **500**     | delete failed         | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## downloadSSHKey

> string downloadSSHKey(xOrgID, id, part)

Download ssh key files by ID (org scoped)

Download &#x60;part&#x3D;public|private|both&#x60; of the keypair. &#x60;both&#x60; returns a zip file.

### Example

```ts
import { Configuration, SshApi } from "@glueops/autoglue-sdk-go";
import type { DownloadSSHKeyRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new SshApi(config);

  const body = {
    // string | Organization UUID
    xOrgID: xOrgID_example,
    // string | SSH Key ID (UUID)
    id: id_example,
    // 'public' | 'private' | 'both' | Which part to download
    part: part_example,
  } satisfies DownloadSSHKeyRequest;

  try {
    const data = await api.downloadSSHKey(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type                        | Description            | Notes                                                   |
| ---------- | --------------------------- | ---------------------- | ------------------------------------------------------- |
| **xOrgID** | `string`                    | Organization UUID      | [Defaults to `undefined`]                               |
| **id**     | `string`                    | SSH Key ID (UUID)      | [Defaults to `undefined`]                               |
| **part**   | `public`, `private`, `both` | Which part to download | [Defaults to `undefined`] [Enum: public, private, both] |

### Return type

**string**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description               | Response headers |
| ----------- | ------------------------- | ---------------- |
| **200**     | file content              | -                |
| **400**     | invalid id / invalid part | -                |
| **401**     | Unauthorized              | -                |
| **403**     | organization required     | -                |
| **404**     | not found                 | -                |
| **500**     | download failed           | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## getSSHKey

> DtoSshRevealResponse getSSHKey(id, xOrgID, reveal)

Get ssh key by ID (org scoped)

Returns public key fields. Append &#x60;?reveal&#x3D;true&#x60; to include the private key PEM.

### Example

```ts
import { Configuration, SshApi } from "@glueops/autoglue-sdk-go";
import type { GetSSHKeyRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new SshApi(config);

  const body = {
    // string | SSH Key ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // boolean | Reveal private key PEM (optional)
    reveal: true,
  } satisfies GetSSHKeyRequest;

  try {
    const data = await api.getSSHKey(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type      | Description            | Notes                                |
| ---------- | --------- | ---------------------- | ------------------------------------ |
| **id**     | `string`  | SSH Key ID (UUID)      | [Defaults to `undefined`]            |
| **xOrgID** | `string`  | Organization UUID      | [Optional] [Defaults to `undefined`] |
| **reveal** | `boolean` | Reveal private key PEM | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoSshRevealResponse**](DtoSshRevealResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | When reveal&#x3D;true | -                |
| **400**     | invalid id            | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **404**     | not found             | -                |
| **500**     | fetch failed          | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## listPublicSshKeys

> Array&lt;DtoSshResponse&gt; listPublicSshKeys(xOrgID)

List ssh keys (org scoped)

Returns ssh keys for the organization in X-Org-ID.

### Example

```ts
import { Configuration, SshApi } from "@glueops/autoglue-sdk-go";
import type { ListPublicSshKeysRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new SshApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies ListPublicSshKeysRequest;

  try {
    const data = await api.listPublicSshKeys(body);
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

[**Array&lt;DtoSshResponse&gt;**](DtoSshResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **500**     | failed to list keys   | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
