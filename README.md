# Autoglue

## Before modifying anything
this repo uses git subtree
Subtree is like “vendor the files” with the option to push/pull history, but it behaves like normal files in the parent—no detached HEADs, no separate checkout.

- Main repo: github.com/glueops/autoglue
- SDK repo: github.com/glueops/autoglue-sdk-go
- SDK Path in main: sdk/go/
```bash
 # one-time: add the external repo as a subtree living at sdk/go/
git remote add sdk-origin git@github.com:glueops/autoglue-sdk-go.git
git subtree add --prefix=sdk/go sdk-origin main --squash
```

After changes in the API:
```bash
# Regenerate Swagger
make swagger

# Regenerate all SDKs - this includes the go and typescript SDKs, as well as the vendored TS SDK consumed by UI
make sdk-all

# update SDK repo from main (after regeneration)
git subtree push --prefix=sdk/go sdk-origin main
```