package main

import (
	"fmt"
	"log"

	"github.com/iancoleman/orderedmap"
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

	data := orderedmap.New()
	// err = utils.LoadJSON(config.Root, &raw)
	if err := utils.LoadJSONOrdered(config.Root, data); err != nil {
		panic(err)
	}

	// fmt.Println(data)

	// for _, key := range obj.Keys() {
	// 	val, _ := obj.Get(key)
	// 	println(key, val.(string))
	// }

	// fmt.Println(data)

	// data := utils.ToOrderedMap(raw)
	normalizer := normalizer.New(
		config.Root,
	)
	result, err := normalizer.Normalize(data, &config)
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
