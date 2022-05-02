package pgdutil

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"repo.pegadaian.co.id/ms-pds/modules/pgdlogger"
)

const (
	letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func CreateFloat64(input float64) *float64 {
	return &input
}

func CreateInt64(input int64) *int64 {
	return &input
}

func CreateInt8(input int8) *int8 {
	return &input
}

func ToJson(obj interface{}) string {
	b, err := json.Marshal(obj)

	if err != nil {
		pgdlogger.Make().Panic(err)
	}

	return string(b)
}

func FloatToString(input float64) string {
	return fmt.Sprintf("%f", input)
}

func RandomStr(n int, arrChecker map[string]bool) string {
	var randString string
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	randString = string(b)

	if arrChecker[randString] {
		randString = string(RandomStr(n, arrChecker))
	}

	arrChecker[randString] = true

	return randString
}

func NowDbBun() time.Time {
	return time.Now().Add(7 * time.Hour)
}

// NowUTC to get real current datetime but UTC format
func NowUTC() time.Time {
	return time.Now().UTC().Add(7 * time.Hour)
}

func InterfaceToMap(obj interface{}) map[string]interface{} {
	var mappedObj map[string]interface{}

	byteObj, err := json.Marshal(obj)

	if err != nil {
		pgdlogger.Make().Panic(err)
	}

	if err := json.Unmarshal(byteObj, &mappedObj); err != nil {
		pgdlogger.Make().Panic(err)
	}

	return mappedObj
}

func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}
