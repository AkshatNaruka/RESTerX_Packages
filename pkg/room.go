package pkg

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"
)

type RoomService struct{}

type CreateRoomRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Description string `json:"description"`
	RoomType    string `json:"roomType"` // free, premium
	DeviceLimit int    `json:"deviceLimit"`
}

type JoinRoomRequest struct {
	RoomID   string `json:"roomId"`
	Password string `json:"password"`
	DeviceID string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
}

type RoomResponse struct {
	Room     *Room   `json:"room"`
	Token    string  `json:"token"`
	RoomID   string  `json:"roomId"`
}

type RoomDataResponse struct {
	Collections  []DBCollection  `json:"collections"`
	Environments []DBEnvironment `json:"environments"`
	Requests     []DBRequest     `json:"requests"`
}

func NewRoomService() *RoomService {
	return &RoomService{}
}

// generateRoomID generates a unique room identifier
func (rs *RoomService) generateRoomID() (string, error) {
	bytes := make([]byte, 6) // 12 character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HashPassword hashes a plain text password
func (rs *RoomService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword compares a hashed password with a plain text password
func (rs *RoomService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CreateRoom creates a new room with a unique room ID
func (rs *RoomService) CreateRoom(req CreateRoomRequest) (*RoomResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("room name is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Generate unique room ID
	roomID, err := rs.generateRoomID()
	if err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := rs.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Set defaults
	deviceLimit := req.DeviceLimit
	if deviceLimit == 0 {
		deviceLimit = 3 // Default to 3 devices for free tier
	}
	
	roomType := req.RoomType
	if roomType == "" {
		roomType = "free"
	}

	// Create room
	room := Room{
		RoomID:      roomID,
		Name:        req.Name,
		Password:    hashedPassword,
		Description: req.Description,
		RoomType:    roomType,
		DeviceLimit: deviceLimit,
		IsActive:    true,
	}

	if err := DB.Create(&room).Error; err != nil {
		return nil, err
	}

	// Generate token for the room
	token, err := rs.GenerateRoomToken(room.ID, roomID)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	room.Password = ""

	return &RoomResponse{
		Room:   &room,
		Token:  token,
		RoomID: roomID,
	}, nil
}

// JoinRoom allows a device to join an existing room
func (rs *RoomService) JoinRoom(req JoinRoomRequest) (*RoomResponse, error) {
	// Validate input
	if req.RoomID == "" {
		return nil, errors.New("room ID is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}
	if req.DeviceID == "" {
		return nil, errors.New("device ID is required")
	}

	// Find room by room ID
	var room Room
	if err := DB.Where("room_id = ?", req.RoomID).First(&room).Error; err != nil {
		return nil, errors.New("room not found")
	}

	if !room.IsActive {
		return nil, errors.New("room is not active")
	}

	// Check password
	if err := rs.CheckPassword(room.Password, req.Password); err != nil {
		return nil, errors.New("invalid password")
	}

	// Check device limit
	var deviceCount int64
	DB.Model(&RoomDevice{}).Where("room_id = ?", room.ID).Count(&deviceCount)
	
	// Check if device already exists
	var existingDevice RoomDevice
	deviceExists := DB.Where("room_id = ? AND device_id = ?", room.ID, req.DeviceID).First(&existingDevice).Error == nil
	
	if !deviceExists && deviceCount >= int64(room.DeviceLimit) {
		return nil, errors.New("device limit reached for this room")
	}

	// Add or update device
	if deviceExists {
		// Update last active time
		DB.Model(&existingDevice).Update("last_active", time.Now())
		if req.DeviceName != "" {
			DB.Model(&existingDevice).Update("device_name", req.DeviceName)
		}
	} else {
		// Add new device
		deviceName := req.DeviceName
		if deviceName == "" {
			deviceName = "Device " + req.DeviceID[:8]
		}
		
		device := RoomDevice{
			RoomID:     room.ID,
			DeviceID:   req.DeviceID,
			DeviceName: deviceName,
			LastActive: time.Now(),
		}
		
		if err := DB.Create(&device).Error; err != nil {
			return nil, err
		}
	}

	// Generate token for the room
	token, err := rs.GenerateRoomToken(room.ID, room.RoomID)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	room.Password = ""

	return &RoomResponse{
		Room:   &room,
		Token:  token,
		RoomID: room.RoomID,
	}, nil
}

// GenerateRoomToken generates a JWT token for room access
func (rs *RoomService) GenerateRoomToken(roomDBID uint, roomID string) (string, error) {
	authService := NewAuthService()
	
	// Create a pseudo-user for the room (using room's ID as user ID)
	// This allows reusing existing JWT infrastructure
	pseudoUser := &User{
		ID:       roomDBID + 1000000, // Offset to avoid conflicts with real users
		Username: "room_" + roomID,
		Email:    roomID + "@room.resterx.com",
		Role:     "room",
	}
	
	return authService.GenerateToken(pseudoUser, roomDBID)
}

// ValidateRoomToken validates a room token and returns room info
func (rs *RoomService) ValidateRoomToken(tokenString string) (*Room, error) {
	authService := NewAuthService()
	claims, err := authService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Role != "room" {
		return nil, errors.New("invalid room token")
	}

	// Extract room ID from workspace ID (which stores room DB ID)
	roomDBID := claims.WorkspaceID
	
	var room Room
	if err := DB.First(&room, roomDBID).Error; err != nil {
		return nil, errors.New("room not found")
	}

	return &room, nil
}

// GetRoomData retrieves all data for a room
func (rs *RoomService) GetRoomData(roomID uint) (*RoomDataResponse, error) {
	var collections []DBCollection
	if err := DB.Preload("Requests").Where("room_id = ?", roomID).Find(&collections).Error; err != nil {
		return nil, err
	}

	var environments []DBEnvironment
	if err := DB.Where("room_id = ?", roomID).Find(&environments).Error; err != nil {
		return nil, err
	}

	// Get all requests for the room's collections
	var requests []DBRequest
	for _, collection := range collections {
		requests = append(requests, collection.Requests...)
	}

	return &RoomDataResponse{
		Collections:  collections,
		Environments: environments,
		Requests:     requests,
	}, nil
}

// SyncRoomData updates room data (collections, environments, requests)
func (rs *RoomService) SyncRoomData(roomID uint, data RoomDataResponse) error {
	tx := DB.Begin()

	// Update or create collections
	for _, collection := range data.Collections {
		collection.RoomID = roomID
		collection.WorkspaceID = 0 // Clear workspace ID for room-based collections
		
		if collection.ID > 0 {
			// Update existing
			if err := tx.Save(&collection).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Create new
			if err := tx.Create(&collection).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Update or create environments
	for _, env := range data.Environments {
		env.RoomID = roomID
		env.WorkspaceID = 0 // Clear workspace ID for room-based environments
		
		if env.ID > 0 {
			// Update existing
			if err := tx.Save(&env).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Create new
			if err := tx.Create(&env).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Update or create requests
	for _, request := range data.Requests {
		if request.ID > 0 {
			// Update existing
			if err := tx.Save(&request).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Create new
			if err := tx.Create(&request).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// GetRoomDevices returns all devices connected to a room
func (rs *RoomService) GetRoomDevices(roomID uint) ([]RoomDevice, error) {
	var devices []RoomDevice
	if err := DB.Where("room_id = ?", roomID).Order("last_active DESC").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// RemoveDevice removes a device from a room
func (rs *RoomService) RemoveDevice(roomID uint, deviceID string) error {
	return DB.Where("room_id = ? AND device_id = ?", roomID, deviceID).Delete(&RoomDevice{}).Error
}

// GetRoomByRoomID retrieves a room by its room ID
func (rs *RoomService) GetRoomByRoomID(roomID string) (*Room, error) {
	var room Room
	if err := DB.Where("room_id = ?", roomID).First(&room).Error; err != nil {
		return nil, errors.New("room not found")
	}
	room.Password = "" // Clear password
	return &room, nil
}

// DeleteRoom deletes a room and all its data
func (rs *RoomService) DeleteRoom(roomID uint, password string) error {
	// Get room
	var room Room
	if err := DB.First(&room, roomID).Error; err != nil {
		return errors.New("room not found")
	}

	// Verify password
	if err := rs.CheckPassword(room.Password, password); err != nil {
		return errors.New("invalid password")
	}

	tx := DB.Begin()

	// Delete all devices
	if err := tx.Where("room_id = ?", roomID).Delete(&RoomDevice{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete all collections and their requests
	if err := tx.Where("room_id = ?", roomID).Delete(&DBCollection{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete all environments
	if err := tx.Where("room_id = ?", roomID).Delete(&DBEnvironment{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete room
	if err := tx.Delete(&room).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
