package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// StorageBackend defines the interface for data storage
type StorageBackend interface {
	// Collection operations
	SaveCollection(collection *StorageCollection) error
	GetCollection(id string) (*StorageCollection, error)
	GetCollectionsByRoom(roomID string) ([]*StorageCollection, error)
	GetAllCollections() ([]*StorageCollection, error)
	DeleteCollection(id string) error

	// Environment operations
	SaveEnvironment(env *StorageEnvironment) error
	GetEnvironment(id string) (*StorageEnvironment, error)
	GetEnvironmentsByRoom(roomID string) ([]*StorageEnvironment, error)
	GetAllEnvironments() ([]*StorageEnvironment, error)
	DeleteEnvironment(id string) error

	// Request operations
	SaveRequest(req *StorageRequest) error
	GetRequest(id string) (*StorageRequest, error)
	GetRequestsByCollection(collectionID string) ([]*StorageRequest, error)
	GetAllRequests() ([]*StorageRequest, error)
	DeleteRequest(id string) error

	// History operations
	SaveHistory(history *StorageHistory) error
	GetHistoryByRoom(roomID string, limit int) ([]*StorageHistory, error)
	GetAllHistory(limit int) ([]*StorageHistory, error)
	DeleteHistory(id string) error

	// Health check
	Ping() error
	Close() error
}

// Storage models for NoSQL/SQL agnostic storage
type StorageCollection struct {
	ID          string                 `json:"id" bson:"_id,omitempty"`
	Name        string                 `json:"name" bson:"name"`
	Description string                 `json:"description" bson:"description"`
	RoomID      string                 `json:"roomId,omitempty" bson:"roomId,omitempty"`
	Requests    []string               `json:"requests" bson:"requests"` // Array of request IDs
	Metadata    map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt" bson:"updatedAt"`
}

type StorageEnvironment struct {
	ID        string                 `json:"id" bson:"_id,omitempty"`
	Name      string                 `json:"name" bson:"name"`
	Variables map[string]interface{} `json:"variables" bson:"variables"`
	RoomID    string                 `json:"roomId,omitempty" bson:"roomId,omitempty"`
	IsActive  bool                   `json:"isActive" bson:"isActive"`
	CreatedAt time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updatedAt"`
}

type StorageRequest struct {
	ID           string                 `json:"id" bson:"_id,omitempty"`
	CollectionID string                 `json:"collectionId" bson:"collectionId"`
	Name         string                 `json:"name" bson:"name"`
	Method       string                 `json:"method" bson:"method"`
	URL          string                 `json:"url" bson:"url"`
	Headers      map[string]interface{} `json:"headers" bson:"headers"`
	Body         string                 `json:"body,omitempty" bson:"body,omitempty"`
	QueryParams  map[string]interface{} `json:"queryParams,omitempty" bson:"queryParams,omitempty"`
	AuthType     string                 `json:"authType,omitempty" bson:"authType,omitempty"`
	AuthData     map[string]interface{} `json:"authData,omitempty" bson:"authData,omitempty"`
	CreatedAt    time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt" bson:"updatedAt"`
}

type StorageHistory struct {
	ID           string                 `json:"id" bson:"_id,omitempty"`
	RoomID       string                 `json:"roomId,omitempty" bson:"roomId,omitempty"`
	Method       string                 `json:"method" bson:"method"`
	URL          string                 `json:"url" bson:"url"`
	StatusCode   int                    `json:"statusCode" bson:"statusCode"`
	ResponseTime int64                  `json:"responseTime" bson:"responseTime"`
	Timestamp    time.Time              `json:"timestamp" bson:"timestamp"`
	Headers      map[string]interface{} `json:"headers,omitempty" bson:"headers,omitempty"`
	Response     string                 `json:"response,omitempty" bson:"response,omitempty"`
}

// MongoStorage implements StorageBackend using MongoDB
type MongoStorage struct {
	client       *mongo.Client
	database     *mongo.Database
	collections  *mongo.Collection
	environments *mongo.Collection
	requests     *mongo.Collection
	history      *mongo.Collection
}

// NewMongoStorage creates a new MongoDB storage backend
func NewMongoStorage(connectionString, dbName string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set default values
	if connectionString == "" {
		connectionString = "mongodb://localhost:27017"
	}
	if dbName == "" {
		dbName = "resterx"
	}

	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(dbName)

	log.Printf("Connected to MongoDB: %s/%s", connectionString, dbName)

	return &MongoStorage{
		client:       client,
		database:     database,
		collections:  database.Collection("collections"),
		environments: database.Collection("environments"),
		requests:     database.Collection("requests"),
		history:      database.Collection("history"),
	}, nil
}

// Collection operations
func (m *MongoStorage) SaveCollection(collection *StorageCollection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.UpdatedAt = time.Now()
	if collection.CreatedAt.IsZero() {
		collection.CreatedAt = time.Now()
	}

	if collection.ID == "" {
		collection.ID = primitive.NewObjectID().Hex()
		_, err := m.collections.InsertOne(ctx, collection)
		return err
	}

	filter := bson.M{"_id": collection.ID}
	update := bson.M{"$set": collection}
	_, err := m.collections.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (m *MongoStorage) GetCollection(id string) (*StorageCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var collection StorageCollection
	filter := bson.M{"_id": id}
	err := m.collections.FindOne(ctx, filter).Decode(&collection)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (m *MongoStorage) GetCollectionsByRoom(roomID string) ([]*StorageCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"roomId": roomID}
	cursor, err := m.collections.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var collections []*StorageCollection
	if err := cursor.All(ctx, &collections); err != nil {
		return nil, err
	}
	return collections, nil
}

func (m *MongoStorage) GetAllCollections() ([]*StorageCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.collections.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var collections []*StorageCollection
	if err := cursor.All(ctx, &collections); err != nil {
		return nil, err
	}
	return collections, nil
}

func (m *MongoStorage) DeleteCollection(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := m.collections.DeleteOne(ctx, filter)
	return err
}

// Environment operations
func (m *MongoStorage) SaveEnvironment(env *StorageEnvironment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	env.UpdatedAt = time.Now()
	if env.CreatedAt.IsZero() {
		env.CreatedAt = time.Now()
	}

	if env.ID == "" {
		env.ID = primitive.NewObjectID().Hex()
		_, err := m.environments.InsertOne(ctx, env)
		return err
	}

	filter := bson.M{"_id": env.ID}
	update := bson.M{"$set": env}
	_, err := m.environments.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (m *MongoStorage) GetEnvironment(id string) (*StorageEnvironment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var env StorageEnvironment
	filter := bson.M{"_id": id}
	err := m.environments.FindOne(ctx, filter).Decode(&env)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

func (m *MongoStorage) GetEnvironmentsByRoom(roomID string) ([]*StorageEnvironment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"roomId": roomID}
	cursor, err := m.environments.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var environments []*StorageEnvironment
	if err := cursor.All(ctx, &environments); err != nil {
		return nil, err
	}
	return environments, nil
}

func (m *MongoStorage) GetAllEnvironments() ([]*StorageEnvironment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.environments.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var environments []*StorageEnvironment
	if err := cursor.All(ctx, &environments); err != nil {
		return nil, err
	}
	return environments, nil
}

func (m *MongoStorage) DeleteEnvironment(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := m.environments.DeleteOne(ctx, filter)
	return err
}

// Request operations
func (m *MongoStorage) SaveRequest(req *StorageRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req.UpdatedAt = time.Now()
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}

	if req.ID == "" {
		req.ID = primitive.NewObjectID().Hex()
		_, err := m.requests.InsertOne(ctx, req)
		return err
	}

	filter := bson.M{"_id": req.ID}
	update := bson.M{"$set": req}
	_, err := m.requests.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (m *MongoStorage) GetRequest(id string) (*StorageRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req StorageRequest
	filter := bson.M{"_id": id}
	err := m.requests.FindOne(ctx, filter).Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (m *MongoStorage) GetRequestsByCollection(collectionID string) ([]*StorageRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"collectionId": collectionID}
	cursor, err := m.requests.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []*StorageRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}

func (m *MongoStorage) GetAllRequests() ([]*StorageRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := m.requests.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var requests []*StorageRequest
	if err := cursor.All(ctx, &requests); err != nil {
		return nil, err
	}
	return requests, nil
}

func (m *MongoStorage) DeleteRequest(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := m.requests.DeleteOne(ctx, filter)
	return err
}

// History operations
func (m *MongoStorage) SaveHistory(history *StorageHistory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if history.Timestamp.IsZero() {
		history.Timestamp = time.Now()
	}

	if history.ID == "" {
		history.ID = primitive.NewObjectID().Hex()
		_, err := m.history.InsertOne(ctx, history)
		return err
	}

	filter := bson.M{"_id": history.ID}
	update := bson.M{"$set": history}
	_, err := m.history.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (m *MongoStorage) GetHistoryByRoom(roomID string, limit int) ([]*StorageHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 100
	}

	filter := bson.M{"roomId": roomID}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit))
	cursor, err := m.history.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []*StorageHistory
	if err := cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func (m *MongoStorage) GetAllHistory(limit int) ([]*StorageHistory, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 100
	}

	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit))
	cursor, err := m.history.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []*StorageHistory
	if err := cursor.All(ctx, &history); err != nil {
		return nil, err
	}
	return history, nil
}

func (m *MongoStorage) DeleteHistory(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	_, err := m.history.DeleteOne(ctx, filter)
	return err
}

func (m *MongoStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Ping(ctx, nil)
}

func (m *MongoStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

// SQLiteStorage implements StorageBackend using SQLite with GORM
type SQLiteStorage struct {
	db *gorm.DB
}

// SQLite models for GORM
type SQLiteCollection struct {
	ID          string    `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:""`
	RoomID      string    `gorm:"index"`
	Requests    string    `gorm:"type:text"` // JSON array of request IDs
	Metadata    string    `gorm:"type:text"` // JSON metadata
	CreatedAt   time.Time `gorm:""`
	UpdatedAt   time.Time `gorm:""`
}

type SQLiteEnvironment struct {
	ID        string    `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Variables string    `gorm:"type:text"` // JSON variables
	RoomID    string    `gorm:"index"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time `gorm:""`
	UpdatedAt time.Time `gorm:""`
}

type SQLiteRequest struct {
	ID           string    `gorm:"primaryKey"`
	CollectionID string    `gorm:"index"`
	Name         string    `gorm:"not null"`
	Method       string    `gorm:"not null"`
	URL          string    `gorm:"not null"`
	Headers      string    `gorm:"type:text"` // JSON headers
	Body         string    `gorm:"type:text"`
	QueryParams  string    `gorm:"type:text"` // JSON query params
	AuthType     string    `gorm:""`
	AuthData     string    `gorm:"type:text"` // JSON auth data
	CreatedAt    time.Time `gorm:""`
	UpdatedAt    time.Time `gorm:""`
}

type SQLiteHistory struct {
	ID           string    `gorm:"primaryKey"`
	RoomID       string    `gorm:"index"`
	Method       string    `gorm:"not null"`
	URL          string    `gorm:"not null"`
	StatusCode   int       `gorm:""`
	ResponseTime int64     `gorm:""`
	Timestamp    time.Time `gorm:"index"`
	Headers      string    `gorm:"type:text"` // JSON headers
	Response     string    `gorm:"type:text"`
}

func (SQLiteCollection) TableName() string  { return "storage_collections" }
func (SQLiteEnvironment) TableName() string { return "storage_environments" }
func (SQLiteRequest) TableName() string     { return "storage_requests" }
func (SQLiteHistory) TableName() string     { return "storage_history" }

// NewSQLiteStorage creates a new SQLite storage backend
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	if dbPath == "" {
		dbPath = "resterx_storage.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	// Auto-migrate tables
	if err := db.AutoMigrate(
		&SQLiteCollection{},
		&SQLiteEnvironment{},
		&SQLiteRequest{},
		&SQLiteHistory{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate SQLite tables: %w", err)
	}

	log.Printf("Connected to SQLite: %s", dbPath)

	return &SQLiteStorage{db: db}, nil
}

// Helper functions for JSON marshaling/unmarshaling
func toJSON(v interface{}) string {
	if v == nil {
		return "{}"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func fromJSON(s string, v interface{}) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}

// Collection operations for SQLite
func (s *SQLiteStorage) SaveCollection(collection *StorageCollection) error {
	collection.UpdatedAt = time.Now()
	if collection.CreatedAt.IsZero() {
		collection.CreatedAt = time.Now()
	}
	if collection.ID == "" {
		collection.ID = primitive.NewObjectID().Hex()
	}

	sqlCollection := SQLiteCollection{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		RoomID:      collection.RoomID,
		Requests:    toJSON(collection.Requests),
		Metadata:    toJSON(collection.Metadata),
		CreatedAt:   collection.CreatedAt,
		UpdatedAt:   collection.UpdatedAt,
	}

	return s.db.Save(&sqlCollection).Error
}

func (s *SQLiteStorage) GetCollection(id string) (*StorageCollection, error) {
	var sqlCollection SQLiteCollection
	if err := s.db.First(&sqlCollection, "id = ?", id).Error; err != nil {
		return nil, err
	}

	collection := &StorageCollection{
		ID:          sqlCollection.ID,
		Name:        sqlCollection.Name,
		Description: sqlCollection.Description,
		RoomID:      sqlCollection.RoomID,
		CreatedAt:   sqlCollection.CreatedAt,
		UpdatedAt:   sqlCollection.UpdatedAt,
	}

	fromJSON(sqlCollection.Requests, &collection.Requests)
	fromJSON(sqlCollection.Metadata, &collection.Metadata)

	return collection, nil
}

func (s *SQLiteStorage) GetCollectionsByRoom(roomID string) ([]*StorageCollection, error) {
	var sqlCollections []SQLiteCollection
	if err := s.db.Where("room_id = ?", roomID).Find(&sqlCollections).Error; err != nil {
		return nil, err
	}

	collections := make([]*StorageCollection, 0, len(sqlCollections))
	for _, sc := range sqlCollections {
		c := &StorageCollection{
			ID:          sc.ID,
			Name:        sc.Name,
			Description: sc.Description,
			RoomID:      sc.RoomID,
			CreatedAt:   sc.CreatedAt,
			UpdatedAt:   sc.UpdatedAt,
		}
		fromJSON(sc.Requests, &c.Requests)
		fromJSON(sc.Metadata, &c.Metadata)
		collections = append(collections, c)
	}

	return collections, nil
}

func (s *SQLiteStorage) GetAllCollections() ([]*StorageCollection, error) {
	var sqlCollections []SQLiteCollection
	if err := s.db.Find(&sqlCollections).Error; err != nil {
		return nil, err
	}

	collections := make([]*StorageCollection, 0, len(sqlCollections))
	for _, sc := range sqlCollections {
		c := &StorageCollection{
			ID:          sc.ID,
			Name:        sc.Name,
			Description: sc.Description,
			RoomID:      sc.RoomID,
			CreatedAt:   sc.CreatedAt,
			UpdatedAt:   sc.UpdatedAt,
		}
		fromJSON(sc.Requests, &c.Requests)
		fromJSON(sc.Metadata, &c.Metadata)
		collections = append(collections, c)
	}

	return collections, nil
}

func (s *SQLiteStorage) DeleteCollection(id string) error {
	return s.db.Delete(&SQLiteCollection{}, "id = ?", id).Error
}

// Environment operations for SQLite
func (s *SQLiteStorage) SaveEnvironment(env *StorageEnvironment) error {
	env.UpdatedAt = time.Now()
	if env.CreatedAt.IsZero() {
		env.CreatedAt = time.Now()
	}
	if env.ID == "" {
		env.ID = primitive.NewObjectID().Hex()
	}

	sqlEnv := SQLiteEnvironment{
		ID:        env.ID,
		Name:      env.Name,
		Variables: toJSON(env.Variables),
		RoomID:    env.RoomID,
		IsActive:  env.IsActive,
		CreatedAt: env.CreatedAt,
		UpdatedAt: env.UpdatedAt,
	}

	return s.db.Save(&sqlEnv).Error
}

func (s *SQLiteStorage) GetEnvironment(id string) (*StorageEnvironment, error) {
	var sqlEnv SQLiteEnvironment
	if err := s.db.First(&sqlEnv, "id = ?", id).Error; err != nil {
		return nil, err
	}

	env := &StorageEnvironment{
		ID:        sqlEnv.ID,
		Name:      sqlEnv.Name,
		RoomID:    sqlEnv.RoomID,
		IsActive:  sqlEnv.IsActive,
		CreatedAt: sqlEnv.CreatedAt,
		UpdatedAt: sqlEnv.UpdatedAt,
	}

	fromJSON(sqlEnv.Variables, &env.Variables)

	return env, nil
}

func (s *SQLiteStorage) GetEnvironmentsByRoom(roomID string) ([]*StorageEnvironment, error) {
	var sqlEnvs []SQLiteEnvironment
	if err := s.db.Where("room_id = ?", roomID).Find(&sqlEnvs).Error; err != nil {
		return nil, err
	}

	environments := make([]*StorageEnvironment, 0, len(sqlEnvs))
	for _, se := range sqlEnvs {
		e := &StorageEnvironment{
			ID:        se.ID,
			Name:      se.Name,
			RoomID:    se.RoomID,
			IsActive:  se.IsActive,
			CreatedAt: se.CreatedAt,
			UpdatedAt: se.UpdatedAt,
		}
		fromJSON(se.Variables, &e.Variables)
		environments = append(environments, e)
	}

	return environments, nil
}

func (s *SQLiteStorage) GetAllEnvironments() ([]*StorageEnvironment, error) {
	var sqlEnvs []SQLiteEnvironment
	if err := s.db.Find(&sqlEnvs).Error; err != nil {
		return nil, err
	}

	environments := make([]*StorageEnvironment, 0, len(sqlEnvs))
	for _, se := range sqlEnvs {
		e := &StorageEnvironment{
			ID:        se.ID,
			Name:      se.Name,
			RoomID:    se.RoomID,
			IsActive:  se.IsActive,
			CreatedAt: se.CreatedAt,
			UpdatedAt: se.UpdatedAt,
		}
		fromJSON(se.Variables, &e.Variables)
		environments = append(environments, e)
	}

	return environments, nil
}

func (s *SQLiteStorage) DeleteEnvironment(id string) error {
	return s.db.Delete(&SQLiteEnvironment{}, "id = ?", id).Error
}

// Request operations for SQLite
func (s *SQLiteStorage) SaveRequest(req *StorageRequest) error {
	req.UpdatedAt = time.Now()
	if req.CreatedAt.IsZero() {
		req.CreatedAt = time.Now()
	}
	if req.ID == "" {
		req.ID = primitive.NewObjectID().Hex()
	}

	sqlReq := SQLiteRequest{
		ID:           req.ID,
		CollectionID: req.CollectionID,
		Name:         req.Name,
		Method:       req.Method,
		URL:          req.URL,
		Headers:      toJSON(req.Headers),
		Body:         req.Body,
		QueryParams:  toJSON(req.QueryParams),
		AuthType:     req.AuthType,
		AuthData:     toJSON(req.AuthData),
		CreatedAt:    req.CreatedAt,
		UpdatedAt:    req.UpdatedAt,
	}

	return s.db.Save(&sqlReq).Error
}

func (s *SQLiteStorage) GetRequest(id string) (*StorageRequest, error) {
	var sqlReq SQLiteRequest
	if err := s.db.First(&sqlReq, "id = ?", id).Error; err != nil {
		return nil, err
	}

	req := &StorageRequest{
		ID:           sqlReq.ID,
		CollectionID: sqlReq.CollectionID,
		Name:         sqlReq.Name,
		Method:       sqlReq.Method,
		URL:          sqlReq.URL,
		Body:         sqlReq.Body,
		AuthType:     sqlReq.AuthType,
		CreatedAt:    sqlReq.CreatedAt,
		UpdatedAt:    sqlReq.UpdatedAt,
	}

	fromJSON(sqlReq.Headers, &req.Headers)
	fromJSON(sqlReq.QueryParams, &req.QueryParams)
	fromJSON(sqlReq.AuthData, &req.AuthData)

	return req, nil
}

func (s *SQLiteStorage) GetRequestsByCollection(collectionID string) ([]*StorageRequest, error) {
	var sqlReqs []SQLiteRequest
	if err := s.db.Where("collection_id = ?", collectionID).Find(&sqlReqs).Error; err != nil {
		return nil, err
	}

	requests := make([]*StorageRequest, 0, len(sqlReqs))
	for _, sr := range sqlReqs {
		r := &StorageRequest{
			ID:           sr.ID,
			CollectionID: sr.CollectionID,
			Name:         sr.Name,
			Method:       sr.Method,
			URL:          sr.URL,
			Body:         sr.Body,
			AuthType:     sr.AuthType,
			CreatedAt:    sr.CreatedAt,
			UpdatedAt:    sr.UpdatedAt,
		}
		fromJSON(sr.Headers, &r.Headers)
		fromJSON(sr.QueryParams, &r.QueryParams)
		fromJSON(sr.AuthData, &r.AuthData)
		requests = append(requests, r)
	}

	return requests, nil
}

func (s *SQLiteStorage) GetAllRequests() ([]*StorageRequest, error) {
	var sqlReqs []SQLiteRequest
	if err := s.db.Find(&sqlReqs).Error; err != nil {
		return nil, err
	}

	requests := make([]*StorageRequest, 0, len(sqlReqs))
	for _, sr := range sqlReqs {
		r := &StorageRequest{
			ID:           sr.ID,
			CollectionID: sr.CollectionID,
			Name:         sr.Name,
			Method:       sr.Method,
			URL:          sr.URL,
			Body:         sr.Body,
			AuthType:     sr.AuthType,
			CreatedAt:    sr.CreatedAt,
			UpdatedAt:    sr.UpdatedAt,
		}
		fromJSON(sr.Headers, &r.Headers)
		fromJSON(sr.QueryParams, &r.QueryParams)
		fromJSON(sr.AuthData, &r.AuthData)
		requests = append(requests, r)
	}

	return requests, nil
}

func (s *SQLiteStorage) DeleteRequest(id string) error {
	return s.db.Delete(&SQLiteRequest{}, "id = ?", id).Error
}

// History operations for SQLite
func (s *SQLiteStorage) SaveHistory(history *StorageHistory) error {
	if history.Timestamp.IsZero() {
		history.Timestamp = time.Now()
	}
	if history.ID == "" {
		history.ID = primitive.NewObjectID().Hex()
	}

	sqlHistory := SQLiteHistory{
		ID:           history.ID,
		RoomID:       history.RoomID,
		Method:       history.Method,
		URL:          history.URL,
		StatusCode:   history.StatusCode,
		ResponseTime: history.ResponseTime,
		Timestamp:    history.Timestamp,
		Headers:      toJSON(history.Headers),
		Response:     history.Response,
	}

	return s.db.Create(&sqlHistory).Error
}

func (s *SQLiteStorage) GetHistoryByRoom(roomID string, limit int) ([]*StorageHistory, error) {
	if limit <= 0 {
		limit = 100
	}

	var sqlHistories []SQLiteHistory
	if err := s.db.Where("room_id = ?", roomID).Order("timestamp DESC").Limit(limit).Find(&sqlHistories).Error; err != nil {
		return nil, err
	}

	histories := make([]*StorageHistory, 0, len(sqlHistories))
	for _, sh := range sqlHistories {
		h := &StorageHistory{
			ID:           sh.ID,
			RoomID:       sh.RoomID,
			Method:       sh.Method,
			URL:          sh.URL,
			StatusCode:   sh.StatusCode,
			ResponseTime: sh.ResponseTime,
			Timestamp:    sh.Timestamp,
			Response:     sh.Response,
		}
		fromJSON(sh.Headers, &h.Headers)
		histories = append(histories, h)
	}

	return histories, nil
}

func (s *SQLiteStorage) GetAllHistory(limit int) ([]*StorageHistory, error) {
	if limit <= 0 {
		limit = 100
	}

	var sqlHistories []SQLiteHistory
	if err := s.db.Order("timestamp DESC").Limit(limit).Find(&sqlHistories).Error; err != nil {
		return nil, err
	}

	histories := make([]*StorageHistory, 0, len(sqlHistories))
	for _, sh := range sqlHistories {
		h := &StorageHistory{
			ID:           sh.ID,
			RoomID:       sh.RoomID,
			Method:       sh.Method,
			URL:          sh.URL,
			StatusCode:   sh.StatusCode,
			ResponseTime: sh.ResponseTime,
			Timestamp:    sh.Timestamp,
			Response:     sh.Response,
		}
		fromJSON(sh.Headers, &h.Headers)
		histories = append(histories, h)
	}

	return histories, nil
}

func (s *SQLiteStorage) DeleteHistory(id string) error {
	return s.db.Delete(&SQLiteHistory{}, "id = ?", id).Error
}

func (s *SQLiteStorage) Ping() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func (s *SQLiteStorage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// StorageManager handles storage with fallback support
type StorageManager struct {
	primary   StorageBackend
	fallback  StorageBackend
	usePrimary bool
}

// NewStorageManager creates a new storage manager with MongoDB primary and SQLite fallback
func NewStorageManager() (*StorageManager, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	mongoDBName := os.Getenv("MONGODB_DATABASE")
	sqlitePath := os.Getenv("SQLITE_PATH")

	if sqlitePath == "" {
		sqlitePath = "resterx_storage.db"
	}

	var primary StorageBackend
	var usePrimary bool

	// Try MongoDB first if URI is provided
	if mongoURI != "" {
		mongo, err := NewMongoStorage(mongoURI, mongoDBName)
		if err == nil {
			log.Println("Using MongoDB as primary storage")
			primary = mongo
			usePrimary = true
		} else {
			log.Printf("Failed to connect to MongoDB: %v", err)
			log.Println("Falling back to SQLite")
		}
	}

	// Create SQLite as fallback
	sqlite, err := NewSQLiteStorage(sqlitePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize fallback storage: %w", err)
	}

	// If MongoDB failed or not configured, use SQLite as primary
	if primary == nil {
		log.Println("Using SQLite as primary storage")
		primary = sqlite
		usePrimary = true
	}

	return &StorageManager{
		primary:    primary,
		fallback:   sqlite,
		usePrimary: usePrimary,
	}, nil
}

// Wrapper methods that try primary then fallback
func (sm *StorageManager) SaveCollection(collection *StorageCollection) error {
	err := sm.primary.SaveCollection(collection)
	if err != nil && sm.fallback != nil {
		log.Printf("Primary storage failed, using fallback: %v", err)
		return sm.fallback.SaveCollection(collection)
	}
	return err
}

func (sm *StorageManager) GetCollection(id string) (*StorageCollection, error) {
	result, err := sm.primary.GetCollection(id)
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetCollection(id)
	}
	return result, err
}

func (sm *StorageManager) GetCollectionsByRoom(roomID string) ([]*StorageCollection, error) {
	result, err := sm.primary.GetCollectionsByRoom(roomID)
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetCollectionsByRoom(roomID)
	}
	return result, err
}

func (sm *StorageManager) GetAllCollections() ([]*StorageCollection, error) {
	result, err := sm.primary.GetAllCollections()
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetAllCollections()
	}
	return result, err
}

func (sm *StorageManager) DeleteCollection(id string) error {
	err := sm.primary.DeleteCollection(id)
	if err != nil && sm.fallback != nil {
		return sm.fallback.DeleteCollection(id)
	}
	return err
}

func (sm *StorageManager) SaveEnvironment(env *StorageEnvironment) error {
	err := sm.primary.SaveEnvironment(env)
	if err != nil && sm.fallback != nil {
		return sm.fallback.SaveEnvironment(env)
	}
	return err
}

func (sm *StorageManager) GetEnvironment(id string) (*StorageEnvironment, error) {
	result, err := sm.primary.GetEnvironment(id)
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetEnvironment(id)
	}
	return result, err
}

func (sm *StorageManager) GetEnvironmentsByRoom(roomID string) ([]*StorageEnvironment, error) {
	result, err := sm.primary.GetEnvironmentsByRoom(roomID)
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetEnvironmentsByRoom(roomID)
	}
	return result, err
}

func (sm *StorageManager) GetAllEnvironments() ([]*StorageEnvironment, error) {
	result, err := sm.primary.GetAllEnvironments()
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetAllEnvironments()
	}
	return result, err
}

func (sm *StorageManager) DeleteEnvironment(id string) error {
	err := sm.primary.DeleteEnvironment(id)
	if err != nil && sm.fallback != nil {
		return sm.fallback.DeleteEnvironment(id)
	}
	return err
}

func (sm *StorageManager) SaveRequest(req *StorageRequest) error {
	err := sm.primary.SaveRequest(req)
	if err != nil && sm.fallback != nil {
		return sm.fallback.SaveRequest(req)
	}
	return err
}

func (sm *StorageManager) GetRequest(id string) (*StorageRequest, error) {
	result, err := sm.primary.GetRequest(id)
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetRequest(id)
	}
	return result, err
}

func (sm *StorageManager) GetRequestsByCollection(collectionID string) ([]*StorageRequest, error) {
	result, err := sm.primary.GetRequestsByCollection(collectionID)
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetRequestsByCollection(collectionID)
	}
	return result, err
}

func (sm *StorageManager) GetAllRequests() ([]*StorageRequest, error) {
	result, err := sm.primary.GetAllRequests()
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetAllRequests()
	}
	return result, err
}

func (sm *StorageManager) DeleteRequest(id string) error {
	err := sm.primary.DeleteRequest(id)
	if err != nil && sm.fallback != nil {
		return sm.fallback.DeleteRequest(id)
	}
	return err
}

func (sm *StorageManager) SaveHistory(history *StorageHistory) error {
	err := sm.primary.SaveHistory(history)
	if err != nil && sm.fallback != nil {
		return sm.fallback.SaveHistory(history)
	}
	return err
}

func (sm *StorageManager) GetHistoryByRoom(roomID string, limit int) ([]*StorageHistory, error) {
	result, err := sm.primary.GetHistoryByRoom(roomID, limit)
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetHistoryByRoom(roomID, limit)
	}
	return result, err
}

func (sm *StorageManager) GetAllHistory(limit int) ([]*StorageHistory, error) {
	result, err := sm.primary.GetAllHistory(limit)
	if err != nil && sm.fallback != nil {
		return sm.fallback.GetAllHistory(limit)
	}
	return result, err
}

func (sm *StorageManager) DeleteHistory(id string) error {
	err := sm.primary.DeleteHistory(id)
	if err != nil && sm.fallback != nil {
		return sm.fallback.DeleteHistory(id)
	}
	return err
}

func (sm *StorageManager) Ping() error {
	return sm.primary.Ping()
}

func (sm *StorageManager) Close() error {
	var errs []error
	if err := sm.primary.Close(); err != nil {
		errs = append(errs, err)
	}
	if sm.fallback != nil {
		if err := sm.fallback.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.New("failed to close storage backends")
	}
	return nil
}

// Global storage manager instance
var StorageMgr *StorageManager

// InitStorage initializes the storage manager
func InitStorage() error {
	var err error
	StorageMgr, err = NewStorageManager()
	return err
}
