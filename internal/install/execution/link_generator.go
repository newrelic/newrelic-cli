package execution

type LinkGenerator interface {
	GenerateExplorerLink(filter string) string
	GenerateEntityLink(entityGUID string) string
	GenerateRedirectURL(status InstallStatus) string
}
