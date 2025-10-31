-- Spotify Clone - Triggers, Procedures, and Functions
-- This file contains advanced MySQL database objects for the music catalog

-- ======================
-- UTILITY TABLES
-- ======================

-- Create a table to track album statistics (needed for triggers)
CREATE TABLE IF NOT EXISTS album_stats (
    album_id INT PRIMARY KEY,
    track_count INT DEFAULT 0,
    total_duration INT DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE
);

-- Create a table to track play counts (for trigger demo)
CREATE TABLE IF NOT EXISTS track_stats (
    track_id INT PRIMARY KEY,
    play_count INT DEFAULT 0,
    last_played TIMESTAMP NULL,
    FOREIGN KEY (track_id) REFERENCES tracks(id) ON DELETE CASCADE
);

-- ======================
-- TRIGGER 1: Update Album Stats When Track is Added
-- ======================
DELIMITER $$

DROP TRIGGER IF EXISTS after_track_insert$$
CREATE TRIGGER after_track_insert
AFTER INSERT ON tracks
FOR EACH ROW
BEGIN
    -- Insert or update album stats
    INSERT INTO album_stats (album_id, track_count, total_duration)
    VALUES (NEW.album_id, 1, NEW.duration)
    ON DUPLICATE KEY UPDATE
        track_count = track_count + 1,
        total_duration = total_duration + NEW.duration;

    -- Initialize track stats
    INSERT INTO track_stats (track_id, play_count)
    VALUES (NEW.id, 0);
END$$

DELIMITER ;

-- ======================
-- TRIGGER 2: Update Album Stats When Track is Deleted
-- ======================
DELIMITER $$

DROP TRIGGER IF EXISTS after_track_delete$$
CREATE TRIGGER after_track_delete
AFTER DELETE ON tracks
FOR EACH ROW
BEGIN
    -- Update album stats
    UPDATE album_stats
    SET track_count = track_count - 1,
        total_duration = total_duration - OLD.duration
    WHERE album_id = OLD.album_id;

    -- Delete track stats
    DELETE FROM track_stats WHERE track_id = OLD.id;
END$$

DELIMITER ;

-- ======================
-- FUNCTION 1: Calculate Album Total Duration
-- ======================
DELIMITER $$

DROP FUNCTION IF EXISTS get_album_duration$$
CREATE FUNCTION get_album_duration(album_id_param INT)
RETURNS INT
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE total INT;

    SELECT COALESCE(SUM(duration), 0)
    INTO total
    FROM tracks
    WHERE album_id = album_id_param;

    RETURN total;
END$$

DELIMITER ;

-- ======================
-- PROCEDURE 1: Add Track with Validation
-- ======================
DELIMITER $$

DROP PROCEDURE IF EXISTS add_track$$
CREATE PROCEDURE add_track(
    IN p_title VARCHAR(255),
    IN p_artist_id INT,
    IN p_album_id INT,
    IN p_duration INT,
    IN p_genre VARCHAR(100),
    IN p_release_date DATE,
    IN p_file_url VARCHAR(500),
    IN p_cover_url VARCHAR(500),
    OUT p_track_id INT,
    OUT p_status VARCHAR(100)
)
BEGIN
    DECLARE artist_exists INT;
    DECLARE album_exists INT;
    DECLARE album_artist_id INT;

    -- Validate artist exists
    SELECT COUNT(*) INTO artist_exists
    FROM artists
    WHERE id = p_artist_id;

    IF artist_exists = 0 THEN
        SET p_status = 'ERROR: Artist does not exist';
        SET p_track_id = NULL;
    ELSE
        -- Validate album exists and belongs to the artist
        SELECT COUNT(*), MAX(artist_id) INTO album_exists, album_artist_id
        FROM albums
        WHERE id = p_album_id;

        IF album_exists = 0 THEN
            SET p_status = 'ERROR: Album does not exist';
            SET p_track_id = NULL;
        ELSEIF album_artist_id != p_artist_id THEN
            SET p_status = 'ERROR: Album does not belong to the specified artist';
            SET p_track_id = NULL;
        ELSE
            -- Insert the track
            INSERT INTO tracks (title, artist_id, album_id, duration, genre, release_date, file_url, cover_url)
            VALUES (p_title, p_artist_id, p_album_id, p_duration, p_genre, p_release_date, p_file_url, p_cover_url);

            SET p_track_id = LAST_INSERT_ID();
            SET p_status = 'SUCCESS: Track added successfully';
        END IF;
    END IF;
END$$

DELIMITER ;

-- ======================
-- PROCEDURE 2: Get Artist Statistics
-- ======================
DELIMITER $$

DROP PROCEDURE IF EXISTS get_artist_stats$$
CREATE PROCEDURE get_artist_stats(IN p_artist_id INT)
BEGIN
    DECLARE artist_name_var VARCHAR(255);

    -- Get artist name
    SELECT name INTO artist_name_var
    FROM artists
    WHERE id = p_artist_id;

    -- Return comprehensive statistics
    SELECT
        p_artist_id AS artist_id,
        artist_name_var AS artist_name,
        COUNT(DISTINCT a.id) AS total_albums,
        COUNT(DISTINCT t.id) AS total_tracks,
        COALESCE(SUM(t.duration), 0) AS total_duration_seconds,
        ROUND(COALESCE(SUM(t.duration), 0) / 60.0, 2) AS total_duration_minutes,
        COUNT(DISTINCT t.genre) AS unique_genres,
        GROUP_CONCAT(DISTINCT t.genre ORDER BY t.genre SEPARATOR ', ') AS genres,
        MIN(t.release_date) AS first_release,
        MAX(t.release_date) AS latest_release,
        COALESCE(AVG(ts.play_count), 0) AS avg_plays_per_track
    FROM artists ar
    LEFT JOIN albums a ON ar.id = a.artist_id
    LEFT JOIN tracks t ON ar.id = t.artist_id
    LEFT JOIN track_stats ts ON t.id = ts.track_id
    WHERE ar.id = p_artist_id
    GROUP BY ar.id, ar.name;
END$$

DELIMITER ;

-- ======================
-- INITIALIZE EXISTING DATA
-- ======================

-- Initialize album_stats for existing tracks
INSERT INTO album_stats (album_id, track_count, total_duration)
SELECT album_id, COUNT(*) as track_count, SUM(duration) as total_duration
FROM tracks
GROUP BY album_id
ON DUPLICATE KEY UPDATE
    track_count = VALUES(track_count),
    total_duration = VALUES(total_duration);

-- Initialize track_stats for existing tracks
INSERT INTO track_stats (track_id, play_count)
SELECT id, 0
FROM tracks
ON DUPLICATE KEY UPDATE play_count = play_count;

-- ======================
-- USAGE EXAMPLES
-- ======================

-- Example 1: Add a new track using the stored procedure
-- CALL add_track(
--     'New Song Title',
--     1,
--     1,
--     240,
--     'Pop',
--     '2025-01-01',
--     'https://audio.example.com/new-song.mp3',
--     'https://i.scdn.co/image/cover.jpg',
--     @new_track_id,
--     @status
-- );
-- SELECT @new_track_id AS track_id, @status AS status;

-- Example 2: Get statistics for an artist
-- CALL get_artist_stats(1);

-- Example 3: Use the function to get album duration
-- SELECT id, title, get_album_duration(id) AS total_duration_seconds
-- FROM albums;

-- Example 4: View album statistics
-- SELECT
--     a.id,
--     a.title,
--     ast.track_count,
--     ast.total_duration,
--     ROUND(ast.total_duration / 60.0, 2) AS duration_minutes
-- FROM albums a
-- LEFT JOIN album_stats ast ON a.id = ast.album_id;

