package ibanlib

import (
	"math/big"
	"regexp"

	"github.com/thewolfnl/numericvalues"
)

// ValidateString : Validate given string
func ValidateString(str string) bool {
	iban := ConvertStringToIban(str)
	return ValidateIban(iban)
}

// ValidateIban : Validate given IBAN object
func ValidateIban(iban *IBAN) bool {
	if iban == nil {
		return false
	}
	return iban.validate()
}

func (iban *IBAN) validate() bool {
	// Validate check digits
	if iban.getCheckDigits() < 2 || iban.getCheckDigits() > 98 {
		return false
	}

	// Validate country specific structure
	if iban.validateCountrySpecificStructure() == false {
		return false
	}

	return iban.validateCheckDigits()
}

func validateStructure(str string) bool {
	return regexp.MustCompile(`^[A-Z]{2}\d{2}[A-Z0-9]{11,30}$`).MatchString(str)
}

func (iban *IBAN) validateCountrySpecificStructure() bool {
	countrySpecificStructure := countries[iban.CountryCode]
	if iban.validateIBANLength(countrySpecificStructure.Length) == false {
		return false
	}
	if iban.validateBBAN(countrySpecificStructure.BBANFormat) == false {
		return false
	}
	return true
}

func (iban *IBAN) validateIBANLength(length int) bool {
	return len(iban.Code()) == length
}

func (iban *IBAN) validateBBAN(BBANFormat string) bool {
	return regexp.MustCompile(BBANFormat).MatchString(iban.BBAN)
}

func (iban *IBAN) validateCheckDigits() bool {
	// Move the four initial characters to the end of the string
	account := iban.BBAN + iban.CountryCode + iban.CheckDigits

	// Replace each letter in the string with two digits, thereby expanding the string, where A = 10, B = 11, ..., Z = 35
	account = numericvalues.LettersToNumericValue(account)

	// Interpret the string as a big integer
	accountnumber, success := new(big.Int).SetString(account, 10)
	if !success {
		return false
	}

	modulo := new(big.Int).SetInt64(97)
	remainder := new(big.Int).Mod(accountnumber, modulo)

	// Check if module is equal to 1
	return remainder.Int64() == 1
}
