package common

// Config --
type Config struct {
	NodeID int64

	StorageDriver string
	FileStorage   FileStore
}

// FileStore --
type FileStore struct {
	PathPrefix string
}

// S3Store --
type S3Store struct {
}
