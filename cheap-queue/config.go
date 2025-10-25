package cheapqueue

import (
	"os"
	"path/filepath"
)

// getTempDir retorna o diretório temporário correto para cada OS
func getTempDir() string {
	// Windows: C:\Users\Username\AppData\Local\Temp
	// Linux/macOS: /tmp
	return os.TempDir()
}

// getQueueCacheDir retorna o diretório onde as mensagens serão persistidas
func (c *CheapQueueEngine) getQueueCacheDir() string {
	// Você pode customizar isso para usar um diretório específico
	// Por exemplo: ~/.cheap-queue/ ou ./cache/
	return getTempDir()
}

// cleanupOldFiles limpa arquivos antigos do projeto (opcional)
func (c *CheapQueueEngine) cleanupOldFiles() error {
	dir := c.getQueueCacheDir()
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	prefix := c.projectId + "_"
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if filepath.HasPrefix(filename, prefix) {
			// Você pode adicionar lógica aqui para limpar arquivos muito antigos
			// baseado em timestamp, por exemplo
		}
	}
	return nil
}
