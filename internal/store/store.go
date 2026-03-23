package store

import (
	"fmt"
	"os"
	"sync"
)

// TemplateStore interface for loading templates.
type TemplateStore interface {
	ReadTemplate(name string) ([]byte, error)
	ListTemplates() ([]string, error)
}

// MemoryStore is an in-memory template store.
type MemoryStore struct {
	templates map[string]string
	mu        sync.RWMutex
}

// NewMemoryStore creates a new in-memory template store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		templates: make(map[string]string),
	}
}

// Register adds a template to store.
func (ms *MemoryStore) Register(name string, content string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.templates[name] = content
}

// ReadTemplate reads a template from store.
func (ms *MemoryStore) ReadTemplate(name string) ([]byte, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	content, ok := ms.templates[name]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", name)
	}
	return []byte(content), nil
}

// ListTemplates lists all templates in the store.
func (ms *MemoryStore) ListTemplates() ([]string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	names := make([]string, 0, len(ms.templates))
	for name := range ms.templates {
		names = append(names, name)
	}
	return names, nil
}

// FileSystemStore loads templates from the filesystem.
type FileSystemStore struct {
	baseDir string
}

// NewFileSystemStore creates a new filesystem template store.
func NewFileSystemStore(baseDir string) *FileSystemStore {
	return &FileSystemStore{baseDir: baseDir}
}

// ReadTemplate reads a template from the filesystem.
func (fs *FileSystemStore) ReadTemplate(name string) ([]byte, error) {
	path := fs.baseDir + "/" + name
	if path[len(path)-5:] != ".wisp" {
		path += ".wisp"
	}
	return os.ReadFile(path)
}

// ListTemplates lists all templates in the base directory.
func (fs *FileSystemStore) ListTemplates() ([]string, error) {
	entries, err := os.ReadDir(fs.baseDir)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
