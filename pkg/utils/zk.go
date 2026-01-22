package utils

import (
	stdBytes "bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"unicode/utf8"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/frontend"
)

const VerificationCodeSize = 6
const InputFrRepresentationSize = 38

const RunesPerElement = 7
const BufferSize = 4
const ElementByteSize = 32

func StringToFrElements(input string) ([InputFrRepresentationSize]fr.Element, error) {
	runeCount := utf8.RuneCountInString(input)

	runes := make([]rune, runeCount)

	ind := 0
	for _, runeValue := range input {
		runes[ind] = runeValue
		ind++
	}

	var elements [InputFrRepresentationSize]fr.Element

	ind = 0

	for i := 0; i < runeCount; i += RunesPerElement {
		var elementBytes [ElementByteSize]byte

		offset := 0
		for j := i; j < i+RunesPerElement; j++ {
			if j == runeCount {
				break
			}

			binary.LittleEndian.PutUint32(elementBytes[offset:offset+BufferSize], uint32(runes[j]))

			offset += BufferSize
		}

		element, err := fr.LittleEndian.Element(&elementBytes)
		if err != nil {
			return [38]fr.Element{}, fmt.Errorf("error while converting input to FR elements: %e", err)
		}

		elements[ind] = element
		ind++
	}

	for ; ind < InputFrRepresentationSize; ind++ {
		elements[ind] = fr.NewElement(0)
	}

	return elements, nil
}

func StringToCircuitVariables(
	input string,
) ([InputFrRepresentationSize]frontend.Variable, error) {
	elements, err := StringToFrElements(input)
	if err != nil {
		return [InputFrRepresentationSize]frontend.Variable{}, err
	}

	var circuitVariables [InputFrRepresentationSize]frontend.Variable

	for i := 0; i < InputFrRepresentationSize; i++ {
		circuitVariables[i] = elements[i]
	}

	return circuitVariables, nil
}

func ConvertCodeToCircuitVariables(code string) ([6]frontend.Variable, error) {
	codeAsCircuitVariables := [6]frontend.Variable{}

	for i := 0; i < len(code); i++ {
		currByte := code[i]

		if currByte >= '0' && currByte <= '9' {
			currByte -= '0'
		} else if currByte >= 'a' && currByte <= 'f' {
			currByte = currByte - 'a' + 10
		} else {
			return codeAsCircuitVariables, fmt.Errorf("invalid character at index %d of the verification code", i)
		}

		codeAsCircuitVariables[i] = currByte
	}

	return codeAsCircuitVariables, nil
}

func Equal(fst []byte, snd []byte) bool {
	if len(fst) != len(snd) {
		return false
	}

	for i := 0; i < len(fst); i++ {
		if fst[i] != snd[i] {
			return false
		}
	}

	return true
}

func WriteBytes[T io.WriterTo](key T) []byte {
	var byteWriter stdBytes.Buffer

	_, err := key.WriteTo(&byteWriter)
	if err != nil {
		log.Fatalf("Error while writing key to byte buffer: %e", err)
	}

	return byteWriter.Bytes()
}

func ReadBytes[T io.ReaderFrom](key T, keyBytes []byte) {
	byteBuffer := stdBytes.NewBuffer([]byte{})

	_, err := byteBuffer.Write(keyBytes)
	if err != nil {
		log.Fatalf("Error while writing bytes to byte buffer: %e", err)
	}

	_, err = key.ReadFrom(byteBuffer)
	if err != nil {
		log.Fatalf("Error while decoding key from bytes: %e", err)
	}
}
