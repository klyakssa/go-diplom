package luhn

import "strings"

func Valid(numberString string) bool {
	card := strings.ReplaceAll(numberString, " ", "")
	sum := 0
	nDigits := len(card)
	parity := nDigits % 2
	for i := 0; i <= nDigits-1; i++ {
		digit := int(card[i] - '0')
		if digit < 0 || digit > 9 {
			return false
		}
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}

//  function checkLuhn(string purportedCC) {
//      int sum := 0
//      int nDigits := length(purportedCC)
//      int parity := nDigits modulus 2
//      for i from 0 to nDigits - 1 {
//          int digit := integer(purportedCC[i])
//          if i modulus 2 = parity
//              digit := digit × 2
//              if digit > 9
//                  digit := digit - 9
//          sum := sum + digit
//      }
//      return (sum modulus 10) = 0
//  }
