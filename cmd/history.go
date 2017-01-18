package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/spf13/cobra"
	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"
	"github.com/wallix/awless/revision/repo"
)

func init() {
	RootCmd.AddCommand(historyCmd)
}

var historyCmd = &cobra.Command{
	Use:     "history",
	Aliases: []string{"who"},
	Short:   "Show your infrastucture history",

	RunE: func(cmd *cobra.Command, args []string) error {
		if !repo.IsGitInstalled() {
			fmt.Printf("No history available. You need to install git")
			os.Exit(0)
		}

		rep, err := repo.NewRepo()
		exitOn(err)

		rev, err := rep.LoadRev("c0be154d7a438107fc0eea79752cf8d59da100d8")
		exitOn(err)

		preRev, err := rep.LoadRev("087fa2f2581e0742824bdba5f35d8f866fe80b70")
		exitOn(err)

		root, err := node.NewNodeFromStrings("/region", config.GetDefaultRegion())
		if err != nil {
			return err
		}

		infraDiff, err := rdf.NewHierarchicalDiffer().Run(root, preRev.Infra, rev.Infra)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("------ INFRA ------")
		infraDiff.FullGraph().VisitDepthFirst(root, printWithDiff)

		return nil
	},
}

func printWithDiff(g *rdf.Graph, n *node.Node, distance int) {
	var lit *literal.Literal
	diff, err := g.TriplesForSubjectPredicate(n, rdf.DiffPredicate)
	if len(diff) > 0 && err == nil {
		lit, _ = diff[0].Object().Literal()
	}

	var tabs bytes.Buffer
	for i := 0; i < distance; i++ {
		tabs.WriteByte('\t')
	}

	switch lit {
	case rdf.ExtraLiteral:
		color.Set(color.FgGreen)
		fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
		color.Unset()
	case rdf.MissingLiteral:
		color.Set(color.FgRed)
		fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
		color.Unset()
	default:
		fmt.Fprintf(os.Stdout, "%s%s, %s\n", tabs.String(), n.Type(), n.ID())
	}
}
