package api

type PushSubscriptionRequest struct {
	Endpoint string `json:"endpoint"`
	Keys     PushSubscriptionKeys `json:"keys"`
}

type PushSubscriptionKeys struct {
	P256DH string `json:"p256dh"`
	Auth   string `json:"auth"`
}
