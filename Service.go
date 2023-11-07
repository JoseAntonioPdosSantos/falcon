package falcon

type Service interface {
	GetByKey(key string) string
}
