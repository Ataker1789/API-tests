package advertisement_helper

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/internal/utils"
)

// CreateMultipleTestAdvertisements создаёт n тестовых объявлений и возвращает их ID.
func CreateMultipleTestAdvertisements(
	t *testing.T,
	token string,
	n int,
) []string {
	var ids []string
	for i := 0; i < n; i++ {
		title := "AUTO_TEST_" + utils.RandomString(8)
		body := myAdvertisements.CreateAdvertisement(
			t,
			token,
			title,
			"Description for "+title,
			100,
			2,
			[]string{"../../testdata/square.png"},
			http.StatusCreated,
		)
		adID := gjson.Get(body, "id").String()
		require.NotEmpty(t, adID, "ID созданного объявления не должен быть пустым")
		ids = append(ids, adID)
		utils.LogWithLabelAndTimeStamp("Helper", fmt.Sprintf("Создано объявление %s с ID %s", title, adID))
	}
	return ids
}
