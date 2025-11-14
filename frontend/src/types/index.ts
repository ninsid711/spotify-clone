export interface Track {
  id: number;
  title: string;
  artist_id: number;
  artist_name: string;
  album_id: number;
  album_name: string;
  duration: number;
  genre: string;
  release_date: string;
  file_url: string;
  cover_url: string;
  created_at: string;
}

export interface Artist {
  id: number;
  name: string;
  bio: string;
  image_url: string;
  created_at: string;
}

export interface Album {
  id: number;
  title: string;
  artist_id: number;
  artist_name?: string;
  release_date: string;
  cover_url: string;
  created_at: string;
}

export interface User {
  id: string;
  email: string;
  username: string;
  display_name: string;
  profile_picture_url?: string;
  preferences: UserPreferences;
  listening_history: ListeningHistory[];
  favorite_artists: number[];
  favorite_genres: string[];
  created_at: string;
  updated_at: string;
}

export interface UserPreferences {
  theme: 'light' | 'dark';
  language: string;
  explicit_content: boolean;
  preferred_genres: string[];
}

export interface ListeningHistory {
  track_id: number;
  played_at: string;
  duration: number;
  completed: boolean;
}

export interface Playlist {
  id: string;
  user_id: string;
  name: string;
  description: string;
  track_ids: number[];
  tracks?: Track[];
  is_public: boolean;
  cover_url: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface SearchResults {
  tracks: Track[];
  artists: Artist[];
  albums: Album[];
}

