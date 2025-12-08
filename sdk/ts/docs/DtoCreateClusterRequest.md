# DtoCreateClusterRequest

## Properties

| Name               | Type   |
| ------------------ | ------ |
| `cluster_provider` | string |
| `docker_image`     | string |
| `docker_tag`       | string |
| `name`             | string |
| `region`           | string |

## Example

```typescript
import type { DtoCreateClusterRequest } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  cluster_provider: null,
  docker_image: null,
  docker_tag: null,
  name: null,
  region: null,
} satisfies DtoCreateClusterRequest;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoCreateClusterRequest;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
