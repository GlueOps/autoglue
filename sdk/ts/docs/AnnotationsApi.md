# AnnotationsApi

All URIs are relative to _http://localhost:8080/api/v1_

| Method                                                     | HTTP request                 | Description                       |
| ---------------------------------------------------------- | ---------------------------- | --------------------------------- |
| [**createAnnotation**](AnnotationsApi.md#createannotation) | **POST** /annotations        | Create annotation (org scoped)    |
| [**deleteAnnotation**](AnnotationsApi.md#deleteannotation) | **DELETE** /annotations/{id} | Delete annotation (org scoped)    |
| [**getAnnotation**](AnnotationsApi.md#getannotation)       | **GET** /annotations/{id}    | Get annotation by ID (org scoped) |
| [**listAnnotations**](AnnotationsApi.md#listannotations)   | **GET** /annotations         | List annotations (org scoped)     |
| [**updateAnnotation**](AnnotationsApi.md#updateannotation) | **PATCH** /annotations/{id}  | Update annotation (org scoped)    |

## createAnnotation

> DtoAnnotationResponse createAnnotation(body, xOrgID)

Create annotation (org scoped)

Creates an annotation.

### Example

```ts
import {
  Configuration,
  AnnotationsApi,
} from '@glueops/autoglue-sdk-go';
import type { CreateAnnotationRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new AnnotationsApi(config);

  const body = {
    // DtoCreateAnnotationRequest | Annotation payload
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies CreateAnnotationRequest;

  try {
    const data = await api.createAnnotation(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type                                                        | Description        | Notes                                |
| ---------- | ----------------------------------------------------------- | ------------------ | ------------------------------------ |
| **body**   | [DtoCreateAnnotationRequest](DtoCreateAnnotationRequest.md) | Annotation payload |                                      |
| **xOrgID** | `string`                                                    | Organization UUID  | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoAnnotationResponse**](DtoAnnotationResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                   | Response headers |
| ----------- | ----------------------------- | ---------------- |
| **201**     | Created                       | -                |
| **400**     | invalid json / missing fields | -                |
| **401**     | Unauthorized                  | -                |
| **403**     | organization required         | -                |
| **500**     | create failed                 | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## deleteAnnotation

> string deleteAnnotation(id, xOrgID)

Delete annotation (org scoped)

Permanently deletes the annotation.

### Example

```ts
import { Configuration, AnnotationsApi } from "@glueops/autoglue-sdk-go";
import type { DeleteAnnotationRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new AnnotationsApi(config);

  const body = {
    // string | Annotation ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies DeleteAnnotationRequest;

  try {
    const data = await api.deleteAnnotation(body);
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
| **id**     | `string` | Annotation ID (UUID) | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID    | [Optional] [Defaults to `undefined`] |

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

## getAnnotation

> DtoAnnotationResponse getAnnotation(id, xOrgID)

Get annotation by ID (org scoped)

Returns one annotation. Add &#x60;include&#x3D;node_pools&#x60; to include node pools.

### Example

```ts
import { Configuration, AnnotationsApi } from "@glueops/autoglue-sdk-go";
import type { GetAnnotationRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new AnnotationsApi(config);

  const body = {
    // string | Annotation ID (UUID)
    id: id_example,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies GetAnnotationRequest;

  try {
    const data = await api.getAnnotation(body);
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
| **id**     | `string` | Annotation ID (UUID) | [Defaults to `undefined`]            |
| **xOrgID** | `string` | Organization UUID    | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoAnnotationResponse**](DtoAnnotationResponse.md)

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

## listAnnotations

> Array&lt;DtoAnnotationResponse&gt; listAnnotations(xOrgID, key, value, q)

List annotations (org scoped)

Returns annotations for the organization in X-Org-ID. Filters: &#x60;key&#x60;, &#x60;value&#x60;, and &#x60;q&#x60; (key contains). Add &#x60;include&#x3D;node_pools&#x60; to include linked node pools.

### Example

```ts
import { Configuration, AnnotationsApi } from "@glueops/autoglue-sdk-go";
import type { ListAnnotationsRequest } from "@glueops/autoglue-sdk-go";

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
  const api = new AnnotationsApi(config);

  const body = {
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
    // string | Exact key (optional)
    key: key_example,
    // string | Exact value (optional)
    value: value_example,
    // string | key contains (case-insensitive) (optional)
    q: q_example,
  } satisfies ListAnnotationsRequest;

  try {
    const data = await api.listAnnotations(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type     | Description                     | Notes                                |
| ---------- | -------- | ------------------------------- | ------------------------------------ |
| **xOrgID** | `string` | Organization UUID               | [Optional] [Defaults to `undefined`] |
| **key**    | `string` | Exact key                       | [Optional] [Defaults to `undefined`] |
| **value**  | `string` | Exact value                     | [Optional] [Defaults to `undefined`] |
| **q**      | `string` | key contains (case-insensitive) | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;DtoAnnotationResponse&gt;**](DtoAnnotationResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`

### HTTP response details

| Status code | Description                | Response headers |
| ----------- | -------------------------- | ---------------- |
| **200**     | OK                         | -                |
| **401**     | Unauthorized               | -                |
| **403**     | organization required      | -                |
| **500**     | failed to list annotations | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

## updateAnnotation

> DtoAnnotationResponse updateAnnotation(id, body, xOrgID)

Update annotation (org scoped)

Partially update annotation fields.

### Example

```ts
import {
  Configuration,
  AnnotationsApi,
} from '@glueops/autoglue-sdk-go';
import type { UpdateAnnotationRequest } from '@glueops/autoglue-sdk-go';

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
  const api = new AnnotationsApi(config);

  const body = {
    // string | Annotation ID (UUID)
    id: id_example,
    // DtoUpdateAnnotationRequest | Fields to update
    body: ...,
    // string | Organization UUID (optional)
    xOrgID: xOrgID_example,
  } satisfies UpdateAnnotationRequest;

  try {
    const data = await api.updateAnnotation(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

| Name       | Type                                                        | Description          | Notes                                |
| ---------- | ----------------------------------------------------------- | -------------------- | ------------------------------------ |
| **id**     | `string`                                                    | Annotation ID (UUID) | [Defaults to `undefined`]            |
| **body**   | [DtoUpdateAnnotationRequest](DtoUpdateAnnotationRequest.md) | Fields to update     |                                      |
| **xOrgID** | `string`                                                    | Organization UUID    | [Optional] [Defaults to `undefined`] |

### Return type

[**DtoAnnotationResponse**](DtoAnnotationResponse.md)

### Authorization

[OrgKeyAuth](../README.md#OrgKeyAuth), [OrgSecretAuth](../README.md#OrgSecretAuth), [BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`

### HTTP response details

| Status code | Description               | Response headers |
| ----------- | ------------------------- | ---------------- |
| **200**     | OK                        | -                |
| **400**     | invalid id / invalid json | -                |
| **401**     | Unauthorized              | -                |
| **403**     | organization required     | -                |
| **404**     | not found                 | -                |
| **500**     | update failed             | -                |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
