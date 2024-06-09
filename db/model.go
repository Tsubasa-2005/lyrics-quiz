package db

type QuizCounter struct {
	ID            uint `gorm:"primary_key"`
	QuizCounterID int  `gorm:"not null"`
	Count         int  `gorm:"not null"`
}

type LyricsCounter struct {
	ID              uint `gorm:"primary_key"`
	LyricsCounterID uint `gorm:"not null"`
	Count           int  `gorm:"not null"`
}

type Lyrics struct {
	ID       uint   `gorm:"primary_key"`
	LyricsID int    `gorm:"not null unique"`
	Lyrics   string `gorm:"not null"`
}

type Answer struct {
	ID        uint   `gorm:"primary_key"`
	AnswerID  uint   `gorm:"not null unique"`
	MusicName string `gorm:"not null"`
}
