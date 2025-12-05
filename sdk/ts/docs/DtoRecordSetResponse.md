# DtoRecordSetResponse

## Properties

| Name          | Type   |
| ------------- | ------ |
| `created_at`  | string |
| `domain_id`   | string |
| `fingerprint` | string |
| `id`          | string |
| `last_error`  | string |
| `name`        | string |
| `owner`       | string |
| `status`      | string |
| `ttl`         | number |
| `type`        | string |
| `updated_at`  | string |
| `values`      | object |

## Example

```typescript
import type { DtoRecordSetResponse } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  created_at: null,
  domain_id: null,
  fingerprint: null,
  id: null,
  last_error: null,
  name: null,
  owner: null,
  status: null,
  ttl: null,
  type: null,
  updated_at: null,
  values: null,
} satisfies DtoRecordSetResponse;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoRecordSetResponse;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
