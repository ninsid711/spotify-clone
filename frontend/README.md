# Spotify Clone Frontend

A modern, responsive React TypeScript frontend for the Spotify Clone application.

## Features

- 🎵 Browse and search music tracks
- 📋 Create and manage playlists
- 🔐 User authentication (login/register)
- 🎯 Personalized recommendations
- 📱 Fully responsive design
- 🎨 Spotify-inspired UI

## Tech Stack

- **React 18** with TypeScript
- **React Router** for navigation
- **Axios** for API calls
- **Context API** for state management
- **CSS3** for styling

## Getting Started

### Prerequisites

- Node.js (v14 or higher)
- npm or yarn
- Backend API running on http://localhost:8080

### Installation

1. Install dependencies:
```bash
npm install
```

2. Create a `.env` file in the frontend directory:
```
REACT_APP_API_URL=http://localhost:8080/api/v1
```

3. Start the development server:
```bash
npm start
```

The app will open at [http://localhost:3000](http://localhost:3000)

## Available Scripts

- `npm start` - Runs the app in development mode
- `npm test` - Launches the test runner
- `npm run build` - Builds the app for production
- `npm run eject` - Ejects from Create React App (one-way operation)

## Project Structure

```
frontend/
├── public/              # Static files
├── src/
│   ├── components/      # Reusable components
│   │   ├── Navbar.tsx
│   │   └── TrackCard.tsx
│   ├── context/         # React Context providers
│   │   └── AuthContext.tsx
│   ├── pages/           # Page components
│   │   ├── Home.tsx
│   │   ├── Login.tsx
│   │   ├── Register.tsx
│   │   ├── Playlists.tsx
│   │   ├── PlaylistDetail.tsx
│   │   └── Recommendations.tsx
│   ├── services/        # API services
│   │   └── api.ts
│   ├── types/           # TypeScript types
│   │   └── index.ts
│   ├── App.tsx          # Main app component
│   └── index.tsx        # Entry point
└── package.json
```

## Features Overview

### Home Page
- Browse all tracks with pagination
- Search tracks by title or artist
- Filter by genre
- View trending tracks

### Authentication
- Register with email, username, and favorite genres
- Login with email and password
- Persistent authentication with JWT tokens

### Playlists
- Create custom playlists
- Add/remove tracks from playlists
- View playlist details
- Delete playlists

### Recommendations
- Personalized track recommendations
- Based on listening history and preferences
- Genre-based recommendations

## API Integration

The frontend communicates with the backend API using Axios. All API calls are centralized in `src/services/api.ts`:

- **Auth API**: Register, Login
- **Tracks API**: Get tracks, search, record plays
- **Artists API**: Get artists, artist stats
- **Albums API**: Get albums, album stats
- **Playlists API**: CRUD operations for playlists
- **Recommendations API**: Get personalized and trending tracks

## Styling

The app uses custom CSS with a dark theme inspired by Spotify:

- CSS variables for consistent theming
- Responsive design with mobile-first approach
- Smooth transitions and hover effects
- Grid and Flexbox layouts

## Contributing

1. Create a feature branch
2. Make your changes
3. Test thoroughly
4. Submit a pull request

## License

MIT License

