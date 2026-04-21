package social

type SocialAccount struct {
	ID             uint
	UserID         uint
	Provider       string
	ProviderUserID string
	AvatarURL      string
}
