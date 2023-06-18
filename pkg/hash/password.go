package hash

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	GenerateFromPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword string, password string) error
}

type HashPassword struct {
	cost int
}

func NewHash() *HashPassword {
	return &HashPassword{
		cost: bcrypt.DefaultCost,
	}
}
func (h *HashPassword) GenerateFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *HashPassword) CompareHashAndPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
