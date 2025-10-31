-- Spotify Clone Seed Data
-- This file contains sample data for the music catalog in MySQL

-- Insert Genres
INSERT INTO genres (name) VALUES
('Pop'),
('Rock'),
('Hip Hop'),
('R&B'),
('Electronic'),
('Jazz'),
('Classical'),
('Country'),
('Indie'),
('Metal');

-- Insert Artists
INSERT INTO artists (name, bio, image_url) VALUES
('The Weeknd', 'Canadian singer, songwriter, and record producer known for his distinctive voice and dark R&B style.', 'https://i.scdn.co/image/ab6761610000e5eb214f3cf1cbe7139c1e26ffbb'),
('Taylor Swift', 'American singer-songwriter known for narrative songwriting and genre versatility.', 'https://i.scdn.co/image/ab6761610000e5eb859e4c14fa59296c8649e0e4'),
('Drake', 'Canadian rapper, singer, and songwriter, one of the best-selling music artists worldwide.', 'https://i.scdn.co/image/ab6761610000e5eb4293385d324db8558179afd9'),
('Billie Eilish', 'American singer and songwriter known for her unique voice and dark pop style.', 'https://i.scdn.co/image/ab6761610000e5eb6a8f2a838b6c8c5cd64d1b1c'),
('Ed Sheeran', 'English singer-songwriter known for his heartfelt ballads and pop hits.', 'https://i.scdn.co/image/ab6761610000e5eb3bcef85e105dfc42399a26d0'),
('Daft Punk', 'French electronic music duo known for their innovative approach to electronic music.', 'https://i.scdn.co/image/ab6761610000e5eb0672b0e8d72a7e8be4e5e6a6'),
('Kendrick Lamar', 'American rapper and songwriter known for his complex lyricism and social commentary.', 'https://i.scdn.co/image/ab6761610000e5eb437b9e2a82505b3d93ff1022'),
('Adele', 'British singer-songwriter known for her powerful voice and emotional ballads.', 'https://i.scdn.co/image/ab6761610000e5eb4293385d324db8558179afd9'),
('Arctic Monkeys', 'English rock band known for their indie rock sound and witty lyrics.', 'https://i.scdn.co/image/ab6761610000e5eb7da39dea0a72f581535fb11f'),
('Post Malone', 'American rapper, singer, and songwriter known for blending genres.', 'https://i.scdn.co/image/ab6761610000e5ebb6c4e0d8c3f5e6c5e9a9c8b8');

-- Insert Albums
INSERT INTO albums (title, artist_id, release_date, cover_url) VALUES
('After Hours', 1, '2020-03-20', 'https://i.scdn.co/image/ab67616d0000b2738863bc11d2aa12b54f5aeb36'),
('1989', 2, '2014-10-27', 'https://i.scdn.co/image/ab67616d0000b273cd9db6e7cdb5b8c51e9e7c2e'),
('Certified Lover Boy', 3, '2021-09-03', 'https://i.scdn.co/image/ab67616d0000b273cd945b4e3de57edd28481a3f'),
('Happier Than Ever', 4, '2021-07-30', 'https://i.scdn.co/image/ab67616d0000b273d8e1e1d8f5e0a5f5c5e5e5e5'),
('Divide', 5, '2017-03-03', 'https://i.scdn.co/image/ab67616d0000b273ba5db46f4b838ef6027e6f96'),
('Random Access Memories', 6, '2013-05-17', 'https://i.scdn.co/image/ab67616d0000b273d8a1e0e8f5e0a5f5c5e5e5e5'),
('DAMN.', 7, '2017-04-14', 'https://i.scdn.co/image/ab67616d0000b273e0e0e0e8f5e0a5f5c5e5e5e5'),
('25', 8, '2015-11-20', 'https://i.scdn.co/image/ab67616d0000b273d8e8e8e8f5e0a5f5c5e5e5e5'),
('AM', 9, '2013-09-09', 'https://i.scdn.co/image/ab67616d0000b273d8f8f8f8f5e0a5f5c5e5e5e5'),
('Hollywood\'s Bleeding', 10, '2019-09-06', 'https://i.scdn.co/image/ab67616d0000b273d8g8g8g8f5e0a5f5c5e5e5e5');

-- Insert Tracks
INSERT INTO tracks (title, artist_id, album_id, duration, genre, release_date, file_url, cover_url) VALUES
-- The Weeknd - After Hours
('Blinding Lights', 1, 1, 200, 'Pop', '2020-03-20', 'https://audio.example.com/blinding-lights.mp3', 'https://i.scdn.co/image/ab67616d0000b2738863bc11d2aa12b54f5aeb36'),
('Save Your Tears', 1, 1, 215, 'Pop', '2020-03-20', 'https://audio.example.com/save-your-tears.mp3', 'https://i.scdn.co/image/ab67616d0000b2738863bc11d2aa12b54f5aeb36'),
('In Your Eyes', 1, 1, 237, 'R&B', '2020-03-20', 'https://audio.example.com/in-your-eyes.mp3', 'https://i.scdn.co/image/ab67616d0000b2738863bc11d2aa12b54f5aeb36'),

-- Taylor Swift - 1989
('Shake It Off', 2, 2, 219, 'Pop', '2014-10-27', 'https://audio.example.com/shake-it-off.mp3', 'https://i.scdn.co/image/ab67616d0000b273cd9db6e7cdb5b8c51e9e7c2e'),
('Blank Space', 2, 2, 231, 'Pop', '2014-10-27', 'https://audio.example.com/blank-space.mp3', 'https://i.scdn.co/image/ab67616d0000b273cd9db6e7cdb5b8c51e9e7c2e'),
('Style', 2, 2, 231, 'Pop', '2014-10-27', 'https://audio.example.com/style.mp3', 'https://i.scdn.co/image/ab67616d0000b273cd9db6e7cdb5b8c51e9e7c2e'),

-- Drake - Certified Lover Boy
('Way 2 Sexy', 3, 3, 257, 'Hip Hop', '2021-09-03', 'https://audio.example.com/way-2-sexy.mp3', 'https://i.scdn.co/image/ab67616d0000b273cd945b4e3de57edd28481a3f'),
('Girls Want Girls', 3, 3, 248, 'Hip Hop', '2021-09-03', 'https://audio.example.com/girls-want-girls.mp3', 'https://i.scdn.co/image/ab67616d0000b273cd945b4e3de57edd28481a3f'),
('Champagne Poetry', 3, 3, 307, 'Hip Hop', '2021-09-03', 'https://audio.example.com/champagne-poetry.mp3', 'https://i.scdn.co/image/ab67616d0000b273cd945b4e3de57edd28481a3f'),

-- Billie Eilish - Happier Than Ever
('Happier Than Ever', 4, 4, 298, 'Pop', '2021-07-30', 'https://audio.example.com/happier-than-ever.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8e1e1d8f5e0a5f5c5e5e5e5'),
('Therefore I Am', 4, 4, 174, 'Pop', '2021-07-30', 'https://audio.example.com/therefore-i-am.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8e1e1d8f5e0a5f5c5e5e5e5'),
('My Future', 4, 4, 210, 'Pop', '2021-07-30', 'https://audio.example.com/my-future.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8e1e1d8f5e0a5f5c5e5e5e5'),

-- Ed Sheeran - Divide
('Shape of You', 5, 5, 233, 'Pop', '2017-03-03', 'https://audio.example.com/shape-of-you.mp3', 'https://i.scdn.co/image/ab67616d0000b273ba5db46f4b838ef6027e6f96'),
('Perfect', 5, 5, 263, 'Pop', '2017-03-03', 'https://audio.example.com/perfect.mp3', 'https://i.scdn.co/image/ab67616d0000b273ba5db46f4b838ef6027e6f96'),
('Castle on the Hill', 5, 5, 261, 'Pop', '2017-03-03', 'https://audio.example.com/castle-on-the-hill.mp3', 'https://i.scdn.co/image/ab67616d0000b273ba5db46f4b838ef6027e6f96'),

-- Daft Punk - Random Access Memories
('Get Lucky', 6, 6, 369, 'Electronic', '2013-05-17', 'https://audio.example.com/get-lucky.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8a1e0e8f5e0a5f5c5e5e5e5'),
('Instant Crush', 6, 6, 337, 'Electronic', '2013-05-17', 'https://audio.example.com/instant-crush.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8a1e0e8f5e0a5f5c5e5e5e5'),
('Lose Yourself to Dance', 6, 6, 353, 'Electronic', '2013-05-17', 'https://audio.example.com/lose-yourself-to-dance.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8a1e0e8f5e0a5f5c5e5e5e5'),

-- Kendrick Lamar - DAMN.
('HUMBLE.', 7, 7, 177, 'Hip Hop', '2017-04-14', 'https://audio.example.com/humble.mp3', 'https://i.scdn.co/image/ab67616d0000b273e0e0e0e8f5e0a5f5c5e5e5e5'),
('DNA.', 7, 7, 185, 'Hip Hop', '2017-04-14', 'https://audio.example.com/dna.mp3', 'https://i.scdn.co/image/ab67616d0000b273e0e0e0e8f5e0a5f5c5e5e5e5'),
('LOYALTY.', 7, 7, 227, 'Hip Hop', '2017-04-14', 'https://audio.example.com/loyalty.mp3', 'https://i.scdn.co/image/ab67616d0000b273e0e0e0e8f5e0a5f5c5e5e5e5'),

-- Adele - 25
('Hello', 8, 8, 295, 'Pop', '2015-11-20', 'https://audio.example.com/hello.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8e8e8e8f5e0a5f5c5e5e5e5'),
('When We Were Young', 8, 8, 290, 'Pop', '2015-11-20', 'https://audio.example.com/when-we-were-young.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8e8e8e8f5e0a5f5c5e5e5e5'),
('Send My Love', 8, 8, 223, 'Pop', '2015-11-20', 'https://audio.example.com/send-my-love.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8e8e8e8f5e0a5f5c5e5e5e5'),

-- Arctic Monkeys - AM
('Do I Wanna Know?', 9, 9, 272, 'Rock', '2013-09-09', 'https://audio.example.com/do-i-wanna-know.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8f8f8f8f5e0a5f5c5e5e5e5'),
('R U Mine?', 9, 9, 200, 'Rock', '2013-09-09', 'https://audio.example.com/r-u-mine.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8f8f8f8f5e0a5f5c5e5e5e5'),
('Why\'d You Only Call Me When You\'re High?', 9, 9, 161, 'Rock', '2013-09-09', 'https://audio.example.com/whyd-you-only-call-me.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8f8f8f8f5e0a5f5c5e5e5e5'),

-- Post Malone - Hollywood's Bleeding
('Circles', 10, 10, 215, 'Pop', '2019-09-06', 'https://audio.example.com/circles.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8g8g8g8f5e0a5f5c5e5e5e5'),
('Sunflower', 10, 10, 158, 'Pop', '2019-09-06', 'https://audio.example.com/sunflower.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8g8g8g8f5e0a5f5c5e5e5e5'),
('Goodbyes', 10, 10, 175, 'Hip Hop', '2019-09-06', 'https://audio.example.com/goodbyes.mp3', 'https://i.scdn.co/image/ab67616d0000b273d8g8g8g8f5e0a5f5c5e5e5e5');

