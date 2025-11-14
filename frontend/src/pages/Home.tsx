import React, { useState, useEffect } from 'react';
import { tracksAPI, recommendationsAPI } from '../services/api';
import { Track } from '../types';
import TrackCard from '../components/TrackCard';
import './Home.css';

const Home: React.FC = () => {
  const [tracks, setTracks] = useState<Track[]>([]);
  const [trendingTracks, setTrendingTracks] = useState<Track[]>([]);
  const [selectedGenre, setSelectedGenre] = useState<string>('');
  const [searchQuery, setSearchQuery] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [page, setPage] = useState(1);

  const GENRES = ['Pop', 'Rock', 'Hip Hop', 'Jazz', 'Classical', 'Electronic', 'Country', 'R&B', 'Indie', 'Metal'];

  useEffect(() => {
    loadTracks();
    loadTrendingTracks();
  }, [page, selectedGenre, searchQuery]);

  const loadTracks = async () => {
    setIsLoading(true);
    try {
      const params: any = { page, limit: 20 };
      if (selectedGenre) params.genre = selectedGenre;
      if (searchQuery) params.search = searchQuery;

      const response = await tracksAPI.getTracks(params);
      setTracks(response.data.tracks || []);
    } catch (error) {
      console.error('Failed to load tracks:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const loadTrendingTracks = async () => {
    try {
      const response = await recommendationsAPI.getTrending();
      setTrendingTracks(response.data.tracks || []);
    } catch (error) {
      console.error('Failed to load trending tracks:', error);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setPage(1);
    loadTracks();
  };

  const handleGenreChange = (genre: string) => {
    setSelectedGenre(genre === selectedGenre ? '' : genre);
    setPage(1);
  };

  return (
    <div className="home-container">
      <div className="home-header">
        <h1>Discover Music</h1>
        <form onSubmit={handleSearch} className="search-form">
          <input
            type="text"
            placeholder="Search for tracks, artists..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="search-input"
          />
          <button type="submit" className="btn-search">Search</button>
        </form>
      </div>

      {/* Genre Filter */}
      <div className="genre-filter">
        <h3>Filter by Genre</h3>
        <div className="genre-pills">
          {GENRES.map((genre) => (
            <button
              key={genre}
              className={`genre-pill ${selectedGenre === genre ? 'active' : ''}`}
              onClick={() => handleGenreChange(genre)}
            >
              {genre}
            </button>
          ))}
        </div>
      </div>

      {/* Trending Section */}
      {!searchQuery && !selectedGenre && trendingTracks.length > 0 && (
        <section className="tracks-section">
          <h2>🔥 Trending Now</h2>
          <div className="tracks-grid">
            {trendingTracks.map((track) => (
              <TrackCard key={track.id} track={track} />
            ))}
          </div>
        </section>
      )}

      {/* All Tracks Section */}
      <section className="tracks-section">
        <h2>{searchQuery ? 'Search Results' : selectedGenre ? `${selectedGenre} Tracks` : 'All Tracks'}</h2>
        {isLoading ? (
          <div className="loading">Loading tracks...</div>
        ) : tracks.length === 0 ? (
          <div className="no-results">No tracks found</div>
        ) : (
          <>
            <div className="tracks-grid">
              {tracks.map((track) => (
                <TrackCard key={track.id} track={track} />
              ))}
            </div>
            <div className="pagination">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="btn-pagination"
              >
                Previous
              </button>
              <span className="page-info">Page {page}</span>
              <button
                onClick={() => setPage((p) => p + 1)}
                disabled={tracks.length < 20}
                className="btn-pagination"
              >
                Next
              </button>
            </div>
          </>
        )}
      </section>
    </div>
  );
};

export default Home;

