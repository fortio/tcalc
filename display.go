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
	unicodeString string = "Unicode: "
	italicPrefix  string = "\x1b[3m"
)

var (
	GREEN = tcolor.Green.Foreground()
	RED   = tcolor.Red.Foreground()
)

func bitString(cur, prev int) string {
	switch [2]int{cur, prev} {
	case [2]int{0, 1}:
		return RED + "0" + tcolor.Reset
	case [2]int{1, 0}:
		return GREEN + "1" + tcolor.Reset
	default:
		return strconv.Itoa(cur)
	}
}

func binaryDisplayStrings(num, prev int64) []string {
	var rows [4][4][]string
	var j, k, w int
	for i := 63; i > -1; i-- {
		value := (int(((1 << i) & num) >> i))
		value = max(value, -value)
		prevValue := (int(((1 << i) & prev) >> i))
		prevValue = max(prevValue, -prevValue)
		valueString := bitString((value), (prevValue))
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
	return hexString + fmt.Sprintf("%X\n", uint64(num))
}

func displayString(num, prev int64, err error) []string {
	display := append([]string{
		"",
		unicodeDisplayString(num),
		decimalDisplayString(num),
		uintDisplayString(num),
		hexDisplayString(num),
	},
		binaryDisplayStrings(num, prev)...,
	)
	if err != nil {
		display[0] = RED + "Last input was invalid" + tcolor.Reset
	}
	return display
}

func unicodeDisplayString(num int64) string {
	switch num {
	case 12:
		return "Unicode: "
	case 7:
		return "Unicode: "
	case 10:
		return "Unicode: \\n"
	case 11:
		return "Unicode: \\r"
	default:
		return "Unicode: " + string(rune(num))
	}
}
