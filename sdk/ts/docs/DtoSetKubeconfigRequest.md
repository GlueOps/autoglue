# DtoSetKubeconfigRequest

## Properties

| Name         | Type   |
| ------------ | ------ |
| `kubeconfig` | string |

## Example

```typescript
import type { DtoSetKubeconfigRequest } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  kubeconfig: null,
} satisfies DtoSetKubeconfigRequest;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoSetKubeconfigRequest;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
