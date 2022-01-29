package language

import (
	"server-monitoring/databases"
	"server-monitoring/domain/model"
	loogers "server-monitoring/utils/looger"

)

const (
	querySelectCategory     = "SELECT * FROM `language`;"
	querySelectLanguageBYId = "SELECT * FROM `language` WHERE id = ?;"
	queryInsertCategory     = "INSERT INTO `language` (`language`,code,language_code,is_rtl,icon_name)VALUES (?,?,?,?,?);"
)

func (l *Language) FindFist() error {
	stmt, err := databases.Client.Prepare(querySelectLanguageBYId)
	if err != nil {
		loogers.Error("Error while preparing find language", err)
		return model.ErrCode
	}
	defer stmt.Close()
	row := stmt.QueryRow(l.Id)
	if err := row.Scan(&l.Id, &l.Language, &l.Code, &l.LanguageCode, &l.IsRtl, &l.IconName); err != nil {
		loogers.Error("Error while preparing find language", err)
		return model.StandardizeError(err)
	}
	return nil
}

func (l *Language) FindAll() ([]Language, error) {
	stmt, err := databases.Client.Prepare(querySelectCategory)
	if err != nil {
		loogers.Error("Error while preparing find language", err)
		return nil, model.ErrCode
	}
	defer stmt.Close()
	rows, errQuery := stmt.Query()
	if errQuery != nil {
		loogers.Error("Error while preparing find language", err)
		return nil, model.ErrCode
	}
	defer rows.Close()

	var languages []Language
	for rows.Next() {
		var language Language
		if err := rows.Scan(&language.Id, &language.Language, &language.Code,&language.LanguageCode, &language.IsRtl, &language.IconName); err != nil {
			continue
		}
		languages = append(languages, language)
	}
	return languages, nil
}

func (l *Language) Insert() error {
	stmt, err := databases.Client.Prepare(queryInsertCategory)
	if err != nil {
		loogers.Error("Error while preparing insert language", err)
		return model.ErrCode
	}
	defer stmt.Close()

	result, errInsert := stmt.Exec(l.Language, l.Code)

	if errInsert != nil {
		loogers.Error("Error while preparing insert language", err)
		return model.ErrCode
	}
	lastid, err := result.LastInsertId()
	if err != nil {
		loogers.Error("Error while get last id", err)
		return model.ErrCode
	}
	l.Id = int(lastid)
	return nil
}
