package serializers

import (
	"fmt"
	"time"

	flatbuffers "github.com/google/flatbuffers/go"

	flatbuffers_generated "github.com/tomotakashimizu/go-serialization-benchmarks/go_serialization_benchmarks/flatbuffers"
	"github.com/tomotakashimizu/go-serialization-benchmarks/internal/models"
)

// FlatBuffersSerializer implements Serializer interface for FlatBuffers
type FlatBuffersSerializer struct{}

// NewFlatBuffersSerializer creates a new FlatBuffersSerializer
func NewFlatBuffersSerializer() *FlatBuffersSerializer {
	return &FlatBuffersSerializer{}
}

// Name returns the name of the serializer
func (f *FlatBuffersSerializer) Name() string {
	return "FlatBuffers"
}

// Marshal serializes a User to FlatBuffers bytes
func (f *FlatBuffersSerializer) Marshal(user models.User) ([]byte, error) {
	builder := flatbuffers.NewBuilder(1024)

	// Create UserList with single user
	userOffset, err := f.convertUserToFlatBuffer(builder, user)
	if err != nil {
		return nil, err
	}

	// Create UserList
	flatbuffers_generated.UserListStartUsersVector(builder, 1)
	builder.PrependUOffsetT(userOffset)
	usersVector := builder.EndVector(1)

	flatbuffers_generated.UserListStart(builder)
	flatbuffers_generated.UserListAddUsers(builder, usersVector)
	userList := flatbuffers_generated.UserListEnd(builder)

	builder.Finish(userList)
	return builder.FinishedBytes(), nil
}

// Unmarshal deserializes FlatBuffers bytes to a User
func (f *FlatBuffersSerializer) Unmarshal(data []byte) (models.User, error) {
	userList := flatbuffers_generated.GetRootAsUserList(data, 0)

	if userList.UsersLength() == 0 {
		return models.User{}, fmt.Errorf("no users in flatbuffer data")
	}

	fbUser := new(flatbuffers_generated.User)
	if !userList.Users(fbUser, 0) {
		return models.User{}, fmt.Errorf("failed to get user from flatbuffer")
	}

	return f.convertFlatBufferToUser(fbUser)
}

// MarshalUsers serializes a collection of Users to FlatBuffers bytes
func (f *FlatBuffersSerializer) MarshalUsers(users models.Users) ([]byte, error) {
	builder := flatbuffers.NewBuilder(1024 * len(users))

	// Convert all users to FlatBuffer objects
	userOffsets := make([]flatbuffers.UOffsetT, len(users))
	for i, user := range users {
		offset, err := f.convertUserToFlatBuffer(builder, user)
		if err != nil {
			return nil, err
		}
		userOffsets[i] = offset
	}

	// Create UserList
	flatbuffers_generated.UserListStartUsersVector(builder, len(userOffsets))
	for i := len(userOffsets) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(userOffsets[i])
	}
	usersVector := builder.EndVector(len(userOffsets))

	flatbuffers_generated.UserListStart(builder)
	flatbuffers_generated.UserListAddUsers(builder, usersVector)
	userList := flatbuffers_generated.UserListEnd(builder)

	builder.Finish(userList)
	return builder.FinishedBytes(), nil
}

// UnmarshalUsers deserializes FlatBuffers bytes to a collection of Users
func (f *FlatBuffersSerializer) UnmarshalUsers(data []byte) (models.Users, error) {
	userList := flatbuffers_generated.GetRootAsUserList(data, 0)

	users := make(models.Users, userList.UsersLength())
	fbUser := new(flatbuffers_generated.User)

	for i := 0; i < userList.UsersLength(); i++ {
		if !userList.Users(fbUser, i) {
			return nil, fmt.Errorf("failed to get user %d from flatbuffer", i)
		}

		user, err := f.convertFlatBufferToUser(fbUser)
		if err != nil {
			return nil, err
		}
		users[i] = user
	}

	return users, nil
}

// convertUserToFlatBuffer converts a models.User to FlatBuffer format
func (f *FlatBuffersSerializer) convertUserToFlatBuffer(builder *flatbuffers.Builder, user models.User) (flatbuffers.UOffsetT, error) {
	// Create all nested objects first (deepest first)

	// Convert Metadata first
	var metadataVector flatbuffers.UOffsetT
	if len(user.Metadata) > 0 {
		metadataOffsets := make([]flatbuffers.UOffsetT, 0, len(user.Metadata))
		for key, value := range user.Metadata {
			metadataOffset := f.convertMetadataEntryToFlatBuffer(builder, key, value)
			metadataOffsets = append(metadataOffsets, metadataOffset)
		}
		flatbuffers_generated.UserStartMetadataVector(builder, len(metadataOffsets))
		for i := len(metadataOffsets) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(metadataOffsets[i])
		}
		metadataVector = builder.EndVector(len(metadataOffsets))
	}

	// Convert Tags
	var tagsVector flatbuffers.UOffsetT
	if len(user.Tags) > 0 {
		tagOffsets := make([]flatbuffers.UOffsetT, len(user.Tags))
		for i, tag := range user.Tags {
			tagOffsets[i] = builder.CreateString(tag)
		}
		flatbuffers_generated.UserStartTagsVector(builder, len(tagOffsets))
		for i := len(tagOffsets) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(tagOffsets[i])
		}
		tagsVector = builder.EndVector(len(tagOffsets))
	}

	// Convert Profile
	profileOffset, err := f.convertProfileToFlatBuffer(builder, user.Profile)
	if err != nil {
		return 0, err
	}

	// Convert Settings
	settingsOffset, err := f.convertSettingsToFlatBuffer(builder, user.Settings)
	if err != nil {
		return 0, err
	}

	// Convert strings last
	nameOffset := builder.CreateString(user.Name)
	emailOffset := builder.CreateString(user.Email)

	// Create User
	flatbuffers_generated.UserStart(builder)
	flatbuffers_generated.UserAddId(builder, user.ID)
	flatbuffers_generated.UserAddName(builder, nameOffset)
	flatbuffers_generated.UserAddEmail(builder, emailOffset)
	flatbuffers_generated.UserAddAge(builder, int32(user.Age))
	flatbuffers_generated.UserAddIsActive(builder, user.IsActive)
	flatbuffers_generated.UserAddProfile(builder, profileOffset)
	flatbuffers_generated.UserAddSettings(builder, settingsOffset)
	if len(user.Tags) > 0 {
		flatbuffers_generated.UserAddTags(builder, tagsVector)
	}
	if len(user.Metadata) > 0 {
		flatbuffers_generated.UserAddMetadata(builder, metadataVector)
	}
	flatbuffers_generated.UserAddCreatedAt(builder, user.CreatedAt.UnixNano())

	return flatbuffers_generated.UserEnd(builder), nil
}

// convertProfileToFlatBuffer converts a models.Profile to FlatBuffer format
func (f *FlatBuffersSerializer) convertProfileToFlatBuffer(builder *flatbuffers.Builder, profile models.Profile) (flatbuffers.UOffsetT, error) {
	// Convert SocialLinks first
	var socialLinksVector flatbuffers.UOffsetT
	if len(profile.SocialLinks) > 0 {
		linkOffsets := make([]flatbuffers.UOffsetT, len(profile.SocialLinks))
		for i, link := range profile.SocialLinks {
			platformOffset := builder.CreateString(link.Platform)
			urlOffset := builder.CreateString(link.URL)

			flatbuffers_generated.LinkStart(builder)
			flatbuffers_generated.LinkAddPlatform(builder, platformOffset)
			flatbuffers_generated.LinkAddUrl(builder, urlOffset)
			linkOffsets[i] = flatbuffers_generated.LinkEnd(builder)
		}
		flatbuffers_generated.ProfileStartSocialLinksVector(builder, len(linkOffsets))
		for i := len(linkOffsets) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(linkOffsets[i])
		}
		socialLinksVector = builder.EndVector(len(linkOffsets))
	}

	// Convert Preferences
	preferencesOffset, err := f.convertPreferencesToFlatBuffer(builder, profile.Preferences)
	if err != nil {
		return 0, err
	}

	// Create strings last
	firstNameOffset := builder.CreateString(profile.FirstName)
	lastNameOffset := builder.CreateString(profile.LastName)
	bioOffset := builder.CreateString(profile.Bio)
	avatarOffset := builder.CreateString(profile.Avatar)

	flatbuffers_generated.ProfileStart(builder)
	flatbuffers_generated.ProfileAddFirstName(builder, firstNameOffset)
	flatbuffers_generated.ProfileAddLastName(builder, lastNameOffset)
	flatbuffers_generated.ProfileAddBio(builder, bioOffset)
	flatbuffers_generated.ProfileAddAvatar(builder, avatarOffset)
	if len(profile.SocialLinks) > 0 {
		flatbuffers_generated.ProfileAddSocialLinks(builder, socialLinksVector)
	}
	flatbuffers_generated.ProfileAddPreferences(builder, preferencesOffset)

	return flatbuffers_generated.ProfileEnd(builder), nil
}

// convertPreferencesToFlatBuffer converts a models.Preferences to FlatBuffer format
func (f *FlatBuffersSerializer) convertPreferencesToFlatBuffer(builder *flatbuffers.Builder, prefs models.Preferences) (flatbuffers.UOffsetT, error) {
	// Convert Notifications map first
	var notificationsVector flatbuffers.UOffsetT
	if len(prefs.Notifications) > 0 {
		notificationOffsets := make([]flatbuffers.UOffsetT, 0, len(prefs.Notifications))
		for key, value := range prefs.Notifications {
			keyOffset := builder.CreateString(key)

			flatbuffers_generated.NotificationSettingStart(builder)
			flatbuffers_generated.NotificationSettingAddKey(builder, keyOffset)
			flatbuffers_generated.NotificationSettingAddValue(builder, value)
			notificationOffsets = append(notificationOffsets, flatbuffers_generated.NotificationSettingEnd(builder))
		}
		flatbuffers_generated.PreferencesStartNotificationsVector(builder, len(notificationOffsets))
		for i := len(notificationOffsets) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(notificationOffsets[i])
		}
		notificationsVector = builder.EndVector(len(notificationOffsets))
	}

	// Convert Privacy
	privacyOffset := f.convertPrivacySettingsToFlatBuffer(builder, prefs.Privacy)

	// Create strings last
	themeOffset := builder.CreateString(prefs.Theme)
	languageOffset := builder.CreateString(prefs.Language)

	flatbuffers_generated.PreferencesStart(builder)
	flatbuffers_generated.PreferencesAddTheme(builder, themeOffset)
	flatbuffers_generated.PreferencesAddLanguage(builder, languageOffset)
	if len(prefs.Notifications) > 0 {
		flatbuffers_generated.PreferencesAddNotifications(builder, notificationsVector)
	}
	flatbuffers_generated.PreferencesAddPrivacy(builder, privacyOffset)

	return flatbuffers_generated.PreferencesEnd(builder), nil
}

// convertPrivacySettingsToFlatBuffer converts a models.PrivacySettings to FlatBuffer format
func (f *FlatBuffersSerializer) convertPrivacySettingsToFlatBuffer(builder *flatbuffers.Builder, privacy models.PrivacySettings) flatbuffers.UOffsetT {
	flatbuffers_generated.PrivacySettingsStart(builder)
	flatbuffers_generated.PrivacySettingsAddProfilePublic(builder, privacy.ProfilePublic)
	flatbuffers_generated.PrivacySettingsAddEmailVisible(builder, privacy.EmailVisible)
	flatbuffers_generated.PrivacySettingsAddShowActivity(builder, privacy.ShowActivity)
	return flatbuffers_generated.PrivacySettingsEnd(builder)
}

// convertSettingsToFlatBuffer converts a models.Settings to FlatBuffer format
func (f *FlatBuffersSerializer) convertSettingsToFlatBuffer(builder *flatbuffers.Builder, settings models.Settings) (flatbuffers.UOffsetT, error) {
	// Convert Features first
	var featuresVector flatbuffers.UOffsetT
	if len(settings.Features) > 0 {
		featureOffsets := make([]flatbuffers.UOffsetT, len(settings.Features))
		for i, feature := range settings.Features {
			featureOffsets[i] = builder.CreateString(feature)
		}
		flatbuffers_generated.SettingsStartFeaturesVector(builder, len(featureOffsets))
		for i := len(featureOffsets) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(featureOffsets[i])
		}
		featuresVector = builder.EndVector(len(featureOffsets))
	}

	// Convert Limits map
	var limitsVector flatbuffers.UOffsetT
	if len(settings.Limits) > 0 {
		limitOffsets := make([]flatbuffers.UOffsetT, 0, len(settings.Limits))
		for key, value := range settings.Limits {
			keyOffset := builder.CreateString(key)

			flatbuffers_generated.LimitSettingStart(builder)
			flatbuffers_generated.LimitSettingAddKey(builder, keyOffset)
			flatbuffers_generated.LimitSettingAddValue(builder, int32(value))
			limitOffsets = append(limitOffsets, flatbuffers_generated.LimitSettingEnd(builder))
		}
		flatbuffers_generated.SettingsStartLimitsVector(builder, len(limitOffsets))
		for i := len(limitOffsets) - 1; i >= 0; i-- {
			builder.PrependUOffsetT(limitOffsets[i])
		}
		limitsVector = builder.EndVector(len(limitOffsets))
	}

	// Create strings last
	languageOffset := builder.CreateString(settings.Language)
	timezoneOffset := builder.CreateString(settings.TimeZone)

	flatbuffers_generated.SettingsStart(builder)
	flatbuffers_generated.SettingsAddLanguage(builder, languageOffset)
	flatbuffers_generated.SettingsAddTimezone(builder, timezoneOffset)
	if len(settings.Features) > 0 {
		flatbuffers_generated.SettingsAddFeatures(builder, featuresVector)
	}
	if len(settings.Limits) > 0 {
		flatbuffers_generated.SettingsAddLimits(builder, limitsVector)
	}

	return flatbuffers_generated.SettingsEnd(builder), nil
}

// convertMetadataEntryToFlatBuffer converts a metadata key-value pair to FlatBuffer format
func (f *FlatBuffersSerializer) convertMetadataEntryToFlatBuffer(builder *flatbuffers.Builder, key string, value interface{}) flatbuffers.UOffsetT {
	keyOffset := builder.CreateString(key)
	var stringValueOffset flatbuffers.UOffsetT

	// Create string values first if needed
	switch v := value.(type) {
	case string:
		stringValueOffset = builder.CreateString(v)
	default:
		// Fallback to string representation for unknown types
		if v != nil && fmt.Sprintf("%T", v) != "int" && fmt.Sprintf("%T", v) != "bool" && fmt.Sprintf("%T", v) != "float32" && fmt.Sprintf("%T", v) != "float64" {
			stringValueOffset = builder.CreateString(fmt.Sprintf("%v", v))
		}
	}

	flatbuffers_generated.MetadataEntryStart(builder)
	flatbuffers_generated.MetadataEntryAddKey(builder, keyOffset)

	switch v := value.(type) {
	case string:
		flatbuffers_generated.MetadataEntryAddStringValue(builder, stringValueOffset)
		flatbuffers_generated.MetadataEntryAddValueType(builder, 0) // string
	case int:
		flatbuffers_generated.MetadataEntryAddIntValue(builder, int32(v))
		flatbuffers_generated.MetadataEntryAddValueType(builder, 1) // int
	case bool:
		flatbuffers_generated.MetadataEntryAddBoolValue(builder, v)
		flatbuffers_generated.MetadataEntryAddValueType(builder, 2) // bool
	case float32:
		flatbuffers_generated.MetadataEntryAddFloatValue(builder, float64(v))
		flatbuffers_generated.MetadataEntryAddValueType(builder, 3) // float
	case float64:
		flatbuffers_generated.MetadataEntryAddFloatValue(builder, v)
		flatbuffers_generated.MetadataEntryAddValueType(builder, 3) // float
	default:
		// Fallback to string representation
		flatbuffers_generated.MetadataEntryAddStringValue(builder, stringValueOffset)
		flatbuffers_generated.MetadataEntryAddValueType(builder, 0) // string
	}

	return flatbuffers_generated.MetadataEntryEnd(builder)
}

// convertFlatBufferToUser converts a FlatBuffer User to models.User
func (f *FlatBuffersSerializer) convertFlatBufferToUser(fbUser *flatbuffers_generated.User) (models.User, error) {
	user := models.User{
		ID:        fbUser.Id(),
		Name:      string(fbUser.Name()),
		Email:     string(fbUser.Email()),
		Age:       int(fbUser.Age()),
		IsActive:  fbUser.IsActive(),
		CreatedAt: time.Unix(0, fbUser.CreatedAt()),
	}

	// Convert Profile
	fbProfile := fbUser.Profile(nil)
	if fbProfile != nil {
		profile, err := f.convertFlatBufferToProfile(fbProfile)
		if err != nil {
			return models.User{}, err
		}
		user.Profile = profile
	}

	// Convert Settings
	fbSettings := fbUser.Settings(nil)
	if fbSettings != nil {
		settings, err := f.convertFlatBufferToSettings(fbSettings)
		if err != nil {
			return models.User{}, err
		}
		user.Settings = settings
	}

	// Convert Tags
	user.Tags = make([]string, fbUser.TagsLength())
	for i := 0; i < fbUser.TagsLength(); i++ {
		user.Tags[i] = string(fbUser.Tags(i))
	}

	// Convert Metadata
	user.Metadata = make(map[string]interface{})
	fbMetadata := new(flatbuffers_generated.MetadataEntry)
	for i := 0; i < fbUser.MetadataLength(); i++ {
		if fbUser.Metadata(fbMetadata, i) {
			key := string(fbMetadata.Key())
			var value interface{}
			switch fbMetadata.ValueType() {
			case 0: // string
				value = string(fbMetadata.StringValue())
			case 1: // int
				value = int(fbMetadata.IntValue())
			case 2: // bool
				value = fbMetadata.BoolValue()
			case 3: // float
				value = fbMetadata.FloatValue()
			default:
				value = string(fbMetadata.StringValue())
			}
			user.Metadata[key] = value
		}
	}

	return user, nil
}

// convertFlatBufferToProfile converts a FlatBuffer Profile to models.Profile
func (f *FlatBuffersSerializer) convertFlatBufferToProfile(fbProfile *flatbuffers_generated.Profile) (models.Profile, error) {
	profile := models.Profile{
		FirstName: string(fbProfile.FirstName()),
		LastName:  string(fbProfile.LastName()),
		Bio:       string(fbProfile.Bio()),
		Avatar:    string(fbProfile.Avatar()),
	}

	// Convert SocialLinks
	profile.SocialLinks = make([]models.Link, fbProfile.SocialLinksLength())
	fbLink := new(flatbuffers_generated.Link)
	for i := 0; i < fbProfile.SocialLinksLength(); i++ {
		if fbProfile.SocialLinks(fbLink, i) {
			profile.SocialLinks[i] = models.Link{
				Platform: string(fbLink.Platform()),
				URL:      string(fbLink.Url()),
			}
		}
	}

	// Convert Preferences
	fbPreferences := fbProfile.Preferences(nil)
	if fbPreferences != nil {
		preferences, err := f.convertFlatBufferToPreferences(fbPreferences)
		if err != nil {
			return models.Profile{}, err
		}
		profile.Preferences = preferences
	}

	return profile, nil
}

// convertFlatBufferToPreferences converts a FlatBuffer Preferences to models.Preferences
func (f *FlatBuffersSerializer) convertFlatBufferToPreferences(fbPrefs *flatbuffers_generated.Preferences) (models.Preferences, error) {
	prefs := models.Preferences{
		Theme:    string(fbPrefs.Theme()),
		Language: string(fbPrefs.Language()),
	}

	// Convert Notifications
	prefs.Notifications = make(map[string]bool)
	fbNotification := new(flatbuffers_generated.NotificationSetting)
	for i := 0; i < fbPrefs.NotificationsLength(); i++ {
		if fbPrefs.Notifications(fbNotification, i) {
			key := string(fbNotification.Key())
			value := fbNotification.Value()
			prefs.Notifications[key] = value
		}
	}

	// Convert Privacy
	fbPrivacy := fbPrefs.Privacy(nil)
	if fbPrivacy != nil {
		prefs.Privacy = models.PrivacySettings{
			ProfilePublic: fbPrivacy.ProfilePublic(),
			EmailVisible:  fbPrivacy.EmailVisible(),
			ShowActivity:  fbPrivacy.ShowActivity(),
		}
	}

	return prefs, nil
}

// convertFlatBufferToSettings converts a FlatBuffer Settings to models.Settings
func (f *FlatBuffersSerializer) convertFlatBufferToSettings(fbSettings *flatbuffers_generated.Settings) (models.Settings, error) {
	settings := models.Settings{
		Language: string(fbSettings.Language()),
		TimeZone: string(fbSettings.Timezone()),
	}

	// Convert Features
	settings.Features = make([]string, fbSettings.FeaturesLength())
	for i := 0; i < fbSettings.FeaturesLength(); i++ {
		settings.Features[i] = string(fbSettings.Features(i))
	}

	// Convert Limits
	settings.Limits = make(map[string]int)
	fbLimit := new(flatbuffers_generated.LimitSetting)
	for i := 0; i < fbSettings.LimitsLength(); i++ {
		if fbSettings.Limits(fbLimit, i) {
			key := string(fbLimit.Key())
			value := int(fbLimit.Value())
			settings.Limits[key] = value
		}
	}

	return settings, nil
}
