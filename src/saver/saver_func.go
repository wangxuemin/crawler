package saver

import (
	"cspub"
	"env"
	"model"
	"proto"
	"scheduler"
)

func (saver *Saver) entrance(result *cspub.CspubFetchResult) {
	c := result.User_data.(*scheduler.CrawContext)

	entry, err := proto.DecodeEntry(result.Html_body)
	if err != nil {
		env.Log.Warn("get entry [%s] error: %s", result.Target_url, err.Error())
		return
	}

	for _, page := range entry.Page_list {
		page_context := scheduler.CrawContext{
			Level: scheduler.CRAW_LEVEL_PAGE,
			Cpid:  c.Cpid,
		}
		task := &scheduler.RpcTask{
			Target_url: page,
			Context:    page_context,
		}
		go saver.SendTask(task)
	}
}

func (saver *Saver) page(result *cspub.CspubFetchResult) {
	c := result.User_data.(*scheduler.CrawContext)

	page, err := proto.DecodePage(result.Html_body)
	if err != nil {
		env.Log.Warn("get page [%s] error: %s", result.Target_url, err.Error())
		return
	}

	saved_novels, err := model.GetAllNovelsOfCP(c.Cpid)
	if err != nil {
		env.Log.Warn("get all novels of cp [%d] error", c.Cpid)
		return
	}

	//index of novel in db
	saved_novel_map := make(map[int]*model.Novel)
	for _, saved_novel := range saved_novels {
		saved_novel_map[saved_novel.Raw_book_id] = saved_novel
	}

	for _, basic_novel := range *page {
		saved_novel, ok := saved_novel_map[basic_novel.Novel_id]

		if ok {
			if saved_novel.Cp_update_time < basic_novel.Update_time {
				env.Log.Debug("novel [raw:%d] [book_id:%d] not updated, skip",
					saved_novel.Raw_book_id, saved_novel.Book_id)
				continue
			}
		}

		//if not found in old novel or need update
		novel_context := scheduler.CrawContext{
			Level:     scheduler.CRAW_LEVEL_NOVEL,
			Cpid:      c.Cpid,
			Rawbookid: basic_novel.Novel_id,
		}

		task := &scheduler.RpcTask{
			Target_url: basic_novel.Link,
			Context:    novel_context,
		}
		go saver.SendTask(task)
	}
}

func (saver *Saver) novel(result *cspub.CspubFetchResult) {
	c := result.User_data.(*scheduler.CrawContext)

	novel, err := DecodeNovel(result.Html_body)
	if err != nil {
		env.Log.Warn("decode novel [%s] error: %s", result.Target_url, err.Error())
		return
	}

	saved_novel, err := model.GetNovelFromCp(c.Cpid, C.Rawbookid)
	if err != nil {
		env.Log.Warn("decode novel [%d|%d] error: %s", C.Cpid, C.Rawbookid, err.Error())
		return
	}
	if saved_novel == nil {

	}
}
