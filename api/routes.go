package api

import (
	"github.com/glueops/autoglue/api/handlers/authentication"
	"github.com/glueops/autoglue/api/handlers/clusters"
	"github.com/glueops/autoglue/api/handlers/credentials"
	"github.com/glueops/autoglue/api/handlers/health"
	"github.com/glueops/autoglue/api/handlers/nodepools"
	"github.com/glueops/autoglue/api/handlers/nodetaints"
	"github.com/glueops/autoglue/api/handlers/orgs"
	"github.com/glueops/autoglue/api/handlers/servers"
	"github.com/glueops/autoglue/api/handlers/ssh"
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
	clustersRouter.HandleFunc("", clusters.ListClusters).Methods("GET")
	clustersRouter.HandleFunc("", clusters.CreateCluster).Methods("POST")
	clustersRouter.HandleFunc("/{id}", clusters.GetCluster).Methods("GET")
	clustersRouter.HandleFunc("/{id}", clusters.UpdateCluster).Methods("PATCH")
	clustersRouter.HandleFunc("/{id}", clusters.DeleteCluster).Methods("DELETE")
	clustersRouter.HandleFunc("/{id}/node-pools", clusters.ListClusterNodeGroups).Methods("GET")
	clustersRouter.HandleFunc("/{id}/node-pools", clusters.AttachClusterNodeGroups).Methods("POST")
	clustersRouter.HandleFunc("/{id}/node-pools/{nodeGroupId}", clusters.DetachClusterNodeGroup).Methods("DELETE")

	credentialsRouter := v1Router.PathPrefix("/credentials").Subrouter()
	credentialsRouter.Use(auth)
	credentialsRouter.HandleFunc("", credentials.ListCredentials).Methods("GET")
	credentialsRouter.HandleFunc("", credentials.CreateCredential).Methods("POST")
	credentialsRouter.HandleFunc("/{id}", credentials.GetCredential).Methods("GET")
	credentialsRouter.HandleFunc("/{id}", credentials.UpdateCredential).Methods("PATCH")
	credentialsRouter.HandleFunc("/{id}", credentials.DeleteCredential).Methods("DELETE")

	sshRouter := v1Router.PathPrefix("/ssh").Subrouter()
	sshRouter.Use(auth)
	sshRouter.HandleFunc("", ssh.ListPublicKeys).Methods("GET")
	sshRouter.HandleFunc("", ssh.CreateSSHKey).Methods("POST")
	sshRouter.HandleFunc("/{id}", ssh.GetSSHKey).Methods("GET")
	sshRouter.HandleFunc("/{id}", ssh.DeleteSSHKey).Methods("DELETE")
	sshRouter.HandleFunc("/{id}/download", ssh.DownloadSSHKey).Methods("GET")

	serversRouter := v1Router.PathPrefix("/servers").Subrouter()
	serversRouter.Use(auth)
	serversRouter.HandleFunc("", servers.ListServers).Methods("GET")
	serversRouter.HandleFunc("", servers.CreateServer).Methods("POST")
	serversRouter.HandleFunc("/{id}", servers.GetServer).Methods("GET")
	serversRouter.HandleFunc("/{id}", servers.UpdateServer).Methods("PATCH")
	serversRouter.HandleFunc("/{id}", servers.DeleteServer).Methods("DELETE")

	nodeTaintsRouter := v1Router.PathPrefix("/node-taints").Subrouter()
	nodeTaintsRouter.Use(auth)
	nodeTaintsRouter.HandleFunc("", nodetaints.ListNodeTaints).Methods("GET")
	nodeTaintsRouter.HandleFunc("", nodetaints.CreateNodeTaint).Methods("POST")
	nodeTaintsRouter.HandleFunc("/{id}", nodetaints.GetNodeTaint).Methods("GET")
	nodeTaintsRouter.HandleFunc("/{id}", nodetaints.UpdateNodeTaint).Methods("PATCH")
	nodeTaintsRouter.HandleFunc("/{id}", nodetaints.DeleteNodeTaint).Methods("DELETE")

	nodePoolsRouter := v1Router.PathPrefix("/node-pools").Subrouter()
	nodePoolsRouter.Use(auth)
	nodePoolsRouter.HandleFunc("", nodepools.ListNodePools).Methods("GET")
	nodePoolsRouter.HandleFunc("", nodepools.CreateNodePool).Methods("POST")
	nodePoolsRouter.HandleFunc("/{id}", nodepools.GetNodePool).Methods("GET")
	nodePoolsRouter.HandleFunc("/{id}", nodepools.UpdateNodeGroup).Methods("PATCH")
	nodePoolsRouter.HandleFunc("/{id}", nodepools.DeleteNodeGroup).Methods("DELETE")
	nodePoolsRouter.HandleFunc("/{id}/servers", nodepools.ListNodeGroupServers).Methods("GET")
	nodePoolsRouter.HandleFunc("/{id}/servers", nodepools.AttachNodeGroupServers).Methods("POST")
	nodePoolsRouter.HandleFunc("/{id}/servers/{serverId}", nodepools.DetachNodeGroupServer).Methods("DELETE")
}
