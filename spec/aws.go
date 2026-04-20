package spec

// AWSResources contains AWS-specific resource bindings.
type AWSResources struct {
	// RDS instances
	RDS []RDSInstance `json:"rds,omitempty"`

	// DynamoDB tables
	DynamoDB []DynamoDBTable `json:"dynamodb,omitempty"`

	// SQS queues
	SQS []SQSQueue `json:"sqs,omitempty"`

	// SNS topics
	SNS []SNSTopic `json:"sns,omitempty"`

	// S3 buckets
	S3 []S3Bucket `json:"s3,omitempty"`

	// Bedrock models
	Bedrock []BedrockModel `json:"bedrock,omitempty"`

	// Lambda functions
	Lambda []LambdaFunction `json:"lambda,omitempty"`

	// ElastiCache clusters (Redis/Memcached)
	ElastiCache []ElastiCacheCluster `json:"elasticache,omitempty"`

	// Kinesis streams
	Kinesis []KinesisStream `json:"kinesis,omitempty"`

	// EventBridge event buses
	EventBridge []EventBridgeBus `json:"eventbridge,omitempty"`

	// SecretsManager secrets
	SecretsManager []SecretsManagerSecret `json:"secretsmanager,omitempty"`

	// API Gateway APIs
	APIGateway []APIGatewayAPI `json:"apigateway,omitempty"`

	// CloudFront distributions
	CloudFront []CloudFrontDistribution `json:"cloudfront,omitempty"`

	// VPC binding
	VPC string `json:"vpc,omitempty"`

	// Subnets
	Subnets []string `json:"subnets,omitempty"`

	// SecurityGroups
	SecurityGroups []string `json:"securityGroups,omitempty"`
}

// RDSInstance represents an RDS database instance.
type RDSInstance struct {
	// Name or identifier of the RDS instance
	Name string `json:"name"`

	// Engine (mysql, postgres, aurora-mysql, etc.)
	Engine string `json:"engine,omitempty"`

	// Port number
	Port int `json:"port,omitempty"`
}

// DynamoDBTable represents a DynamoDB table.
type DynamoDBTable struct {
	// Name of the table
	Name string `json:"name"`
}

// SQSQueue represents an SQS queue.
type SQSQueue struct {
	// Name of the queue
	Name string `json:"name"`
}

// SNSTopic represents an SNS topic.
type SNSTopic struct {
	// Name of the topic
	Name string `json:"name"`
}

// S3Bucket represents an S3 bucket.
type S3Bucket struct {
	// Name of the bucket
	Name string `json:"name"`
}

// BedrockModel represents an AWS Bedrock model.
type BedrockModel struct {
	// ModelID is the Bedrock model identifier
	ModelID string `json:"modelId"`
}

// LambdaFunction represents an AWS Lambda function.
type LambdaFunction struct {
	// Name of the function
	Name string `json:"name"`

	// Runtime (e.g., "go1.x", "python3.9", "nodejs18.x")
	Runtime string `json:"runtime,omitempty"`
}

// ElastiCacheCluster represents an ElastiCache cluster.
type ElastiCacheCluster struct {
	// Name of the cluster
	Name string `json:"name"`

	// Engine (redis, memcached)
	Engine string `json:"engine,omitempty"`

	// Port number
	Port int `json:"port,omitempty"`
}

// KinesisStream represents a Kinesis data stream.
type KinesisStream struct {
	// Name of the stream
	Name string `json:"name"`
}

// EventBridgeBus represents an EventBridge event bus.
type EventBridgeBus struct {
	// Name of the event bus
	Name string `json:"name"`
}

// SecretsManagerSecret represents a Secrets Manager secret.
type SecretsManagerSecret struct {
	// Name of the secret
	Name string `json:"name"`
}

// APIGatewayAPI represents an API Gateway API.
type APIGatewayAPI struct {
	// Name of the API
	Name string `json:"name"`

	// Type (REST, HTTP, WebSocket)
	Type string `json:"type,omitempty"`
}

// CloudFrontDistribution represents a CloudFront distribution.
type CloudFrontDistribution struct {
	// ID of the distribution
	ID string `json:"id"`

	// Domain name
	Domain string `json:"domain,omitempty"`
}
