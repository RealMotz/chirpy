package main

import "testing"

func TestCleanBod(t *testing.T) {
	t.Run("Test case # 1", func(t *testing.T) {
		const body = "This is a kerfuffle opinion I need to share with the world"
		const expected = "This is a **** opinion I need to share with the world"
		actual := cleanBody(body)
		if actual != expected {
			t.Errorf("Expected: %s", expected)
			t.Errorf("Actual: %s", actual)
		}
	})

	t.Run("Test case # 2", func(t *testing.T) {
		const body = "This kerfuffle is a sharbert of a fornax"
		const expected = "This **** is a **** of a ****"
		actual := cleanBody(body)
		if actual != expected {
			t.Errorf("Expected: %s", expected)
			t.Errorf("Actual: %s", actual)
		}
	})

	t.Run("Test case # 3", func(t *testing.T) {
		const body = "This kerfuffle! is a sharbert! of a fornax!"
		actual := cleanBody(body)
		if actual != body {
			t.Errorf("Expected: %s", body)
			t.Errorf("Actual: %s", actual)
		}
	})
}
