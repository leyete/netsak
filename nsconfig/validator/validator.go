package validator

import (
	"fmt"
	"regexp"
	_regexp "github.com/mattrbeam/netsak/nsconfig/validator/regexp"

	valid "github.com/asaskevich/govalidator"
	"github.com/badoux/checkmail"
)

func IsURL(str string) bool {
	var result = valid.IsURL(str)
	return result
}

func IsEmail(email string) bool {
	err := checkmail.ValidateFormat(email)
	if err != nil {
		fmt.Println(err)
		return false
	}
	/*err = checkmail.ValidateHost(email)
	if err != nil {
		fmt.Println(err)
		return false
	}*/
	return true
}

func match(text string, regex *regexp.Regexp) []string {
	parsed := regex.FindAllString(text, -1)
	return parsed
}

// Date finds all date strings
func Date(text string) []string {
	return match(text, _regexp.DateRegex)
}

// Time finds all time strings
func Time(text string) []string {
	return match(text, _regexp.TimeRegex)
}

// Phones finds all phone numbers
func Phones(text string) []string {
	return match(text, _regexp.PhoneRegex)
}

// PhonesWithExts finds all phone numbers with ext
func PhonesWithExts(text string) []string {
	return match(text, _regexp.PhonesWithExtsRegex)
}

// Links finds all link strings
func Links(text string) []string {
	return match(text, _regexp.LinkRegex)
}

// Emails finds all email strings
func Emails(text string) []string {
	return match(text, _regexp.EmailRegex)
}

// IPv4s finds all IPv4 addresses
func IPv4s(text string) []string {
	return match(text, _regexp.IPv4Regex)
}

// IPv6s finds all IPv6 addresses
func IPv6s(text string) []string {
	return match(text, _regexp.IPv6Regex)
}

// IPs finds all IP addresses (both IPv4 and IPv6)
func IPs(text string) []string {
	return match(text, _regexp.IPRegex)
}

// NotKnownPorts finds all not-known port numbers
func NotKnownPorts(text string) []string {
	return match(text, _regexp.NotKnownPortRegex)
}

// Prices finds all price strings
func Prices(text string) []string {
	return match(text, _regexp.PriceRegex)
}

// HexColors finds all hex color values
func HexColors(text string) []string {
	return match(text, _regexp.HexColorRegex)
}

// CreditCards finds all credit card numbers
func CreditCards(text string) []string {
	return match(text, _regexp.CreditCardRegex)
}

// BtcAddresses finds all bitcoin addresses
func BtcAddresses(text string) []string {
	return match(text, _regexp.BtcAddressRegex)
}

// StreetAddresses finds all street addresses
func StreetAddresses(text string) []string {
	return match(text, _regexp.StreetAddressRegex)
}

// ZipCodes finds all zip codes
func ZipCodes(text string) []string {
	return match(text, _regexp.ZipCodeRegex)
}

// PoBoxes finds all po-box strings
func PoBoxes(text string) []string {
	return match(text, _regexp.PoBoxRegex)
}

// SSNs finds all SSN strings
func SSNs(text string) []string {
	return match(text, _regexp.SSNRegex)
}

// MD5Hexes finds all MD5 hex strings
func MD5Hexes(text string) []string {
	return match(text, _regexp.MD5HexRegex)
}

// SHA1Hexes finds all SHA1 hex strings
func SHA1Hexes(text string) []string {
	return match(text, _regexp.SHA1HexRegex)
}

// SHA256Hexes finds all SHA256 hex strings
func SHA256Hexes(text string) []string {
	return match(text, _regexp.SHA256HexRegex)
}

// GUIDs finds all GUID strings
func GUIDs(text string) []string {
	return match(text, _regexp.GUIDRegex)
}

// ISBN13s finds all ISBN13 strings
func ISBN13s(text string) []string {
	return match(text, _regexp.ISBN13Regex)
}

// ISBN10s finds all ISBN10 strings
func ISBN10s(text string) []string {
	return match(text, _regexp.ISBN10Regex)
}

// VISACreditCards finds all VISA credit card numbers
func VISACreditCards(text string) []string {
	return match(text, _regexp.VISACreditCardRegex)
}

// MCCreditCards finds all MasterCard credit card numbers
func MCCreditCards(text string) []string {
	return match(text, _regexp.MCCreditCardRegex)
}

// MACAddresses finds all MAC addresses
func MACAddresses(text string) []string {
	return match(text, _regexp.MACAddressRegex)
}

// IBANs finds all IBAN strings
func IBANs(text string) []string {
	return match(text, _regexp.IBANRegex)
}