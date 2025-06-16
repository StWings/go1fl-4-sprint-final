package daysteps

import (
	"errors"
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	if len(data) == 0 {
		log.Println("parsePackage(data string) invalid data: " + data)
		return 0, 0, errors.New("parsePackage(data string) invalid data: " + data)
	}
	info := strings.Split(data, ",")
	if len(info) != 2 {
		log.Println("parsePackage(data string) invalid data: " + data)
		return 0, 0, errors.New("parsePackage(data string) invalid data: " + data)
	}
	steps, err := strconv.Atoi(info[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		log.Println("invalid data steps")
		return 0, 0, errors.New("invalid data steps")
	}
	duration, err := time.ParseDuration(info[1])
	if err != nil {
		return 0, 0, err
	}
	if duration <= 0 {
		log.Println("invalid duration")
		return 0, 0, errors.New("invalid duration")
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}
	if steps <= 0 {
		return ""
	}
	meters := float64(steps) * stepLength
	distanceInKm := meters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)
}
