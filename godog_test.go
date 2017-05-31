package janus

import (
	"flag"
	"net/url"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/hellofresh/janus/features/bootstrap"
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/config"
	"github.com/hellofresh/janus/pkg/errors"
	"gopkg.in/mgo.v2"
)

var runGoDogTests bool

func init() {
	flag.BoolVar(&runGoDogTests, "godog", false, "Set this flag is you want to run godog BDD tests")
	flag.Bool("random", false, "Randomize features/scenarios execution. Flag is passed to godog")
	flag.Bool("stop-on-failure", false, "Stop processing on first failed scenario.. Flag is passed to godog")
	flag.Parse()
}

func FeatureContext(s *godog.Suite) {
	c, err := config.Load("")
	if nil != err {
		panic(err)
	}

	var apiRepo api.Repository

	dsnURL, err := url.Parse(c.Database.DSN)
	switch dsnURL.Scheme {
	case "mongodb":
		session, err := mgo.Dial(c.Database.DSN)
		if err != nil {
			panic(err)
		}

		session.SetMode(mgo.Monotonic, true)

		apiRepo, err = api.NewMongoAppRepository(session)
		if err != nil {
			panic(err)
		}
	case "file":
		var apiPath = dsnURL.Path + "/apis"

		apiRepo, err = api.NewFileSystemRepository(apiPath)
		if err != nil {
			panic(err)
		}
	default:
		panic(errors.ErrInvalidScheme)
	}

	bootstrap.RegisterRequestContext(s, c.Port, c.Web.Port, c.Web.Credentials)
	bootstrap.RegisterAPIContext(s, c.Web.ReadOnly, apiRepo)
}

func TestMain(m *testing.M) {
	if !runGoDogTests {
		os.Exit(0)
	}

	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
		FeatureContext(s)
	}, godog.Options{
		Format: "progress",
		Paths:  []string{"features"},
	})

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}
