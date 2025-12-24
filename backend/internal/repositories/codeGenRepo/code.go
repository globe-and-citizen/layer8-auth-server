package codeGenRepo

import (
	"globe-and-citizen/layer8/auth-server/backend/pkg/code"
)

type ICodeGeneratorRepository interface {
	GenerateEmailVerificationCode(salt string, userEmail string) (string, error)
	GeneratePhoneVerificationCode(salt string, phoneNumber string) (string, error)
}

type CodeGeneratorRepository struct {
	generator code.ICodeGenerator
}

func NewCodeGenerateRepository(generator code.ICodeGenerator) ICodeGeneratorRepository {
	return &CodeGeneratorRepository{}
}

func (cgr *CodeGeneratorRepository) GenerateEmailVerificationCode(salt string, userEmail string) (string, error) {
	return cgr.generator.GenerateCode(salt, userEmail)
}

func (cgr *CodeGeneratorRepository) GeneratePhoneVerificationCode(salt string, phoneNumber string) (string, error) {
	return cgr.generator.GenerateCode(salt, phoneNumber)
}
