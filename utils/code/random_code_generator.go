package code

import (
	"math/rand"
	"strconv"
	"strings"
)

type RandomCodeGenerator struct {
	verificationCodeSize int
}

func NewRandomCodeGenerator(verificationCodeSize int) *RandomCodeGenerator {
	g := new(RandomCodeGenerator)
	g.verificationCodeSize = verificationCodeSize
	return g
}

func (g *RandomCodeGenerator) GenerateCode(salt string, emailAddress string) (string, error) {
	verificationCode := make([]string, g.verificationCodeSize)
	for i := 0; i < g.verificationCodeSize; i++ {
		verificationCode[i] = strconv.Itoa(rand.Intn(10))
	}
	return strings.Join(verificationCode, ""), nil
}
