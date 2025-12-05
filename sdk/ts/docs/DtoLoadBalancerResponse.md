# DtoLoadBalancerResponse

## Properties

| Name                 | Type   |
| -------------------- | ------ |
| `created_at`         | string |
| `id`                 | string |
| `kind`               | string |
| `name`               | string |
| `organization_id`    | string |
| `private_ip_address` | string |
| `public_ip_address`  | string |
| `updated_at`         | string |

## Example

```typescript
import type { DtoLoadBalancerResponse } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  created_at: null,
  id: null,
  kind: null,
  name: null,
  organization_id: null,
  private_ip_address: null,
  public_ip_address: null,
  updated_at: null,
} satisfies DtoLoadBalancerResponse;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoLoadBalancerResponse;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
