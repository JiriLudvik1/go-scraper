package repository

import (
	"go-scraper/models"
	"time"

	"github.com/patrickmn/go-cache"
)

var c *cache.Cache
var isInitialized = false

func InitializeIgnoreCache() {
	if isInitialized {
		return
	}

	c = cache.New(5*time.Minute, 10*time.Minute)
	isInitialized = true
}

// Adding
func AddListingToIgnoreCache(listing *models.Listing) {
	c.Add(getListingCacheKey(listing), true, cache.DefaultExpiration)
}

func AddUserToIgnoreCache(user *models.User) {
	c.Add(getUserCacheKey(user), true, cache.DefaultExpiration)
}

// Checking
func IsListingIgnored(listing *models.Listing) bool {
	_, found := c.Get(getListingCacheKey(listing))
	return found
}

func IsUserIgnored(user *models.User) bool {
	_, found := c.Get(getUserCacheKey(user))
	return found
}

// Helpers, cache keys
func getListingCacheKey(listing *models.Listing) string {
	return "listing-" + listing.ID
}
func getUserCacheKey(user *models.User) string {
	return "user-" + user.UserName
}
