package services

import (
	"family-calendar-backend/db"
	"family-calendar-backend/db/models"
)

// FindOrCreateUser finds an existing user or creates a new one based on OAuth provider info
func FindOrCreateUser(authProvider, authProviderID, givenName, familyName, email string) (*models.User, error) {
	var user models.User

	// Try to find existing user by auth provider and provider ID
	result := db.DB.Where("auth_provider = ? AND auth_provider_id = ?", authProvider, authProviderID).First(&user)

	if result.Error == nil {
		// User found, update their information in case it changed
		user.GivenName = givenName
		user.FamilyName = familyName
		user.Email = email
		db.DB.Save(&user)
		return &user, nil
	}

	// User not found, create new user
	user = models.User{
		GivenName:      givenName,
		FamilyName:     familyName,
		Email:          email,
		AuthProvider:   authProvider,
		AuthProviderID: authProviderID,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
