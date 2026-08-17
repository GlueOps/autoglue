package auth

import (
	"crypto/subtle"
	"fmt"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AgentKeyPrefix follows the "org_" convention already used for API keys: the
// visible half of a credential says what kind of principal it belongs to, so a
// value spotted in a log or a config file is identifiable without a lookup.
const AgentKeyPrefix = "agt_"

// agentPrefixLen is the leading slice of the key kept in cleartext for display.
// Twelve characters, same as the org keys, which is enough to tell two
// credentials apart in a UI and far short of guessable.
const agentPrefixLen = 12

// AgentCredential is the tuple minted at enrolment: id, key, secret.
//
// Three parts rather than the usual two so that no part of it is derivable from
// the others — the id names the row, the key is the indexed lookup, and the
// secret is the only thing that proves possession. The plaintext Key and Secret
// exist only to be returned to the enrolling bastion once; only the hashes are
// ever persisted.
type AgentCredential struct {
	ID     uuid.UUID
	Key    string
	Secret string

	KeyHash    string
	SecretHash string
	Prefix     string
}

// MintAgentCredential generates a fresh agent credential tuple.
//
// Entropy and hashing follow findOrCreateClusterAutomationKey exactly — 128
// bits for the key, 256 for the secret, SHA256 for the key digest and argon2id
// for the secret — because the two credentials protect the same thing and a
// weaker agent credential would simply become the easier way in.
//
// The id is minted here rather than left to the column default so the caller
// holds the whole tuple before it writes anything: the enrolment response has
// to state the id, and reading it back after an insert is one more place for
// the transaction to be half-done.
func MintAgentCredential() (AgentCredential, error) {
	keySuffix, err := RandomB64URL(16)
	if err != nil {
		return AgentCredential{}, fmt.Errorf("entropy_error: %w", err)
	}
	secret, err := RandomB64URL(32)
	if err != nil {
		return AgentCredential{}, fmt.Errorf("entropy_error: %w", err)
	}

	key := AgentKeyPrefix + keySuffix

	secretHash, err := HashSecretArgon2id(secret)
	if err != nil {
		return AgentCredential{}, fmt.Errorf("hash_error: %w", err)
	}

	prefix := key
	if len(prefix) > agentPrefixLen {
		prefix = prefix[:agentPrefixLen]
	}

	return AgentCredential{
		ID:         uuid.New(),
		Key:        key,
		Secret:     secret,
		KeyHash:    SHA256Hex(key),
		SecretHash: secretHash,
		Prefix:     prefix,
	}, nil
}

// ValidateAgentKeyPair validates the agent credential tuple carried in
// X-Agent-ID / X-Agent-KEY / X-Agent-SECRET.
//
// Every failure returns nil and none of them are distinguishable, so a caller
// holding half a credential cannot use the response to learn which half is
// wrong.
//
// Unlike ValidateOrgKeyPair this predicates on status as well as expiry. That
// validator gets away without a revocation check only because a sweeper deletes
// the keys it issues; an agent row is never deleted — re-enrolment revokes it,
// precisely so its tasks outlive the credential that produced them — so without
// the status predicate a revoked bastion would keep authenticating forever.
func ValidateAgentKeyPair(agentID, agentKey, secret string, db *gorm.DB) *models.Agent {
	if agentID == "" || agentKey == "" || secret == "" {
		return nil
	}
	id, err := uuid.Parse(agentID)
	if err != nil {
		return nil
	}

	// SHA256 on the key gives a plain indexed equality, which is safe precisely
	// because the key is not the secret.
	var a models.Agent
	if err := db.Where(
		"key_hash = ? AND status = ? AND (expires_at IS NULL OR expires_at > ?)",
		SHA256Hex(agentKey), models.AgentStatusActive, time.Now(),
	).First(&a).Error; err != nil {
		return nil
	}

	// Binding the id is what makes this a tuple rather than a pair: the caller
	// must present the whole thing it was issued, so a stale id left over from a
	// previous enrolment fails closed instead of quietly authenticating as
	// whichever row the key happens to name.
	if subtle.ConstantTimeCompare(id[:], a.ID[:]) != 1 {
		return nil
	}

	// argon2id, compared in constant time by ComparePasswordAndHash. Never
	// compare the encoded hashes with ==: that is a length-prefixed string
	// compare and it leaks a prefix at a time.
	if ok, _ := VerifySecretArgon2id(a.SecretHash, secret); !ok {
		return nil
	}

	return &a
}
