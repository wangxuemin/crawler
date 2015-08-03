package model

import (
	"env"
	"proto"
)

type Novel struct {
	Short_id    int    `short_id`
	Book_id     int    `book_id`
	Raw_book_id int    `raw_book_id`
	Dir_url     string `dir_url`
	Cp_id       int    `cp_id`
	//Cp_name   string `cp_name`
	Gid       int    `gid`
	Book_name string `book_name`
	Author_id int    `author_id`
	//Author   string `author`
	Channel           string `channel`
	Category          string `category`
	Chapter_price     int    `chapter_price`
	Tag               string `tag`
	Description       string `description`
	Cp_logo           string `cp_logo`
	Logo              string `logo`
	Cp_exclusive_flag int    `cp_exclusive_flag`
	Save_content      int    `save_content`
	//Language string `language`
	//Format string `format`
	//Roll_number    int `roll_number`
	//Chapter_number int `chapter_number`
	Wordsum        int `word_sum`
	Cp_update_time int `cp_update_time`
	Create_time    int `create_time`
	Update_time    int `update_time`
}

func NewFromProto(url string, cp Cp, author Author,
	novel_info proto.NovelInfo) *Novel {
	novel := &Novel{
		Raw_book_id:       novel_info.Novel_id,
		Dir_url:           url,
		Cp_id:             cp.Cp_id,
		Book_name:         novel_info.Novel_name,
		Gid:               GetGid(novel_info.Novel_name, author.Author_name),
		Author_id:         author.Author_id,
		Channel:           novel_info.Channel,
		Category:          novel_info.Category,
		Tag:               novel_info.Tag,
		Chapter_price:     novel_info.Chapter_price,
		Description:       novel_info.Description,
		Cp_logo:           novel_info.Logo,
		Logo:              novel_info.Logo,
		Cp_exclusive_flag: 0, //TODO wtf
		Save_content:      novel_info.Save_content,
	}

	return novel
}

func GetAllNovelsOfCP(cp_id int) ([]*Novel, error) {
	stmt, err := env.Db.Prepare(
		`SELECT book_id, raw_book_id, cp_update_time
         FROM novel_basic_info
         WHERE cp_id=?`)
	if err != nil {
		env.Log.Warn("prepare error : %s", err.Error())
		return nil, err
	}

	defer stmt.Close()

	var novels []*Novel
	rows, err := stmt.Query(cp_id)
	if err != nil {
		env.Log.Warn("get rows of %d error: %s", cp_id, err.Error())
		return nil, err
	}

	for rows.Next() {
		novel := &Novel{}
		err := rows.Scan(&novel.Book_id, &novel.Raw_book_id, &novel.Cp_update_time)
		if err != nil {
			env.Log.Warn(err.Error())
			return nil, err
		}
		novels = append(novels, novel)
	}

	return novels, nil
}

func GetNovelFromCp(cpid, rawbookid int) (*Novel, *error) {
	stmt, err := env.Db.Prepare(
		`SELECT short_id, book_id, raw_book_id, dir_url, cp_id, gid,
                book_name, author_id, channel, category, tag, description,
                logo, cp_exclusive_flag, save_content,
                word_sum, cp_update_time, create_time, update_time
         FROM novel_basic_info
         WHERE cp_id=? and rawk_book_id=?`)
	if err != nil {
		env.Log.Warn("prepare sql error : %s", err.Error())
		return nil, err
	}
	defer stmt.Close()

	novel := &Novel{}

	err = stmt.QueryRow(cpid, rawbookid).
		Scan(novel.Short_id, novel.Book_id, novel.Raw_book_id, novel.Dir_url, novel.Cp_id,
		novel.Gid, novel.Book_name, novel.Author_id, novel.Channel, novel.Category,
		nove.Tag, novel.Description, novel.Logo, novel.Cp_exclusive_flag, novel.Save_content,
		novel.Word_sum, novel.Cp_update_time, novel.Create_time, novel.Update_time)
	if err != nil {
		env.Log.Warn("get novel[%d:%d] error: %s", cpid, rawbookid, err.Error())
		return nil, err
	}

	return novel, nil
}

func (novel *Novel) Insert() error {
	ret, err := env.Db.Exec(
		`INSERT INTO novel_basic_info(raw_book_id, dir_url, cp_id, gid,
         book_name, author_id, channel, category, tag, description, cp_logo,
         logo, cp_exclusive_flag, save_content, word_sum, update_time,
         create_time, cp_update_time
         VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, from_unixtime(), 
         from_unixtime(), ?)`,
		novel.Raw_book_id, novel.Dir_url, novel.Cp_id, novel.Gid, novel.Book_name,
		novel.Author_id, novel.Channel, novel.Category, novel.Tag, novel.Description,
		novel.Cp_logo, novel.Cp_logo, novel.Cp_exclusive_flag, novel.Save_content,
		novel.Word_sum, novel.Cp_update_time)

	if err != nil {
		env.Log.Warn("[SQL][INSERT][raw_book_id:%d][error:%s]",
			novel.Raw_book_id, err.Error())
		return error
	}

	newid, err := ret.LastInsertId()
	if err != nil {
		env.Log.Warn("[SQL][INSERT][raw_book_id:%d][error:%s]",
			novel.Raw_book_id, err.Error())
		return err
	}

	novel.Book_id = GetGid(short_id)
	ret, err := env.Db.Exec(
		`UPDATE novel_basic_info SET book_id=? WHERE raw_book_id=?`,
		novel.Raw_book_id, novel.Book_id)

	if err != nil {
		env.Log.Warn("[SQL][INSERT][raw_book_id:%d][error:%s]",
			novel.Raw_book_id, err.Error())
		return err
	}

	//TODO check affected rows?
	return nil
}

func (novel *Novel) UpdateEssential() error {
	ret, err := env.Db.Exec(
		`UPDATE novel_basic_info SET dir_url=?, gid=?, book_name=?, author_id=?
         channel=?, category=?, tag=?, description=?, cp_logo=?, cp_exclusive_flag=?,
         save_content=?, word_sum=?, cp_update_time=from_unixtime(), 
         WHERE book_id=?`,
		novel.Dir_url, novel.Gid, novel.Book_name, novel.Author_id, novel.Channel,
		novel.Category, novel.Tag, novel.Description, novel.Cp_logo, novel.Cp_exclusive_flag,
		novel.Save_content, novel.Word_sum, novel.Book_id)

	if err != nil {
		env.Log.Warn("[SQL][UPDATE][book_id:%d][error:%s]",
			novel.Book_id, err.Error())
		return err
	}

	affected, err := ret.RowsAffected()
	if err != nil {
		env.Log.Warn("get update result error : %s", err.Error())
		return err
	}

	if affected != 1 {
		env.Log.Warn("affected [%d] rows, rediculous", affected)
		return errors.New{"unknown upate error"}
	}

	return nil
}
