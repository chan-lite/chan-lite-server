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

// SaveThread - TODO
func SaveThread(ID string, boardString string, threadString string, threadData *ChanThreadPage) error {

	// connect to database
	db, err := Connect()
	if err != nil {
		return err
	}

	// 1.
	// Save each individual post, not including image (use zero for Tim, Tn_W, Tn_H)
	// Include reference to thread number
	// Include check for if thread has already been saved
	// 2. Save the thread number, along with user and board
	// Include check if user has already save this thread before

	// Check if thread has been saved already
	count := 0
	db.Model(&ThreadSaveTest{}).Where("board = ? AND thread = ?", boardString, threadString).Count(&count)
	if count != 0 {
		return errors.New("thread has already been saved")
		// Another user (or maybe the same user) has saved this thread
		// If some user error out, else copy data for new user
	}

	threadRow := ThreadSaveTest{Board: boardString, Thread: threadString, User: ID}
	db.Create(&threadRow)

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

	return nil
}
