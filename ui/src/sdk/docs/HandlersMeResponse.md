
# HandlersMeResponse


## Properties

Name | Type
------------ | -------------
`avatar_url` | string
`created_at` | Date
`display_name` | string
`emails` | [Array&lt;ModelsUserEmail&gt;](ModelsUserEmail.md)
`id` | string
`is_disabled` | boolean
`organizations` | [Array&lt;ModelsOrganization&gt;](ModelsOrganization.md)
`primary_email` | string
`updated_at` | Date

## Example

```typescript
import type { HandlersMeResponse } from '@glueops/autoglue-sdk'

// TODO: Update the object below with actual values
const example = {
  "avatar_url": null,
  "created_at": null,
  "display_name": null,
  "emails": null,
  "id": null,
  "is_disabled": null,
  "organizations": null,
  "primary_email": null,
  "updated_at": null,
} satisfies HandlersMeResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HandlersMeResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


