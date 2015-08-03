package model

import (
	"env"
)

type Site struct {
	Id          int    `id`
	Site_name   string `site_name"`
	Site_entry  string `site_entry`
	Entry_type  int    `entry_type`
	Has_content int    `has_content`
	/*
		Create_time int    `create_time`
		Update_time int    `update_time`

			Send_email  string `send_email`
			Address     string `address`
			Send_sms    string `send_sms`
			Mobile      string `mobile"`
	*/
}

func GetSiteByName(site_name string) (*Site, error) {
	var stmt *sql.Stmt
	site := &Site{}

	if stmt, err := env.Db.Prepare(
		`SELECT id, site_name, site_entry, entry_type, has_content
         FROM novel_site_info
         WHERE site_name=?`); err != nil {
		env.Log.Warn("[SQL][SELECT][site_name:%s][error:%s]", site_name, err.Error())
		return nil, err
	}
	defer stmt.Close()

	if err := stmt.QueryRow(site_name).
		Scan(&site.Id, &site.Site_name, &site.Site_entry,
		&site.Entry_type, &site.Has_content); err != nil {
		env.Log.Warn("[SQL][SELECT][site_name:%s][error:%s]", site_name, err.Error())
		return nil, err
	}

	return site, nil
}

func GetSites() ([]*Site, error) {
	var rows *sql.Rows

	if rows, err := env.Db.Query(
		`SELECT id, site_name, site_entry, entry_type, has_content 
        FROM novel_site_info`); err != nil {
		env.Log.Warn("[SQL][SELECT][error:%s]", err.Error())
		return nil, err
	}
	defer rows.Close()

	var sites []*Site

	for rows.Next() {
		site := &Site{}
		if err := rows.Scan(
			&site.Id, &site.Site_name, &site.Site_entry,
			&site.Entry_type, &site.Has_content); err != nil {
			env.Log.Warn("[SQL][SELECT][error:%s]", err.Error())
			return nil, err
		}

		sites = append(sites, site)
	}

	return sites, nil
}
