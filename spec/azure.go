package spec

// AzureResources contains Azure-specific resource bindings.
type AzureResources struct {
	// CosmosDB accounts
	CosmosDB []CosmosDBAccount `json:"cosmosdb,omitempty"`

	// ServiceBus namespaces
	ServiceBus []ServiceBusNamespace `json:"servicebus,omitempty"`

	// BlobStorage accounts
	BlobStorage []BlobStorageAccount `json:"blobstorage,omitempty"`

	// Functions apps
	Functions []AzureFunction `json:"functions,omitempty"`

	// RedisCache instances
	RedisCache []AzureRedisCache `json:"rediscache,omitempty"`

	// KeyVault instances
	KeyVault []AzureKeyVault `json:"keyvault,omitempty"`

	// EventHubs namespaces
	EventHubs []EventHubsNamespace `json:"eventhubs,omitempty"`

	// SQLDatabase instances
	SQLDatabase []AzureSQLDatabase `json:"sqldatabase,omitempty"`

	// VNet binding
	VNet string `json:"vnet,omitempty"`

	// Subnets
	Subnets []string `json:"subnets,omitempty"`
}

// CosmosDBAccount represents an Azure Cosmos DB account.
type CosmosDBAccount struct {
	// Name of the account
	Name string `json:"name"`

	// API type (sql, mongodb, cassandra, gremlin, table)
	API string `json:"api,omitempty"`
}

// ServiceBusNamespace represents an Azure Service Bus namespace.
type ServiceBusNamespace struct {
	// Name of the namespace
	Name string `json:"name"`

	// Queues in this namespace
	Queues []string `json:"queues,omitempty"`

	// Topics in this namespace
	Topics []string `json:"topics,omitempty"`
}

// BlobStorageAccount represents an Azure Blob Storage account.
type BlobStorageAccount struct {
	// Name of the storage account
	Name string `json:"name"`

	// Containers in this account
	Containers []string `json:"containers,omitempty"`
}

// AzureFunction represents an Azure Functions app.
type AzureFunction struct {
	// Name of the function app
	Name string `json:"name"`

	// Runtime (e.g., "dotnet", "node", "python", "java")
	Runtime string `json:"runtime,omitempty"`
}

// AzureRedisCache represents an Azure Cache for Redis instance.
type AzureRedisCache struct {
	// Name of the cache
	Name string `json:"name"`

	// Port number
	Port int `json:"port,omitempty"`
}

// AzureKeyVault represents an Azure Key Vault.
type AzureKeyVault struct {
	// Name of the vault
	Name string `json:"name"`
}

// EventHubsNamespace represents an Azure Event Hubs namespace.
type EventHubsNamespace struct {
	// Name of the namespace
	Name string `json:"name"`

	// Hubs in this namespace
	Hubs []string `json:"hubs,omitempty"`
}

// AzureSQLDatabase represents an Azure SQL Database.
type AzureSQLDatabase struct {
	// Name of the database
	Name string `json:"name"`

	// Server name
	Server string `json:"server,omitempty"`

	// Port number
	Port int `json:"port,omitempty"`
}
