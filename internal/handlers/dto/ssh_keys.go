package dto

import (
	"github.com/glueops/autoglue/internal/common"
)

type CreateSSHRequest struct {
	Name    string  `json:"name"`
	Comment string  `json:"comment,omitempty" example:"deploy@autoglue"`
	Bits    *int    `json:"bits,omitempty"` // Only for RSA
	Type    *string `json:"type,omitempty"` // "rsa" (default) or "ed25519"
}

type SshResponse struct {
	common.AuditFields
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
	// ServerCount is how many servers reference this key. Zero means nothing in
	// autoglue uses it, which is not the same as unused: the key may be in an
	// authorized_keys on a host autoglue does not track. Absent means the count
	// was not computed for this response, so treat absent and zero differently.
	//
	// A pointer for exactly that reason. Reveal and download describe key
	// material rather than attachment and leave it nil; a plain int would
	// render that as 0 and badge every revealed key as unattached.
	ServerCount         *int   `json:"server_count,omitempty"`
	EncryptedPrivateKey string `json:"-"`
	PrivateIV           string `json:"-"`
	PrivateTag          string `json:"-"`
}

type SshRevealResponse struct {
	SshResponse
	PrivateKey string `json:"private_key"`
}

type SshMaterialJSON struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	// Exactly one of the following will be populated for part=public/private.
	PublicKey  *string `json:"public_key,omitempty"`  // OpenSSH authorized_key (string)
	PrivatePEM *string `json:"private_pem,omitempty"` // PKCS#1/PEM (string)
	// For part=both with mode=json we'll return a base64 zip
	ZipBase64 *string `json:"zip_base64,omitempty"` // base64-encoded zip
	// Suggested filenames (SDKs can save to disk without inferring names)
	Filenames []string `json:"filenames"`
}
