package databases

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	loogers "server-monitoring/utils/looger"

	"log"
	"os"
	"strings"
)

var (
	Client   *sql.DB
	username = "root"
	password = "mitQWERTY"
	host  ="localhost"
	schema= "pool"
)

func init() {

	name, err1 := os.Hostname()
	if err1 != nil {
		//panic(err1)
	}
	if !strings.Contains(name,"DESKTOP") {
		username = "master"
		password = "mitQWERTY"
	}
	fmt.Println("hostname:", name)

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		username, password, host,"3306", schema,
	)
	var err error
	Client, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(); err != nil {
		panic(err)
	}

	_ = mysql.SetLogger(loogers.GetLogger())
	log.Println("database successfully configured")
}
