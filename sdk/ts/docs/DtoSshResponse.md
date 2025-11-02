# DtoSshResponse

## Properties

| Name              | Type   |
| ----------------- | ------ |
| `created_at`      | string |
| `fingerprint`     | string |
| `id`              | string |
| `name`            | string |
| `organization_id` | string |
| `public_key`      | string |
| `updated_at`      | string |

## Example

```typescript
import type { DtoSshResponse } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  created_at: null,
  fingerprint: null,
  id: null,
  name: null,
  organization_id: null,
  public_key: null,
  updated_at: null,
} satisfies DtoSshResponse;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoSshResponse;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
