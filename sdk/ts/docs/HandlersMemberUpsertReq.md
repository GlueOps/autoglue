# HandlersMemberUpsertReq

## Properties

| Name      | Type   |
| --------- | ------ |
| `role`    | string |
| `user_id` | string |

## Example

```typescript
import type { HandlersMemberUpsertReq } from "@glueops/autoglue-sdk-go";

// TODO: Update the object below with actual values
const example = {
  role: member,
  user_id: null,
} satisfies HandlersMemberUpsertReq;

console.log(example);

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example);
console.log(exampleJSON);

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HandlersMemberUpsertReq;
console.log(exampleParsed);
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)
