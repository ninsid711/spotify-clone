import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { playlistsAPI } from '../services/api';
import { Playlist } from '../types';
import './Playlists.css';

const Playlists: React.FC = () => {
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newPlaylist, setNewPlaylist] = useState({
    name: '',
    description: '',
    is_public: true,
  });

  useEffect(() => {
    loadPlaylists();
  }, []);

  const loadPlaylists = async () => {
    setIsLoading(true);
    try {
      const response = await playlistsAPI.getUserPlaylists();
      setPlaylists(response.data.playlists || []);
    } catch (error) {
      console.error('Failed to load playlists:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreatePlaylist = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await playlistsAPI.createPlaylist(newPlaylist);
      setShowCreateModal(false);
      setNewPlaylist({ name: '', description: '', is_public: true });
      loadPlaylists();
    } catch (error) {
      console.error('Failed to create playlist:', error);
    }
  };

  return (
    <div className="playlists-container">
      <div className="playlists-header">
        <h1>Your Playlists</h1>
        <button
          className="btn-primary"
          onClick={() => setShowCreateModal(true)}
        >
          + Create Playlist
        </button>
      </div>

      {isLoading ? (
        <div className="loading">Loading playlists...</div>
      ) : playlists.length === 0 ? (
        <div className="empty-state">
          <p>You don't have any playlists yet.</p>
          <button
            className="btn-primary"
            onClick={() => setShowCreateModal(true)}
          >
            Create your first playlist
          </button>
        </div>
      ) : (
        <div className="playlists-grid">
          {playlists.map((playlist) => (
            <Link
              key={playlist.id}
              to={`/playlists/${playlist.id}`}
              className="playlist-card"
            >
              <div className="playlist-cover">
                {playlist.cover_url ? (
                  <img src={playlist.cover_url} alt={playlist.name} />
                ) : (
                  <div className="playlist-cover-placeholder">
                    <span>🎵</span>
                  </div>
                )}
              </div>
              <div className="playlist-info">
                <h3>{playlist.name}</h3>
                <p>{playlist.description || 'No description'}</p>
                <span className="playlist-meta">
                  {playlist.track_ids?.length || 0} tracks
                  {playlist.is_public ? ' • Public' : ' • Private'}
                </span>
              </div>
            </Link>
          ))}
        </div>
      )}

      {/* Create Playlist Modal */}
      {showCreateModal && (
        <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>Create New Playlist</h2>
            <form onSubmit={handleCreatePlaylist}>
              <div className="form-group">
                <label htmlFor="name">Playlist Name</label>
                <input
                  id="name"
                  type="text"
                  value={newPlaylist.name}
                  onChange={(e) =>
                    setNewPlaylist({ ...newPlaylist, name: e.target.value })
                  }
                  placeholder="My Awesome Playlist"
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="description">Description</label>
                <textarea
                  id="description"
                  value={newPlaylist.description}
                  onChange={(e) =>
                    setNewPlaylist({
                      ...newPlaylist,
                      description: e.target.value,
                    })
                  }
                  placeholder="Describe your playlist..."
                  rows={3}
                />
              </div>

              <div className="form-group checkbox-group">
                <label>
                  <input
                    type="checkbox"
                    checked={newPlaylist.is_public}
                    onChange={(e) =>
                      setNewPlaylist({
                        ...newPlaylist,
                        is_public: e.target.checked,
                      })
                    }
                  />
                  Make playlist public
                </label>
              </div>

              <div className="modal-actions">
                <button
                  type="button"
                  className="btn-secondary"
                  onClick={() => setShowCreateModal(false)}
                >
                  Cancel
                </button>
                <button type="submit" className="btn-primary">
                  Create
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Playlists;

