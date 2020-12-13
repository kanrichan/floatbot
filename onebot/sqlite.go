package onebot

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"yaya/core"
)

func (conf *Yaml) runDB() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[数据库] DB =X=> =X=> Start Error: %v", err)
			WARN("[数据库] DB ==> ==> Sleep")
		}
	}()
	for i, _ := range conf.BotConfs {
		conf.BotConfs[i].DB = openEventDB(conf.BotConfs[i].Bot)
	}
}

func openEventDB(botID int64) *sql.DB {
	CreatePath(AppPath + core.Int2Str(botID))
	db, err := sql.Open("sqlite3", AppPath+core.Int2Str(botID)+"/event.db")
	if err != nil {
		ERROR("[数据库] Open DB ERROR: %v", err)
	}
	table := `
    CREATE TABLE IF NOT EXISTS event (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        self_id INT NULL,
        message_type INT NULL,
        sub_type INT NULL,
        group_id INT NULL,
        user_id INT NULL,
        notice_id INT NULL,
        message TEXT NULL,
        message_num INT NULL,
        message_id INT NULL,
        raw_message BLOB NULL,
        time INT NULL,
        ret INT NULL
    );
    `
	if _, err := db.Exec(table); err != nil {
		ERROR("[数据库] Create DB ERROR: %v", err)
	}
	return db
}

func (xe *XEvent) event2DB(db *sql.DB) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[数据库] DB =X=> =X=> Start Error: %v", err)
			WARN("[数据库] DB ==> ==> Sleep")
		}
	}()
	stmt, err := db.Prepare(`INSERT INTO event(
		self_id,
		message_type,
        sub_type,
        group_id,
        user_id,
        notice_id,
        message,
        message_num,
        message_id,
        raw_message,
        time,
        ret
		) values(?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		ERROR("[数据库] Event =X=> ==> DB ERROR: %v", err)
	}

	res, err := stmt.Exec(
		xe.selfID,
		xe.mseeageType,
		xe.subType,
		xe.groupID,
		xe.userID,
		xe.noticID,
		xe.message,
		xe.messageNum,
		xe.messageID,
		xe.rawMessage,
		xe.time,
		xe.ret,
	)
	if err != nil {
		ERROR("[数据库] Event =X=> ==> DB ERROR: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		ERROR("[数据库] Event =X=> ==> DB ERROR: %v", err)
	}
	xe.cqID = id
}

func db2Mseeage(db *sql.DB, bot int64, id int64) XEvent {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM event where id=%d", id))
	if err != nil {
		ERROR("[数据库] DB =X=> ==> Event ERROR: %v", err)
	}
	defer rows.Close()
	var (
		selfID      int64
		mseeageType int64
		subType     int64
		groupID     int64
		userID      int64
		noticID     int64
		message     string
		messageNum  int64
		messageID   int64
		rawMessage  []byte
		time        int64
		ret         int64
		cqID        int64
	)
	for rows.Next() {
		err = rows.Scan(
			&cqID,
			&selfID,
			&mseeageType,
			&subType,
			&groupID,
			&userID,
			&noticID,
			&message,
			&messageNum,
			&messageID,
			&rawMessage,
			&time,
			&ret,
		)
	}
	if err != nil {
		ERROR("[数据库] DB =X=> ==> Event ERROR: %v", err)
	}
	return XEvent{
		selfID:      selfID,
		mseeageType: mseeageType,
		subType:     subType,
		groupID:     groupID,
		userID:      userID,
		noticID:     noticID,
		message:     message,
		messageNum:  messageNum,
		messageID:   messageID,
		rawMessage:  rawMessage,
		time:        time,
		ret:         ret,
		cqID:        cqID,
	}
}

func (xe *XEvent) xq2cqid(db *sql.DB) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM event where message_num=%d", xe.messageNum))
	if err != nil {
		ERROR("[数据库] DB =X=> ==> Event ERROR: %v", err)
	}
	defer rows.Close()
	var (
		selfID      int64
		mseeageType int64
		subType     int64
		groupID     int64
		userID      int64
		noticID     int64
		message     string
		messageNum  int64
		messageID   int64
		rawMessage  []byte
		time        int64
		ret         int64
		cqID        int64
	)
	for rows.Next() {
		err = rows.Scan(
			&cqID,
			&selfID,
			&mseeageType,
			&subType,
			&groupID,
			&userID,
			&noticID,
			&message,
			&messageNum,
			&messageID,
			&rawMessage,
			&time,
			&ret,
		)
	}
	if err != nil {
		ERROR("[数据库] DB =X=> ==> Event ERROR: %v", err)
	}
	xe.cqID = cqID
}
