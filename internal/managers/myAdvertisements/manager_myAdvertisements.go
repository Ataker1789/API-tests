package myAdvertisements

import (
	"log"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/client/http/myAdvertisements"
	httpHelper "api-tests-template/internal/helpers/http-helper"
	"api-tests-template/internal/utils"
)

// GetMyAdvertisements получает список собственных объявлений
func GetMyAdvertisements(t *testing.T, token string, expectedStatusCode int) string {
	utils.LogWithLabelAndTimeStamp("Manager", "Получение списка объявлений пользователя")

	resp := myAdvertisements.HttpGetMyAdvertisements(t, token)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)

	utils.LogWithLabelAndTimeStamp("Manager", "Список объявлений получен")
	return body
}

// CreateAdvertisement создает новое объявление
func CreateAdvertisement(
	t *testing.T,
	token string,
	title string,
	description string,
	price int,
	quantity int,
	photos []string,
	expectedStatusCode int,
) string {
	utils.LogWithLabelAndTimeStamp("Manager", "Создание объявления: "+title)

	fields := buildAdvertisementFields(title, description, price, quantity)

	resp := myAdvertisements.HttpPostAdvertisement(t, token, fields, photos)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)

	utils.LogWithLabelAndTimeStamp("Manager", "Объявление создано")
	return body
}

// GetAdvertisementByID получает объявление по ID
func GetAdvertisementByID(t *testing.T, id string, expectedStatusCode int) string {
	utils.LogWithLabelAndTimeStamp("Manager", "Получение объявления по ID: "+id)

	resp := myAdvertisements.HttpGetAdvertisementByID(t, id)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)

	return body
}

// SearchAdvertisements ищет объявления по запросу
func SearchAdvertisements(t *testing.T, query string, expectedStatusCode int) string {
	utils.LogWithLabelAndTimeStamp("Manager", "Поиск объявлений: "+query)

	resp := myAdvertisements.HttpSearchAdvertisements(t, query)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)

	return body
}

// GetAdvertisementPhotos получает фотографии объявления
func GetAdvertisementPhotos(t *testing.T, token string, id string, expectedStatusCode int) string {
	utils.LogWithLabelAndTimeStamp("Manager", "Получение фотографий объявления: "+id)

	resp := myAdvertisements.HttpGetAdvertisementPhotos(t, token, id)
	body := httpHelper.AssertStatusCode(t, resp, expectedStatusCode)

	return body
}

// DeleteAdvertisement удаляет объявление по ID
func DeleteAdvertisement(t *testing.T, token string, id string) {
	utils.LogWithLabelAndTimeStamp("Manager", "Удаление объявления: "+id)

	resp := myAdvertisements.HttpDeleteAdvertisement(t, token, id)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		log.Printf("[WARNING] Не удалось удалить объявление %s: статус %d", id, resp.StatusCode)
	} else {
		utils.LogWithLabelAndTimeStamp("Manager", "Объявление удалено: "+id)
	}
}

// buildAdvertisementFields строит map полей для создания объявления
func buildAdvertisementFields(title, description string, price, quantity int) map[string]string {
	fields := map[string]string{}

	if title != "" {
		fields["title"] = title
	}

	if description != "" {
		fields["description"] = description
	}

	if price > 0 {
		fields["price"] = strconv.Itoa(price)
	}

	if quantity > 0 {
		fields["quantity"] = strconv.Itoa(quantity)
	}

	return fields
}
