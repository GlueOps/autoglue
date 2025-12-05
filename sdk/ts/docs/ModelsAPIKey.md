# ModelsAPIKey

## Properties

| Name           | Type    |
| -------------- | ------- |
| `created_at`   | Date    |
| `expires_at`   | Date    |
| `id`           | string  |
| `last_used_at` | Date    |
| `name`         | string  |
| `org_id`       | string  |
| `prefix`       | string  |
| `revoked`      | boolean |
| `scope`        | string  |
| `updated_at`   | Date    |
| `user_id`      | string  |

## Example

```typescript
import type { ModelsAPIKey } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  created_at: null,
  expires_at: null,
  id: null,
  last_used_at: null,
  name: null,
  org_id: null,
  prefix: null,
  revoked: null,
  scope: null,
  updated_at: null,
  user_id: null,
} satisfies ModelsAPIKey;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsAPIKey;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
