package model

import "context"

type FeedRequest struct {
	Ctx         context.Context
	UserProfile *Profile
}
