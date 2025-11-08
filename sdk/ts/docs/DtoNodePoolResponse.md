# DtoNodePoolResponse

## Properties

| Name              | Type                                                           |
| ----------------- | -------------------------------------------------------------- |
| `annotations`     | [Array&lt;DtoAnnotationResponse&gt;](DtoAnnotationResponse.md) |
| `created_at`      | string                                                         |
| `id`              | string                                                         |
| `labels`          | [Array&lt;DtoLabelResponse&gt;](DtoLabelResponse.md)           |
| `name`            | string                                                         |
| `organization_id` | string                                                         |
| `role`            | string                                                         |
| `servers`         | [Array&lt;DtoServerResponse&gt;](DtoServerResponse.md)         |
| `taints`          | [Array&lt;DtoTaintResponse&gt;](DtoTaintResponse.md)           |
| `updated_at`      | string                                                         |

## Example

```typescript
import type { DtoNodePoolResponse } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  annotations: null,
  created_at: null,
  id: null,
  labels: null,
  name: null,
  organization_id: null,
  role: null,
  servers: null,
  taints: null,
  updated_at: null,
} satisfies DtoNodePoolResponse;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoNodePoolResponse;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
