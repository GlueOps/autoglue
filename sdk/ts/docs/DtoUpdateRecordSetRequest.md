# DtoUpdateRecordSetRequest

## Properties

| Name     | Type                |
| -------- | ------------------- |
| `name`   | string              |
| `status` | string              |
| `ttl`    | number              |
| `type`   | string              |
| `values` | Array&lt;string&gt; |

## Example

```typescript
import type { DtoUpdateRecordSetRequest } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  name: null,
  status: null,
  ttl: null,
  type: null,
  values: null,
} satisfies DtoUpdateRecordSetRequest;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoUpdateRecordSetRequest;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
