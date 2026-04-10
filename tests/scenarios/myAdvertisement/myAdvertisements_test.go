package myAdvertisement

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	advertisementHelper "api-tests-template/internal/helpers/advertisement-helper"
	validationHelper "api-tests-template/internal/helpers/validation-helper"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData  models.LoginOkResponse
	createdAds []string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя с кредами из переменных окружения")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	utils.LogWithLabelAndTimeStamp("Setup", "Авторизован пользователь: "+s.loginData.User.Email)
}

func (s *TestSuite) TearDownTest() {
	if len(s.createdAds) > 0 {
		utils.LogWithLabelAndTimeStamp("Cleanup", "Удаление созданных объявлений")
		for _, adID := range s.createdAds {
			myAdvertisements.DeleteAdvertisement(s.T(), s.loginData.Token, adID)
		}
		s.createdAds = []string{}
	}
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// ==================== ПОЗИТИВНЫЕ ТЕСТЫ ====================

/*
Positive: добавление объявления со всеми возможными параметрами.
Проверка, что объявление доступно через ручку GET /advertisement и в поиске по полному вхождению названия объявления.
Проверка, что фотография была загружена на сервер и доступна по ручке /advertisements/{id}/photos.
*/
func (s *TestSuite) TestCreateAdvertisementWithAllParametersPositive() {
	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithAllParametersPositive - START")

	// Подготовка тестовых данных
	title := "AUTO_TEST_" + utils.RandomString(8)
	description := "Описание тестового объявления для автотестов " + utils.RandomString(5)
	price := 12345
	quantity := 5
	photos := []string{"../../testdata/square.png"}

	var adID string

	s.Run("Создание_объявления_со_всеми_параметрами", func() {
		utils.LogWithLabelAndTimeStamp("Step", "POST /advertisement с title, description, price, quantity, photos")

		body := myAdvertisements.CreateAdvertisement(
			s.T(),
			s.loginData.Token,
			title,
			description,
			price,
			quantity,
			photos,
			http.StatusCreated,
		)

		// Валидация ответа
		adID = gjson.Get(body, "id").String()
		require.NotEmpty(s.T(), adID, "ID объявления не должен быть пустым")
		s.createdAds = append(s.createdAds, adID)

		require.Equal(s.T(), title, gjson.Get(body, "title").String(), "title не совпадает")
		require.Equal(s.T(), description, gjson.Get(body, "description").String(), "description не совпадает")
		require.Equal(s.T(), float64(price), gjson.Get(body, "price").Float(), "price не совпадает")
		require.Equal(s.T(), int64(quantity), gjson.Get(body, "quantity").Int(), "quantity не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, gjson.Get(body, "user_id").String(), "user_id не совпадает")
		require.NotEmpty(s.T(), gjson.Get(body, "created_at").String(), "created_at не должен быть пустым")
		require.NotEmpty(s.T(), gjson.Get(body, "updated_at").String(), "updated_at не должен быть пустым")

		utils.LogWithLabelAndTimeStamp("Result", "Объявление создано успешно с ID: "+adID)
	})

	s.Run("Получение_объявления_по_ID_через_GET", func() {
		utils.LogWithLabelAndTimeStamp("Step", "GET /advertisement?id="+adID)

		body := myAdvertisements.GetAdvertisementByID(s.T(), adID, http.StatusOK)

		// Валидация полей объявления
		require.Equal(s.T(), adID, gjson.Get(body, "id").String(), "id не совпадает")
		require.Equal(s.T(), title, gjson.Get(body, "title").String(), "title не совпадает")
		require.Equal(s.T(), description, gjson.Get(body, "description").String(), "description не совпадает")
		require.Equal(s.T(), float64(price), gjson.Get(body, "price").Float(), "price не совпадает")
		require.Equal(s.T(), int64(quantity), gjson.Get(body, "quantity").Int(), "quantity не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, gjson.Get(body, "user_id").String(), "user_id не совпадает")

		utils.LogWithLabelAndTimeStamp("Result", "Объявление успешно получено через GET /advertisement")
	})

	s.Run("Поиск_объявления_по_полному_названию", func() {
		utils.LogWithLabelAndTimeStamp("Step", "GET /advertisements?search="+title)

		body := myAdvertisements.SearchAdvertisements(s.T(), title, http.StatusOK)

		// Проверка наличия объявления в результатах поиска
		items := gjson.Get(body, "items").Array()
		require.NotEmpty(s.T(), items, "Результаты поиска не должны быть пустыми")

		found := false
		for _, item := range items {
			if gjson.Get(item.String(), "id").String() == adID {
				found = true
				require.Equal(s.T(), title, gjson.Get(item.String(), "title").String(), "title в поиске не совпадает")
				break
			}
		}
		require.True(s.T(), found, "Объявление не найдено в результатах поиска по названию")

		utils.LogWithLabelAndTimeStamp("Result", "Объявление найдено в поиске по полному названию")
	})

	s.Run("Проверка_загрузки_фотографий_на_сервер", func() {
		utils.LogWithLabelAndTimeStamp("Step", "GET /advertisements/"+adID+"/photos")

		body := myAdvertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, adID, http.StatusOK)

		// Валидация структуры фотографий
		photosArray := gjson.Parse(body).Array()
		require.NotEmpty(s.T(), photosArray, "Список фотографий не должен быть пустым")
		require.GreaterOrEqual(s.T(), len(photosArray), 1, "Должна быть минимум 1 фотография")

		firstPhoto := photosArray[0]
		photoID := gjson.Get(firstPhoto.String(), "id").String()
		photoURL := gjson.Get(firstPhoto.String(), "url").String()

		require.NotEmpty(s.T(), photoID, "ID фотографии не должен быть пустым")
		require.NotEmpty(s.T(), photoURL, "URL фотографии не должен быть пустым")
		require.Contains(s.T(), photoURL, "http", "URL должен содержать протокол http/https")

		utils.LogWithLabelAndTimeStamp("Result", "Фотография загружена на сервер и доступна по URL: "+photoURL)
	})

	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithAllParametersPositive - PASSED ✓")
}

/*
Positive: Проверка получения собственных объявлений в профиле
*/
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	utils.LogWithLabelAndTimeStamp("Test", "TestGetMyAdvertisementsPositive - START")

	// Создаём два тестовых объявления, чтобы гарантированно были данные
	createdIDs := advertisementHelper.CreateMultipleTestAdvertisements(s.T(), s.loginData.Token, 2)
	s.createdAds = append(s.createdAds, createdIDs...) // добавляем все ID в список на удаление

	var advertisementsBody string
	s.Run("Получение_списка_собственных_объявлений", func() {
		utils.LogWithLabelAndTimeStamp("Step", "GET /my/advertisements")
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK)
		utils.LogWithLabelAndTimeStamp("Result", "Список объявлений получен")
	})

	s.Run("Валидация_структуры_ответа_и_полей_объявлений", func() {
		utils.LogWithLabelAndTimeStamp("Step", "Проверка полей каждого объявления в ответе")

		items := gjson.Get(advertisementsBody, "items").Array()
		require.NotEmpty(s.T(), items, "Должны быть объявления в профиле пользователя")

		utils.LogWithLabelAndTimeStamp("Info", fmt.Sprintf("Найдено объявлений: %d", len(items)))

		for i, item := range items {
			s.T().Logf("Валидация объявления #%d", i+1) // заменили log.Printf на t.Logf

			// Проверка принадлежности текущему пользователю
			require.Equal(s.T(), s.loginData.User.Id, gjson.Get(item.String(), "user_id").String(),
				"user_id должен принадлежать текущему пользователю")

			// Проверка обязательных полей
			require.NotEmpty(s.T(), gjson.Get(item.String(), "id").String(), "id не должен быть пустым")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "title").String(), "title не должен быть пустым")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "description").String(), "description не должен быть пустым")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "created_at").String(), "created_at не должен быть пустым")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "updated_at").String(), "updated_at не должен быть пустым")

			// Проверка наличия фотографий
			photos := gjson.Get(item.String(), "photos").Array()
			require.NotEmpty(s.T(), photos, "У объявления должны быть фотографии")
		}

		utils.LogWithLabelAndTimeStamp("Result", "Все объявления прошли валидацию успешно")
	})

	utils.LogWithLabelAndTimeStamp("Test", "TestGetMyAdvertisementsPositive - PASSED ✓")
}

// ==================== НЕГАТИВНЫЕ ТЕСТЫ ====================

/*
Негативный: попытка создания объявления без авторизации (без Bearer token).
Проверка возврата 401 Unauthorized.
*/
func (s *TestSuite) TestCreateAdvertisementWithoutAuthNegative() {
	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutAuthNegative - START")

	s.Run("Попытка_создания_объявления_без_токена", func() {
		utils.LogWithLabelAndTimeStamp("Step", "POST /advertisement без Authorization header")

		body := myAdvertisements.CreateAdvertisement(
			s.T(),
			"", // пустой токен
			"Test Title",
			"Test Description",
			100,
			1,
			[]string{"../../testdata/square.png"},
			http.StatusUnauthorized,
		)

		utils.LogWithLabelAndTimeStamp("Step", "Валидация ответа 401 Unauthorized")
		validationHelper.ValidateUnauthorizedError(s.T(), body)

		utils.LogWithLabelAndTimeStamp("Result", "Получена корректная ошибка 401 при отсутствии токена")
	})

	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutAuthNegative - PASSED ✓")
}

/*
Негативный: попытка получения своих объявлений с невалидным токеном.
Проверка возврата 401 Unauthorized с корректным сообщением.
*/
func (s *TestSuite) TestGetMyAdvertisementsWithInvalidTokenNegative() {
	utils.LogWithLabelAndTimeStamp("Test", "TestGetMyAdvertisementsWithInvalidTokenNegative - START")

	s.Run("Запрос_с_невалидным_токеном", func() {
		utils.LogWithLabelAndTimeStamp("Step", "GET /my/advertisements с токеном 'invalid_token_123'")

		body := myAdvertisements.GetMyAdvertisements(s.T(), "invalid_token_123", http.StatusUnauthorized)

		utils.LogWithLabelAndTimeStamp("Step", "Валидация ответа 401 и сообщения об ошибке")

		errorField := gjson.Get(body, "error").String()
		messageField := gjson.Get(body, "message").String()

		require.Equal(s.T(), "unauthorized", errorField, "Ошибка должна быть 'unauthorized'")
		require.Equal(s.T(), "Invalid or expired token", messageField, "Сообщение не совпадает")

		utils.LogWithLabelAndTimeStamp("Result", "Получена корректная ошибка: "+errorField)
	})

	utils.LogWithLabelAndTimeStamp("Test", "TestGetMyAdvertisementsWithInvalidTokenNegative - PASSED ✓")
}

/*
Негативный: попытка создания объявления без обязательного поля title.
Согласно Swagger: title required=true
Ожидается 400 Bad Request.
*/
func (s *TestSuite) TestCreateAdvertisementWithoutTitleNegative() {
	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutTitleNegative - START")

	s.Run("Попытка_создания_без_поля_title", func() {
		utils.LogWithLabelAndTimeStamp("Step", "POST /advertisement без обязательного поля 'title'")

		body := myAdvertisements.CreateAdvertisement(
			s.T(),
			s.loginData.Token,
			"", // пустой title
			"Test Description",
			100,
			1,
			[]string{"../../testdata/square.png"},
			http.StatusBadRequest,
		)

		utils.LogWithLabelAndTimeStamp("Step", "Валидация ответа 400 Bad Request")
		validationHelper.ValidateBadRequestError(s.T(), body)

		utils.LogWithLabelAndTimeStamp("Result", "Получена корректная ошибка 400 при отсутствии title")
	})

	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutTitleNegative - PASSED ✓")
}

/*
Негативный: попытка создания объявления без обязательного поля description.
Согласно Swagger: description required=true
Ожидается 400 Bad Request.
*/
func (s *TestSuite) TestCreateAdvertisementWithoutDescriptionNegative() {
	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutDescriptionNegative - START")

	s.Run("Попытка_создания_без_поля_description", func() {
		utils.LogWithLabelAndTimeStamp("Step", "POST /advertisement без обязательного поля 'description'")

		body := myAdvertisements.CreateAdvertisement(
			s.T(),
			s.loginData.Token,
			"Test Title "+utils.RandomString(5),
			"", // пустой description
			100,
			1,
			[]string{"../../testdata/square.png"},
			http.StatusBadRequest,
		)

		utils.LogWithLabelAndTimeStamp("Step", "Валидация ответа 400 Bad Request")
		validationHelper.ValidateBadRequestError(s.T(), body)

		utils.LogWithLabelAndTimeStamp("Result", "Получена корректная ошибка 400 при отсутствии description")
	})

	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutDescriptionNegative - PASSED ✓")
}

/*
Негативный: попытка создания объявления без обязательного поля photos.
Согласно Swagger: photos required=true (1-3 photos required)
Ожидается 400 Bad Request.
*/
func (s *TestSuite) TestCreateAdvertisementWithoutPhotosNegative() {
	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutPhotosNegative - START")

	s.Run("Попытка_создания_без_поля_photos", func() {
		utils.LogWithLabelAndTimeStamp("Step", "POST /advertisement без обязательного поля 'photos'")

		body := myAdvertisements.CreateAdvertisement(
			s.T(),
			s.loginData.Token,
			"Test Title "+utils.RandomString(5),
			"Test Description",
			100,
			1,
			nil, // nil photos
			http.StatusBadRequest,
		)

		utils.LogWithLabelAndTimeStamp("Step", "Валидация ответа 400 Bad Request")
		validationHelper.ValidateBadRequestError(s.T(), body)

		utils.LogWithLabelAndTimeStamp("Result", "Получена корректная ошибка 400 при отсутствии photos")
	})

	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutPhotosNegative - PASSED ✓")
}

/*
Негативный: попытка создания объявления без всех обязательных полей.
Проверка множественной валидации на бэкенде.
Ожидается 400 Bad Request.
*/
func (s *TestSuite) TestCreateAdvertisementWithoutAllRequiredFieldsNegative() {
	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutAllRequiredFieldsNegative - START")

	s.Run("Попытка_создания_без_всех_обязательных_полей", func() {
		utils.LogWithLabelAndTimeStamp("Step", "POST /advertisement без title, description, photos")

		body := myAdvertisements.CreateAdvertisement(
			s.T(),
			s.loginData.Token,
			"", // нет title
			"", // нет description
			100,
			1,
			nil, // нет photos
			http.StatusBadRequest,
		)

		utils.LogWithLabelAndTimeStamp("Step", "Валидация ответа 400 Bad Request")
		validationHelper.ValidateBadRequestError(s.T(), body)

		utils.LogWithLabelAndTimeStamp("Result", "Получена корректная ошибка 400 при отсутствии всех обязательных полей")
	})

	utils.LogWithLabelAndTimeStamp("Test", "TestCreateAdvertisementWithoutAllRequiredFieldsNegative - PASSED ✓")
}
