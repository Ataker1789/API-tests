package path

const (
	AuthLoginPath        = "/auth/login"
	MyAdvertisementsPath = "/my/advertisements"
	AdvertisementPath    = "/advertisement"
	AdvertisementsPath   = "/advertisements"
)

func AdvertisementPhotosPath(id string) string {
	return AdvertisementsPath + "/" + id + "/photos"
}
