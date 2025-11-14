import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { playlistsAPI } from '../services/api';
import { Playlist, Track } from '../types';
import './PlaylistDetail.css';

const PlaylistDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [playlist, setPlaylist] = useState<Playlist | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (id) {
      loadPlaylist();
    }
  }, [id]);

  const loadPlaylist = async () => {
    if (!id) return;
    setIsLoading(true);
    try {
      const response = await playlistsAPI.getPlaylistById(id);
      setPlaylist(response.data.playlist);
      setTracks(response.data.tracks || []);
    } catch (error) {
      console.error('Failed to load playlist:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeletePlaylist = async () => {
    if (!id || !window.confirm('Are you sure you want to delete this playlist?')) {
      return;
    }
    try {
      await playlistsAPI.deletePlaylist(id);
      navigate('/playlists');
    } catch (error) {
      console.error('Failed to delete playlist:', error);
    }
  };

  const handleRemoveTrack = async (trackId: number) => {
    if (!id) return;
    try {
      await playlistsAPI.removeTrackFromPlaylist(id, trackId);
      loadPlaylist();
    } catch (error) {
      console.error('Failed to remove track:', error);
    }
  };

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  if (isLoading) {
    return <div className="loading">Loading playlist...</div>;
  }

  if (!playlist) {
    return <div className="error">Playlist not found</div>;
  }

  const totalDuration = tracks.reduce((sum, track) => sum + track.duration, 0);

  return (
    <div className="playlist-detail-container">
      <div className="playlist-header">
        <div className="playlist-cover-large">
          {playlist.cover_url ? (
            <img src={playlist.cover_url} alt={playlist.name} />
          ) : (
            <div className="playlist-cover-placeholder">
              <span>🎵</span>
            </div>
          )}
        </div>
        <div className="playlist-info-large">
          <span className="playlist-type">Playlist</span>
          <h1>{playlist.name}</h1>
          <p className="playlist-description">{playlist.description}</p>
          <div className="playlist-meta">
            <span>{tracks.length} songs</span>
            <span>•</span>
            <span>{Math.floor(totalDuration / 60)} min</span>
            <span>•</span>
            <span>{playlist.is_public ? 'Public' : 'Private'}</span>
          </div>
        </div>
      </div>

      <div className="playlist-actions">
        <button className="btn-primary">▶ Play All</button>
        <button className="btn-danger" onClick={handleDeletePlaylist}>
          Delete Playlist
        </button>
      </div>

      {tracks.length === 0 ? (
        <div className="empty-playlist">
          <p>This playlist is empty. Add some tracks!</p>
        </div>
      ) : (
        <div className="tracks-table">
          <div className="tracks-table-header">
            <div className="track-number">#</div>
            <div className="track-title">Title</div>
            <div className="track-artist">Artist</div>
            <div className="track-album">Album</div>
            <div className="track-duration">Duration</div>
            <div className="track-actions">Actions</div>
          </div>
          {tracks.map((track, index) => (
            <div key={track.id} className="tracks-table-row">
              <div className="track-number">{index + 1}</div>
              <div className="track-title">
                <img
                  src={track.cover_url || '/placeholder.png'}
                  alt={track.title}
                  className="track-thumbnail"
                />
                <span>{track.title}</span>
              </div>
              <div className="track-artist">{track.artist_name}</div>
              <div className="track-album">{track.album_name}</div>
              <div className="track-duration">{formatDuration(track.duration)}</div>
              <div className="track-actions">
                <button
                  className="btn-icon"
                  onClick={() => handleRemoveTrack(track.id)}
                  title="Remove from playlist"
                >
                  ✕
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default PlaylistDetail;

