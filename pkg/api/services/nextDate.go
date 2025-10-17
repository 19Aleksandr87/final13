package services

import (
	"errors"
	"sort"
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
	case "w":
		str, err := week(d, nowI, dstartI)
		return str, err
	case "m":
		str, err := month(d, nowI, dstartI)
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
		dstartI = dstartI.AddDate(0, 0, d[0][0])
		if nowI.Before(dstartI) {
			return string(dstartI.Format(FormatDate)), nil
		}
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

func month(d [][]int, nowI, dstartI time.Time) (string, error) {
	//проверяем коррктность данных
	if len(d) < 1 || len(d) > 2 {
		return "", errors.New("некорректные данные: для месяца")
	}
	//проверяем коррктность дней
	for _, i := range d[0] {
		if i < (-2) || i > 31 {
			return "", errors.New("некорректные данные: для дней")
		}
	}
	//проверяем коррктность месяцев
	if len(d) == 2 {
		sort.Ints(d[1])
		for _, i := range d[1] {
			if i < 0 || i > 12 {
				return "", errors.New("некорректные данные: для месяца")
			}
		}
	} else {
		d = append(d, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	}
	// определяем максимальное количество дней в месяце
	var dayMax = map[int]int{
		1: 31, 2: 28, 3: 31, 4: 30, 5: 31, 6: 30, 7: 31, 8: 31, 9: 30, 10: 31, 11: 30, 12: 31,
	}
	// проверяем на наличие несуществующих 31 дней в месяце
	control(d, dayMax)
	// приравниваем текуцую дату к дате отсчета если она меньше
	if dstartI.Before(nowI) {
		dstartI = nowI
	}
	// генератор следующей даты
	dstartI = generator(d[0], d[1], dayMax, dstartI)

	return string(dstartI.Format(FormatDate)), nil
}

func control(d [][]int, max map[int]int) error {
	max[2] = 29
	for _, i := range d[1] {
		for _, j := range d[0] {
			if max[i] >= j {
				return nil
			}
		}
	}
	return errors.New("указаны несуществующие дни в месяце")
}

func generator(days []int, months []int, max map[int]int, dstartI time.Time) time.Time {
	day := int(dstartI.Day())
	month := int(dstartI.Month())
	year := int(dstartI.Year())
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 {
		max[2] = 29
	}

	for _, m := range months {
		if month > m {
			continue
		}

		data := daysM(days, max[m])
		if len(data) == 0 {
			continue
		}

		if month < m {
			monthT := time.Month(m)
			return time.Date(year, monthT, data[0], 0, 0, 0, 0, time.UTC)
		}
		for _, i := range data {
			if day > i {
				continue
			}
			if day < i {
				monthT := time.Month(month)
				return time.Date(year, monthT, i, 0, 0, 0, 0, time.UTC)
			}
		}
		continue
	}

	dstartI = time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	generator(days, months, max, dstartI)
	return dstartI
}

func daysM(data []int, max int) []int {
	for i := 0; i < len(data); i++ {
		if data[i] == -1 {
			data = append(data, max)
		}
		if data[i] == -2 {
			data = append(data, max-1)
		}
		if data[i] > max {
			data = append(data[:i], data[i+1:]...)
			continue
		}
	}
	sort.Ints(data)
	return data
}

func week(d [][]int, nowI, dstartI time.Time) (string, error) {
	if len(d) != 1 {
		return "", errors.New("некорректные данные: для недели")
	}
	for _, i := range d[0] {
		if (i < 1) || (i > 7) {
			return "", errors.New("указан несуществующий день недели")
		}
	}

	sort.Ints(d[0])

	var w = map[string]int{
		"Monday": 1, "Tuesday": 2, "Wednesda": 3, "Thursday": 4,
		"Friday": 5, "Saturday": 6, "Sunday": 7,
	}

	if dstartI.Before(nowI) {
		dstartI = nowI
	}
	weekday := dstartI.Weekday().String()
	for _, i := range d[0] {
		if w[weekday] < i {
			dstartI = dstartI.AddDate(0, 0, i-w[weekday])
			return string(dstartI.Format(FormatDate)), nil
		}
	}
	dstartI = dstartI.AddDate(0, 0, (7 - w[weekday] + d[0][0]))
	return string(dstartI.Format(FormatDate)), nil
}
