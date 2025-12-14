package codeGenRepo

import (
	"globe-and-citizen/layer8/auth-server/pkg/code"
)

type ICodeGeneratorRepository interface {
	GenerateVerificationCode(salt string, userEmail string) (string, error)
}

type CodeGeneratorRepository struct {
	generator code.ICodeGenerator
}

func NewCodeGenerateRepository(generator code.ICodeGenerator) ICodeGeneratorRepository {
	return &CodeGeneratorRepository{}
}

func (cgr *CodeGeneratorRepository) GenerateVerificationCode(salt string, userEmail string) (string, error) {
	return cgr.generator.GenerateCode(salt, userEmail)
}
