package model

import (
	"env"
)

type Cp struct {
	Cp_id   int    `cp_id`
	Cp_name string `cp_name`
}

func GetCpBySite(site *Site) (*Cp, error) {
	stmt, err := env.Db.Prepare(
		`SELECT cp_id, cp_name 
         FROM novel_cp_info
         WHERE cp_name=?`)
	if err != nil {
		env.Log.Warn("prepare error: %s", err.Error())
		return nil, err
	}

	defer stmt.Close()
	cp := &Cp{}
	err = stmt.QueryRow(site.Site_name).Scan(&cp.Cp_id, &cp.Cp_name)
	if err != nil {
		env.Log.Warn("get cp %s info error: %s", site.Site_name, err.Error())
		return nil, err
	}

	return cp, nil

}
