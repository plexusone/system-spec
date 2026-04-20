package spec

// CloudflareResources contains Cloudflare resource bindings.
type CloudflareResources struct {
	// Zone is the DNS zone
	Zone string `json:"zone,omitempty"`

	// Workers lists Cloudflare Workers
	Workers []CloudflareWorker `json:"workers,omitempty"`

	// R2Buckets lists R2 storage buckets
	R2Buckets []R2Bucket `json:"r2Buckets,omitempty"`

	// KV lists KV namespaces
	KV []CloudflareKV `json:"kv,omitempty"`

	// D1 lists D1 databases
	D1 []CloudflareD1 `json:"d1,omitempty"`

	// Queues lists Cloudflare Queues
	Queues []CloudflareQueue `json:"queues,omitempty"`

	// DurableObjects lists Durable Object namespaces
	DurableObjects []CloudflareDurableObject `json:"durableObjects,omitempty"`
}

// CloudflareWorker represents a Cloudflare Worker.
type CloudflareWorker struct {
	// Name of the worker
	Name string `json:"name"`

	// Route pattern
	Route string `json:"route,omitempty"`
}

// R2Bucket represents a Cloudflare R2 bucket.
type R2Bucket struct {
	// Name of the bucket
	Name string `json:"name"`
}

// CloudflareKV represents a Cloudflare KV namespace.
type CloudflareKV struct {
	// Name of the KV namespace
	Name string `json:"name"`
}

// CloudflareD1 represents a Cloudflare D1 database.
type CloudflareD1 struct {
	// Name of the database
	Name string `json:"name"`
}

// CloudflareQueue represents a Cloudflare Queue.
type CloudflareQueue struct {
	// Name of the queue
	Name string `json:"name"`
}

// CloudflareDurableObject represents a Durable Object namespace.
type CloudflareDurableObject struct {
	// Name of the Durable Object class
	Name string `json:"name"`
}
