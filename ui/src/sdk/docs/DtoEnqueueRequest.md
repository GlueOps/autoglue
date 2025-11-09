
# DtoEnqueueRequest


## Properties

Name | Type
------------ | -------------
`payload` | object
`queue` | string
`run_at` | string
`type` | string

## Example

```typescript
import type { DtoEnqueueRequest } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "payload": null,
  "queue": default,
  "run_at": 2025-11-05T08:00:00Z,
  "type": email.send,
} satisfies DtoEnqueueRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoEnqueueRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


