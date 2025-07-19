package organizer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrizkifadil26/medix/internal/organize"
	util "github.com/mrizkifadil26/medix/utils"
)

const configPath = "config/organize-icons.json"

func RunPreview() {
	fmt.Println("🔍 Previewing scattered icon organization...")

	// 1. Load config
	var cfg organize.OrganizeConfig
	err := util.LoadJSON(configPath, &cfg)
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 1. Load raw.json and build slug→type map
	rawEntries, err := organize.LoadRawMetadata(cfg.RawMetadataPath)
	if err != nil {
		fmt.Printf("❌ Failed to load raw metadata: %v\n", err)
		os.Exit(1)
	}
	slugMap := organize.BuildSlugMap(rawEntries)

	// 2. Load source dirs from config
	icons, err := organize.LoadScatteredIcons(cfg.Sources, cfg.ExcludeDirs)
	if err != nil {
		fmt.Printf("❌ Failed to load scattered icons: %v\n", err)
		os.Exit(1)
	}

	// 4. Build move plan (source → target)
	movePlan := organize.BuildMovePlan(icons, slugMap, cfg.OutputBase)

	// 5. Show preview
	for _, plan := range movePlan {
		sourceDir := filepath.Base(filepath.Dir(plan.SourcePath)) // e.g., "Icons"
		sourceFile := filepath.Base(plan.SourcePath)              // e.g., "1917.ico"
		source := fmt.Sprintf("../%s/%s", sourceDir, sourceFile)  // e.g., ../Icons/1917.ico

		targetGroup := filepath.Base(plan.Group)                   // e.g., "Action"
		targetFile := filepath.Base(plan.TargetPath)               // e.g., "1917.ico"
		target := fmt.Sprintf("../%s/%s", targetGroup, targetFile) // e.g., ../Action/1917.ico

		switch {
		// case !plan.Matched:
		// 	fmt.Printf("🔴 %s → ❌ no match\n", source)
		case plan.Duplicate:
			fmt.Printf("⚠️  %s → %s (duplicate detected)\n", source, target)
		default:
			fmt.Printf("🟢 %s → %s\n", source, target)
		}
	}

	// 6. Save preview plan
	f, err := os.Create(cfg.PlanPath)
	if err == nil {
		_ = json.NewEncoder(f).Encode(movePlan)
		f.Close()
	}

	fmt.Println("✅ Preview complete.")
}
