# DtoCreateSSHRequest

## Properties

| Name      | Type   |
| --------- | ------ |
| `bits`    | number |
| `comment` | string |
| `name`    | string |
| `type`    | string |

## Example

```typescript
import type { DtoCreateSSHRequest } from '@glueops/autoglue-sdk'

// TODO: Update the object below with actual values
const example = {
  "bits": null,
  "comment": deploy@autoglue,
  "name": null,
  "type": null,
} satisfies DtoCreateSSHRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoCreateSSHRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
