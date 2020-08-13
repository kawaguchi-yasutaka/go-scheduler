package scheduler

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var minuteAbsoluteRegExp = `^(([0-9]|[1-5][0-9])|([0-9]|[1-5][0-9])-([0-9]|[1-5][0-9]))(,([0-9]|[1-5][0-9])|,([0-9]|[1-5][0-9])-([0-9]|[1-5][0-9]))*$|^\*$`
var minuteIntervalRegExp = `^(\*|([0-9]|[1-5][0-9])-([0-9]|[1-5][0-9]))/([0-9]|[1-5][0-9])`
var minuteAnyValueRegExp = `^\*$`

var errorInvalidMinuteFormat = errors.New("invalid cron minute format")

func ValidateMinuteFormat(str string) bool {
	r := regexp.MustCompile(fmt.Sprintf(
		"%s|%s|%s",
		minuteAnyValueRegExp,
		minuteAbsoluteRegExp,
		minuteIntervalRegExp),
	)
	return r.MatchString(str)
}


func ParseMinute(str string) (error ,cronField) {
	result := strings.Split(str,",")
	cronFiledMap := map[int64]int64{}
	for _, v := range result {
		if strings.Contains(v, "/") {
			slice := strings.Split(v, "/")
			from := 0
			end := 59
			if strings.Contains(slice[0], "-") {
				slice := strings.Split(slice[0], "-")
				f,err := strconv.Atoi(slice[0])
				if err != nil {
					return err,nil
				}
				e, err := strconv.Atoi(slice[1])
				if err != nil {
					return err,nil
				}
				from = f
				end = e
			}
			per, err := strconv.Atoi(slice[1])
			if err != nil {
				return err, nil
			}
			for i := from;i <=end;i++ {
				if i % per == 0 {
					cronFiledMap[int64(i)] = int64(i)
					}
			}
		} else if strings.Contains(v, "-") {
			slice := strings.Split(v, "-")
			from,err := strconv.Atoi(slice[0])
			if err != nil {
				return err,nil
			}
			end, err := strconv.Atoi(slice[1])
			if err != nil {
				return err,nil
			}
			//if from > end {
			//}
			for i := from;i<=end;i++ {
				cronFiledMap[int64(i)] = int64(i)
			}
		} else if v == "*"{
			for i := 0;i <=59;i++ {
				cronFiledMap[int64(i)] = int64(i)
			}
		} else {
			a, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err, nil
			}
			cronFiledMap[a] = a
		}
	}
	var field cronField
	for _, v := range cronFiledMap {field = append(field, v)}
	sort.Slice(field, func(i, j int) bool {return field[i] < field[j]})
	return nil, field
}

