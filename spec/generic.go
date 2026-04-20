package spec

// GenericResources contains provider-agnostic resource bindings.
// Use these for self-hosted or multi-cloud deployments.
type GenericResources struct {
	// Redis instances
	Redis []RedisInstance `json:"redis,omitempty"`

	// Kafka clusters
	Kafka []KafkaCluster `json:"kafka,omitempty"`

	// Elasticsearch clusters
	Elasticsearch []ElasticsearchCluster `json:"elasticsearch,omitempty"`

	// MongoDB instances
	MongoDB []MongoDBInstance `json:"mongodb,omitempty"`

	// PostgreSQL databases
	PostgreSQL []PostgreSQLDatabase `json:"postgresql,omitempty"`

	// MySQL databases
	MySQL []MySQLDatabase `json:"mysql,omitempty"`

	// RabbitMQ instances
	RabbitMQ []RabbitMQInstance `json:"rabbitmq,omitempty"`

	// Cassandra clusters
	Cassandra []CassandraCluster `json:"cassandra,omitempty"`

	// Memcached instances
	Memcached []MemcachedInstance `json:"memcached,omitempty"`

	// Minio (S3-compatible) instances
	Minio []MinioInstance `json:"minio,omitempty"`

	// NATS servers
	NATS []NATSServer `json:"nats,omitempty"`

	// etcd clusters
	Etcd []EtcdCluster `json:"etcd,omitempty"`

	// Vault (HashiCorp) instances
	Vault []VaultInstance `json:"vault,omitempty"`

	// Consul instances
	Consul []ConsulInstance `json:"consul,omitempty"`
}

// RedisInstance represents a Redis instance.
type RedisInstance struct {
	// Name or identifier
	Name string `json:"name"`

	// Host address
	Host string `json:"host,omitempty"`

	// Port number (default: 6379)
	Port int `json:"port,omitempty"`
}

// KafkaCluster represents a Kafka cluster.
type KafkaCluster struct {
	// Name of the cluster
	Name string `json:"name"`

	// Bootstrap servers (comma-separated)
	BootstrapServers string `json:"bootstrapServers,omitempty"`
}

// ElasticsearchCluster represents an Elasticsearch cluster.
type ElasticsearchCluster struct {
	// Name of the cluster
	Name string `json:"name"`

	// Host address
	Host string `json:"host,omitempty"`

	// Port number (default: 9200)
	Port int `json:"port,omitempty"`
}

// MongoDBInstance represents a MongoDB instance.
type MongoDBInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// Host address
	Host string `json:"host,omitempty"`

	// Port number (default: 27017)
	Port int `json:"port,omitempty"`

	// Database name
	Database string `json:"database,omitempty"`
}

// PostgreSQLDatabase represents a PostgreSQL database.
type PostgreSQLDatabase struct {
	// Name of the database
	Name string `json:"name"`

	// Host address
	Host string `json:"host,omitempty"`

	// Port number (default: 5432)
	Port int `json:"port,omitempty"`
}

// MySQLDatabase represents a MySQL database.
type MySQLDatabase struct {
	// Name of the database
	Name string `json:"name"`

	// Host address
	Host string `json:"host,omitempty"`

	// Port number (default: 3306)
	Port int `json:"port,omitempty"`
}

// RabbitMQInstance represents a RabbitMQ instance.
type RabbitMQInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// Host address
	Host string `json:"host,omitempty"`

	// Port number (default: 5672)
	Port int `json:"port,omitempty"`

	// Vhost
	Vhost string `json:"vhost,omitempty"`
}

// CassandraCluster represents a Cassandra cluster.
type CassandraCluster struct {
	// Name of the cluster
	Name string `json:"name"`

	// Contact points (comma-separated hosts)
	ContactPoints string `json:"contactPoints,omitempty"`

	// Port number (default: 9042)
	Port int `json:"port,omitempty"`

	// Keyspace
	Keyspace string `json:"keyspace,omitempty"`
}

// MemcachedInstance represents a Memcached instance.
type MemcachedInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// Host address
	Host string `json:"host,omitempty"`

	// Port number (default: 11211)
	Port int `json:"port,omitempty"`
}

// MinioInstance represents a Minio (S3-compatible) instance.
type MinioInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// Endpoint URL
	Endpoint string `json:"endpoint,omitempty"`

	// Buckets
	Buckets []string `json:"buckets,omitempty"`
}

// NATSServer represents a NATS server.
type NATSServer struct {
	// Name of the server
	Name string `json:"name"`

	// URL (e.g., "nats://localhost:4222")
	URL string `json:"url,omitempty"`
}

// EtcdCluster represents an etcd cluster.
type EtcdCluster struct {
	// Name of the cluster
	Name string `json:"name"`

	// Endpoints (comma-separated)
	Endpoints string `json:"endpoints,omitempty"`
}

// VaultInstance represents a HashiCorp Vault instance.
type VaultInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// Address (e.g., "https://vault.example.com:8200")
	Address string `json:"address,omitempty"`
}

// ConsulInstance represents a HashiCorp Consul instance.
type ConsulInstance struct {
	// Name of the instance
	Name string `json:"name"`

	// Address (e.g., "consul.example.com:8500")
	Address string `json:"address,omitempty"`
}
