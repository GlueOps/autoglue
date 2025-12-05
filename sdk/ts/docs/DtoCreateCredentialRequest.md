# DtoCreateCredentialRequest

## Properties

| Name                  | Type   |
| --------------------- | ------ |
| `account_id`          | string |
| `credential_provider` | string |
| `kind`                | string |
| `name`                | string |
| `region`              | string |
| `schema_version`      | number |
| `scope`               | object |
| `scope_kind`          | string |
| `scope_version`       | number |
| `secret`              | object |

## Example

```typescript
import type { DtoCreateCredentialRequest } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  account_id: null,
  credential_provider: null,
  kind: null,
  name: null,
  region: null,
  schema_version: null,
  scope: null,
  scope_kind: null,
  scope_version: null,
  secret: null,
} satisfies DtoCreateCredentialRequest;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoCreateCredentialRequest;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
