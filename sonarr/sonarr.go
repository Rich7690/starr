package sonarr

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"golift.io/starr"
)

// Sonarr contains all the methods to interact with a Sonarr server.
type Sonarr struct {
	starr.APIer
}

// Filter values are integers. Given names for ease of discovery.
//nolint:lll
// https://github.com/Sonarr/Sonarr/blob/0cb8d93069d6310abd39ee2fe73219e17aa83fe6/src/NzbDrone.Core/History/EpisodeHistory.cs#L34-L41
const (
	FilterUnknown starr.Filtering = iota
	FilterGrabbed
	FilterSeriesFolderImported
	FilterDownloadFolderImported
	FilterDownloadFailed
	FilterDeleted
	FilterRenamed
	FilterImportFailed
)

// New returns a Sonarr object used to interact with the Sonarr API.
func New(config *starr.Config) *Sonarr {
	if config.Client == nil {
		//nolint:exhaustivestruct,gosec
		config.Client = &http.Client{
			Timeout: config.Timeout.Duration,
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.ValidSSL},
			},
		}
	}

	if config.Debugf == nil {
		config.Debugf = func(string, ...interface{}) {}
	}

	return &Sonarr{APIer: config}
}

// GetQueue returns a single page from the Sonarr Queue (processing, but not yet imported).
// WARNING: 12/30/2021 - this method changed.
// If you need control over the page, use sonarr.GetQueuePage().
// This function simply returns the number of queue records desired,
// up to the number of records present in the application.
// It grabs records in (paginated) batches of perPage, and concatenates
// them into one list.  Passing zero for records will return all of them.
func (s *Sonarr) GetQueue(records, perPage int) (*Queue, error) {
	return s.GetQueueContext(context.Background(), records, perPage)
}

func (s *Sonarr) GetQueueContext(ctx context.Context, records, perPage int) (*Queue, error) {
	queue := &Queue{Records: []*QueueRecord{}}
	perPage = starr.SetPerPage(records, perPage)

	for page := 1; ; page++ {
		curr, err := s.GetQueuePageContext(ctx, &starr.Req{PageSize: perPage, Page: page})
		if err != nil {
			return nil, err
		}

		queue.Records = append(queue.Records, curr.Records...)

		if len(queue.Records) >= curr.TotalRecords ||
			(len(queue.Records) >= records && records != 0) ||
			len(curr.Records) == 0 {
			queue.PageSize = curr.TotalRecords
			queue.TotalRecords = curr.TotalRecords
			queue.SortDirection = curr.SortDirection
			queue.SortKey = curr.SortKey

			break
		}

		perPage = starr.AdjustPerPage(records, curr.TotalRecords, len(queue.Records), perPage)
	}

	return queue, nil
}

// GetQueuePage returns a single page from the Sonarr Queue.
// The page size and number is configurable with the input request parameters.
func (s *Sonarr) GetQueuePage(params *starr.Req) (*Queue, error) {
	return s.GetQueuePageContext(context.Background(), params)
}

func (s *Sonarr) GetQueuePageContext(ctx context.Context, params *starr.Req) (*Queue, error) {
	var queue Queue

	params.CheckSet("sortKey", "timeleft")
	params.CheckSet("includeUnknownSeriesItems", "true")

	err := s.GetInto(ctx, "v3/queue", params.Params(), &queue)
	if err != nil {
		return nil, fmt.Errorf("api.Get(queue): %w", err)
	}

	return &queue, nil
}

// GetLanguageProfiles returns all configured language profiles.
func (s *Sonarr) GetLanguageProfiles() ([]*LanguageProfile, error) {
	return s.GetLanguageProfilesContext(context.Background())
}

func (s *Sonarr) GetLanguageProfilesContext(ctx context.Context) ([]*LanguageProfile, error) {
	var profiles []*LanguageProfile

	err := s.GetInto(ctx, "v3/languageprofile", nil, &profiles)
	if err != nil {
		return nil, fmt.Errorf("api.Get(languageprofile): %w", err)
	}

	return profiles, nil
}

// GetRootFolders returns all configured root folders.
func (s *Sonarr) GetRootFolders() ([]*RootFolder, error) {
	return s.GetRootFoldersContext(context.Background())
}

func (s *Sonarr) GetRootFoldersContext(ctx context.Context) ([]*RootFolder, error) {
	var folders []*RootFolder

	err := s.GetInto(ctx, "v3/rootfolder", nil, &folders)
	if err != nil {
		return nil, fmt.Errorf("api.Get(rootfolder): %w", err)
	}

	return folders, nil
}
