
# DtoPageJob


## Properties

Name | Type
------------ | -------------
`items` | [Array&lt;DtoJob&gt;](DtoJob.md)
`page` | number
`page_size` | number
`total` | number

## Example

```typescript
import type { DtoPageJob } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "items": null,
  "page": null,
  "page_size": null,
  "total": null,
} satisfies DtoPageJob

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DtoPageJob
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


