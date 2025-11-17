import type React from "react"

declare global {
  namespace JSX {
    interface IntrinsicElements {
      "rapi-doc": React.DetailedHTMLProps<React.HTMLAttributes<HTMLElement>, HTMLElement> & {
        "spec-url"?: string
        "render-style"?: string
        theme?: string
        "show-header"?: string | boolean
        "persist-auth"?: string | boolean
        "allow-advanced-search"?: string | boolean
        "schema-description-expanded"?: string | boolean
        "allow-schema-description-expand-toggle"?: string | boolean
        "allow-spec-file-download"?: string | boolean
        "allow-spec-file-load"?: string | boolean
        "allow-spec-url-load"?: string | boolean
        "allow-try"?: string | boolean
        "schema-style"?: string
        "fetch-credentials"?: string
        "default-api-server"?: string
        "api-key-name"?: string
        "api-key-location"?: string
        "api-key-value"?: string
      }
    }
  }
}

export {}