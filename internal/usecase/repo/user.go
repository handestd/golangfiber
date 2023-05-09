package repo

import "newanysock/pkg/database/mysql"

func MatchRecord(field string, value interface{}, model interface{}) (bool, modelOut interface{}) {
	rs := mysql.DB.Where(field+" = ?", value).First(model)

	if rs.Error != nil {
		return false, model
	} else if rs.RowsAffected > 0 {
		return true, model
	}
	return false, model
}

func MatchRecordTableName(table, name string, value interface{}, model interface{}) bool {
	rs := mysql.DB.Table(table).Where(name+" = ?", value).First(model)
	if rs.RowsAffected > 0 {
		return true
	}
	return false
}
