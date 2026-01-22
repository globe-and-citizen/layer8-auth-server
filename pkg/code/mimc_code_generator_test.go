package code

import (
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

const salt = "ajdjsjsaafktyowqqrtgpowrkdkdkfak"

func TestGenerateCode_EmailWithASCIICharacters(t *testing.T) {
	email := "myemail@gmail.com"

	generator := NewMIMCCodeGenerator()
	code, err := generator.GenerateCode(salt, email)

	assert.Nil(t, err)
	assert.True(t, len(code) == utils.VerificationCodeSize)
}

func TestGenerateCode_EmailWithChineseCharacters(t *testing.T) {
	email := "用户@例子.广告"

	generator := NewMIMCCodeGenerator()
	code, err := generator.GenerateCode(salt, email)

	assert.Nil(t, err)
	assert.True(t, len(code) == utils.VerificationCodeSize)
}

func TestGenerateCode_EmailWithGermanCharacters(t *testing.T) {
	email := "Dörte@Sörensen.example.com"

	generator := NewMIMCCodeGenerator()
	code, err := generator.GenerateCode(salt, email)

	assert.Nil(t, err)
	assert.True(t, len(code) == utils.VerificationCodeSize)
}
