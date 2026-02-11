package zkRepo

import (
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"globe-and-citizen/layer8/auth-server/pkg/zk"
)

type IZkRepository interface {
	GenerateProof(salt string, target string, verificationCode string) ([]byte, uint, error)
	VerifyProof(verificationCode string, salt string, proofBytes []byte) error
}

type ZkRepository struct {
	proofProcessor zk.IProofProcessor
}

func (z ZkRepository) GenerateProof(salt string, target string, verificationCode string) ([]byte, uint, error) {
	proof, id, err := z.proofProcessor.GenerateProof(target, salt, verificationCode)
	return proof, id, utils.StackError(err)
}

func (z ZkRepository) VerifyProof(verificationCode string, salt string, proofBytes []byte) error {
	return utils.StackError(z.proofProcessor.VerifyProof(verificationCode, salt, proofBytes))
}

func NewZkRepository(zkProofProcessor zk.IProofProcessor) IZkRepository {
	return &ZkRepository{
		proofProcessor: zkProofProcessor,
	}
}
