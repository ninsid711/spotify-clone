import React, { useState, useEffect } from 'react';
import { recommendationsAPI } from '../services/api';
import { Track } from '../types';
import TrackCard from '../components/TrackCard';
import './Recommendations.css';

const Recommendations: React.FC = () => {
  const [personalTracks, setPersonalTracks] = useState<Track[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadRecommendations();
  }, []);

  const loadRecommendations = async () => {
    setIsLoading(true);
    try {
      const response = await recommendationsAPI.getPersonalRecommendations();
      setPersonalTracks(response.data.tracks || []);
    } catch (error) {
      console.error('Failed to load recommendations:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="recommendations-container">
      <div className="recommendations-header">
        <h1>Made For You</h1>
        <p>Personalized recommendations based on your listening history</p>
      </div>

      {isLoading ? (
        <div className="loading">Loading recommendations...</div>
      ) : personalTracks.length === 0 ? (
        <div className="empty-state">
          <p>Start listening to tracks to get personalized recommendations!</p>
        </div>
      ) : (
        <section className="tracks-section">
          <h2>✨ Your Personalized Mix</h2>
          <div className="tracks-grid">
            {personalTracks.map((track) => (
              <TrackCard key={track.id} track={track} />
            ))}
          </div>
        </section>
      )}
    </div>
  );
};

export default Recommendations;

