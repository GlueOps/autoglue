
# DtoAnnotationResponse


## Properties

Name | Type
------------ | -------------
`created_at` | string
`id` | string
`key` | string
`organization_id` | string
`updated_at` | string
`value` | string

## Example

```typescript
import type { DtoAnnotationResponse } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "created_at": null,
  "id": null,
  "key": null,
  "organization_id": null,
  "updated_at": null,
  "value": null,
} satisfies DtoAnnotationResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoAnnotationResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


