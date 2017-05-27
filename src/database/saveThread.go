package database

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// ChanThreadPage - TODO
type ChanThreadPage struct {
	Posts []chanThreadPost
}

type chanThreadPost struct {
	No       int64
	Now      string
	Name     string
	Com      string
	Filename string
	Ext      string
	W        int64
	H        int64
	Tn_W     int64
	Tn_H     int64
	Tim      int64
	Time     int64
}

type threadSave struct {
	chanThreadPost
	user uint
}

// ThreadSaveTest - TODO
type ThreadSaveTest struct {
	gorm.Model
	Board  string `gorm:"size:255"`
	Thread string `gorm:"size:255"`
	User   string
}

// PostSaveTest - TODO
type PostSaveTest struct {
	gorm.Model
	Thread   string `gorm:"size:255"`
	No       int64
	Now      string
	Name     string
	Com      string
	Filename string
	Ext      string
	Tim      int64
	Time     int64
}

// GetSavedThreads - Returns a unqiue array<string> of boards saved by the user
//									 in the form of `g`, not `/g/`
func GetSavedThreads(userID string) ([]string, error) {
	var saved []string

	// Connect to the database.
	db, err := Connect()
	if err != nil {
		return saved, err
	}

	rows, err := db.Table("thread_save_tests").Select("board").Where("user = ?", userID).Group("thread").Rows()
	if err != nil {
		return saved, err
	}

	for rows.Next() {
		var board string
		rows.Scan(&board)
		saved = append(saved, board)
	}

	defer db.Close()

	return saved, nil
}

// SaveThread - TODO
func SaveThread(ID string, boardString string, threadString string, threadData *ChanThreadPage) error {
	// Connect to the database.
	db, err := Connect()
	if err != nil {
		return err
	}
	// Check if the current user has saved
	// the current thread before.
	count := 0
	db.Model(&ThreadSaveTest{}).Where(`
		board = ? AND thread = ? AND user = ?
	`, boardString, threadString, ID).Count(&count)
	if count != 0 {
		// The current user has saved the thread already.
		// Bail out!!!
		return errors.New("you have already saved this thread")
	}
	// Check if the thread has been saved already by a different user.
	// If it has lets just create the link and bail out.
	// Do NOT add posts to the database.
	db.Model(&ThreadSaveTest{}).Where(`
		board = ? AND thread = ?
	`, boardString, threadString).Count(&count)
	// Time to add the link
	threadRow := ThreadSaveTest{Board: boardString, Thread: threadString, User: ID}
	db.Create(&threadRow)
	// Check if we should bail or bote.
	if count != 0 {
		// This thread has already been saved by another user.
		// Do not copy posts, just create thread link for user.
		return nil
	}
	// We need to add the posts.
	// Add posts async, this is a heavy task.
	// There can be up to 300 posts.
	go func(posts []chanThreadPost) {
		for i := 0; i < len(posts); i++ {
			go func(post chanThreadPost) {
				row := PostSaveTest{
					Thread:   threadString,
					No:       post.No,
					Now:      post.Now,
					Name:     post.Name,
					Com:      post.Com,
					Filename: post.Filename,
					Ext:      post.Ext,
					Tim:      post.Tim,
					Time:     post.Time,
				}
				db.Create(&row)
				if i == len(posts)-1 {
					defer db.Close()
				}
			}(posts[i])
		}
	}(threadData.Posts)
	// Return eagerly.
	return nil
}
