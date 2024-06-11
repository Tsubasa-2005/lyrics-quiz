-- Create table for QuizManager
CREATE TABLE quiz_manager (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL UNIQUE,
    the_number_of_questions INTEGER NOT NULL,
    quiz_count INTEGER NOT NULL,
    lyrics_count INTEGER NOT NULL,
    status TEXT NOT NULL,
    type TEXT NOT NULL
);

-- Create table for Lyrics
CREATE TABLE lyrics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_manager_id TEXT NOT NULL,
    question_number INTEGER NOT NULL,
    count INTEGER NOT NULL,
    lyrics TEXT NOT NULL,
    FOREIGN KEY (quiz_manager_id) REFERENCES quiz_manager(user_id)
);

-- Create table for Answer
CREATE TABLE answer (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_manager_id TEXT NOT NULL,
    question_number INTEGER NOT NULL,
    music_name TEXT NOT NULL,
    FOREIGN KEY (quiz_manager_id) REFERENCES quiz_manager(user_id)
);

-- Create table for Artist
CREATE TABLE artist (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    quiz_manager_id TEXT NOT NULL,
    artist TEXT NOT NULL,
    FOREIGN KEY (quiz_manager_id) REFERENCES quiz_manager(user_id)
);
