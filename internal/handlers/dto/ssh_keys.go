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
	Name                string `json:"name"`
	PublicKey           string `json:"public_key"`
	Fingerprint         string `json:"fingerprint"`
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
