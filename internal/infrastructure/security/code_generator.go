package security

import (
	"math/rand"
	"strconv"
)

type RandomCodeGenerator struct{}

func NewRandomCodeGenerator() *RandomCodeGenerator {
	return &RandomCodeGenerator{}
}

func (g *RandomCodeGenerator) GenerateRecoveryCode() string {
	number := rand.Intn(900000) + 100000
	return strconv.Itoa(number)
}
