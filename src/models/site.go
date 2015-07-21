package model

import (
	"database/sql"
)

var (
	db
)

type Site struct {
	Id          int    `id`
	Site_name   string `site_name"`
	Site_entry  string `site_entry`
	Entry_type  int    `entry_type`
	Has_content int    `has_content`
	Create_time int    `create_time`
	Update_time int    `update_time`
	send_email  int    `send_email`
	address     string `address`
	send_sms    int    `send_sms`
	mobile      string `mobile"`
}

func (*Site) NewSites(has_content int) {

}
