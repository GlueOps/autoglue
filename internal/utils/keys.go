package utils

import (
	"errors"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/spf13/viper"
)

func getMasterKey() ([]byte, error) {
	return decodeB64(viper.GetString("master_key"))
}

func getOrCreateTenantKey(orgID string) ([]byte, error) {
	var orgKey models.OrganizationKey
	db.DB.Where("organization_id = ?", orgID).First(&orgKey)

	if orgKey.ID != 0 {
		encryptedKey, _ := decodeB64(orgKey.EncryptedKey)
		iv, _ := decodeB64(orgKey.IV)
		tag, _ := decodeB64(orgKey.Tag)
		masterKey, _ := getMasterKey()
		return decryptAESGCM(encryptedKey, masterKey, iv, tag)
	}

	tenantKey, _ := randomBytes(32)
	masterKey, _ := getMasterKey()
	encrypted, iv, tag, _ := encryptAESGCM(tenantKey, masterKey)

	orgKey = models.OrganizationKey{
		OrganizationID: orgID,
		EncryptedKey:   encodeB64(encrypted),
		IV:             encodeB64(iv),
		Tag:            encodeB64(tag),
	}
	db.DB.Create(&orgKey)

	return tenantKey, nil
}

func encryptWithTenantKey(plainText, orgID, provider string) error {
	tenantKey, err := getOrCreateTenantKey(orgID)
	if err != nil {
		return err
	}

	ciphertext, iv, tag, err := encryptAESGCM([]byte(plainText), tenantKey)
	if err != nil {
		return err
	}

	var cred models.Credential
	db.DB.Where("organization_id = ? AND provider = ?", orgID, provider).First(&cred)
	cred.OrganizationID = orgID
	cred.Provider = provider
	cred.EncryptedData = encodeB64(ciphertext)
	cred.IV = encodeB64(iv)
	cred.Tag = encodeB64(tag)

	if cred.ID != 0 {
		db.DB.Save(&cred)
	} else {
		db.DB.Create(&cred)
	}
	return nil
}

func decryptCredentials(orgID, provider string) (string, error) {
	var cred models.Credential
	if err := db.DB.Where("organization_id = ? AND provider = ?", orgID, provider).First(&cred).Error; err != nil {
		return "", errors.New("credential not found")
	}

	tenantKey, err := getOrCreateTenantKey(orgID)
	if err != nil {
		return "", err
	}

	ciphertext, _ := decodeB64(cred.EncryptedData)
	iv, _ := decodeB64(cred.IV)
	tag, _ := decodeB64(cred.Tag)

	plain, err := decryptAESGCM(ciphertext, tenantKey, iv, tag)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
