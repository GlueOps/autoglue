
# HandlersUserAPIKeyOut


## Properties

Name | Type
------------ | -------------
`created_at` | string
`expires_at` | string
`id` | string
`last_used_at` | string
`name` | string
`plain` | string
`scope` | string

## Example

```typescript
import type { HandlersUserAPIKeyOut } from '@glueops/autoglue-sdk'

// TODO: Update the object below with actual values
const example = {
  "created_at": null,
  "expires_at": null,
  "id": null,
  "last_used_at": null,
  "name": null,
  "plain": null,
  "scope": null,
} satisfies HandlersUserAPIKeyOut

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HandlersUserAPIKeyOut
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


