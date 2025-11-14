import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if available
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  register: (data: {
    email: string;
    password: string;
    username: string;
    display_name: string;
    genres?: string[];
  }) => api.post('/auth/register', data),

  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
};

// Tracks API
export const tracksAPI = {
  getTracks: (params?: { page?: number; limit?: number; genre?: string; search?: string }) =>
    api.get('/tracks', { params }),

  getTrackById: (id: number) => api.get(`/tracks/${id}`),

  getSimilarTracks: (id: number) => api.get(`/tracks/${id}/similar`),

  addTrack: (data: any) => api.post('/tracks/add', data),

  recordPlay: (id: number) => api.post(`/tracks/${id}/play`),
};

// Artists API
export const artistsAPI = {
  getArtists: () => api.get('/artists'),

  getArtistById: (id: number) => api.get(`/artists/${id}`),

  getArtistStats: (id: number) => api.get(`/artists/${id}/stats`),
};

// Albums API
export const albumsAPI = {
  getAlbums: () => api.get('/albums'),

  getAlbumStats: (id: number) => api.get(`/albums/${id}/stats`),

  getAlbumDuration: (id: number) => api.get(`/albums/${id}/duration`),
};

// Playlists API
export const playlistsAPI = {
  createPlaylist: (data: { name: string; description: string; is_public: boolean }) =>
    api.post('/playlists', data),

  getUserPlaylists: () => api.get('/playlists'),

  getPlaylistById: (id: string) => api.get(`/playlists/${id}`),

  updatePlaylist: (id: string, data: any) => api.put(`/playlists/${id}`, data),

  deletePlaylist: (id: string) => api.delete(`/playlists/${id}`),

  addTrackToPlaylist: (id: string, trackId: number) =>
    api.post(`/playlists/${id}/tracks`, { track_id: trackId }),

  removeTrackFromPlaylist: (id: string, trackId: number) =>
    api.delete(`/playlists/${id}/tracks/${trackId}`),
};

// Recommendations API
export const recommendationsAPI = {
  getTrending: () => api.get('/recommendations/trending'),

  getGenreRecommendations: (genre: string) =>
    api.get(`/recommendations/genre/${genre}`),

  getPersonalRecommendations: () => api.get('/recommendations'),
};

// Profile API
export const profileAPI = {
  getProfile: () => api.get('/profile'),

  updatePreferences: (data: any) => api.put('/profile/preferences', data),
};

// Search API
export const searchAPI = {
  search: (query: string) => api.get('/search', { params: { q: query } }),
};

export default api;

