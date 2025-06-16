package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	info := strings.Split(data, ",")
	if len(info) != 3 {
		return 0, "", 0, errors.New("invalid data")
	}
	steps, err := strconv.Atoi(info[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("invalid data steps")
	}
	duration, err := time.ParseDuration(info[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("invalid duration")
	}

	return steps, info[1], duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	distance := (float64(steps) * stepLength) / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	speed := distance / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	if steps <= 0 {
		return "", errors.New("invalid data steps")
	}
	if weight <= 0 || height <= 0 {
		return "", errors.New("вес и рост должны быть больше 0")
	}
	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	switch activity {
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity,
			duration.Hours(),
			distance,
			speed,
			calories), nil
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity,
			duration.Hours(),
			distance,
			speed,
			calories), nil
	default:
		return "", errors.New("неизвестный тип тренировки")

	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, errors.New("invalid data steps")
	}
	if weight <= 0 {
		return 0, errors.New("invalid data weight")
	}
	if height <= 0 {
		return 0, errors.New("invalid data height")
	}
	if duration <= 0 {
		return 0, errors.New("invalid data duration")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return (weight * meanSpeed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, errors.New("invalid data steps")
	}
	if weight <= 0 {
		return 0, errors.New("invalid data weight")
	}
	if height <= 0 {
		return 0, errors.New("invalid data height")
	}
	if duration <= 0 {
		return 0, errors.New("invalid data duration")
	}
	minutes := duration.Minutes()
	meanSpeed := meanSpeed(steps, height, duration)
	calories := (weight * meanSpeed * minutes) / minInH

	return calories * walkingCaloriesCoefficient, nil
}
