package execution

type LinkGenerator interface {
	GenerateExplorerLink(status InstallStatus) string
	GenerateEntityLink(entityGUID string) string
	GenerateLoggingLink(entityGUID string) string
	GenerateFleetLink(entityGUID string) string
	GenerateGuidedInstallDocLink() string
	GenerateRedirectURL(status InstallStatus) string
}
