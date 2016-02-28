package backend

import (
	"github.com/mindcastio/podcast-feed"

	"github.com/mindcastio/mindcastio/backend/util"
)

func podcastDetailsToMetadata(podcast *podcast.PodcastDetails) *PodcastMetadata {
	meta := PodcastMetadata{
		podcast.Uid,
		podcast.Title,
		podcast.Subtitle,
		podcast.Url,
		podcast.Feed,
		podcast.Description,
		podcast.Published,
		podcast.Language,
		podcast.Image,
		podcast.Owner.Name,
		podcast.Owner.Email,
		"",
		0,
		0,
		0,
		0,
		util.Timestamp(),
		0,
	}
	return &meta
}

func episodeDetailsToMetadata(episode *podcast.EpisodeDetails, puid string) *EpisodeMetadata {
	meta := EpisodeMetadata{
		episode.Uid,
		episode.Title,
		episode.Url,
		episode.Description,
		episode.Published,
		episode.Duration,
		episode.Author,
		episode.Content.Url,
		episode.Content.Type,
		episode.Content.Size,
		puid,
		0,
		util.Timestamp(),
		0,
	}
	return &meta
}

func podcastMetadataToSummary(p *PodcastMetadata) PodcastSummary {
	return PodcastSummary{
		p.Uid,
		p.Title,
		p.OwnerName,
		p.Description,
		p.Url,
		p.Feed,
		p.ImageUrl,
		p.Published,
	}
}

func podcastMetadataToSearch(p *PodcastMetadata) PodcastMetadataSearch {
	return PodcastMetadataSearch{
		p.Uid,
		p.Title,
		p.Subtitle,
		p.Description,
		p.Published,
		p.Language,
		p.OwnerName,
		p.OwnerEmail,
		p.Tags,
	}
}

func episodeMetadataToSearch(e *EpisodeMetadata) EpisodeMetadataSearch {
	return EpisodeMetadataSearch{
		e.Uid,
		e.Title,
		e.Url,
		e.Description,
		e.Published,
		e.Author,
		e.PodcastUid,
	}
}
