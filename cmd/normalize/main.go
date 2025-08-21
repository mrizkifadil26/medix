package main

import (
	"fmt"
	"log"

	"github.com/mrizkifadil26/medix/normalizer"
	_ "github.com/mrizkifadil26/medix/normalizer/actions/extractor"
	_ "github.com/mrizkifadil26/medix/normalizer/actions/formatter"
	_ "github.com/mrizkifadil26/medix/normalizer/actions/replacer"
	_ "github.com/mrizkifadil26/medix/normalizer/actions/transformer"
	"github.com/mrizkifadil26/medix/utils"
)

func main() {
	args, err := normalizer.ParseCLI()
	if err != nil {
		log.Fatalf("Error parsing CLI: %v", err)
	}

	var config normalizer.Config
	if args.ConfigPath != nil {
		config, err = utils.LoadConfig[normalizer.Config](*args.ConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config file: %v", err)
		}
	}

	if err := utils.MergeInto(&config, &args.Config, utils.MergeOptions{
		Overwrite: true,
		Recursive: true,
	}); err != nil {
		log.Fatalf("Failed to merge CLI config: %v", err)
	}

	data := utils.NewOrderedMap[string, any]()
	if err := utils.LoadJSON(config.Root, data); err != nil {
		panic(err)
	}

	n := normalizer.New(&config)
	result, err := n.Normalize(data)
	if err != nil {
		fmt.Println(err)
	}

	// ContinueOnError := false
	// result, err := normalizer.Process(
	// data,
	// config,
	// )
	// result := []string{}

	// Always try to write output, even if errors occurred
	if config.OutputPath != "" {
		if err := utils.WriteJSON(config.OutputPath, result); err != nil {
			log.Fatalf("Failed to write output: %v", err)
		}
	}

	// if err != nil {
	// 	if ContinueOnError {
	// 		fmt.Println("✅ Process completed with errors. Check output for details.")
	// 	} else {
	// 		log.Fatalf("❌ Process failed: %v", err)
	// 	}
	// } else {
	// 	fmt.Println("✅ Process completed")
	// }
}
