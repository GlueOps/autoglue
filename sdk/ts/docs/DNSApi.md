# DNSApi

All URIs are relative to *https://autoglue.glueopshosted.com/api/v1*

| Method                                           | HTTP request                              | Description                                                                                  |
| ------------------------------------------------ | ----------------------------------------- | -------------------------------------------------------------------------------------------- |
| [**createDomain**](DNSApi.md#createdomain)       | **POST** /dns/domains                     | Create a domain (org scoped)                                                                 |
| [**createRecordSet**](DNSApi.md#createrecordset) | **POST** /dns/domains/{domain_id}/records | Create a record set (pending; Archer will UPSERT to Route 53)                                |
| [**deleteDomain**](DNSApi.md#deletedomain)       | **DELETE** /dns/domains/{id}              | Delete a domain                                                                              |
| [**deleteRecordSet**](DNSApi.md#deleterecordset) | **DELETE** /dns/records/{id}              | Delete a record set (API removes row; worker can optionally handle external deletion policy) |
| [**getDomain**](DNSApi.md#getdomain)             | **GET** /dns/domains/{id}                 | Get a domain (org scoped)                                                                    |
| [**listDomains**](DNSApi.md#listdomains)         | **GET** /dns/domains                      | List domains (org scoped)                                                                    |
| [**listRecordSets**](DNSApi.md#listrecordsets)   | **GET** /dns/domains/{domain_id}/records  | List record sets for a domain                                                                |
| [**updateDomain**](DNSApi.md#updatedomain)       | **PATCH** /dns/domains/{id}               | Update a domain (org scoped)                                                                 |
| [**updateRecordSet**](DNSApi.md#updaterecordset) | **PATCH** /dns/records/{id}               | Update a record set (flips to pending for reconciliation)                                    |

## createDomain

> DtoDomainResponse createDomain(dtoCreateDomainRequest, xOrgID)

Create a domain (org scoped)

Creates a domain bound to a Route 53 scoped credential. Archer will backfill ZoneID if omitted.

### Example

```ts
import {
  Configuration,
  DNSApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateDomainRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new DNSApi(config);

  const body = {
    // DtoCreateDomainRequest | Domain payload
    dtoCreateDomainRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateDomainRequest;

  try {
    const data = await api.createDomain(body);
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
| **dtoCreateDomainRequest** | [DtoCreateDomainRequest](DtoCreateDomainRequest.md) | Domain payload    |                                      |
| **xOrgID**                 | `string`                                            | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoDomainResponse**](DtoDomainResponse.md)

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
| **401**     | Unauthorized          | -                |
| **403**     | organization required | -                |
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## createRecordSet

> DtoRecordSetResponse createRecordSet(domainId, dtoCreateRecordSetRequest, xOrgID)

Create a record set (pending; Archer will UPSERT to Route 53)

### Example

```ts
import {
  Configuration,
  DNSApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateRecordSetRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new DNSApi(config);

  const body = {
    // string | Domain ID (UUID)
    domainId: domainId_example,
    // DtoCreateRecordSetRequest | Record set payload
    dtoCreateRecordSetRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateRecordSetRequest;

  try {
    const data = await api.createRecordSet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                          | Type                                                      | Description        | Notes                                |
| ----------------------------- | --------------------------------------------------------- | ------------------ | ------------------------------------ |
| **domainId**                  | `string`                                                  | Domain ID (UUID)   | [Defaults to `undefined`]            |
| **dtoCreateRecordSetRequest** | [DtoCreateRecordSetRequest](DtoCreateRecordSetRequest.md) | Record set payload |                                      |
| **xOrgID**                    | `string`                                                  | Organization UUID  | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoRecordSetResponse**](DtoRecordSetResponse.md)

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

## deleteDomain

> deleteDomain(id, xOrgID)

Delete a domain

### Example

```ts
import { Configuration, DNSApi } from "@glueops/autoglue-sdk-go";
import type { DeleteDomainRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new DNSApi(config);

  const body = {
    // string | Domain ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteDomainRequest;

  try {
    const data = await api.deleteDomain(body);
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
| **id**     | `string` | Domain ID (UUID)  | [Defaults to `undefined`]            |
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
| **403**     | organization required | -                |
| **404**     | not found             | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## deleteRecordSet

> deleteRecordSet(id, xOrgID)

Delete a record set (API removes row; worker can optionally handle external deletion policy)

### Example

```ts
import { Configuration, DNSApi } from "@glueops/autoglue-sdk-go";
import type { DeleteRecordSetRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new DNSApi(config);

  const body = {
    // string | Record Set ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteRecordSetRequest;

  try {
    const data = await api.deleteRecordSet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type     | Description          | Notes                                |
| ---------- | -------- | -------------------- | ------------------------------------ |
| **id**     | `string` | Record Set ID (UUID) | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID    | [Optional] [Defaults to `undefined`] |

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

## getDomain

> DtoDomainResponse getDomain(id, xOrgID)

Get a domain (org scoped)

### Example

```ts
import { Configuration, DNSApi } from "@glueops/autoglue-sdk-go";
import type { GetDomainRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new DNSApi(config);

  const body = {
    // string | Domain ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies GetDomainRequest;

  try {
    const data = await api.getDomain(body);
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
| **id**     | `string` | Domain ID (UUID)  | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoDomainResponse**](DtoDomainResponse.md)

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
| **404**     | not found             | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## listDomains

> Array&lt;DtoDomainResponse&gt; listDomains(xOrgID, domainName, status, q)

List domains (org scoped)

Returns domains for X-Org-ID. Filters: &#x60;domain_name&#x60;, &#x60;status&#x60;, &#x60;q&#x60; (contains).

### Example

```ts
import { Configuration, DNSApi } from "@glueops/autoglue-sdk-go";
import type { ListDomainsRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new DNSApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // string | Exact domain name (lowercase, no trailing dot) (optional)
    domainName: domainName_example,
    // string | pending|provisioning|ready|failed (optional)
    status: status_example,
    // string | Domain contains (case-insensitive) (optional)
    q: q_example,
  } satisfies ListDomainsRequest;

  try {
    const data = await api.listDomains(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name           | Type     | Description                                    | Notes                                |
| -------------- | -------- | ---------------------------------------------- | ------------------------------------ | ----- | ------ | ------------------------------------ |
| **xOrgID**     | `string` | Organization UUID                              | [Optional] [Defaults to `undefined`] |
| **domainName** | `string` | Exact domain name (lowercase, no trailing dot) | [Optional] [Defaults to `undefined`] |
| **status**     | `string` | pending                                        | provisioning                         | ready | failed | [Optional] [Defaults to `undefined`] |
| **q**          | `string` | Domain contains (case-insensitive)             | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoDomainResponse&gt;**](DtoDomainResponse.md)

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
| **500**     | db error              | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## listRecordSets

> Array&lt;DtoRecordSetResponse&gt; listRecordSets(domainId, xOrgID, name, type, status)

List record sets for a domain

Filters: &#x60;name&#x60;, &#x60;type&#x60;, &#x60;status&#x60;.

### Example

```ts
import { Configuration, DNSApi } from "@glueops/autoglue-sdk-go";
import type { ListRecordSetsRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new DNSApi(config);

  const body = {
    // string | Domain ID (UUID)
    domainId: domainId_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // string | Exact relative name or FQDN (server normalizes) (optional)
    name: name_example,
    // string | RR type (A, AAAA, CNAME, TXT, MX, NS, SRV, CAA) (optional)
    type: type_example,
    // string | pending|provisioning|ready|failed (optional)
    status: status_example,
  } satisfies ListRecordSetsRequest;

  try {
    const data = await api.listRecordSets(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name         | Type     | Description                                     | Notes                                |
| ------------ | -------- | ----------------------------------------------- | ------------------------------------ | ----- | ------ | ------------------------------------ |
| **domainId** | `string` | Domain ID (UUID)                                | [Defaults to `undefined`]            |
| **xOrgID**   | `string` | Organization UUID                               | [Optional] [Defaults to `undefined`] |
| **name**     | `string` | Exact relative name or FQDN (server normalizes) | [Optional] [Defaults to `undefined`] |
| **type**     | `string` | RR type (A, AAAA, CNAME, TXT, MX, NS, SRV, CAA) | [Optional] [Defaults to `undefined`] |
| **status**   | `string` | pending                                         | provisioning                         | ready | failed | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoRecordSetResponse&gt;**](DtoRecordSetResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description           | Response headers |
| ----------- | --------------------- | ---------------- |
| **200**     | OK                    | -                |
| **403**     | organization required | -                |
| **404**     | domain not found      | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## updateDomain

> DtoDomainResponse updateDomain(id, dtoUpdateDomainRequest, xOrgID)

Update a domain (org scoped)

### Example

```ts
import {
  Configuration,
  DNSApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateDomainRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new DNSApi(config);

  const body = {
    // string | Domain ID (UUID)
    id: id_example,
    // DtoUpdateDomainRequest | Fields to update
    dtoUpdateDomainRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateDomainRequest;

  try {
    const data = await api.updateDomain(body);
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
| **id**                     | `string`                                            | Domain ID (UUID)  | [Defaults to `undefined`]            |
| **dtoUpdateDomainRequest** | [DtoUpdateDomainRequest](DtoUpdateDomainRequest.md) | Fields to update  |                                      |
| **xOrgID**                 | `string`                                            | Organization UUID | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoDomainResponse**](DtoDomainResponse.md)

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

## updateRecordSet

> DtoRecordSetResponse updateRecordSet(id, dtoUpdateRecordSetRequest, xOrgID)

Update a record set (flips to pending for reconciliation)

### Example

```ts
import {
  Configuration,
  DNSApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateRecordSetRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new DNSApi(config);

  const body = {
    // string | Record Set ID (UUID)
    id: id_example,
    // DtoUpdateRecordSetRequest | Fields to update
    dtoUpdateRecordSetRequest: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateRecordSetRequest;

  try {
    const data = await api.updateRecordSet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name                          | Type                                                      | Description          | Notes                                |
| ----------------------------- | --------------------------------------------------------- | -------------------- | ------------------------------------ |
| **id**                        | `string`                                                  | Record Set ID (UUID) | [Defaults to `undefined`]            |
| **dtoUpdateRecordSetRequest** | [DtoUpdateRecordSetRequest](DtoUpdateRecordSetRequest.md) | Fields to update     |                                      |
| **xOrgID**                    | `string`                                                  | Organization UUID    | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoRecordSetResponse**](DtoRecordSetResponse.md)

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
