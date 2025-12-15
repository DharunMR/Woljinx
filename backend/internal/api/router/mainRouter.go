package router

import (
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MainRouter(client *mongo.Client) *http.ServeMux {

	mRouter := MovieRoute(client)
	mRouter.Handle("/", UserRoute(client))
	return mRouter

}
