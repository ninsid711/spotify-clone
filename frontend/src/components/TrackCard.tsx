import React, { useState } from 'react';
import { Track } from '../types';
import { tracksAPI, playlistsAPI } from '../services/api';
import { useAuth } from '../context/AuthContext';
import './TrackCard.css';

interface TrackCardProps {
  track: Track;
}

const TrackCard: React.FC<TrackCardProps> = ({ track }) => {
  const { token } = useAuth();
  const [showAddToPlaylist, setShowAddToPlaylist] = useState(false);
  const [playlists, setPlaylists] = useState<any[]>([]);

  const formatDuration = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const handlePlay = async () => {
    if (token) {
      try {
        await tracksAPI.recordPlay(track.id);
        console.log('Play recorded');
      } catch (error) {
        console.error('Failed to record play:', error);
      }
    }
    // In a real app, you would actually play the track here
    console.log('Playing:', track.title);
  };

  const handleAddToPlaylist = async () => {
    if (!token) {
      alert('Please login to add tracks to playlists');
      return;
    }
    try {
      const response = await playlistsAPI.getUserPlaylists();
      setPlaylists(response.data.playlists || []);
      setShowAddToPlaylist(true);
    } catch (error) {
      console.error('Failed to load playlists:', error);
    }
  };

  const addToPlaylist = async (playlistId: string) => {
    try {
      await playlistsAPI.addTrackToPlaylist(playlistId, track.id);
      setShowAddToPlaylist(false);
      alert('Track added to playlist!');
    } catch (error) {
      console.error('Failed to add track to playlist:', error);
      alert('Failed to add track to playlist');
    }
  };

  return (
    <div className="track-card">
      <div className="track-cover">
        <img
          src={track.cover_url || '/placeholder.png'}
          alt={track.title}
          onError={(e) => {
            (e.target as HTMLImageElement).src = '/placeholder.png';
          }}
        />
        <div className="track-overlay">
          <button className="play-button" onClick={handlePlay}>
            ▶
          </button>
        </div>
      </div>
      <div className="track-info">
        <h3 className="track-title">{track.title}</h3>
        <p className="track-artist">{track.artist_name}</p>
        <p className="track-album">{track.album_name}</p>
        <div className="track-meta">
          <span className="track-genre">{track.genre}</span>
          <span className="track-duration">{formatDuration(track.duration)}</span>
        </div>
      </div>
      <div className="track-actions">
        <button
          className="btn-icon"
          onClick={handleAddToPlaylist}
          title="Add to playlist"
        >
          +
        </button>
      </div>

      {/* Add to Playlist Modal */}
      {showAddToPlaylist && (
        <div
          className="modal-overlay"
          onClick={() => setShowAddToPlaylist(false)}
        >
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>Add to Playlist</h3>
            {playlists.length === 0 ? (
              <p>You don't have any playlists yet.</p>
            ) : (
              <div className="playlist-list">
                {playlists.map((playlist) => (
                  <button
                    key={playlist.id}
                    className="playlist-item"
                    onClick={() => addToPlaylist(playlist.id)}
                  >
                    {playlist.name}
                  </button>
                ))}
              </div>
            )}
            <button
              className="btn-secondary"
              onClick={() => setShowAddToPlaylist(false)}
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default TrackCard;

