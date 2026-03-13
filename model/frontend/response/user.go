package response

import "server/model/frontend"

type LoginResponse struct {
	User      frontend.User `json:"user"`
	Token     string        `json:"token"`
	ExpiresAt int64         `json:"expiresAt"`
}
