/**
@Author: wei-g
@Date:   2021/1/8 3:52 下午
@Description:
*/

package cmd

func sparrowFlags() {
	RootCmd.Flags().StringP("gin.addr", "a", ":8080", "Run attaches to a http.Server and starts listening and serving HTTP requests")
	RootCmd.Flags().StringP("gin.mode", "m", "release", "SetMode sets gin mode according to input string")
	RootCmd.Flags().Bool("gin.no_route", true, "Set the default NoRoute handler")
	RootCmd.Flags().Bool("gin.no_method", true, "Set the default NoMethod handler")
}
