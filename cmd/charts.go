package cmd

import (
	autopatch "github.com/dingobar/autopatch/autopatch"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CheckCharts(cmd *cobra.Command, args []string) {
	logrus.Infoln("Checking for chart updates...")
	var chartConfig []autopatch.ChartConfig
	err := viper.UnmarshalKey("charts", &chartConfig)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	errs := autopatch.LoopChartsAndCheck(chartConfig)

	for _, err := range errs {
		logrus.Error(err)
	}

	check, err := cmd.Flags().GetBool("check")
	if err != nil {
		logrus.Fatalf("Unexpected error %s", err)
	}
	if len(errs) > 0 && check {
		logrus.Errorf("FAILED - %d charts have new versions", len(errs))
		logrus.Exit(1)
	} else {
		logrus.Info("OK")
	}
}

// ChartsCmd represents the charts command
var ChartsCmd = &cobra.Command{
	Use:   "charts",
	Short: "Checks the current version of the chart released and compares to the latest available",
	Long: `This command automatically checks for the latest version of the application in the chart
	repo, and compares it with the desired release in the given namespace.`,
	Run: CheckCharts,
}

func init() {
	rootCmd.AddCommand(ChartsCmd)
	ChartsCmd.Flags().BoolP("check", "c", false, "os.exit(1) if there is any pending update")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chartsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chartsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
