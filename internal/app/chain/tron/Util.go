package tron

import (
	"fmt"
	"log"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

const (
	LocalnetRPCEndpoint = "http://localhost:8899"
	TestnetRPCEndpoint  = "https://api.shasta.trongrid.io"
	MainnetRPCEndpoint  = "https://api.trongrid.io"
)

// ConvertToBigInt takes an interface{} and converts it to *big.Int
func ConvertToBigInt(value interface{}) (*big.Int, error) {
	bigIntValue := new(big.Int)

	switch v := reflect.ValueOf(value); v.Kind() {
	case reflect.String:
		// If value is a string, try to parse it as a big integer
		_, ok := bigIntValue.SetString(value.(string), 10)
		if !ok {
			return nil, fmt.Errorf("invalid string format")
		}
	case reflect.Int, reflect.Int64:
		// If value is an int or int64, convert it to big.Int
		bigIntValue.SetInt64(v.Int())
	case reflect.Float64:
		// If value is a float64, convert to string first and then big.Int
		str := strconv.FormatFloat(value.(float64), 'f', -1, 64)
		_, ok := bigIntValue.SetString(str, 10)
		if !ok {
			return nil, fmt.Errorf("invalid float format")
		}
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}

	return bigIntValue, nil
}

// ConvertToReadableAmount takes a balance of any type and precision, and converts it to a human-readable form.
func ConvertToReadableAmount(balance interface{}, precision interface{}) string {
	// Convert balance and precision to *big.Int
	balanceBigInt, err := ConvertToBigInt(balance)
	if err != nil {
		log.Fatalf("invalid balance: %v", err)
		return ""
	}

	precisionBigInt, err := ConvertToBigInt(precision)
	if err != nil {
		log.Fatalf("invalid precision: %v", err)
		return ""
	}

	// If precision is 0, return balance as string
	if precisionBigInt.Cmp(big.NewInt(0)) == 0 {
		return balanceBigInt.String()
	}

	// Create 10^precision as *big.Int
	precisionFactor := new(big.Int).Exp(big.NewInt(10), precisionBigInt, nil)

	// Perform division: balance / 10^precision
	readableAmount := new(big.Rat).SetFrac(balanceBigInt, precisionFactor)

	// Convert to string with decimal places
	prec := precisionBigInt.Int64() // Convert precision back to int64 for formatting
	readableAmountStr := readableAmount.FloatString(int(prec))
	if strings.Contains(readableAmountStr, ".") {
		readableAmountStr = strings.TrimRight(strings.TrimRight(readableAmountStr, "0"), ".")
	}
	return readableAmountStr
}
