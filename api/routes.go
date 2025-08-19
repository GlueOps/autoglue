package api

import (
	"github.com/glueops/autoglue/api/handlers/authentication"
	"github.com/glueops/autoglue/api/handlers/clusters"
	"github.com/glueops/autoglue/api/handlers/credentials"
	"github.com/glueops/autoglue/api/handlers/health"
	"github.com/glueops/autoglue/api/handlers/orgs"
	"github.com/glueops/autoglue/api/middleware"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func RegisterRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/healthz", health.Check).Methods("GET")

	v1Router := apiRouter.PathPrefix("/v1").Subrouter()

	jwtSecret := viper.GetString("authentication.jwt_secret")
	auth := middleware.AuthMiddleware(jwtSecret)

	authRouter := v1Router.PathPrefix("/authentication").Subrouter()
	authRouter.HandleFunc("/login", authentication.Login).Methods("POST")
	authRouter.HandleFunc("/register", authentication.Register).Methods("POST")

	authPrivate := v1Router.PathPrefix("/authentication").Subrouter()
	authPrivate.Use(auth)
	authPrivate.HandleFunc("/refresh", authentication.Refresh).Methods("POST")
	authPrivate.HandleFunc("/logout", authentication.Logout).Methods("POST")
	authPrivate.HandleFunc("/me", authentication.Me).Methods("GET")

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

	credentialsRouter := v1Router.PathPrefix("/credentials").Subrouter()
	credentialsRouter.Use(auth)
	credentialsRouter.HandleFunc("", credentials.ListCredentials).Methods("GET")
	credentialsRouter.HandleFunc("", credentials.CreateCredential).Methods("POST")
	credentialsRouter.HandleFunc("/{id}", credentials.GetCredential).Methods("GET")
	credentialsRouter.HandleFunc("/{id}", credentials.UpdateCredential).Methods("PATCH")
	credentialsRouter.HandleFunc("/{id}", credentials.DeleteCredential).Methods("DELETE")

}
