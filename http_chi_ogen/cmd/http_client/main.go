package main

import (
	"context"
	"errors"
	"log"

	"github.com/brianvoe/gofakeit/v7"

	weatherV1 "github.com/arlengur/go-micros/http_chi_ogen/pkg/openapi/weather/v1"
)

const (
	serverURL       = "http://localhost:8080"
	defaultCityName = "Moscow"
	defaultMinTemp  = -10
	defaultMaxTemp  = 40
)

func main() {
	ctx := context.Background()

	// Инициализация Ogen-клиента
	client, err := weatherV1.NewClient(serverURL)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании клиента: %v", err)
	}

	log.Println("=== Тестирование API для работы с данными о погоде ===")
	log.Println()

	// 1. Пытаемся получить данные о погоде (которых пока нет)
	log.Printf("🌦️ Получение данных о погоде для города %s\n", defaultCityName)
	log.Println("===================================================")

	weatherResp, err := client.GetWeatherByCity(ctx, weatherV1.GetWeatherByCityParams{
		City: defaultCityName,
	})
	// Проверяем статус ошибки - если 404, значит данных просто нет
	if err != nil {
		// Проверяем, что ошибка содержит информацию о статусе 404
		var errResp *weatherV1.GenericErrorStatusCode
		if errors.As(err, &errResp) && errResp.StatusCode == 404 {
			log.Printf("ℹ️ Данные о погоде для города %s не найдены\n", defaultCityName)
			return
		}

		log.Printf("❌ Ошибка при получении погоды: %v\n", err)
		return
	}

	log.Printf("Данные о погоде для города %s: %+v\n", defaultCityName, weatherResp)

	// 2. Обновляем данные о погоде
	log.Printf("🔄 Обновление данных о погоде для города %s\n", defaultCityName)
	log.Println("=====================================================")

	// Создаем запрос на обновление погоды
	updateRequest := &weatherV1.UpdateWeatherRequest{
		Temperature: gofakeit.Float32Range(defaultMinTemp, defaultMaxTemp),
	}

	updatedWeather, err := client.UpdateWeatherByCity(ctx, updateRequest, weatherV1.UpdateWeatherByCityParams{
		City: defaultCityName,
	})
	if err != nil {
		log.Printf("❌ Ошибка при обновлении погоды: %v\n", err)
		return
	}
	log.Printf("✅ Данные о погоде обновлены: %+v\n", updatedWeather)

	// 3. Получаем обновленные данные о погоде
	log.Printf("🌦️ Получение обновленных данных о погоде для города %s\n", defaultCityName)
	log.Println("===========================================================")

	weatherResp, err = client.GetWeatherByCity(ctx, weatherV1.GetWeatherByCityParams{
		City: defaultCityName,
	})
	if err != nil {
		log.Printf("❌ Ошибка при получении погоды: %v\n", err)
		return
	}

	log.Printf("✅ Получены данные о погоде: %+v\n", weatherResp)
	log.Println("Тестирование завершено успешно!")
}
