package auth

import "testing"

func TestCheckPassword(t *testing.T) {
    password := "mySecret123"

    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("failed to hash password: %v", err)
    }

    t.Run("correct password should be true", func(t *testing.T) {
        result, _ := CheckPasswordHash(password, hash)
        if result != true {
            t.Errorf("expected true, got false")
        }
    })

    t.Run("wrong password should be false", func(t *testing.T) {
        result, _ := CheckPasswordHash("wrongPassword", hash)
        if result != false {
            t.Errorf("expected false, got true")
        }
    })
}

