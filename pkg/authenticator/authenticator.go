package authenticator

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"

	"github.com/docker-slim/oauth2-authenticator/pkg/browser"
)

type Flags struct {
	Verbose   bool
	Browser   string
	Provider  string
	Proto     string
	Address   string
	Port      int
	Callback  string
	ID        string
	Secret    string
	ScopeList string
	Scopes    []string
	ShowUser  bool
}

func GetFlags() *Flags {
	var input Flags

	flag.BoolVar(&input.Verbose, "verbose", false, "verbose output")
	flag.StringVar(&input.Browser, "browser", "default", "'default', one of the supported browsers or 'manual'")
	flag.StringVar(&input.Provider, "provider", "github", "one of the supported oauth providers or 'manual'")
	flag.StringVar(&input.Proto, "proto", "http", "server protocol - http or https")
	flag.StringVar(&input.Address, "address", "localhost", "server address")
	flag.IntVar(&input.Port, "port", 3000, "server port")
	flag.StringVar(&input.Callback, "callback", "/callback", "callback handler path")
	flag.StringVar(&input.ID, "id", "", "oauth app client ID")
	flag.StringVar(&input.Secret, "secret", "", "oauth app client secret")
	flag.StringVar(&input.ScopeList, "scopes", "", "comman separated list of scopes")
	flag.BoolVar(&input.ShowUser, "showuser", false, "show user information")

	flag.Parse()
	if len(input.ScopeList) > 0 {
		input.Scopes = strings.Split(input.ScopeList, ",")
	}

	if len(input.ID) == 0 {
		input.ID = os.Getenv("CLIENT_ID")
	}

	if len(input.Secret) == 0 {
		input.Secret = os.Getenv("CLIENT_SECRET")
	}

	if !browser.IsValid(input.Browser) {
		panic(fmt.Errorf("Invalid browser type"))
	}

	return &input
}

type App struct {
	Input       *Flags
	CallbackURL string
	AuthURL     string
	AccessToken string
	UserInfo    *goth.User
	Provider    goth.Provider
}

func New(input *Flags) *App {
	if input == nil {
		input = GetFlags()
	}

	callbackURL := fmt.Sprintf("%s://%s:%d%s",
		input.Proto, input.Address, input.Port, input.Callback)

	app := App{
		Input:       input,
		CallbackURL: callbackURL,
	}

	app.initProvider()
	return &app
}

func (app *App) initProvider() {
	switch app.Input.Provider {
	case "github":
		app.Provider = github.New(app.Input.ID, app.Input.Secret, app.CallbackURL, app.Input.Scopes...)
	case "gitlab":
		app.Provider = gitlab.New(app.Input.ID, app.Input.Secret, app.CallbackURL, app.Input.Scopes...)
	case "bitbucket":
		app.Provider = bitbucket.New(app.Input.ID, app.Input.Secret, app.CallbackURL, app.Input.Scopes...)
	case "google":
		app.Provider = google.New(app.Input.ID, app.Input.Secret, app.CallbackURL, app.Input.Scopes...)
	case "amazon":
		app.Provider = amazon.New(app.Input.ID, app.Input.Secret, app.CallbackURL, app.Input.Scopes...)
	default:
		panic(fmt.Errorf("Unknown provider"))
	}
}

func (app *App) Run() {
	state := newState()
	sess, err := app.Provider.BeginAuth(state)
	if err != nil {
		panic(err)
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		panic(err)
	}

	app.AuthURL = url

	if app.Input.Browser == "manual" {
		fmt.Printf("Open this URL in a browser: %s\n", url)
	} else {
		if app.Input.Verbose {
			fmt.Printf("Opening auth URL (%s)...\n", url)
		}

		err = browser.OpenURL(app.AuthURL, app.Input.Browser, 0)
		if err != nil {
			panic(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	http.HandleFunc(app.Input.Callback,
		func(w http.ResponseWriter, r *http.Request) {
			defer wg.Done()

			if callbackState := r.URL.Query().Get("state"); callbackState != state {
				http.Error(w, fmt.Sprintf("bad state: %s", callbackState), http.StatusUnauthorized)
				return
			}

			params := r.URL.Query()
			if params.Encode() == "" && r.Method == "POST" {
				r.ParseForm()
				params = r.Form
			}

			accessToken, err := sess.Authorize(app.Provider, params)
			if err != nil {
				panic(err)
			}

			app.AccessToken = accessToken
			w.Write([]byte(accessToken))

			if app.Input.ShowUser {
				user, err := app.Provider.FetchUser(sess)
				if err != nil {
					fmt.Printf("Error getting user info - %v\n", err)
					return
				}

				app.UserInfo = &user
			}
		})

	server := http.Server{
		Addr: fmt.Sprintf(":%d", app.Input.Port),
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	wg.Wait()
	if err := server.Shutdown(context.Background()); err != nil {
		panic(err)
	}
}

func (app *App) UserJSON() string {
	if app.UserInfo != nil {
		encoded, err := json.MarshalIndent(app.UserInfo, "", "  ")
		if err != nil {
			fmt.Printf("JSON error - %v\n", err)
			return ""
		}

		return string(encoded)
	}

	return ""
}

func newState() string {
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(nonceBytes)
}
