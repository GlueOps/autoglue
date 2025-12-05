# HandlersVersionResponse

## Properties

| Name         | Type    |
| ------------ | ------- |
| `built`      | string  |
| `builtBy`    | string  |
| `commit`     | string  |
| `commitTime` | string  |
| `go`         | string  |
| `goArch`     | string  |
| `goOS`       | string  |
| `modified`   | boolean |
| `revision`   | string  |
| `vcs`        | string  |
| `version`    | string  |

## Example

```typescript
import type { HandlersVersionResponse } from '@glueops/autoglue-sdk-go'

// TODO: Update the object below with actual values
const example = {
  "built": 2025-11-08T12:34:56Z,
  "builtBy": ci,
  "commit": a1b2c3d,
  "commitTime": 2025-11-08T12:31:00Z,
  "go": go1.23.3,
  "goArch": amd64,
  "goOS": linux,
  "modified": false,
  "revision": a1b2c3d4e5f6abcdef,
  "vcs": git,
  "version": 1.4.2,
} satisfies HandlersVersionResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HandlersVersionResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
