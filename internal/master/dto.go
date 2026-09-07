package master

type ProfileDTO struct {
	UserID            int64   `json:"userId"`
	City              string  `json:"city,omitempty"`
	Bio               string  `json:"bio,omitempty"`
	AvatarURL         *string `json:"avatarUrl,omitempty"`
	RatingAvg         float64 `json:"ratingAvg"`
	RatingCount       int32   `json:"ratingCount"`
	SpecializationIDs []int64 `json:"specializationIds"`
}

type UpdateProfileRequest struct {
	City              string  `json:"city"`
	Bio               string  `json:"bio"`
	AvatarURL         *string `json:"avatarUrl,omitempty"`
	SpecializationIDs []int64 `json:"specializationIds"`
}

func ToProfileDTO(p Profile, specializationIDs []int64) ProfileDTO {
	return ProfileDTO{
		UserID:            p.UserID,
		City:              p.City,
		Bio:               p.Bio,
		AvatarURL:         p.AvatarURL,
		RatingAvg:         p.RatingAvg,
		RatingCount:       p.RatingCount,
		SpecializationIDs: specializationIDs,
	}
}
