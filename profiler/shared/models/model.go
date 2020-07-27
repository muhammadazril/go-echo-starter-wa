package models

// Config Lucent App
const (
	ConfigAppName               string = "app.name"
	ConfigAppDebug              string = "app.debug"
	ConfigAppMaintenance        string = "app.maintenance"
	ConfigAppTimezone           string = "app.timezone"
	ConfigAppPort               string = "app.port"
	ConfigAppJWTSecret          string = "app.jwt.secret"
	ConfigAppJWTLifetime        string = "app.jwt.lifetime"
	ConfigAppKeyPrivate         string = "app.keys.private_key"
	ConfigAppKeyPublic          string = "app.keys.public_key"
	ConfigAppVersion            string = "app.version"
	ConfigAppActiveLinkLifetime string = "app.activation_link.lifetime"
)

// Config Redis
const (
	ConfigRedisPassword string = "redis.password"
	ConfigRedisHost     string = "redis.host"
	ConfigRedisPort     string = "redis.port"
	ConfigRedisDatabase string = "redis.database"
)

// Config MongoDB
const (
	ConfigMongoPassword string = "mongodb.password"
	ConfigMongoUsername string = "mongodb.username"
	ConfigMongoHost     string = "mongodb.host"
	ConfigMongoPort     string = "mongodb.port"
	ConfigMongoDatabase string = "mongodb.database"
	ConfigMongoReplica  string = "mongodb.replica_set"
)

// Config for Kafka
const (
	ConfigKafkaHosts    string = "kafka.hosts"
	ConfigKafkaUsername string = "kafka.username"
	ConfigKafkaPassword string = "kafka.password"
)
