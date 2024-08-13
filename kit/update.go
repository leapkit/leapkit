package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

func update() error {
	fmt.Println("🔍 Looking for the latest version...")

	response, err := http.Get("https://api.github.com/repos/leapkit/leapkit/tags")
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %w", err)
	}

	defer response.Body.Close()

	type tag struct {
		Name string `json:"name"`
	}

	var tags []tag
	err = json.NewDecoder(response.Body).Decode(&tags)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if len(tags) == 0 {
		return fmt.Errorf("😞 no tags found")
	}

	//get only the ones that contain the key kit
	var kitTags []tag
	for _, tag := range tags {
		if strings.Contains(tag.Name, "kit") {
			kitTags = append(kitTags, tag)
		}
	}

	if len(kitTags) == 0 {
		return fmt.Errorf("😞 no tags found")
	}

	latestTag := kitTags[0].Name
	fmt.Println("😎 Found the latest version:", latestTag)
	fmt.Println("👀 Updating leapkit/kit to the latest version...")
	version := strings.Split(latestTag, "/")[1]

	cmd := exec.Command("go", "install", "github.com/leapkit/leapkit/kit@"+version)
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Failed to update leapkit/kit to the latest version.")
		return err
	}

	fmt.Println("✅ Updated leapkit/kit to the latest version.")
	return nil
}
