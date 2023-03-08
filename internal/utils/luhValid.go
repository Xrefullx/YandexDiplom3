package utils

func LuhValid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luh int
	for i := 0; number > 0; i++ {
		cursor := number % 10
		if i%2 == 0 { // even
			cursor = cursor * 2
			if cursor > 9 {
				cursor = cursor%10 + cursor/10
			}
		}
		luh += cursor
		number = number / 10
	}
	return luh % 10
}
