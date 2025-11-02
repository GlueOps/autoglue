package keys

func Decrypt(encKeyB64, enc string) ([]byte, error) {
	return decryptAESGCM(encKeyB64, enc)
}
