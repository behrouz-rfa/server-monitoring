package visits

import (
	"server-monitoring/databases"
	"server-monitoring/domain/model"
	loogers "server-monitoring/utils/looger"
)

const queryInsertVisited = "INSERT INTO visited(url,ip,user_id)VALUES (?,?,?);"

func (l *Visit) Insert() error {
	stmt, err := databases.Client.Prepare(queryInsertVisited)
	if err != nil {
		loogers.Error("Error while preparing insert visited", err)
		return model.ErrCode
	}
	defer stmt.Close()

	_, errInsert := stmt.Exec(l.Url, l.Ip, l.UserId)

	if errInsert != nil {
		loogers.Error("Error while  cv insert visited", err)
		return model.ErrCode
	}


	return nil
}
