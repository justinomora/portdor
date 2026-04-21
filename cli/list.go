package cli

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered services",
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := apiClient.ListServices()
		if err != nil {
			return err
		}
		if len(services) == 0 {
			fmt.Println("No services registered.")
			return nil
		}
		services = sortByProject(services)
		fmt.Printf("%-20s %-15s %-8s %-10s\n",
			"NAME", "PROJECT", "PORT", "STATUS")
		fmt.Println(strings.Repeat("─", 57))
		for _, svc := range services {
			port := "—"
			if p, ok := svc["port"].(float64); ok && p > 0 {
				port = fmt.Sprintf("%d", int(p))
			}
			project := "ungrouped"
			if p, ok := svc["project"].(string); ok && p != "" {
				project = p
			}
			fmt.Printf("%-20s %-15s %-8s %-10s\n",
				svc["name"], project, port, svc["status"])
		}
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status <name>",
	Short: "Show detailed status for a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		services, err := apiClient.ListServices()
		if err != nil {
			return err
		}
		for _, svc := range services {
			if svc["name"] == args[0] {
				data, _ := json.MarshalIndent(svc, "", "  ")
				fmt.Println(string(data))
				return nil
			}
		}
		return fmt.Errorf("service '%s' not found", args[0])
	},
}

func sortByProject(services []map[string]any) []map[string]any {
	sort.Slice(services, func(i, j int) bool {
		pi, _ := services[i]["project"].(string)
		pj, _ := services[j]["project"].(string)
		if pi != pj {
			return pi < pj
		}
		ni, _ := services[i]["name"].(string)
		nj, _ := services[j]["name"].(string)
		return ni < nj
	})
	return services
}

func init() {
	rootCmd.AddCommand(listCmd, statusCmd)
}
