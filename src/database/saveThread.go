package database

type chanThreadPage struct {
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

// SaveThread - TODO
func SaveThread(ID string, threadData chanThreadPage) error {

	// TODO
	// Save each individual post with reference to thread numner
	// The ref will also have user ID data.
	// e z p z
	save := new(threadSave)
	save.No = threadData.No

	db, err := Connect()
	if err != nil {
		return err
	}

	db.Where("ID = ?", ID).First(&save)

	defer db.Close()

	return err
}
