package main

import (
	"encoding/json"
	"server-monitoring/router"
	"server-monitoring/services"
	"server-monitoring/services/mongoservices"
	"server-monitoring/shared/database"
	"server-monitoring/shared/email"
	"server-monitoring/shared/recaptcha"
	"server-monitoring/shared/server"
	"server-monitoring/shared/session"
	"server-monitoring/shared/translate"
	"server-monitoring/shared/view"
	"server-monitoring/shared/view/plugin"
	"sync"

	"log"
	"os"
	"runtime"
	"server-monitoring/shared/jsonconfig"
)

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	//
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)
	//	app.StartApplication()

	// Configure the session cookie store
	session.Configure(config.Session)

	//if os.Getenv("MONGODB_URI") != "" {
	//	config.Database.MongoDB.Database = os.Getenv("MONGODB_URI")
	//}

	// Connect to database
	database.Connect(config.Database)

	translate.Config("fa_IR")
	//websockets.Setup()
	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)
	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.RootAdmin, config.Template.RootFront, config.Template.Children, config.Template.ChildrenAdmin, config.Template.ChildrenFront)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.Ranges(),
		plugin.Add(),
		plugin.GB(),
		plugin.Uptime(),
		plugin.Sub(),
		plugin.PrettyTime(),
		plugin.Translate(config.View),
		recaptcha.Plugin(),
	)

	//cron
	//c := cron.New()
	////c.AddFunc("*/1 * * * *", func() {
	////	log.Print("[Job 1]Every minute job\n")
	////})
	//c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	//c.Start()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Run(router.LoadHTTP(), router.LoadHTTPS(), config.Server)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		//http_log.Run()
		services.Run()
	}()
	wg.Add(1)
	go func() {
		mongoservices.Run()
	}()
	wg.Wait()
}

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database  database.Info   `json:"Database"`
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
