package main

import (
	"database/sql"
	"fmt"
)

type recipe struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Ingredients string `json:"ingredients"`
	Description string `json:"description"`
}

func (r *recipe) getRecipe(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT title, ingredients, description FROM recipes WHERE id=%d", r.ID)
	return db.QueryRow(statement).Scan(&r.Title, &r.Ingredients, &r.Description)
}

func (r *recipe) updateRecipe(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE recipes SET title='%s', ingredients='%s', description='%s' WHERE id=%d", r.Title, r.Ingredients, r.Description, r.ID)
	_, err := db.Exec(statement)
	return err
}

func (r *recipe) deleteRecipe(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM recipes WHERE id=%d", r.ID)
	_, err := db.Exec(statement)
	return err
}

func (r *recipe) createRecipe(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO recipes(title, ingredients, description) VALUES('%s', '%s', '%s')", r.Title, r.Ingredients, r.Description)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&r.ID)
	if err != nil {
		return err
	}
	return nil
}

func getRecipes(db *sql.DB) ([]recipe, error) {
	statement := "SELECT id, title, ingredients, description FROM recipes"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	recipes := []recipe{}
	for rows.Next() {
		var r recipe
		if err := rows.Scan(&r.ID, &r.Title, &r.Ingredients, &r.Description); err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}
	return recipes, nil
}
