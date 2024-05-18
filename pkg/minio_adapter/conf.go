package minio_adapter

type MinioConfig struct {
	ServerURL  string `long:"url"`
	AccessKey  string `long:"access"`
	SecretKey  string `long:"secret"`
	BucketName string `long:"bucket"`
	TestFile   string `long:"test-file"`
}
