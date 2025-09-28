package helper

import (
	"strconv"
)

func DurationFromStartAndEnd(startTime, endTime string) (string, error) {
	iStartTime, err := strconv.Atoi(startTime)
	if err != nil {
		return "", err
	}

	iEndTime, err := strconv.Atoi(endTime)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(iEndTime - iStartTime), nil
}
