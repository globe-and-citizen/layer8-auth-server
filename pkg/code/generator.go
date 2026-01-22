package code

type ICodeGenerator interface {
	GenerateCode(salt string, input string) (string, error)
}
