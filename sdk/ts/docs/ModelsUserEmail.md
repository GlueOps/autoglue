# ModelsUserEmail

## Properties

| Name          | Type                        |
| ------------- | --------------------------- |
| `created_at`  | Date                        |
| `email`       | string                      |
| `id`          | string                      |
| `is_primary`  | boolean                     |
| `is_verified` | boolean                     |
| `updated_at`  | Date                        |
| `user`        | [ModelsUser](ModelsUser.md) |
| `user_id`     | string                      |

## Example

```typescript
import type { ModelsUserEmail } from "@glueops/autoglue-sdk";

// TODO: Update the object below with actual values
const example = {
  created_at: null,
  email: null,
  id: null,
  is_primary: null,
  is_verified: null,
  updated_at: null,
  user: null,
  user_id: null,
} satisfies ModelsUserEmail;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsUserEmail;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
