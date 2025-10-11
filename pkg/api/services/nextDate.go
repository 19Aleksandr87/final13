package services

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now, dstart, repeat string) (string, error) {
	var d = [][]int{}

	//подготовка строки
	repeat = strings.TrimSpace(strings.ToLower(repeat))
	if len(repeat) == 0 {
		return "", errors.New("пустая строка")
	}

	//разделенние строки
	rep := strings.Split(repeat, " ")
	if ((rep[0] != "y") && (len(rep) < 2)) || (len(rep) > 3) {
		return "", errors.New("некорректные данные " + repeat)
	}
	// преобразование данных в int
	for i, str := range rep {
		if i == 0 {
			continue
		}
		dto := []int{}
		st := strings.Split(str, ",")
		for _, s := range st {
			num, err := strconv.Atoi(s)
			if err != nil {
				return "", errors.New("strconv.Atoi err " + repeat)
			}
			dto = append(dto, num)
		}
		d = append(d, dto)
	}
	// преобразование данных в time
	nowI, err := time.Parse(FormatDate, now)
	if err != nil {
		return "",
			errors.New("time.Parse err " + now)
	}
	dstartI, err := time.Parse(FormatDate, dstart)
	if err != nil {
		return "", errors.New("time.Parse err " + dstart)
	}

	switch rep[0] {
	case "d":
		str, err := day(d, nowI, dstartI)
		return str, err

	case "y":
		str, err := year(d, nowI, dstartI)
		return str, err

	default:
		return "", errors.New("неверный формат " + repeat)
	}
}

func day(d [][]int, nowI, dstartI time.Time) (string, error) {
	if len(d) != 1 {
		return "", errors.New("некорректные данные: для дней")
	}
	if (d[0][0] < 0) || (d[0][0] > 400) {
		return "", errors.New("количество дней не может быть меньше 0 или более 400")
	}
	if nowI.Equal(dstartI) {
		return string(dstartI.AddDate(0, 0, 0).Format(FormatDate)), nil
	}
	if nowI.Before(dstartI) {
		return string(dstartI.AddDate(0, 0, d[0][0]).Format(FormatDate)), nil
	}

	for {
		if nowI.Before(dstartI) {
			return string(dstartI.Format(FormatDate)), nil
		}
		dstartI = dstartI.AddDate(0, 0, d[0][0])
	}
}

func year(d [][]int, nowI, dstartI time.Time) (string, error) {
	if len(d) > 0 {
		return "", errors.New("некорректные данные для года")
	}
	if nowI.Before(dstartI) {
		return string(dstartI.AddDate(1, 0, 0).Format(FormatDate)), nil
	}
	if nowI.Equal(dstartI) {
		return string(dstartI.AddDate(0, 0, 0).Format(FormatDate)), nil
	}
	for {
		if nowI.Before(dstartI) {
			return string(dstartI.Format(FormatDate)), nil
		}
		dstartI = dstartI.AddDate(1, 0, 0)
	}
}

/*func month(d [][]int, nowI, dstartI time.Time) (string, error) {

}
*/ /*
func week(d [][]int, nowI, dstartI time.Time) (string, error) {
	if len(d)!=1{
		return "", errors.New("некорректные данные: для недели")
	}
	if (d[0][0] < 1) || (d[0][0] > 7) {
		return "", errors.New("указан несуществующий день недели")
	}
	var w = map[string]int{
		"Monday":   1,
		"Tuesday":  2,
		"Wednesda": 3,
		"Thursday": 4,
		"Friday":   5,
		"Saturday": 6,
		"Sunday":   7,
	}

	b := nowI.Compare(dstartI)
	weekday := nowI.Weekday()
	switch b{
	case -1:

		return string(dstartI.AddDate(0, 0, d[0][0]).Format(FormatDate)), nil
	}

	return "", nil
}
*/
