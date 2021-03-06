package celeritas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/frangar97/celeritas/render"
	"github.com/frangar97/celeritas/session"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Celeritas struct {
	AppName   string
	Debug     bool
	Version   string
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	RoothPath string
	Routes    *chi.Mux
	Render    *render.Render
	JetViews  *jet.Set
	Session   *scs.SessionManager
	config    config
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
}

func (c *Celeritas) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := c.Init(pathConfig)

	if err != nil {
		return err
	}

	err = c.checkDotEnv(rootPath)

	if err != nil {
		return err
	}

	//leemos el .env
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	//creamos logs
	infoLog, errorLog := c.startLoggers()
	c.InfoLog = infoLog
	c.ErrorLog = errorLog

	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = version
	c.RoothPath = rootPath
	c.Routes = c.routes().(*chi.Mux)

	c.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSISTS"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
	}

	sess := session.Session{
		CookieLifetime: c.config.cookie.lifetime,
		CookiePersist:  c.config.cookie.persist,
		CookieName:     c.config.cookie.name,
		SessionType:    c.config.sessionType,
		CookineDomain:  c.config.cookie.domain,
	}

	c.Session = sess.InitSession()

	views := jet.NewSet(jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rootPath)), jet.InDevelopmentMode())

	c.JetViews = views

	c.createRenderer()

	return nil
}

func (c *Celeritas) Init(p initPaths) error {
	root := p.rootPath

	for _, path := range p.folderNames {
		//creamos folder si no existe
		err := c.CreateDirIfNotExist(root + "/" + path)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Celeritas) ListenAndServe() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     c.ErrorLog,
		Handler:      c.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	c.InfoLog.Printf("Servidor iniciado en el puerto %s", os.Getenv("PORT"))
	err := srv.ListenAndServe()
	c.ErrorLog.Fatal(err)
}

func (c *Celeritas) checkDotEnv(path string) error {
	err := c.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))

	if err != nil {
		return err
	}

	return nil
}

func (c *Celeritas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	infoLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (c *Celeritas) createRenderer() {
	myRenderer := render.Render{
		Renderer:  c.config.renderer,
		RoothPath: c.RoothPath,
		Port:      c.config.port,
		JetViews:  c.JetViews,
	}
	c.Render = &myRenderer

}
