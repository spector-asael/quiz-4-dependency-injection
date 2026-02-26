package handler

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
)

func (a *ApplicationDependencies)Routes() http.Handler  {

   // setup a new router
   router := httprouter.New()
   router.NotFound = http.HandlerFunc(a.notFoundResponse)
   router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)
   // setup routes
   (router.HandlerFunc(http.MethodPost, "/v1/transfer", a.transferHandler))

   return a.recoverPanic(router)      
  
}
