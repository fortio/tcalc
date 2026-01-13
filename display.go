package main

import (
	"fmt"
	"strconv"
	"strings"

	"fortio.org/terminal/ansipixels/tcolor"
)

const (
	decimalString string = "Decimal: "
	hexString     string = "Hex: "
	binaryString  string = "Binary: \n"
)

func binaryDisplayStrings(num int64) []string {
	var rows [4][4][]string
	var j, k, w int
	for i := 63; i > -1; i-- {
		value := (int(((1 << i) & num) >> i))
		value = max(value, -value)
		valueString := strconv.Itoa(value)
		if rows[j][k] == nil { //nolint:gosec // we are doing some math to ensure we stay in bounds
			rows[j][k] = make([]string, 4)
		}
		rows[j][k][w] = valueString
		w = (w + 1) % 4
		if w != 0 {
			continue
		}
		k = (k + 1) % 4
		if k != 0 {
			continue
		}
		j = (j + 1) % 4
	}
	display := []string{binaryString}
	for i := range 4 {
		displayValue := strconv.Itoa((64 - (16 * i)))
		var inner []string
		for j := range 4 {
			inner = append(inner, strings.Join(rows[i][j], " "))
		}
		innerString := strings.Join(inner, "  ")

		display = append(display, displayValue+": "+innerString)
	}
	return display
}

func decimalDisplayString(num int64) string {
	return decimalString + strconv.Itoa(int(num)) + "\n"
}

func uintDisplayString(num int64) string {
	//nolint:gosec // we just want to display the unsigned representation of our number
	return "Unsigned " + decimalString + strconv.FormatUint((uint64(num)), 10)
}

func hexDisplayString(num int64) string {
	//nolint:gosec // I think it makes the most sense to display the hex value as unsigned
	return hexString + fmt.Sprintf("%x\n", uint64(num))
}

func displayString(num int64, err error) []string {
	display := append([]string{
		"",
		ASCII(num),
		decimalDisplayString(num),
		uintDisplayString(num),
		hexDisplayString(num),
	},
		binaryDisplayStrings(num)...)
	if err != nil {
		display[0] = tcolor.Red.Foreground() + "Last input was invalid" + tcolor.Reset
	}
	return display
}

func ASCII(num int64) string {
	switch num {
	case 12:
		return "ASCII: "
	case 7:
		return "ASCII: "
	case 10:
		return "ASCII: \\n"
	case 11:
		return "ASCII: \\r"
	default:
		return "ASCII: " + string(rune(num))
	}
}
