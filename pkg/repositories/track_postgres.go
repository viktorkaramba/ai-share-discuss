package repositories

import (
	"PlaylistsSynchronizer.Backend/pkg/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TrackPostgres struct {
	db *sqlx.DB
}

func NewTrackPostgres(db *sqlx.DB) *TrackPostgres {
	return &TrackPostgres{db: db}
}

func (t *TrackPostgres) Create(track models.CreateTrack) (int, error) {
	var trackID int
	query := fmt.Sprintf("INSERT INTO %s (spotifyuri, youtubemusicid, name) values ($1, $2, $3) RETURNING id", tracksTable)
	row1 := t.db.QueryRow(query, track.SpotifyUri, track.YouTubeMusicID, track.Name)
	if err := row1.Scan(&trackID); err != nil {
		return 0, err
	}
	return trackID, nil
}

func (t *TrackPostgres) AddYouTubeMusicTrackPlayList(playListId int, track models.CreateTrack) (int, error) {
	var youTubeTrackID int
	youTubeTrack := fmt.Sprintf("INSERT INTO %s (userID, trackID, playListID, playListYouTubeMusicID) "+
		"values ($1, $2, $3, $4) RETURNING id", youtubeMusicTracks)
	row1 := t.db.QueryRow(youTubeTrack, track.UserID, track.ID, playListId, track.PlayListYouTubeMusicID)
	if err := row1.Scan(&youTubeTrackID); err != nil {
		return 0, err
	}
	return youTubeTrackID, nil
}

func (t *TrackPostgres) AddSpotifyTrackPlayList(playListId int, track models.CreateTrack) (int, error) {
	var playListTrackID int
	playListTrack := fmt.Sprintf("INSERT INTO %s (trackID, playListID) values ($1, $2) RETURNING id", playlistTrackTable)
	row1 := t.db.QueryRow(playListTrack, track.ID, playListId)
	if err := row1.Scan(&playListTrackID); err != nil {
		return 0, err
	}
	return playListTrackID, nil
}

func (t *TrackPostgres) GetAll() ([]models.Track, error) {
	var tracks []models.Track
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id ASC ", tracksTable)
	err := t.db.Select(&tracks, query)
	return tracks, err
}

func (t *TrackPostgres) GetByID(ID int) (*models.Track, error) {
	var track models.Track
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", tracksTable)
	err := t.db.Get(&track, query, ID)
	if track == (models.Track{}) {
		return nil, nil
	}
	return &track, err
}

func (t *TrackPostgres) GetByPlayListTrackID(playListID, trackID int) ([]models.PlayListTrack, error) {
	var tracks []models.PlayListTrack
	query := fmt.Sprintf("SELECT ytm.id, ytm.userid, ytm.trackid, ytm.playlistid, "+
		"ytm.playlistyoutubemusicid, uytm.accesstoken FROM %s ytm JOIN %s uytm "+
		"ON ytm.userid = uytm.userid WHERE ytm.playListID=$1 AND ytm.trackID=$2 ORDER BY id ASC",
		youtubeMusicTracks, userYouTubeMusicTable)
	err := t.db.Select(&tracks, query, playListID, trackID)
	return tracks, err
}

func (t *TrackPostgres) GetByPlayListID(playListID int) ([]models.Track, error) {
	var tracks []models.Track
	query := fmt.Sprintf("SELECT t.id, t.spotifyuri, t.youtubemusicid, t.name "+
		"FROM %s t INNER JOIN %s pt ON t.id = pt.trackid WHERE pt.playlistid=$1 ORDER BY t.id ASC", tracksTable, playlistTrackTable)
	err := t.db.Select(&tracks, query, playListID)
	return tracks, err
}

func (t *TrackPostgres) GetByPlayListIDAndTrackApiID(playListID int, apiID models.ApiTrackID) (*models.Track, error) {
	var track models.Track
	query := fmt.Sprintf("SELECT t.id, t.spotifyUri, t.youTubeMusicID "+
		"FROM %s t INNER JOIN %s pt ON t.id = pt.trackID AND t.spotifyUri=$1 AND t.youTubeMusicID=$2 "+
		"WHERE pt.playListID=$3", tracksTable, playlistTrackTable)
	err := t.db.Get(&track, query, apiID.SpotifyUri, apiID.YouTubeMusic, playListID)
	if track == (models.Track{}) {
		return nil, nil
	}
	return &track, err
}

func (t *TrackPostgres) GetByTrackApiID(apiID models.ApiTrackID) (*models.Track, error) {
	var track models.Track
	query := fmt.Sprintf("SELECT t.id, t.spotifyUri, t.youTubeMusicID "+
		"FROM %s t WHERE t.spotifyUri=$1 AND t.youTubeMusicID=$2", tracksTable)
	err := t.db.Get(&track, query, apiID.SpotifyUri, apiID.YouTubeMusic)
	if track == (models.Track{}) {
		return nil, nil
	}
	return &track, err
}

func (t *TrackPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tracksTable)
	_, err := t.db.Exec(query, id)
	return err
}

func (t *TrackPostgres) DeleteFromPlayList(playListID, trackID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE playListID=$1 AND trackID=$2", playlistTrackTable)
	_, err := t.db.Exec(query, playListID, trackID)
	return err
}

func (t *TrackPostgres) DeleteFromYouTubeMusicPlayList(userID, playListID, trackID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE userID=$1 AND playListID=$2 AND trackID=$3", youtubeMusicTracks)
	_, err := t.db.Exec(query, userID, playListID, trackID)
	return err
}
