package get

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	pkgtypes "github.com/chronicblondiee/searchctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGetNodesCmd() *cobra.Command {
	var (
		roleFilter string
		selector   string
		nameFilter string
		sortBy     string
		desc       bool
		limit      int
		columnsCSV string
	)
	cmd := &cobra.Command{
		Use:     "nodes [NODE_NAME]",
		Short:   "List nodes",
		Long:    "List all nodes with optional filtering, sorting, and custom columns.",
		Aliases: []string{"node", "no"},
		Example: strings.TrimSpace(`
  # Basic list
  searchctl get nodes

  # Filter by role and name substring
  searchctl get nodes --role data --name es-data-

  # Sort by CPU then heap percent (descending) and show top 10
  searchctl get nodes --sort CPU,HEAP.PERCENT --desc --limit 10

  # Choose exact columns
  searchctl get nodes --columns NAME,IP,CPU,HEAP.PERCENT

  # Wide output adds additional load columns automatically
  searchctl get nodes -o wide
        `),
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			nodes, err := c.GetNodes()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting nodes: %v\n", err)
				os.Exit(1)
			}
			// Filters
			filtered := filterNodes(nodes, roleFilter, selector, nameFilter)

			// Sorting
			if sortBy != "" {
				sortNodes(filtered, sortBy, desc)
			}

			// Limiting
			if limit > 0 && limit < len(filtered) {
				filtered = filtered[:limit]
			}

			// Columns
			cols := defaultColumns()
			if viper.GetString("output") == "wide" {
				cols = append(cols, "LOAD_5M", "LOAD_15M")
			}
			if columnsCSV != "" {
				cols = parseColumns(columnsCSV)
			}

			// Convert to interface{} slice for formatting with deterministic column order
			data := make([]interface{}, len(filtered))
			pref := strings.Join(cols, ",")
			for i, node := range filtered {
				row := map[string]interface{}{"__columns": pref}
				for _, col := range cols {
					row[col] = valueForColumn(col, node)
				}
				data[i] = row
			}

			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&roleFilter, "role", "", "Filter by node role (substring match, e.g. data, master, ingest)")
	cmd.Flags().StringVar(&selector, "selector", "", "Reserved for attribute filtering (key=value[,key=value])")
	cmd.Flags().StringVar(&nameFilter, "name", "", "Filter by node name or IP substring")
	cmd.Flags().StringVar(&sortBy, "sort", "", "Comma-separated sort columns (case-insensitive). Common: NAME, IP, CPU, HEAP.PERCENT, RAM.PERCENT, LOAD_1M, LOAD_5M, LOAD_15M")
	cmd.Flags().BoolVar(&desc, "desc", false, "Sort in descending order")
	cmd.Flags().IntVar(&limit, "limit", 0, "Limit number of rows after filtering and sorting")
	cmd.Flags().StringVar(&columnsCSV, "columns", "", "Override table columns (CSV). Default: NAME,HOST,IP,HEAP.PERCENT,RAM.PERCENT,CPU,LOAD_1M,ROLE,MASTER. With -o wide: adds LOAD_5M,LOAD_15M")

	return cmd
}

func defaultColumns() []string {
	return []string{"NAME", "HOST", "IP", "HEAP.PERCENT", "RAM.PERCENT", "CPU", "LOAD_1M", "ROLE", "MASTER"}
}

func parseColumns(csv string) []string {
	parts := strings.Split(csv, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, strings.ToUpper(p))
		}
	}
	return out
}

func filterNodes(nodes []pkgtypes.Node, roleFilter, selector, nameFilter string) []pkgtypes.Node {
	filtered := make([]pkgtypes.Node, 0, len(nodes))
	roleFilter = strings.ToLower(strings.TrimSpace(roleFilter))
	nameFilter = strings.ToLower(strings.TrimSpace(nameFilter))
	// selector parsing is a no-op for now; attributes are not available from cat/nodes
	for _, n := range nodes {
		if roleFilter != "" && !matchesRole(n, roleFilter) {
			continue
		}
		if nameFilter != "" && !strings.Contains(strings.ToLower(n.Name+" "+n.IP), nameFilter) {
			continue
		}
		// Attributes not available from cat/nodes; keep selector as no-op for now
		filtered = append(filtered, n)
	}
	return filtered
}

// sortNodes supports multi-column sort: "A,B,C". Uses numeric sort when value is numeric.
func sortNodes(nodes []pkgtypes.Node, sortBy string, desc bool) {
	columns := parseColumns(sortBy)
	less := func(i, j int) bool {
		for _, col := range columns {
			vi := fmt.Sprintf("%v", valueForColumn(col, nodes[i]))
			vj := fmt.Sprintf("%v", valueForColumn(col, nodes[j]))
			if vi == vj {
				continue
			}
			if ni, err := strconv.ParseFloat(vi, 64); err == nil {
				if nj, err2 := strconv.ParseFloat(vj, 64); err2 == nil {
					if desc {
						return ni > nj
					}
					return ni < nj
				}
			}
			if desc {
				return vi > vj
			}
			return vi < vj
		}
		return false
	}
	sort.SliceStable(nodes, less)
}

// Minimal selector support scaffold
func parseSelector(sel string) map[string]string {
	m := make(map[string]string)
	for _, pair := range strings.Split(sel, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return m
}

// valueForColumn maps column name to value from node
func valueForColumn(col string, n pkgtypes.Node) interface{} {
	switch strings.ToUpper(col) {
	case "NAME":
		return n.Name
	case "HOST":
		return n.Host
	case "IP":
		return n.IP
	case "HEAP.PERCENT":
		return n.HeapPercent
	case "RAM.PERCENT":
		return n.RAMPercent
	case "CPU":
		return n.CPU
	case "LOAD_1M":
		return n.Load1m
	case "LOAD_5M":
		return n.Load5m
	case "LOAD_15M":
		return n.Load15m
	case "ROLE":
		return n.NodeRole
	case "MASTER":
		return n.Master
	default:
		return ""
	}
}

// matchesRole supports friendly role aliases and one-letter role codes from _cat/nodes
func matchesRole(n pkgtypes.Node, roleFilter string) bool {
	rf := strings.ToLower(strings.TrimSpace(roleFilter))
	if rf == "" {
		return true
	}
	roles := strings.ToLower(n.NodeRole)
	switch rf {
	case "master", "m":
		if n.Master == "*" {
			return true
		}
		return strings.Contains(roles, "m")
	case "data", "d":
		return strings.Contains(roles, "d")
	case "ingest", "i":
		return strings.Contains(roles, "i")
	case "ml", "machine_learning":
		return strings.Contains(roles, "ml")
	case "transform", "t":
		return strings.Contains(roles, "t")
	case "remote", "remote_cluster_client", "r":
		return strings.Contains(roles, "r")
	case "voting", "voting_only", "v":
		return strings.Contains(roles, "v")
	case "coordinating", "coord", "-":
		return roles == "-" || roles == ""
	default:
		return strings.Contains(roles, rf)
	}
}

// Alias to reduce long type name in helpers
//
