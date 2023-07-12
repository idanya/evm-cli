package directory

//go:generate mockery --name DirectoryClient
type DirectoryClient interface {
	LookupFunction(hash string) (string, error)
	LookupEvent(hash string) (string, error)
}
