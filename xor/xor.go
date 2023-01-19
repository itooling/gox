package xor

func EncryptXor(src, key []byte) []byte {
	ml, kl := len(src), len(key)
	res := make([]byte, 0)
	for i := 0; i < ml; i++ {
		res = append(res, src[i]^key[i%kl])
	}
	return res
}

func DecryptXor(src, key []byte) []byte {
	ml, kl := len(src), len(key)
	res := make([]byte, 0)
	for i := 0; i < ml; i++ {
		res = append(res, src[i]^key[i%kl])
	}
	return res
}

func EncryptXorString(src, key string) string {
	ml, kl := len(src), len(key)
	var res string
	for i := 0; i < ml; i++ {
		res += string(src[i] ^ key[i%kl])
	}
	return res
}

func DecryptXorString(src, key string) string {
	ml, kl := len(src), len(key)
	var res string
	for i := 0; i < ml; i++ {
		res += string(src[i] ^ key[i%kl])
	}
	return res
}
