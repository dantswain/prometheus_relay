package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "prometheus_relay [prometheus_url] [post_to_url]",
	Short: "Relay prometheus data",
	Long:  `Relay prometheus data`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		prometheusClientUrl := args[0]
		postToUrl := args[1]

		client, err := api.NewClient(api.Config{
			Address: prometheusClientUrl,
		})
		if err != nil {
			fmt.Printf("Error creating client: %v\n", err)
			os.Exit(1)
		}

		v1api := v1.NewAPI(client)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		result, warnings, err := v1api.Query(ctx, "ambient_wx_temperature", time.Now())
		if err != nil {
			fmt.Printf("Error querying Prometheus: %v\n", err)
			os.Exit(1)
		}
		if len(warnings) > 0 {
			fmt.Printf("Warnings: %v\n", warnings)
		}

		url := postToUrl
		for _, v := range result.(model.Vector) {
			location := v.Metric[model.LabelName("location")]
			value := v.Value

			vals := map[string]interface{}{"name": location, "value": value}
			str, _ := json.Marshal(vals)
			_, err := http.Post(url, "application/json", bytes.NewBuffer(str))
			if err != nil {
				fmt.Printf("Error posting: %v\n", err)
				os.Exit(1)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// common args
}
