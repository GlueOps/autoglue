# ArcherAdminApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**adminCancelArcherJob**](ArcherAdminApi.md#admincancelarcherjob) | **POST** /admin/archer/jobs/{id}/cancel | Cancel an Archer job (admin) |
| [**adminEnqueueArcherJob**](ArcherAdminApi.md#adminenqueuearcherjob) | **POST** /admin/archer/jobs | Enqueue a new Archer job (admin) |
| [**adminListArcherJobs**](ArcherAdminApi.md#adminlistarcherjobs) | **GET** /admin/archer/jobs | List Archer jobs (admin) |
| [**adminListArcherQueues**](ArcherAdminApi.md#adminlistarcherqueues) | **GET** /admin/archer/queues | List Archer queues (admin) |
| [**adminRetryArcherJob**](ArcherAdminApi.md#adminretryarcherjob) | **POST** /admin/archer/jobs/{id}/retry | Retry a failed/canceled Archer job (admin) |



## adminCancelArcherJob

> DtoJob adminCancelArcherJob(id)

Cancel an Archer job (admin)

Set job status to canceled if cancellable. For running jobs, this only affects future picks; wire to Archer if you need active kill.

### Example

```ts
import {
  Configuration,
  ArcherAdminApi,
} from '@glueops/autoglue-sdk-go';
import type { AdminCancelArcherJobRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new ArcherAdminApi(config);

  const body = {
    // string | Job ID
    id: id_example,
  } satisfies AdminCancelArcherJobRequest;

  try {
    const data = await api.adminCancelArcherJob(body);
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
| **id** | `string` | Job ID | [Defaults to `undefined`] |

### Return type

[**DtoJob**](DtoJob.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | invalid job or not cancellable |  -  |
| **401** | Unauthorized |  -  |
| **403** | forbidden |  -  |
| **404** | not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## adminEnqueueArcherJob

> DtoJob adminEnqueueArcherJob(body)

Enqueue a new Archer job (admin)

Create a job immediately or schedule it for the future via &#x60;run_at&#x60;.

### Example

```ts
import {
  Configuration,
  ArcherAdminApi,
} from '@glueops/autoglue-sdk-go';
import type { AdminEnqueueArcherJobRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new ArcherAdminApi(config);

  const body = {
    // object | Job parameters
    body: Object,
  } satisfies AdminEnqueueArcherJobRequest;

  try {
    const data = await api.adminEnqueueArcherJob(body);
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
| **body** | `object` | Job parameters | |

### Return type

[**DtoJob**](DtoJob.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | invalid json or missing fields |  -  |
| **401** | Unauthorized |  -  |
| **403** | forbidden |  -  |
| **500** | internal error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## adminListArcherJobs

> DtoPageJob adminListArcherJobs(status, queue, q, page, pageSize)

List Archer jobs (admin)

Paginated background jobs with optional filters. Search &#x60;q&#x60; may match id, type, error, payload (implementation-dependent).

### Example

```ts
import {
  Configuration,
  ArcherAdminApi,
} from '@glueops/autoglue-sdk-go';
import type { AdminListArcherJobsRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new ArcherAdminApi(config);

  const body = {
    // 'queued' | 'running' | 'succeeded' | 'failed' | 'canceled' | 'retrying' | 'scheduled' | Filter by status (optional)
    status: status_example,
    // string | Filter by queue name / worker name (optional)
    queue: queue_example,
    // string | Free-text search (optional)
    q: q_example,
    // number | Page number (optional)
    page: 56,
    // number | Items per page (optional)
    pageSize: 56,
  } satisfies AdminListArcherJobsRequest;

  try {
    const data = await api.adminListArcherJobs(body);
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
| **status** | `queued`, `running`, `succeeded`, `failed`, `canceled`, `retrying`, `scheduled` | Filter by status | [Optional] [Defaults to `undefined`] [Enum: queued, running, succeeded, failed, canceled, retrying, scheduled] |
| **queue** | `string` | Filter by queue name / worker name | [Optional] [Defaults to `undefined`] |
| **q** | `string` | Free-text search | [Optional] [Defaults to `undefined`] |
| **page** | `number` | Page number | [Optional] [Defaults to `1`] |
| **pageSize** | `number` | Items per page | [Optional] [Defaults to `25`] |

### Return type

[**DtoPageJob**](DtoPageJob.md)

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
| **403** | forbidden |  -  |
| **500** | internal error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## adminListArcherQueues

> Array&lt;DtoQueueInfo&gt; adminListArcherQueues()

List Archer queues (admin)

Summary metrics per queue (pending, running, failed, scheduled).

### Example

```ts
import {
  Configuration,
  ArcherAdminApi,
} from '@glueops/autoglue-sdk-go';
import type { AdminListArcherQueuesRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new ArcherAdminApi(config);

  try {
    const data = await api.adminListArcherQueues();
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

[**Array&lt;DtoQueueInfo&gt;**](DtoQueueInfo.md)

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
| **403** | forbidden |  -  |
| **500** | internal error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## adminRetryArcherJob

> DtoJob adminRetryArcherJob(id)

Retry a failed/canceled Archer job (admin)

Marks the job retriable (DB flip). Swap this for an Archer admin call if you expose one.

### Example

```ts
import {
  Configuration,
  ArcherAdminApi,
} from '@glueops/autoglue-sdk-go';
import type { AdminRetryArcherJobRequest } from '@glueops/autoglue-sdk-go';

async function example() {
  console.log("🚀 Testing @glueops/autoglue-sdk-go SDK...");
  const config = new Configuration({ 
    // To configure API key authorization: BearerAuth
    apiKey: "YOUR API KEY",
  });
  const api = new ArcherAdminApi(config);

  const body = {
    // string | Job ID
    id: id_example,
  } satisfies AdminRetryArcherJobRequest;

  try {
    const data = await api.adminRetryArcherJob(body);
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
| **id** | `string` | Job ID | [Defaults to `undefined`] |

### Return type

[**DtoJob**](DtoJob.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **400** | invalid job or not eligible |  -  |
| **401** | Unauthorized |  -  |
| **403** | forbidden |  -  |
| **404** | not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

