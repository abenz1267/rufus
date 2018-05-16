# Rufus, Web-Framework

## Disclaimer

I've built Rufus for personal usage but decided to open-source it, in case someone else finds it useful. Feedback and other proposals are welcomed.

[Example project](https://github.com/abenz1267/rufus_example)

## Features

* uses [Chi](https://github.com/go-chi/chi) as router
* can render json, normal SSR websites and SSR templates for isomorphic websites
* handles translation
* automatic SSL via LetsEncrypt for production
* self-signed certs for development
* custom middleware:
  * in-memory response cache (you can plugin your own as well)
  * logger ( [zerolog](https://github.com/rs/zerolog) )
  * redirect www to no-www
  * content-type setting for json, SSR or "stripped" SSR

## To-Do

* probably more testing
* A/B Testing middleware

## Install

`go get -u github.com/abenz1267/rufus`

## Usage

Code examples are taken from the [example project](https://github.com/abenz1267/rufus_example)

### Preperation

To start using Rufus you need a few basic files:

* config.json
* translation.json, if you want to use the translation feature

You can find examples in this project or the example project linked at the top.

### Loading up Rufus

```
app := &rufus.App{}

if err := app.LoadConfigAndRouter(); err != nil {
	log.Fatal(err)
}
```

### Using shipped in-memory response cache

```
cache := &rufus.Cache{}
app.Middleware.Cache = cache
```

### Registering routes

```
// main.go

go registerRoutes(app)
```

```
// routes.go

func registerRoutes(app *rufus.App) {
    for i := 0; i < app.Translation.Amount; i++ {
	language := <-app.RoutesSender

	if language == "default" {
	language = app.Language
	}

	r := chi.NewRouter()
	app.Router.PrependMiddleware(r, app.Server, app.CSPPolicy)

	r.NotFound(handlers.NotFound{App: app, Language:language.Get)

        handlers.Index{App: app, Language: language}.GetRoutes(r)
        handlers.About{App: app, Language: language}.GetRoutes(r)

	app.RoutesReceiver <- r
    }
}
```

### Example handler

```
// handlers/about.go

type About struct {
	App         *rufus.App `json:"-"`
	Language    string     `json:"language,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
}

// GetRoutes registers routes for About handler
func (h About) GetRoutes(r *chi.Mux) {
	var url strings.Builder
	url.WriteString("/")

	if h.App.Translation.Amount > 1 {
		url.WriteString(h.App.TranslateURL("about", h.Language))
	} else {
		url.WriteString("about")
	}

	r.Get(url.String(), h.get)
}

func (h About) get(w http.ResponseWriter, r *http.Request) {
	h.Title = "About Page"
	h.Description = "This is a description"

    resp := rufus.Response{Status: http.StatusOK, TemplateFile: "about", Data: h}

	h.App.Response = resp

	h.App.Render(w, r)
}
```

## Translation

Rufus does a few things:

* handles translated routes
* creates URLs based on the language
* comes with two simple translation functions that can be used in templates as well as in functions
  * translate
    * template: `{{ translate "phrase" "de" }}`
    * function: `App.Translation.Translate("phrase", "de")`
  * translateURL

When a translation file is present it will also change the way routes are registered automatically, f.e.

without translation file: `https://myapp.com/home`

with translation file: `https://myapp.com/de/startseite` or `https://myapp.com/en/home`
