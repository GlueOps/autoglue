# DtoAttachLoadBalancerRequest

## Properties

| Name               | Type   |
| ------------------ | ------ |
| `load_balancer_id` | string |

## Example

```typescript
import type { DtoAttachLoadBalancerRequest } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  load_balancer_id: null,
} satisfies DtoAttachLoadBalancerRequest;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoAttachLoadBalancerRequest;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
