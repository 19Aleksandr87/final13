package services

import (
	"encoding/json"
	"errors"
	"time"
)

func Check(r *json.Decoder) (string, Content, error) {
	var content Content
	err := r.Decode(&content)

	if err != nil {
		return "", Content{}, err
	}
	if content.Title == "" {
		return "", Content{}, errors.New("укажите заголовок задачи")
	}

	tNow := time.Now()
	var dstart time.Time

	if content.Date == "" {
		dstart = time.Now()
	} else {
		dstart, err = time.Parse(FormatDate, content.Date)
		if err != nil {
			return "", Content{}, err
		}
	}
	var str string
	if content.Repeat == "" {
		if tNow.Before(dstart) {
			str = dstart.Format(FormatDate)
		} else {

			str = tNow.Format(FormatDate)
		}
	} else {
		str, err = NextDate(tNow.Format(FormatDate), dstart.Format(FormatDate), content.Repeat)
	}
	return str, content, err
}
