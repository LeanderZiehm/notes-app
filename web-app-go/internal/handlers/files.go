package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var filesDir string

func SetFilesDir(dir string) {
	filesDir = dir
}

type FileResponse struct {
	Name        string `json:"name"`
	FullName    string `json:"fullName"`
	DisplayName string `json:"displayName"`
}

// getSafeFilePath validates filename and prevents path traversal
func getSafeFilePath(filename string) (string, string, error) {
	safeName := filepath.Base(filename)
	if !strings.HasSuffix(safeName, ".txt") {
		safeName += ".txt"
	}
	fullPath := filepath.Join(filesDir, safeName)
	resolvedPath, _ := filepath.Abs(fullPath)
	resolvedDir, _ := filepath.Abs(filesDir)

	if !strings.HasPrefix(resolvedPath, resolvedDir+string(os.PathSeparator)) && resolvedPath != resolvedDir {
		return "", "", os.ErrInvalid
	}
	return resolvedPath, safeName, nil
}

// ListFiles godoc
// @Summary List all files
// @Description Returns a list of all text files
// @Tags files
// @Produce json
// @Success 200 {array} FileResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/files [get]
func ListFiles(c *gin.Context) {
	files, err := ioutil.ReadDir(filesDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error reading files"})
		return
	}

	var response []FileResponse
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".txt") {
			name := strings.TrimSuffix(f.Name(), ".txt")
			response = append(response, FileResponse{
				Name:        name,
				FullName:    f.Name(),
				DisplayName: name,
			})
		}
	}
	c.JSON(http.StatusOK, response)
}

// GetFile godoc
// @Summary Get file content
// @Description Returns content of a text file
// @Tags files
// @Produce json
// @Param filename path string true "Filename"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files/{filename} [get]
func GetFile(c *gin.Context) {
	filename := c.Param("filename")
	fullPath, _, err := getSafeFilePath(filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid file path"})
		return
	}

	content, err := ioutil.ReadFile(fullPath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, map[string]interface{}{"error": "File not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error reading file"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"content": string(content)})
}

// SaveFile godoc
// @Summary Save file content
// @Description Overwrites content of a text file
// @Tags files
// @Accept json
// @Produce json
// @Param filename path string true "Filename"
// @Param content body map[string]string true "Content"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files/{filename} [post]
func SaveFile(c *gin.Context) {
	filename := c.Param("filename")
	var body map[string]string
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid JSON"})
		return
	}

	fullPath, safeFilename, err := getSafeFilePath(filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid file path"})
		return
	}

	content := body["content"]
	if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error saving file"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"success": true, "message": "File saved successfully", "filename": strings.TrimSuffix(safeFilename, ".txt")})
}

// DeleteFile godoc
// @Summary Delete file
// @Description Deletes a text file
// @Tags files
// @Param filename path string true "Filename"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files/{filename} [delete]
func DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	fullPath, _, err := getSafeFilePath(filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid file path"})
		return
	}

	if err := os.Remove(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, map[string]interface{}{"error": "File not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error deleting file"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"success": true, "message": "File deleted successfully"})
}

// CreateFile godoc
// @Summary Create a new file
// @Description Creates a new text file
// @Tags files
// @Accept json
// @Produce json
// @Param file body map[string]string true "Filename and optional content"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/files [post]
func CreateFile(c *gin.Context) {
	var body map[string]string
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid JSON"})
		return
	}

	filename, ok := body["filename"]
	if !ok || filename == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Filename is required"})
		return
	}

	content := body["content"]

	fullPath, safeFilename, err := getSafeFilePath(filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid file path"})
		return
	}

	// Check existence
	if _, err := os.Stat(fullPath); err == nil {
		c.JSON(http.StatusConflict, map[string]interface{}{"error": "File already exists"})
		return
	}

	if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error creating file"})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"success": true, "message": "File created successfully", "filename": strings.TrimSuffix(safeFilename, ".txt")})
}

