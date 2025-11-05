# DtoQueueInfo

## Properties

| Name        | Type   |
| ----------- | ------ |
| `failed`    | number |
| `name`      | string |
| `pending`   | number |
| `running`   | number |
| `scheduled` | number |

## Example

```typescript
import type { DtoQueueInfo } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "failed": 5,
  "name": default,
  "pending": 42,
  "running": 3,
  "scheduled": 7,
} satisfies DtoQueueInfo

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoQueueInfo
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
