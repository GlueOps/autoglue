# DtoCreateLoadBalancerRequest

## Properties

| Name                 | Type   |
| -------------------- | ------ |
| `kind`               | string |
| `name`               | string |
| `private_ip_address` | string |
| `public_ip_address`  | string |

## Example

```typescript
import type { DtoCreateLoadBalancerRequest } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "kind": public,
  "name": glueops,
  "private_ip_address": 192.168.0.2,
  "public_ip_address": 8.8.8.8,
} satisfies DtoCreateLoadBalancerRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoCreateLoadBalancerRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
