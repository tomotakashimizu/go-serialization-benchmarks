package models

import (
	"fmt"
	"math/rand"
	"time"
)

// Users represents a collection of User objects
type Users []User

// User represents a user with nested structures
type User struct {
	ID        int64                  `json:"id" msgpack:"id" cbor:"id"`
	Name      string                 `json:"name" msgpack:"name" cbor:"name"`
	Email     string                 `json:"email" msgpack:"email" cbor:"email"`
	Age       int                    `json:"age" msgpack:"age" cbor:"age"`
	IsActive  bool                   `json:"is_active" msgpack:"is_active" cbor:"is_active"`
	Profile   Profile                `json:"profile" msgpack:"profile" cbor:"profile"`
	Settings  Settings               `json:"settings" msgpack:"settings" cbor:"settings"`
	Tags      []string               `json:"tags" msgpack:"tags" cbor:"tags"`
	Metadata  map[string]interface{} `json:"metadata" msgpack:"metadata" cbor:"metadata"`
	CreatedAt time.Time              `json:"created_at" msgpack:"created_at" cbor:"created_at"`
}

// Profile represents user profile information (2nd layer)
type Profile struct {
	FirstName   string      `json:"first_name" msgpack:"first_name" cbor:"first_name"`
	LastName    string      `json:"last_name" msgpack:"last_name" cbor:"last_name"`
	Bio         string      `json:"bio" msgpack:"bio" cbor:"bio"`
	Avatar      string      `json:"avatar" msgpack:"avatar" cbor:"avatar"`
	SocialLinks []Link      `json:"social_links" msgpack:"social_links" cbor:"social_links"`
	Preferences Preferences `json:"preferences" msgpack:"preferences" cbor:"preferences"`
}

// Link represents a social media link (3rd layer)
type Link struct {
	Platform string `json:"platform" msgpack:"platform" cbor:"platform"`
	URL      string `json:"url" msgpack:"url" cbor:"url"`
}

// Preferences represents user preferences (3rd layer)
type Preferences struct {
	Theme         string          `json:"theme" msgpack:"theme" cbor:"theme"`
	Language      string          `json:"language" msgpack:"language" cbor:"language"`
	Notifications map[string]bool `json:"notifications" msgpack:"notifications" cbor:"notifications"`
	Privacy       PrivacySettings `json:"privacy" msgpack:"privacy" cbor:"privacy"`
}

// PrivacySettings represents privacy settings (4th layer for deeper nesting)
type PrivacySettings struct {
	ProfilePublic bool `json:"profile_public" msgpack:"profile_public" cbor:"profile_public"`
	EmailVisible  bool `json:"email_visible" msgpack:"email_visible" cbor:"email_visible"`
	ShowActivity  bool `json:"show_activity" msgpack:"show_activity" cbor:"show_activity"`
}

// Settings represents user application settings (2nd layer)
type Settings struct {
	Language string         `json:"language" msgpack:"language" cbor:"language"`
	TimeZone string         `json:"timezone" msgpack:"timezone" cbor:"timezone"`
	Features []string       `json:"features" msgpack:"features" cbor:"features"`
	Limits   map[string]int `json:"limits" msgpack:"limits" cbor:"limits"`
}

// GenerateTestUsers generates a specified number of test users
func GenerateTestUsers(count int) Users {
	users := make(Users, count)

	platforms := []string{"Twitter", "GitHub", "LinkedIn", "Instagram", "Facebook"}
	themes := []string{"dark", "light", "auto", "contrast"}
	languages := []string{"en", "ja", "fr", "de", "es", "zh", "ko"}
	timezones := []string{"UTC", "JST", "PST", "EST", "CET", "IST", "CST"}
	features := []string{"premium", "beta", "notifications", "analytics", "export", "api", "integration"}

	for i := 0; i < count; i++ {
		user := User{
			ID:        int64(i + 1),
			Name:      fmt.Sprintf("User%d", i+1),
			Email:     fmt.Sprintf("user%d@example.com", i+1),
			Age:       rand.Intn(60) + 18,
			IsActive:  rand.Float32() > 0.2, // 80% active
			CreatedAt: time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour),
		}

		// Generate profile
		user.Profile = Profile{
			FirstName: fmt.Sprintf("First%d", i+1),
			LastName:  fmt.Sprintf("Last%d", i+1),
			Bio:       fmt.Sprintf("This is a bio for user %d with some additional details", i+1),
			Avatar:    fmt.Sprintf("https://example.com/avatars/user%d.jpg", i+1),
		}

		// Generate social links (0-4 links)
		linkCount := rand.Intn(5)
		user.Profile.SocialLinks = make([]Link, linkCount)
		for j := 0; j < linkCount; j++ {
			platform := platforms[rand.Intn(len(platforms))]
			user.Profile.SocialLinks[j] = Link{
				Platform: platform,
				URL:      fmt.Sprintf("https://%s.com/user%d", platform, i+1),
			}
		}

		// Generate preferences
		user.Profile.Preferences = Preferences{
			Theme:    themes[rand.Intn(len(themes))],
			Language: languages[rand.Intn(len(languages))],
			Notifications: map[string]bool{
				"email":   rand.Float32() > 0.5,
				"push":    rand.Float32() > 0.5,
				"sms":     rand.Float32() > 0.3,
				"desktop": rand.Float32() > 0.4,
				"weekly":  rand.Float32() > 0.6,
			},
			Privacy: PrivacySettings{
				ProfilePublic: rand.Float32() > 0.3,
				EmailVisible:  rand.Float32() > 0.7,
				ShowActivity:  rand.Float32() > 0.4,
			},
		}

		// Generate settings
		user.Settings = Settings{
			Language: languages[rand.Intn(len(languages))],
			TimeZone: timezones[rand.Intn(len(timezones))],
			Features: generateRandomFeatures(features, rand.Intn(4)+1),
			Limits: map[string]int{
				"api_calls":    rand.Intn(1000) + 100,
				"storage_mb":   rand.Intn(1000) + 100,
				"connections":  rand.Intn(50) + 10,
				"bandwidth_mb": rand.Intn(500) + 50,
			},
		}

		// Generate tags (0-6 tags)
		tagCount := rand.Intn(7)
		user.Tags = make([]string, tagCount)
		for j := 0; j < tagCount; j++ {
			user.Tags[j] = fmt.Sprintf("tag%d", rand.Intn(30)+1)
		}

		// Generate metadata (0-4 entries)
		metadataCount := rand.Intn(5)
		user.Metadata = make(map[string]interface{})
		for j := 0; j < metadataCount; j++ {
			key := fmt.Sprintf("meta%d", j+1)
			switch rand.Intn(4) {
			case 0:
				user.Metadata[key] = fmt.Sprintf("value%d", j+1)
			case 1:
				user.Metadata[key] = rand.Intn(1000)
			case 2:
				user.Metadata[key] = rand.Float32() > 0.5
			case 3:
				user.Metadata[key] = rand.Float64() * 100
			}
		}

		users[i] = user
	}

	return users
}

func generateRandomFeatures(features []string, count int) []string {
	if count > len(features) {
		count = len(features)
	}

	selected := make([]string, 0, count)
	used := make(map[int]bool)

	for len(selected) < count {
		idx := rand.Intn(len(features))
		if !used[idx] {
			selected = append(selected, features[idx])
			used[idx] = true
		}
	}

	return selected
}
