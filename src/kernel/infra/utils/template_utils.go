package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadEmailTemplate(templateName string, data map[string]interface{}) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	templatePath := filepath.Join(currentDir, "src", "kernel", "app", "emails", templateName)

	content, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", templateName, err)
	}

	templateContent := string(content)

	if data != nil {
		templateContent = replaceTemplateMarkers(templateContent, data)
	}

	return templateContent, nil
}

func ReadEmailTemplateWithBase(contentTemplateName string, locale string, data map[string]interface{}) (string, error) {
	baseTemplateName := fmt.Sprintf("base_%s.html", locale)
	baseContent, err := ReadEmailTemplate(baseTemplateName, nil)
	if err != nil {
		return "", fmt.Errorf("failed to read base template: %w", err)
	}

	contentTemplate, err := ReadEmailTemplate(contentTemplateName, nil)
	if err != nil {
		return "", fmt.Errorf("failed to read content template: %w", err)
	}

	combinedTemplate := strings.ReplaceAll(baseContent, "{{content}}", contentTemplate)

	if data != nil {
		combinedTemplate = replaceTemplateMarkers(combinedTemplate, data)
	}

	return combinedTemplate, nil
}

func replaceTemplateMarkers(template string, data map[string]interface{}) string {
	result := template
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

func GetInvitationTemplateName(locale string) string {
	if locale == "fr" {
		return "invitation_content_fr.html"
	}
	return "invitation_content_en.html"
}

func GetTemplatePath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	return filepath.Join(currentDir, "src", "kernel", "app", "emails"), nil
}

func ListAvailableTemplates() ([]string, error) {
	templatePath, err := GetTemplatePath()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".html") {
			templates = append(templates, file.Name())
		}
	}

	return templates, nil
}
