package execution

type LinkGenerator interface {
	GenerateExplorerLink(status InstallStatus) string
	GenerateEntityLink(entityGUID string) string
	GenerateRedirectURL(status InstallStatus) string
}
