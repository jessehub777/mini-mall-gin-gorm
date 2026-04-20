package password

import "golang.org/x/crypto/bcrypt"

// Hash 对明文密码进行哈希处理。
func Hash(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Verify 校验明文密码与哈希密码是否匹配。
func Verify(rawPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	return err == nil
}
