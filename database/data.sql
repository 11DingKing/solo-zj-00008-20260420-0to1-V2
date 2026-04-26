INSERT INTO users (username, email, password_hash) VALUES 
('demo_user', 'demo@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'),
('test_user', 'test@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy');

INSERT INTO songs (name, artist, album, duration, cover_url, audio_file_url) VALUES
('Bohemian Rhapsody', 'Queen', 'A Night at the Opera', 355, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Queen%20Bohemian%20Rhapsody%20album%20cover%20classic%20rock&image_size=square', 'https://example.com/audio/bohemian.mp3'),
('Stairway to Heaven', 'Led Zeppelin', 'Led Zeppelin IV', 482, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Led%20Zeppelin%20IV%20album%20cover%20classic%20rock&image_size=square', 'https://example.com/audio/stairway.mp3'),
('Hotel California', 'Eagles', 'Hotel California', 391, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Eagles%20Hotel%20California%20album%20cover%20sunset%20motel&image_size=square', 'https://example.com/audio/hotel.mp3'),
('Smells Like Teen Spirit', 'Nirvana', 'Nevermind', 301, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Nirvana%20Nevermind%20album%20cover%20baby%20swimming&image_size=square', 'https://example.com/audio/teen.mp3'),
('Imagine', 'John Lennon', 'Imagine', 183, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=John%20Lennon%20Imagine%20album%20cover%20peaceful%20white&image_size=square', 'https://example.com/audio/imagine.mp3'),
('Like a Rolling Stone', 'Bob Dylan', 'Highway 61 Revisited', 379, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Bob%20Dylan%20Highway%2061%20album%20cover%20folk%20rock&image_size=square', 'https://example.com/audio/rolling.mp3'),
('Purple Haze', 'Jimi Hendrix', 'Are You Experienced', 171, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Jimi%20Hendrix%20purple%20haze%20psychedelic%20guitar&image_size=square', 'https://example.com/audio/purple.mp3'),
('Yesterday', 'The Beatles', 'Help!', 125, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=The%20Beatles%20Yesterday%20album%20cover%20black%20and%20white&image_size=square', 'https://example.com/audio/yesterday.mp3'),
('Respect', 'Aretha Franklin', 'I Never Loved a Man the Way I Love You', 147, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Aretha%20Franklin%20Respect%20soul%20music%20golden&image_size=square', 'https://example.com/audio/respect.mp3'),
('Good Vibrations', 'The Beach Boys', 'Smiley Smile', 235, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Beach%20Boys%20Good%20Vibrations%20sunshine%20beach&image_size=square', 'https://example.com/audio/vibrations.mp3'),
('Sweet Child O Mine', 'Guns N Roses', 'Appetite for Destruction', 356, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Guns%20N%20Roses%20rock%20album%20cover%20skull&image_size=square', 'https://example.com/audio/sweet.mp3'),
('Billie Jean', 'Michael Jackson', 'Thriller', 294, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Michael%20Jackson%20Thriller%20album%20cover%20pop%20icon&image_size=square', 'https://example.com/audio/billie.mp3'),
('Hey Jude', 'The Beatles', 'Hey Jude', 431, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=The%20Beatles%20Hey%20Jude%20album%20cover%20colorful&image_size=square', 'https://example.com/audio/jude.mp3'),
('Every Breath You Take', 'The Police', 'Synchronicity', 254, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=The%20Police%20Synchronicity%20album%20cover%20blue&image_size=square', 'https://example.com/audio/breath.mp3'),
('Rolling in the Deep', 'Adele', '21', 228, 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=Adele%2021%20album%20cover%20soul%20singer&image_size=square', 'https://example.com/audio/rolling2.mp3');

INSERT INTO playlists (name, description, is_public, owner_id) VALUES
('Classic Rock Essentials', 'The greatest classic rock songs of all time', true, 1),
('Chill Vibes', 'Relaxing songs for a quiet evening', true, 2),
('Workout Mix', 'High energy songs for your workout', false, 1),
('Road Trip', 'Perfect songs for a long drive', true, 1),
('Morning Coffee', 'Songs to start your day right', true, 2);

INSERT INTO playlist_songs (playlist_id, song_id, position) VALUES
(1, 1, 0),
(1, 2, 1),
(1, 3, 2),
(1, 4, 3),
(1, 5, 4),
(1, 6, 5),
(1, 7, 6),
(1, 8, 7),
(2, 5, 0),
(2, 8, 1),
(2, 9, 2),
(2, 15, 3),
(3, 1, 0),
(3, 4, 1),
(3, 7, 2),
(3, 12, 3),
(3, 11, 4),
(4, 1, 0),
(4, 2, 1),
(4, 3, 2),
(4, 11, 3),
(4, 13, 4),
(5, 5, 0),
(5, 8, 1),
(5, 10, 2),
(5, 14, 3);
