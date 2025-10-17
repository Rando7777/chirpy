package auth

import(
	"github.com/alexedwards/argon2id"	
)

func HashPassword(pw string) (string, error) {
	hash, err := argon2id.CreateHash(pw, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(pw string, hash string) (bool, error){
	match, err := argon2id.ComparePasswordAndHash(pw, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
