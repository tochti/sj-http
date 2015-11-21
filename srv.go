package main

import (
	"fmt"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tochti/gin-angular-kauth"
	"github.com/tochti/gin-static"
	"github.com/tochti/sj-lib"
	"github.com/tochti/smem"
)

func main() {
	app, err := sj.NewApp("sj")
	handleErr(err)

	store := smem.NewStore()
	userStore := sj.NewUserStore(app.DB)

	signedIn := kauth.SignedIn(&store)
	signIn := kauth.SignIn(&store, userStore)

	newSeries := sj.NewAppHandler(app, sj.NewSeriesHandler)
	readSeries := sj.NewAppHandler(app, sj.ReadSeriesHandler)
	removeSeries := sj.NewAppHandler(app, sj.RemoveSeriesHandler)

	appendSeries := sj.NewAppHandler(app, sj.AppendSeriesListHandler)
	readSeriesList := sj.NewAppHandler(app, sj.ReadSeriesListHandler)

	lastWatchedList := sj.NewAppHandler(app, sj.LastWatchedListHandler)
	updatedLastWatched := sj.NewAppHandler(app, sj.UpdateLastWatchedHandler)

	newUser := sj.NewAppHandler(app, sj.NewUserHandler)

	publicDir := app.Specs.PublicDir
	htmlDir := path.Join(publicDir, "html")

	srv := gin.New()
	srv.POST("/Series", signedIn(newSeries))
	srv.GET("/Series/:id", signedIn(readSeries))
	srv.DELETE("/Series/:id", signedIn(removeSeries))
	srv.GET("/SignIn/:name/:password", signIn)
	srv.POST("/User", newUser)
	srv.PATCH("/AppendSeries", signedIn(appendSeries))
	srv.GET("/ReadSeriesList", signedIn(readSeriesList))
	srv.POST("/LastWatched", signedIn(updatedLastWatched))
	srv.GET("/LastWatchedList", signedIn(lastWatchedList))

	srv.Use(ginstatic.Serve("/", ginstatic.LocalFile(htmlDir, false)))
	srv.Static("/public", publicDir)
	srv.Static("/images", app.Specs.ImageDir)

	addr := fmt.Sprintf("%v:%v", app.Specs.Host, app.Specs.Port)
	err = srv.Run(addr)
	handleErr(err)

	os.Exit(0)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
