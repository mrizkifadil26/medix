package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/organizer"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	flags := organizer.Parse()

	var input model.MediaOutput
	if err := utils.LoadJSON(flags.InputPath, &input); err != nil {
		fmt.Fprintln(os.Stderr, "❌ Failed to load input JSON:", err)
		os.Exit(1)
	}

	if err := utils.EnsureDir(flags.TargetDir); err != nil {
		fmt.Fprintln(os.Stderr, "❌ Failed to create target directory:", err)
		os.Exit(1)
	}

	switch flags.Mode {
	case "preview":
		result := organizer.Preview(input, flags.IconSources, flags.TargetDir)
		if err := utils.WriteJSON(flags.ReportPath, result); err != nil {
			fmt.Fprintln(os.Stderr, "❌ Failed to save preview result:", err)
			os.Exit(1)
		}

		if flags.ExportPath != "" {
			if err := utils.WriteJSON(flags.ExportPath, result); err != nil {
				fmt.Fprintln(os.Stderr, "❌ Failed to save preview result:", err)
				os.Exit(1)
			}
		}

		fmt.Println("✅ Preview completed:", flags.TargetDir)

	case "apply":
		var result organizer.OrganizeResult

		if utils.FileExists(flags.ReportPath) {
			if err := utils.LoadJSON(flags.ReportPath, &result); err != nil {
				fmt.Fprintln(os.Stderr, "❌ Failed to load preview result for apply:", err)
				os.Exit(1)
			}
		} else {
			// Fallback to preview
			fmt.Println("⚠️ JSON file not found, running preview...")
			utils.LoadJSON(flags.InputPath, &input)

			result = organizer.Preview(input, flags.IconSources, flags.TargetDir)
			if err := utils.WriteJSON(flags.ReportPath, result); err != nil {
				fmt.Fprintln(os.Stderr, "❌ Failed to save preview result:", err)
				os.Exit(1)
			}

			if flags.ExportPath != "" {
				if err := utils.WriteJSON(flags.ExportPath, result); err != nil {
					fmt.Fprintln(os.Stderr, "❌ Failed to save preview result:", err)
					os.Exit(1)
				}
			}
		}

		// Apply
		organizer.Apply(result)
		fmt.Println("✅ Apply completed:", flags.TargetDir)

	default:
		log.Fatalf("❌ Unknown mode: %s", flags.Mode)
		os.Exit(1)
	}
}
