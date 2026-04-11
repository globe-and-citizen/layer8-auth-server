package codeGenRepo

import (
	"globe-and-citizen/layer8/auth-server/pkg/code"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
)

// fixme: Email and Phone verification codes are currently identical and persistent for each user!!!
type ICodeGeneratorRepository interface {
	GenerateEmailVerificationCode(salt string, userEmail string) (string, error)
	GeneratePhoneVerificationCode(salt string, phoneNumber string) (string, error)
}

type CodeGeneratorRepository struct {
	generator code.ICodeGenerator
}

func NewCodeGenerateRepository(generator code.ICodeGenerator) ICodeGeneratorRepository {
	return &CodeGeneratorRepository{
		generator: generator,
	}
}

func (cgr *CodeGeneratorRepository) GenerateEmailVerificationCode(salt string, userEmail string) (string, error) {
	verCode, err := cgr.generator.GenerateCode(salt, userEmail)
	return verCode, utils.StackError(err)
}

func (cgr *CodeGeneratorRepository) GeneratePhoneVerificationCode(salt string, phoneNumber string) (string, error) {
	verCode, err := cgr.generator.GenerateCode(salt, phoneNumber)
	return verCode, utils.StackError(err)
}
