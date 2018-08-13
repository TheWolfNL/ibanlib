package ibanlib

import (
	"errors"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/thewolfnl/numericvalues"
)

// IBAN struct
type IBAN struct {
	// Country code
	CountryCode string

	// Check digits
	CheckDigits string

	// Country specific bban part
	BBAN string
}

// Code : Full IBAN
func (iban *IBAN) Code() string {
	return iban.CountryCode + iban.CheckDigits + iban.BBAN
}

// PrettyCode : Full IBAN prettyfied for display or printing
func (iban *IBAN) PrettyCode() string {
	return strings.Join(chunkString(iban.Code(), 4), " ")
}

// ConvertStringToIban : Convert string to IBAN object
func ConvertStringToIban(str string) *IBAN {
	// Remove spaces and force uppercase
	str = strings.ToUpper(strings.Replace(str, " ", "", -1))

	if validateStructure(str) {
		return &IBAN{
			CountryCode: str[0:2],
			CheckDigits: str[2:4],
			BBAN:        str[4:],
		}
	}
	return nil
}

// ConvertBBANStringToIban : Convert BBAN string to IBAN object
func ConvertBBANStringToIban(bban string, country string) *IBAN {
	// Remove spaces and force uppercase
	bban = strings.ToUpper(strings.Replace(bban, " ", "", -1))
	country = strings.ToUpper(country)

	alphanumeric := regexp.MustCompile(`^[A-Z0-9]+$`)
	if alphanumeric.MatchString(bban) == false || alphanumeric.MatchString(country) == false {
		return nil
	}

	iban := IBAN{
		CountryCode: country,
		BBAN:        bban,
	}
	digits, err := iban.calculateCheckDigits()
	if err != nil {
		return nil
	}
	iban.CheckDigits = digits

	if iban.validate() {
		return &iban
	}
	return nil
}

func (iban *IBAN) getCheckDigits() int {
	digits, err := strconv.ParseInt(iban.CheckDigits, 10, 64)
	if err != nil {
		log.Println(err)
	}
	return int(digits)
}

func (iban *IBAN) calculateCheckDigits() (string, error) {
	account := iban.BBAN + iban.CountryCode

	account = numericvalues.LettersToNumericValue(account)

	account += "00"

	// Interpret the string as a big integer
	accountnumber, success := new(big.Int).SetString(account, 10)
	if !success {
		return "", errors.New("Problem")
	}

	modulo := new(big.Int).SetInt64(97)
	remainder := new(big.Int).Mod(accountnumber, modulo)

	digits := int(98 - remainder.Int64())

	return strconv.Itoa(digits), nil
}

func chunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}
