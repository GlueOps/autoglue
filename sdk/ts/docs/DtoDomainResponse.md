# DtoDomainResponse

## Properties

| Name              | Type   |
| ----------------- | ------ |
| `created_at`      | string |
| `credential_id`   | string |
| `domain_name`     | string |
| `id`              | string |
| `last_error`      | string |
| `organization_id` | string |
| `status`          | string |
| `updated_at`      | string |
| `zone_id`         | string |

## Example

```typescript
import type { DtoDomainResponse } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  created_at: null,
  credential_id: null,
  domain_name: null,
  id: null,
  last_error: null,
  organization_id: null,
  status: null,
  updated_at: null,
  zone_id: null,
} satisfies DtoDomainResponse;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoDomainResponse;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
