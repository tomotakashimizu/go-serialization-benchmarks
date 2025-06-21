package serializers

import (
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
	pb "github.com/tomotakashimizu/go-serialization-benchmarks/proto"
)

// ProtobufSerializer implements Serializer interface for protobuf
type ProtobufSerializer struct{}

// NewProtobufSerializer creates a new ProtobufSerializer
func NewProtobufSerializer() *ProtobufSerializer {
	return &ProtobufSerializer{}
}

// Name returns the name of the serializer
func (p *ProtobufSerializer) Name() string {
	return "Protobuf"
}

// Marshal serializes a User to Protocol Buffer bytes
func (p *ProtobufSerializer) Marshal(user models.User) ([]byte, error) {
	pbUser, err := p.convertUserToProto(user)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(pbUser)
}

// Unmarshal deserializes Protocol Buffer bytes to a User
func (p *ProtobufSerializer) Unmarshal(data []byte) (models.User, error) {
	var pbUser pb.User
	if err := proto.Unmarshal(data, &pbUser); err != nil {
		return models.User{}, err
	}
	return p.convertUserFromProto(&pbUser)
}

// MarshalUsers serializes a collection of Users to Protocol Buffer bytes
func (p *ProtobufSerializer) MarshalUsers(users models.Users) ([]byte, error) {
	pbUserList := &pb.UserList{
		Users: make([]*pb.User, len(users)),
	}

	for i, user := range users {
		pbUser, err := p.convertUserToProto(user)
		if err != nil {
			return nil, err
		}
		pbUserList.Users[i] = pbUser
	}

	return proto.Marshal(pbUserList)
}

// UnmarshalUsers deserializes Protocol Buffer bytes to a collection of Users
func (p *ProtobufSerializer) UnmarshalUsers(data []byte) (models.Users, error) {
	var pbUserList pb.UserList
	if err := proto.Unmarshal(data, &pbUserList); err != nil {
		return nil, err
	}

	users := make(models.Users, len(pbUserList.Users))
	for i, pbUser := range pbUserList.Users {
		user, err := p.convertUserFromProto(pbUser)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return users, nil
}

// convertUserToProto converts models.User to pb.User
func (p *ProtobufSerializer) convertUserToProto(user models.User) (*pb.User, error) {
	// Handle metadata by converting to JSON strings
	metadata := make(map[string]string)
	for k, v := range user.Metadata {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata value: %w", err)
		}
		metadata[k] = string(jsonBytes)
	}

	pbUser := &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       int32(user.Age),
		IsActive:  user.IsActive,
		Tags:      user.Tags,
		Metadata:  metadata,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	// Convert profile
	if !p.isEmptyProfile(user.Profile) {
		pbProfile, err := p.convertProfileToProto(user.Profile)
		if err != nil {
			return nil, err
		}
		pbUser.Profile = pbProfile
	}

	// Convert settings
	if !p.isEmptySettings(user.Settings) {
		pbSettings, err := p.convertSettingsToProto(user.Settings)
		if err != nil {
			return nil, err
		}
		pbUser.Settings = pbSettings
	}

	return pbUser, nil
}

// convertUserFromProto converts pb.User to models.User
func (p *ProtobufSerializer) convertUserFromProto(pbUser *pb.User) (models.User, error) {
	var createdAt time.Time
	if pbUser.CreatedAt != nil {
		createdAt = pbUser.CreatedAt.AsTime()
	}

	// Handle metadata by converting from JSON strings
	metadata := make(map[string]interface{})
	for k, v := range pbUser.Metadata {
		var value interface{}
		if err := json.Unmarshal([]byte(v), &value); err != nil {
			return models.User{}, fmt.Errorf("failed to unmarshal metadata value: %w", err)
		}
		metadata[k] = value
	}

	user := models.User{
		ID:        pbUser.Id,
		Name:      pbUser.Name,
		Email:     pbUser.Email,
		Age:       int(pbUser.Age),
		IsActive:  pbUser.IsActive,
		Tags:      pbUser.Tags,
		Metadata:  metadata,
		CreatedAt: createdAt,
	}

	// Convert profile
	if pbUser.Profile != nil {
		profile, err := p.convertProfileFromProto(pbUser.Profile)
		if err != nil {
			return models.User{}, err
		}
		user.Profile = profile
	}

	// Convert settings
	if pbUser.Settings != nil {
		settings, err := p.convertSettingsFromProto(pbUser.Settings)
		if err != nil {
			return models.User{}, err
		}
		user.Settings = settings
	}

	return user, nil
}

// convertProfileToProto converts models.Profile to pb.Profile
func (p *ProtobufSerializer) convertProfileToProto(profile models.Profile) (*pb.Profile, error) {
	pbProfile := &pb.Profile{
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Bio:       profile.Bio,
		Avatar:    profile.Avatar,
	}

	// Convert social links
	if len(profile.SocialLinks) > 0 {
		pbProfile.SocialLinks = make([]*pb.Link, len(profile.SocialLinks))
		for i, link := range profile.SocialLinks {
			pbProfile.SocialLinks[i] = &pb.Link{
				Platform: link.Platform,
				Url:      link.URL,
			}
		}
	}

	// Convert preferences
	if !p.isEmptyPreferences(profile.Preferences) {
		pbPreferences, err := p.convertPreferencesToProto(profile.Preferences)
		if err != nil {
			return nil, err
		}
		pbProfile.Preferences = pbPreferences
	}

	return pbProfile, nil
}

// convertProfileFromProto converts pb.Profile to models.Profile
func (p *ProtobufSerializer) convertProfileFromProto(pbProfile *pb.Profile) (models.Profile, error) {
	profile := models.Profile{
		FirstName: pbProfile.FirstName,
		LastName:  pbProfile.LastName,
		Bio:       pbProfile.Bio,
		Avatar:    pbProfile.Avatar,
	}

	// Convert social links
	if len(pbProfile.SocialLinks) > 0 {
		profile.SocialLinks = make([]models.Link, len(pbProfile.SocialLinks))
		for i, pbLink := range pbProfile.SocialLinks {
			profile.SocialLinks[i] = models.Link{
				Platform: pbLink.Platform,
				URL:      pbLink.Url,
			}
		}
	}

	// Convert preferences
	if pbProfile.Preferences != nil {
		preferences, err := p.convertPreferencesFromProto(pbProfile.Preferences)
		if err != nil {
			return models.Profile{}, err
		}
		profile.Preferences = preferences
	}

	return profile, nil
}

// convertPreferencesToProto converts models.Preferences to pb.Preferences
func (p *ProtobufSerializer) convertPreferencesToProto(preferences models.Preferences) (*pb.Preferences, error) {
	pbPreferences := &pb.Preferences{
		Theme:         preferences.Theme,
		Language:      preferences.Language,
		Notifications: preferences.Notifications,
	}

	// Convert privacy settings
	if preferences.Privacy != (models.PrivacySettings{}) {
		pbPreferences.Privacy = &pb.PrivacySettings{
			ProfilePublic: preferences.Privacy.ProfilePublic,
			EmailVisible:  preferences.Privacy.EmailVisible,
			ShowActivity:  preferences.Privacy.ShowActivity,
		}
	}

	return pbPreferences, nil
}

// convertPreferencesFromProto converts pb.Preferences to models.Preferences
func (p *ProtobufSerializer) convertPreferencesFromProto(pbPreferences *pb.Preferences) (models.Preferences, error) {
	preferences := models.Preferences{
		Theme:         pbPreferences.Theme,
		Language:      pbPreferences.Language,
		Notifications: pbPreferences.Notifications,
	}

	// Convert privacy settings
	if pbPreferences.Privacy != nil {
		preferences.Privacy = models.PrivacySettings{
			ProfilePublic: pbPreferences.Privacy.ProfilePublic,
			EmailVisible:  pbPreferences.Privacy.EmailVisible,
			ShowActivity:  pbPreferences.Privacy.ShowActivity,
		}
	}

	return preferences, nil
}

// convertSettingsToProto converts models.Settings to pb.Settings
func (p *ProtobufSerializer) convertSettingsToProto(settings models.Settings) (*pb.Settings, error) {
	// Convert int to int32 for limits
	limits := make(map[string]int32)
	for k, v := range settings.Limits {
		limits[k] = int32(v)
	}

	return &pb.Settings{
		Language: settings.Language,
		Timezone: settings.TimeZone,
		Features: settings.Features,
		Limits:   limits,
	}, nil
}

// convertSettingsFromProto converts pb.Settings to models.Settings
func (p *ProtobufSerializer) convertSettingsFromProto(pbSettings *pb.Settings) (models.Settings, error) {
	// Convert int32 to int for limits
	limits := make(map[string]int)
	for k, v := range pbSettings.Limits {
		limits[k] = int(v)
	}

	return models.Settings{
		Language: pbSettings.Language,
		TimeZone: pbSettings.Timezone,
		Features: pbSettings.Features,
		Limits:   limits,
	}, nil
}

// Helper methods to check if structs are empty

// isEmptyProfile checks if a Profile struct is empty
func (p *ProtobufSerializer) isEmptyProfile(profile models.Profile) bool {
	return profile.FirstName == "" && profile.LastName == "" && profile.Bio == "" && profile.Avatar == "" &&
		len(profile.SocialLinks) == 0 && p.isEmptyPreferences(profile.Preferences)
}

// isEmptySettings checks if a Settings struct is empty
func (p *ProtobufSerializer) isEmptySettings(settings models.Settings) bool {
	return settings.Language == "" && settings.TimeZone == "" &&
		len(settings.Features) == 0 && len(settings.Limits) == 0
}

// isEmptyPreferences checks if a Preferences struct is empty
func (p *ProtobufSerializer) isEmptyPreferences(preferences models.Preferences) bool {
	return preferences.Theme == "" && preferences.Language == "" &&
		len(preferences.Notifications) == 0 && p.isEmptyPrivacySettings(preferences.Privacy)
}

// isEmptyPrivacySettings checks if a PrivacySettings struct is empty
func (p *ProtobufSerializer) isEmptyPrivacySettings(privacy models.PrivacySettings) bool {
	return !privacy.ProfilePublic && !privacy.EmailVisible && !privacy.ShowActivity
}
