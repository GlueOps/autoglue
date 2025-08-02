package api

import (
	"github.com/glueops/autoglue/api/handlers/authentication"
	"github.com/glueops/autoglue/api/handlers/clusters"
	"github.com/glueops/autoglue/api/handlers/orgs"
	"github.com/glueops/autoglue/api/middleware"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func RegisterRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()

	authRouter := v1Router.PathPrefix("/authentication").Subrouter()
	authRouter.HandleFunc("/login", authentication.Login).Methods("POST")
	authRouter.HandleFunc("/register", authentication.Register).Methods("POST")
	authRouter.HandleFunc("/refresh", authentication.Refresh).Methods("POST")
	authRouter.HandleFunc("/logout", authentication.Logout).Methods("POST")

	jwtSecret := viper.GetString("authentication.jwt_secret")
	auth := middleware.AuthMiddleware(jwtSecret)

	orgsRouter := v1Router.PathPrefix("/orgs").Subrouter()
	orgsRouter.Use(auth)
	orgsRouter.HandleFunc("", orgs.ListOrganizations).Methods("GET")
	orgsRouter.HandleFunc("", orgs.CreateOrganization).Methods("POST")
	orgsRouter.HandleFunc("/switch", orgs.SwitchOrganization).Methods("POST")
	orgsRouter.HandleFunc("/invite", orgs.InviteMember).Methods("POST")
	orgsRouter.HandleFunc("/join", orgs.JoinOrganization).Methods("POST")
	orgsRouter.HandleFunc("/members", orgs.ListMembers).Methods("GET")
	orgsRouter.HandleFunc("/members/{userId}", orgs.DeleteMember).Methods("DELETE")
	orgsRouter.HandleFunc("/{orgId}", orgs.UpdateOrganization).Methods("PATCH")
	orgsRouter.HandleFunc("/{orgId}", orgs.DeleteOrganization).Methods("DELETE")

	clustersRouter := v1Router.PathPrefix("/clusters").Subrouter()
	clustersRouter.Use(auth)
	clustersRouter.HandleFunc("", clusters.CreateCluster).Methods("POST")
	clustersRouter.HandleFunc("", clusters.GetClusters).Methods("GET")

}
