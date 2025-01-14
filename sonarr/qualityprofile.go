package sonarr

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"golift.io/starr"
)

// GetQualityProfiles returns all configured quality profiles.
func (s *Sonarr) GetQualityProfiles() ([]*QualityProfile, error) {
	return s.GetQualityProfilesContext(context.Background())
}

func (s *Sonarr) GetQualityProfilesContext(ctx context.Context) ([]*QualityProfile, error) {
	var profiles []*QualityProfile

	err := s.GetInto(ctx, "v3/qualityprofile", nil, &profiles)
	if err != nil {
		return nil, fmt.Errorf("api.Get(qualityprofile): %w", err)
	}

	return profiles, nil
}

// AddQualityProfile updates a quality profile in place.
func (s *Sonarr) AddQualityProfile(profile *QualityProfile) (int64, error) {
	return s.AddQualityProfileContext(context.Background(), profile)
}

func (s *Sonarr) AddQualityProfileContext(ctx context.Context, profile *QualityProfile) (int64, error) {
	post, err := json.Marshal(profile)
	if err != nil {
		return 0, fmt.Errorf("json.Marshal(profile): %w", err)
	}

	var output QualityProfile

	err = s.PostInto(ctx, "v3/qualityProfile", nil, post, &output)
	if err != nil {
		return 0, fmt.Errorf("api.Post(qualityProfile): %w", err)
	}

	return output.ID, nil
}

// UpdateQualityProfile updates a quality profile in place.
func (s *Sonarr) UpdateQualityProfile(profile *QualityProfile) error {
	return s.UpdateQualityProfileContext(context.Background(), profile)
}

func (s *Sonarr) UpdateQualityProfileContext(ctx context.Context, profile *QualityProfile) error {
	put, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("json.Marshal(profile): %w", err)
	}

	_, err = s.Put(ctx, "v3/qualityProfile/"+strconv.FormatInt(profile.ID, starr.Base10), nil, put)
	if err != nil {
		return fmt.Errorf("api.Put(qualityProfile): %w", err)
	}

	return nil
}
