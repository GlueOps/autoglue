
# DtoJob


## Properties

Name | Type
------------ | -------------
`attempts` | number
`created_at` | string
`id` | string
`last_error` | string
`max_attempts` | number
`payload` | object
`queue` | string
`run_at` | string
`status` | [DtoJobStatus](DtoJobStatus.md)
`type` | string
`updated_at` | string

## Example

```typescript
import type { DtoJob } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "attempts": 0,
  "created_at": 2025-11-04T09:30:00Z,
  "id": 01HF7SZK8Z8WG1M3J7S2Z8M2N6,
  "last_error": error message,
  "max_attempts": 3,
  "payload": null,
  "queue": default,
  "run_at": 2025-11-04T09:30:00Z,
  "status": null,
  "type": email.send,
  "updated_at": 2025-11-04T09:30:00Z,
} satisfies DtoJob

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoJob
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


