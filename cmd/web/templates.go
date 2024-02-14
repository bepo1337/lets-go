package main

import "letsgo.bepo1337/internal/models"

type ViewTemplateData struct {
	Snippet *models.Snippet
}

type HomeTemplateData struct {
	Snippets []*models.Snippet
}
