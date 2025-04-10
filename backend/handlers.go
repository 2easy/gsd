package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Project struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Position  float64 `json:"position"`
	CreatedAt string  `json:"created_at"`
	Deadline  string  `json:"deadline,omitempty"`
}

// Request body for position update
type UpdatePositionRequest struct {
	Position float64 `json:"position"`
}

type UpdateProjectRequest struct {
	Position float64 `json:"position,omitempty"`
	Deadline *string `json:"deadline"`
}

type NextAction struct {
	ID          string  `json:"id"`
	Action      string  `json:"action"`
	ProjectID   string  `json:"project_id,omitempty"`
	URL         string  `json:"url,omitempty"`
	Size        string  `json:"size,omitempty"`
	Energy      string  `json:"energy,omitempty"`
	CreatedAt   string  `json:"created_at"`
	CompletedAt string  `json:"completed_at,omitempty"`
	Position    float64 `json:"position"`
}

type UpdateNextActionRequest struct {
	Action      string  `json:"action,omitempty"`
	ProjectID   string  `json:"project_id,omitempty"`
	URL         string  `json:"url,omitempty"`
	Size        string  `json:"size,omitempty"`
	Energy      string  `json:"energy,omitempty"`
	CompletedAt string  `json:"completed_at,omitempty"`
	Position    float64 `json:"position,omitempty"`
}

type InboxItem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	URL         string `json:"url,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type UpdateInboxItemRequest struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
	CreatedAt   string `json:"created_at"`
}

func GetProjects(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, position, deadline FROM projects ORDER BY position")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var projects []Project
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Position, &project.Deadline); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		projects = append(projects, project)
	}

	if len(projects) == 0 {
		projects = []Project{}
	}

	c.JSON(http.StatusOK, projects)
}

func CreateProject(c *gin.Context) {
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		log.Println("Error decoding request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if project.ID == "" {
		project.ID = uuid.New().String()
	}

	// Set creation time
	project.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	log.Println("Inserting project into database with ID:", project.ID, "and Name:", project.Name)

	// Get max position
	var maxPosition sql.NullFloat64
	err := db.QueryRow("SELECT MAX(position) FROM projects").Scan(&maxPosition)
	if err != nil {
		log.Println("Error getting max position:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !maxPosition.Valid {
		project.Position = 1.0
	} else {
		project.Position = maxPosition.Float64 + 1.0
	}

	_, err = db.Exec("INSERT INTO projects (id, name, position, deadline, created_at) VALUES (?, ?, ?, ?, ?)",
		project.ID, project.Name, project.Position, project.Deadline, project.CreatedAt)
	if err != nil {
		log.Println("Error inserting into database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

func UpdateProject(c *gin.Context) {
	projectID := c.Param("id")

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE projects SET"
	var params []interface{}
	var setFields []string

	// Check if position was included
	if req.Position != 0 {
		setFields = append(setFields, " position = ?")
		params = append(params, req.Position)
	}

	// Check if deadline was explicitly included (even if null)
	if req.Deadline != nil {
		setFields = append(setFields, " deadline = ?")
		params = append(params, req.Deadline)
	}

	if len(setFields) == 0 {
		log.Println("No fields to update in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Combine all set fields and add WHERE clause
	for i := 0; i < len(setFields)-1; i++ {
		query += setFields[i] + ","
	}
	query += setFields[len(setFields)-1] + " WHERE id = ?"
	params = append(params, projectID)

	log.Printf("Executing query: %s with params: %v", query, params)
	result, err := db.Exec(query, params...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		log.Printf("No project found with ID: %s", projectID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var project Project
	err = db.QueryRow("SELECT id, name, position, deadline FROM projects WHERE id = ?", projectID).
		Scan(&project.ID, &project.Name, &project.Position, &project.Deadline)
	if err != nil {
		log.Printf("Error fetching updated project: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	projectID := c.Param("id")

	result, err := db.Exec("DELETE FROM projects WHERE id = ?", projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.Status(http.StatusOK)
}

func GetNextActions(c *gin.Context) {
	rows, err := db.Query("SELECT id, action, project_id, url, size, energy, created_at, completed_at, position FROM next_actions ORDER BY position")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var actions []NextAction
	for rows.Next() {
		var action NextAction
		var size, energy, projectID, url, completedAt sql.NullString

		if err := rows.Scan(&action.ID, &action.Action, &projectID, &url, &size,
			&energy, &action.CreatedAt, &completedAt, &action.Position); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Only set the fields if they are not NULL
		if size.Valid {
			action.Size = size.String
		}
		if energy.Valid {
			action.Energy = energy.String
		}
		if projectID.Valid {
			action.ProjectID = projectID.String
		}
		if url.Valid {
			action.URL = url.String
		}
		if completedAt.Valid {
			action.CompletedAt = completedAt.String
		}

		actions = append(actions, action)
	}

	if len(actions) == 0 {
		actions = []NextAction{}
	}
	c.JSON(http.StatusOK, actions)
}

func CreateNextAction(c *gin.Context) {
	var action NextAction
	if err := c.ShouldBindJSON(&action); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if action.ID == "" {
		action.ID = uuid.New().String()
	}

	// Get max position
	var maxPosition sql.NullFloat64
	err := db.QueryRow("SELECT MAX(position) FROM next_actions").Scan(&maxPosition)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !maxPosition.Valid {
		action.Position = 1.0
	} else {
		action.Position = maxPosition.Float64 + 1.0
	}

	// Set creation time
	action.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	// Handle NULL values
	var sizeParam, energyParam interface{}
	if action.Size == "" {
		sizeParam = nil
	} else {
		sizeParam = action.Size
	}
	if action.Energy == "" {
		energyParam = nil
	} else {
		energyParam = action.Energy
	}

	_, err = db.Exec(`
		INSERT INTO next_actions (id, action, project_id, url, size, energy, created_at, completed_at, position) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		action.ID, action.Action, action.ProjectID, action.URL, sizeParam,
		energyParam, action.CreatedAt, nil, action.Position)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, action)
}

func UpdateNextAction(c *gin.Context) {
	actionID := c.Param("id")

	// Get the raw JSON to check which fields were actually included in the request
	var rawJson map[string]interface{}
	if err := c.ShouldBindJSON(&rawJson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE next_actions SET"
	var params []interface{}
	var setFields []string

	// Only update fields that were explicitly included in the request
	if _, exists := rawJson["action"]; exists {
		setFields = append(setFields, " action = ?")
		params = append(params, rawJson["action"])
	}
	if _, exists := rawJson["project_id"]; exists {
		setFields = append(setFields, " project_id = ?")
		params = append(params, rawJson["project_id"])
	}
	if _, exists := rawJson["url"]; exists {
		setFields = append(setFields, " url = ?")
		params = append(params, rawJson["url"])
	}
	if _, exists := rawJson["size"]; exists {
		setFields = append(setFields, " size = ?")
		params = append(params, rawJson["size"])
	}
	if _, exists := rawJson["energy"]; exists {
		setFields = append(setFields, " energy = ?")
		params = append(params, rawJson["energy"])
	}
	if _, exists := rawJson["completed_at"]; exists {
		setFields = append(setFields, " completed_at = ?")
		// If completed_at is explicitly null, we want to remove the completion
		if rawJson["completed_at"] == nil {
			params = append(params, nil)
		} else {
			params = append(params, rawJson["completed_at"])
		}
	}
	if _, exists := rawJson["position"]; exists {
		setFields = append(setFields, " position = ?")
		params = append(params, rawJson["position"])
	}

	if len(setFields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Combine all set fields and add WHERE clause
	for i := 0; i < len(setFields)-1; i++ {
		query += setFields[i] + ","
	}
	query += setFields[len(setFields)-1] + " WHERE id = ?"
	params = append(params, actionID)

	result, err := db.Exec(query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Next action not found"})
		return
	}

	// Fetch the updated action data
	var action NextAction
	var size, energy, projectID, url, completedAt sql.NullString
	err = db.QueryRow(`
        SELECT id, action, project_id, url, size, energy, created_at, completed_at, position 
        FROM next_actions WHERE id = ?`, actionID).
		Scan(&action.ID, &action.Action, &projectID, &url, &size,
			&energy, &action.CreatedAt, &completedAt, &action.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated action"})
		return
	}

	// Set the nullable fields if they are valid
	if size.Valid {
		action.Size = size.String
	}
	if energy.Valid {
		action.Energy = energy.String
	}
	if projectID.Valid {
		action.ProjectID = projectID.String
	}
	if url.Valid {
		action.URL = url.String
	}
	if completedAt.Valid {
		action.CompletedAt = completedAt.String
	}

	c.JSON(http.StatusOK, action)
}

func DeleteNextAction(c *gin.Context) {
	actionID := c.Param("id")
	result, err := db.Exec("DELETE FROM next_actions WHERE id = ?", actionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Next action not found"})
		return
	}

	c.Status(http.StatusOK)
}

func GetInboxItems(c *gin.Context) {
	rows, err := db.Query("SELECT id, description, url, created_at FROM inbox WHERE state IS NULL")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []InboxItem
	for rows.Next() {
		var item InboxItem
		if err := rows.Scan(&item.ID, &item.Description, &item.URL, &item.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		items = []InboxItem{}
	}

	c.JSON(http.StatusOK, items)
}

func CreateInboxItem(c *gin.Context) {
	var item InboxItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if item.ID == "" {
		item.ID = uuid.New().String()
	}

	// Set creation time
	item.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	_, err := db.Exec("INSERT INTO inbox (id, description, url, created_at) VALUES (?, ?, ?, ?)", item.ID, item.Description, item.URL, item.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.BroadcastUpdate(map[string]any{
		"type": "inbox_item_created",
		"data": item,
	})

	c.JSON(http.StatusOK, item)
}

func DeleteInboxItem(c *gin.Context) {
	itemID := c.Param("id")

	result, err := db.Exec("UPDATE inbox SET state = 'deleted' WHERE id = ?", itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inbox item not found"})
		return
	}

	c.Status(http.StatusOK)
}
