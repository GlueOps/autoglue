# DtoUpdateDomainRequest

## Properties

| Name            | Type   |
| --------------- | ------ |
| `credential_id` | string |
| `domain_name`   | string |
| `status`        | string |
| `zone_id`       | string |

## Example

```typescript
import type { DtoUpdateDomainRequest } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  credential_id: null,
  domain_name: null,
  status: null,
  zone_id: null,
} satisfies DtoUpdateDomainRequest;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoUpdateDomainRequest;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
