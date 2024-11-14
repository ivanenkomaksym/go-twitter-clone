package repositories

import "twitter-clone/internal/models"

// ContainsTag checks if a given tag is present in the tags slice
func ContainsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// ContainsTweet checks if a given tweet ID is present in the tweets slice
func ContainsTweet(tweets []models.Tweet, tweetId string) bool {
	for _, t := range tweets {
		if t.ID == tweetId {
			return true
		}
	}
	return false
}
