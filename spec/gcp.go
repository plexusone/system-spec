package spec

// GCPResources contains GCP-specific resource bindings.
type GCPResources struct {
	// CloudSQL instances
	CloudSQL []CloudSQLInstance `json:"cloudsql,omitempty"`

	// PubSub topics
	PubSub []PubSubTopic `json:"pubsub,omitempty"`

	// GCS buckets
	GCS []GCSBucket `json:"gcs,omitempty"`

	// CloudRun services
	CloudRun []CloudRunService `json:"cloudrun,omitempty"`

	// CloudFunctions
	CloudFunctions []CloudFunction `json:"cloudfunctions,omitempty"`

	// Memorystore instances (Redis)
	Memorystore []MemorystoreInstance `json:"memorystore,omitempty"`

	// BigQuery datasets
	BigQuery []BigQueryDataset `json:"bigquery,omitempty"`

	// SecretManager secrets
	SecretManager []GCPSecret `json:"secretmanager,omitempty"`

	// Firestore databases
	Firestore []FirestoreDatabase `json:"firestore,omitempty"`

	// Spanner instances
	Spanner []SpannerInstance `json:"spanner,omitempty"`

	// VPC network
	VPCNetwork string `json:"vpcNetwork,omitempty"`
}

// CloudSQLInstance represents a Cloud SQL instance.
type CloudSQLInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// DatabaseType (mysql, postgres)
	DatabaseType string `json:"databaseType,omitempty"`
}

// PubSubTopic represents a Pub/Sub topic.
type PubSubTopic struct {
	// Name of the topic
	Name string `json:"name"`
}

// GCSBucket represents a GCS bucket.
type GCSBucket struct {
	// Name of the bucket
	Name string `json:"name"`
}

// CloudRunService represents a Cloud Run service.
type CloudRunService struct {
	// Name of the service
	Name string `json:"name"`

	// Region where deployed
	Region string `json:"region,omitempty"`
}

// CloudFunction represents a Cloud Function.
type CloudFunction struct {
	// Name of the function
	Name string `json:"name"`

	// Runtime (e.g., "go121", "python311", "nodejs18")
	Runtime string `json:"runtime,omitempty"`
}

// MemorystoreInstance represents a Memorystore (Redis) instance.
type MemorystoreInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// Port number
	Port int `json:"port,omitempty"`
}

// BigQueryDataset represents a BigQuery dataset.
type BigQueryDataset struct {
	// Name of the dataset
	Name string `json:"name"`
}

// GCPSecret represents a Secret Manager secret.
type GCPSecret struct {
	// Name of the secret
	Name string `json:"name"`
}

// FirestoreDatabase represents a Firestore database.
type FirestoreDatabase struct {
	// Name of the database (default: "(default)")
	Name string `json:"name"`
}

// SpannerInstance represents a Cloud Spanner instance.
type SpannerInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// Database name
	Database string `json:"database,omitempty"`
}
