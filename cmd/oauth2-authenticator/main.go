package main

import (
	"fmt"

	"github.com/docker-slim/oauth2-authenticator/pkg/authenticator"
)

func main() {
	app := authenticator.New(nil)
	app.Run()

	if app.Input.Verbose {
		fmt.Println("ACCESS TOKEN:")
	}

	fmt.Printf("%s\n", app.AccessToken)

	if app.Input.ShowUser {
		if app.Input.Verbose {
			fmt.Println("USER INFO:")
		}

		fmt.Printf("%s\n", app.UserJSON())
	}
}

/*
type flags struct {
	verbose bool
	browser string
	provider string
	proto string
	address string
	port int
	callback string
	id string
	secret string
	scopeList string
	scopes []string
	showUser bool
}

func getFlags() flags {
	var input flags

	flag.BoolVar(&input.verbose, "verbose", false, "verbose output")
	flag.StringVar(&input.browser, "browser", "default", "'default', one of the supported browsers or 'manual'")
	flag.StringVar(&input.provider, "provider", "github", "one of the supported oauth providers or 'manual'")
    flag.StringVar(&input.proto, "proto", "http", "server protocol - http or https")
    flag.StringVar(&input.address, "address", "localhost", "server address")
    flag.IntVar(&input.port, "port", 3000, "server port")
    flag.StringVar(&input.callback, "callback", "/callback", "callback handler path")
    flag.StringVar(&input.id, "id", "", "oauth app client ID")
    flag.StringVar(&input.secret, "secret", "", "oauth app client secret")
    flag.StringVar(&input.scopeList, "scopes", "", "comman separated list of scopes")
    flag.BoolVar(&input.showUser, "showuser", false, "show user information")

    flag.Parse()
    if len(input.scopeList) > 0 {
    	input.scopes = strings.Split(input.scopeList,",")
	}

	if len(input.id) == 0 {
		input.id = os.Getenv("CLIENT_ID")
	}

	if len(input.secret) == 0 {
		input.secret = os.Getenv("CLIENT_SECRET")
	}

    if !browser.IsValid(input.browser) {
    	panic(fmt.Errorf("Invalid browser type"))
    }

    return input
}

func main() {
  	input := getFlags()

	callbackURL := fmt.Sprintf("%s://%s:%d%s",
		input.proto, input.address, input.port, input.callback)

	providers := map[string] goth.Provider {
		"github": github.New(input.id, input.secret, callbackURL, input.scopes...),
		"gitlab": gitlab.New(input.id, input.secret, callbackURL, input.scopes...),
		"bitbucket": bitbucket.New(input.id, input.secret, callbackURL, input.scopes...),
		"google": google.New(input.id, input.secret, callbackURL, input.scopes...),
		"amazon": amazon.New(input.id, input.secret, callbackURL, input.scopes...),
	}

	if _, found := providers[input.provider]; !found {
		panic(fmt.Errorf("Unknown provider"))
	}

	state := newState()
	sess, err := providers[input.provider].BeginAuth(state)
	if err != nil {
		panic(err)
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		panic(err)
	}

	if input.browser == "manual" {
		fmt.Printf("Open this URL in a browser: %s\n", url)
	} else {
		if input.verbose {
			fmt.Printf("Opening auth URL (%s)...\n", url)
		}

		err = browser.OpenURL(url, input.browser, 0)
		if err != nil {
			panic(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)

	http.HandleFunc(input.callback,
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

		accessToken, err := sess.Authorize(providers[input.provider], params)
		if err != nil {
			panic(err)
		}

		if input.verbose {
			fmt.Println("ACCESS TOKEN:")
		}

		fmt.Printf("%s\n", accessToken)
		w.Write([]byte(accessToken))

		if input.showUser {
			user, err := providers[input.provider].FetchUser(sess)
			if err != nil {
				fmt.Printf("Error getting user info - %v\n", err)
				return
			}

			userStr, err := json.MarshalIndent(user, "", "  ")
			if err != nil {
				fmt.Printf("JSON error - %v\n", err)
				return
			}

			fmt.Printf("USER INFO:\n%s\n", userStr)
		}
	})

	server := http.Server{
		Addr: fmt.Sprintf(":%d", input.port),
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

func newState() string{
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(nonceBytes)
}
*/
