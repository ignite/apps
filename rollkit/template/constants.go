package template

const (
	ServerAddCommandsWithStartCmdOptions = "AddCommandsWithStartCmdOptions"

	RollkitV0XStartHandler = "rollserv.StartHandler"
	RollkitV1XStartHandler = "abciserver.StartHandler"
	RollkitServerOptions   = `server.AddCommandsWithStartCmdOptions(
			rootCmd,
			app.DefaultNodeHome,
			newApp, appExport,
			server.StartCmdOptions{
				AddFlags: func(cmd *cobra.Command) {
					abciserver.AddFlags(cmd)
				},
				StartCommandHandler: abciserver.StartHandler(),
			})`
)
