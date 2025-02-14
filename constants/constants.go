package constants

const (
	SuccessMessage      = "success"
	ErrFailedBadRequest = "Data tidak sesuai"
	ErrServerError      = "Terjadi kesalahan pada server"
)

var MappingClient = map[string]string{
	"fastcampus_ecommerce": "ini_secret_key",
}
