# CredentialsApi

All URIs are relative to *https://autoglue.glueopshosted.com/api/v1*

| Method                                                     | HTTP request                      | Description                                     |
| ---------------------------------------------------------- | --------------------------------- | ----------------------------------------------- |
| [**createCredential**](CredentialsApi.md#createcredential) | **POST** /credentials             | Create a credential (encrypts secret)           |
| [**deleteCredential**](CredentialsApi.md#deletecredential) | **DELETE** /credentials/{id}      | Delete credential                               |
| [**getCredential**](CredentialsApi.md#getcredential)       | **GET** /credentials/{id}         | Get credential by ID (metadata only)            |
| [**listCredentials**](CredentialsApi.md#listcredentials)   | **GET** /credentials              | List credentials (metadata only)                |
| [**revealCredential**](CredentialsApi.md#revealcredential) | **POST** /credentials/{id}/reveal | Reveal decrypted secret (one-time read)         |
| [**updateCredential**](CredentialsApi.md#updatecredential) | **PATCH** /credentials/{id}       | Update credential metadata and/or rotate secret |

## createCredential

> DtoCredentialOut createCredential(dtoCreateCredentialRequest, xOrgID)

Create a credential (encrypts secret)

### Example

```ts
import {
  Configuration,
  CredentialsApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateCredentialRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new CredentialsApi(config);

  const body = {
    // DtoCreateCredentialRequest | Credential payload
    dtoCreateCredentialRequest: ...,
    // string | Organization ID (UUID) (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateCredentialRequest;

  try {
    const data = await api.createCredential(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                           | Type                                                        | Description            | Notes                                |
| ------------------------------ | ----------------------------------------------------------- | ---------------------- | ------------------------------------ |
| **dtoCreateCredentialRequest** | [DtoCreateCredentialRequest](DtoCreateCredentialRequest.md) | Credential payload     |                                      |
| **xOrgID**                     | `string`                                                    | Organization ID (UUID) | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoCredentialOut**](DtoCredentialOut.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **201**     | Created               | -                |
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **500**     | internal server error | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## deleteCredential

> deleteCredential(id, xOrgID)

Delete credential

### Example

```ts
import { Configuration, CredentialsApi } from "@glueops/autoglue-sdk-go";
import type { DeleteCredentialRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new CredentialsApi(config);

  const body = {
    // string | Credential ID (UUID)
    id: id_example,
    // string | Organization ID (UUID) (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteCredentialRequest;

  try {
    const data = await api.deleteCredential(body);
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
| **id**     | `string` | Credential ID (UUID)   | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization ID (UUID) | [Optional] [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description | Response headers |
| ----------- | ----------- | ---------------- |
| **204**     | No Content  | -                |
| **404**     | not found   | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## getCredential

> DtoCredentialOut getCredential(id, xOrgID)

Get credential by ID (metadata only)

### Example

```ts
import { Configuration, CredentialsApi } from "@glueops/autoglue-sdk-go";
import type { GetCredentialRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new CredentialsApi(config);

  const body = {
    // string | Credential ID (UUID)
    id: id_example,
    // string | Organization ID (UUID) (optional)
    xOrgID: xOrgID_example,
  } satisfies GetCredentialRequest;

  try {
    const data = await api.getCredential(body);
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
| **id**     | `string` | Credential ID (UUID)   | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization ID (UUID) | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoCredentialOut**](DtoCredentialOut.md)

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
| **500**     | internal server error | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## listCredentials

> Array&lt;DtoCredentialOut&gt; listCredentials(xOrgID, credentialProvider, kind, scopeKind)

List credentials (metadata only)

Returns credential metadata for the current org. Secrets are never returned.

### Example

```ts
import { Configuration, CredentialsApi } from "@glueops/autoglue-sdk-go";
import type { ListCredentialsRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new CredentialsApi(config);

  const body = {
    // string | Organization ID (UUID) (optional)
    xOrgID: xOrgID_example,
    // string | Filter by provider (e.g., aws) (optional)
    credentialProvider: credentialProvider_example,
    // string | Filter by kind (e.g., aws_access_key) (optional)
    kind: kind_example,
    // string | Filter by scope kind (credential_provider/service/resource) (optional)
    scopeKind: scopeKind_example,
  } satisfies ListCredentialsRequest;

  try {
    const data = await api.listCredentials(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                   | Type     | Description                                                 | Notes                                |
| ---------------------- | -------- | ----------------------------------------------------------- | ------------------------------------ |
| **xOrgID**             | `string` | Organization ID (UUID)                                      | [Optional] [Defaults to `undefined`] |
| **credentialProvider** | `string` | Filter by provider (e.g., aws)                              | [Optional] [Defaults to `undefined`] |
| **kind**               | `string` | Filter by kind (e.g., aws_access_key)                       | [Optional] [Defaults to `undefined`] |
| **scopeKind**          | `string` | Filter by scope kind (credential_provider/service/resource) | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoCredentialOut&gt;**](DtoCredentialOut.md)

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
| **500**     | internal server error | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## revealCredential

> { [key: string]: any; } revealCredential(id, xOrgID, body)

Reveal decrypted secret (one-time read)

### Example

```ts
import { Configuration, CredentialsApi } from "@glueops/autoglue-sdk-go";
import type { RevealCredentialRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new CredentialsApi(config);

  const body = {
    // string | Credential ID (UUID)
    id: id_example,
    // string | Organization ID (UUID) (optional)
    xOrgID: xOrgID_example,
    // object (optional)
    body: Object,
  } satisfies RevealCredentialRequest;

  try {
    const data = await api.revealCredential(body);
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
| **id**     | `string` | Credential ID (UUID)   | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization ID (UUID) | [Optional] [Defaults to `undefined`] |
| **body**   | `object` |                        | [Optional]                           |

### Return type

**{ [key: string]: any; }**

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **403**     | organization required | -                |
| **404**     | not found             | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## updateCredential

> DtoCredentialOut updateCredential(id, dtoUpdateCredentialRequest, xOrgID)

Update credential metadata and/or rotate secret

### Example

```ts
import {
  Configuration,
  CredentialsApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateCredentialRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new CredentialsApi(config);

  const body = {
    // string | Credential ID (UUID)
    id: id_example,
    // DtoUpdateCredentialRequest | Fields to update
    dtoUpdateCredentialRequest: ...,
    // string | Organization ID (UUID) (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateCredentialRequest;

  try {
    const data = await api.updateCredential(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                           | Type                                                        | Description            | Notes                                |
| ------------------------------ | ----------------------------------------------------------- | ---------------------- | ------------------------------------ |
| **id**                         | `string`                                                    | Credential ID (UUID)   | [Defaults to `undefined`]            |
| **dtoUpdateCredentialRequest** | [DtoUpdateCredentialRequest](DtoUpdateCredentialRequest.md) | Fields to update       |                                      |
| **xOrgID**                     | `string`                                                    | Organization ID (UUID) | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoCredentialOut**](DtoCredentialOut.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description       | Response headers |
| ----------- | ----------------- | ---------------- |
| **200**     | OK                | -                |
| **403**     | X-Org-ID required | -                |
| **404**     | not found         | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
