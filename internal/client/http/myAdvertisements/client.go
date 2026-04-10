package myAdvertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpGetMyAdvertisements получает список собственных объявлений
func HttpGetMyAdvertisements(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.MyAdvertisementsPath).
		ContentType("application/json").
		Query("limit", "50").
		Expect(t).
		End().Response
}

// HttpPostAdvertisement создает новое объявление
func HttpPostAdvertisement(t *testing.T, token string, fields map[string]string, files []string) *http.Response {
	req := apiRunner.GetRunner().
		Auth(token).
		Create().
		Post(path.AdvertisementPath)

	// Добавляем файлы (photos)
	for _, file := range files {
		req = req.MultipartFile("photos", file)
	}

	// Добавляем поля (title, description, price, quantity)
	for k, v := range fields {
		req = req.MultipartFormData(k, v)
	}

	return req.Expect(t).End().Response
}

// HttpGetAdvertisementByID получает объявление по ID
func HttpGetAdvertisementByID(t *testing.T, id string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpSearchAdvertisements ищет объявления по запросу
func HttpSearchAdvertisements(t *testing.T, query string) *http.Response {
	return apiRunner.GetRunner().Create().Get(path.AdvertisementsPath).
		Query("search", query).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления
func HttpGetAdvertisementPhotos(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPhotosPath(id)).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement удаляет объявление по ID
func HttpDeleteAdvertisement(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
