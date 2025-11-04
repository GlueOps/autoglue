# DtoJob

## Properties

| Name           | Type                            |
| -------------- | ------------------------------- |
| `attempts`     | number                          |
| `created_at`   | string                          |
| `id`           | string                          |
| `last_error`   | string                          |
| `max_attempts` | number                          |
| `payload`      | object                          |
| `queue`        | string                          |
| `run_at`       | string                          |
| `status`       | [DtoJobStatus](DtoJobStatus.md) |
| `type`         | string                          |
| `updated_at`   | string                          |

## Example

```typescript
import type { DtoJob } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  attempts: null,
  created_at: null,
  id: null,
  last_error: null,
  max_attempts: null,
  payload: null,
  queue: null,
  run_at: null,
  status: null,
  type: null,
  updated_at: null,
} satisfies DtoJob;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoJob;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
