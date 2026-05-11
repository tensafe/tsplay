package tsplay_core

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"crypto/rand"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	lua "github.com/yuin/gopher-lua"
)

const (
	zipEncryptedFlag = 0x1
	zipDataDescFlag  = 0x8
	zipUTF8Flag      = 0x800
)

var zipCryptoCRCTable = crc32.MakeTable(crc32.IEEE)

type zipCompressConfig struct {
	FilePath string
	Sources  []string
	Password string
	BaseDir  string
}

type zipExtractConfig struct {
	FilePath  string
	SavePath  string
	Password  string
	Overwrite bool
}

type zipArchiveEntry struct {
	SourcePath string
	Name       string
	Info       fs.FileInfo
	IsDir      bool
}

func zip_compress(L *lua.LState) int {
	config, err := zipCompressConfigFromLua(L)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	result, err := executeZipCompress(config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func zip_extract(L *lua.LState) int {
	config, err := zipExtractConfigFromLua(L)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	result, err := executeZipExtract(config)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(goValueToLua(L, result))
	return 1
}

func zipCompressConfigFromLua(L *lua.LState) (zipCompressConfig, error) {
	if L.GetTop() == 1 {
		if values, ok := luaValueToGo(L.CheckAny(1)).(map[string]any); ok {
			return normalizeZipCompressConfig(values)
		}
	}
	values := map[string]any{
		"file_path": L.CheckString(1),
	}
	if L.GetTop() >= 2 && L.Get(2) != lua.LNil {
		values["source_path"] = luaValueToGo(L.Get(2))
	}
	if L.GetTop() >= 3 && L.Get(3) != lua.LNil {
		third := luaValueToGo(L.Get(3))
		if options, ok := third.(map[string]any); ok {
			for name, value := range options {
				values[name] = value
			}
		} else {
			values["password"] = third
		}
	}
	return normalizeZipCompressConfig(values)
}

func zipExtractConfigFromLua(L *lua.LState) (zipExtractConfig, error) {
	if L.GetTop() == 1 {
		if values, ok := luaValueToGo(L.CheckAny(1)).(map[string]any); ok {
			return normalizeZipExtractConfig(values)
		}
	}
	values := map[string]any{
		"file_path": L.CheckString(1),
		"save_path": L.CheckString(2),
	}
	if L.GetTop() >= 3 && L.Get(3) != lua.LNil {
		third := luaValueToGo(L.Get(3))
		if options, ok := third.(map[string]any); ok {
			for name, value := range options {
				values[name] = value
			}
		} else {
			values["password"] = third
		}
	}
	return normalizeZipExtractConfig(values)
}

func normalizeZipCompressConfig(values map[string]any) (zipCompressConfig, error) {
	config := zipCompressConfig{}
	config.FilePath = firstNonEmptyString(values, "file_path", "archive_path", "output_path", "save_path")
	config.Password = firstNonEmptyString(values, "password")
	config.BaseDir = firstNonEmptyString(values, "base_dir", "root")

	sourceKeys := []string{"source_path", "source", "path", "file", "folder", "files", "folders", "paths", "sources"}
	for _, key := range sourceKeys {
		raw, ok := values[key]
		if !ok || raw == nil {
			continue
		}
		sources, err := zipPathListValue(raw)
		if err != nil {
			return zipCompressConfig{}, fmt.Errorf("zip_compress %s %w", key, err)
		}
		config.Sources = append(config.Sources, sources...)
	}

	if strings.TrimSpace(config.FilePath) == "" {
		return zipCompressConfig{}, fmt.Errorf("zip_compress requires file_path")
	}
	if len(config.Sources) == 0 {
		return zipCompressConfig{}, fmt.Errorf("zip_compress requires source_path, file, files, folder, or folders")
	}
	return config, nil
}

func normalizeZipExtractConfig(values map[string]any) (zipExtractConfig, error) {
	config := zipExtractConfig{
		FilePath:  firstNonEmptyString(values, "file_path", "archive_path", "source_path"),
		SavePath:  firstNonEmptyString(values, "save_path", "output_dir", "dest_dir", "destination"),
		Password:  firstNonEmptyString(values, "password"),
		Overwrite: true,
	}
	if raw, ok := values["overwrite"]; ok && raw != nil {
		overwrite, err := boolParam(raw)
		if err != nil {
			return zipExtractConfig{}, fmt.Errorf("zip_extract overwrite %w", err)
		}
		config.Overwrite = overwrite
	}
	if strings.TrimSpace(config.FilePath) == "" {
		return zipExtractConfig{}, fmt.Errorf("zip_extract requires file_path")
	}
	if strings.TrimSpace(config.SavePath) == "" {
		return zipExtractConfig{}, fmt.Errorf("zip_extract requires save_path")
	}
	return config, nil
}

func firstNonEmptyString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := values[key]
		if !ok || value == nil {
			continue
		}
		if text, ok := value.(string); ok && strings.TrimSpace(text) != "" {
			return text
		}
	}
	return ""
}

func zipPathListValue(value any) ([]string, error) {
	switch typed := value.(type) {
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil, fmt.Errorf("cannot be blank")
		}
		return []string{typed}, nil
	case []string:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			if strings.TrimSpace(item) == "" {
				return nil, fmt.Errorf("must contain non-empty paths")
			}
			items = append(items, item)
		}
		return items, nil
	case []any:
		items := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok || strings.TrimSpace(text) == "" {
				return nil, fmt.Errorf("must be a path string or list of path strings")
			}
			items = append(items, text)
		}
		return items, nil
	default:
		return nil, fmt.Errorf("must be a path string or list of path strings")
	}
}

func executeZipCompress(config zipCompressConfig) (map[string]any, error) {
	entries, err := collectZipArchiveEntries(config)
	if err != nil {
		return nil, err
	}
	if err := writeOutputFile(config.FilePath, nil); err != nil {
		return nil, fmt.Errorf("zip_compress prepare %q: %w", config.FilePath, err)
	}
	file, err := os.Create(config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("zip_compress create %q: %w", config.FilePath, err)
	}

	writer := zip.NewWriter(file)
	fileCount := 0
	dirCount := 0
	for _, entry := range entries {
		if entry.IsDir {
			dirCount++
		} else {
			fileCount++
		}
		if err := addZipArchiveEntry(writer, entry, config.Password); err != nil {
			_ = writer.Close()
			_ = file.Close()
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("zip_compress close archive %q: %w", config.FilePath, err)
	}
	if err := file.Close(); err != nil {
		return nil, fmt.Errorf("zip_compress close file %q: %w", config.FilePath, err)
	}

	info, err := os.Stat(config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("zip_compress stat %q: %w", config.FilePath, err)
	}
	return map[string]any{
		"file_path":   config.FilePath,
		"entries":     len(entries),
		"files":       fileCount,
		"directories": dirCount,
		"bytes":       info.Size(),
		"encrypted":   strings.TrimSpace(config.Password) != "",
	}, nil
}

func collectZipArchiveEntries(config zipCompressConfig) ([]zipArchiveEntry, error) {
	outputAbs, err := filepath.Abs(config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("zip_compress resolve output %q: %w", config.FilePath, err)
	}

	baseAbs := ""
	if strings.TrimSpace(config.BaseDir) != "" {
		baseAbs, err = filepath.Abs(config.BaseDir)
		if err != nil {
			return nil, fmt.Errorf("zip_compress resolve base_dir %q: %w", config.BaseDir, err)
		}
	}

	entries := make([]zipArchiveEntry, 0)
	seen := map[string]string{}
	for _, source := range config.Sources {
		sourceAbs, err := filepath.Abs(source)
		if err != nil {
			return nil, fmt.Errorf("zip_compress resolve source %q: %w", source, err)
		}
		info, err := os.Lstat(sourceAbs)
		if err != nil {
			return nil, fmt.Errorf("zip_compress open source %q: %w", source, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("zip_compress source %q is a symlink; symlinks are not supported", source)
		}
		if info.IsDir() {
			err = filepath.WalkDir(sourceAbs, func(current string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				itemInfo, err := d.Info()
				if err != nil {
					return err
				}
				if itemInfo.Mode()&os.ModeSymlink != 0 {
					return fmt.Errorf("zip_compress source %q is a symlink; symlinks are not supported", current)
				}
				if !itemInfo.IsDir() && !itemInfo.Mode().IsRegular() {
					return fmt.Errorf("zip_compress source %q is not a regular file", current)
				}
				if sameFilePath(current, outputAbs) {
					return nil
				}
				name, err := zipEntryNameForSource(current, sourceAbs, baseAbs, true)
				if err != nil {
					return err
				}
				if name == "" {
					return nil
				}
				return appendZipArchiveEntry(&entries, seen, zipArchiveEntry{
					SourcePath: current,
					Name:       name,
					Info:       itemInfo,
					IsDir:      itemInfo.IsDir(),
				})
			})
			if err != nil {
				return nil, err
			}
			continue
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("zip_compress source %q is not a regular file", source)
		}
		if sameFilePath(sourceAbs, outputAbs) {
			continue
		}
		name, err := zipEntryNameForSource(sourceAbs, sourceAbs, baseAbs, false)
		if err != nil {
			return nil, err
		}
		if err := appendZipArchiveEntry(&entries, seen, zipArchiveEntry{
			SourcePath: sourceAbs,
			Name:       name,
			Info:       info,
			IsDir:      false,
		}); err != nil {
			return nil, err
		}
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("zip_compress found no files or folders to archive")
	}
	return entries, nil
}

func appendZipArchiveEntry(entries *[]zipArchiveEntry, seen map[string]string, entry zipArchiveEntry) error {
	name := entry.Name
	if entry.IsDir && !strings.HasSuffix(name, "/") {
		name += "/"
	}
	if previous, ok := seen[name]; ok && previous != entry.SourcePath {
		return fmt.Errorf("zip_compress duplicate archive entry %q from %q and %q", name, previous, entry.SourcePath)
	}
	seen[name] = entry.SourcePath
	entry.Name = name
	*entries = append(*entries, entry)
	return nil
}

func zipEntryNameForSource(current string, root string, baseAbs string, includeRoot bool) (string, error) {
	var rel string
	var err error
	if baseAbs != "" {
		rel, err = filepath.Rel(baseAbs, current)
		if err != nil {
			return "", fmt.Errorf("zip_compress compute relative path for %q: %w", current, err)
		}
		if rel == "." {
			return "", nil
		}
		if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
			return "", fmt.Errorf("zip_compress source %q is outside base_dir %q", current, baseAbs)
		}
		return sanitizeZipArchiveName(filepath.ToSlash(rel))
	}
	if includeRoot {
		parent := filepath.Dir(root)
		rel, err = filepath.Rel(parent, current)
		if err != nil {
			return "", fmt.Errorf("zip_compress compute relative path for %q: %w", current, err)
		}
		return sanitizeZipArchiveName(filepath.ToSlash(rel))
	}
	return sanitizeZipArchiveName(filepath.Base(current))
}

func sanitizeZipArchiveName(name string) (string, error) {
	name = strings.ReplaceAll(name, "\\", "/")
	name = path.Clean(name)
	name = strings.TrimPrefix(name, "./")
	if name == "." || name == "" {
		return "", nil
	}
	if path.IsAbs(name) || strings.HasPrefix(name, "../") || name == ".." || filepath.VolumeName(name) != "" {
		return "", fmt.Errorf("zip entry name %q is not a relative path", name)
	}
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return "", fmt.Errorf("zip entry name %q contains an unsafe path segment", name)
		}
	}
	return name, nil
}

func addZipArchiveEntry(writer *zip.Writer, entry zipArchiveEntry, password string) error {
	header, err := zip.FileInfoHeader(entry.Info)
	if err != nil {
		return fmt.Errorf("zip_compress header %q: %w", entry.SourcePath, err)
	}
	header.Name = entry.Name
	if utf8.ValidString(header.Name) {
		header.Flags |= zipUTF8Flag
	}
	if entry.IsDir {
		header.Method = zip.Store
		header.Name = strings.TrimSuffix(header.Name, "/") + "/"
		_, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("zip_compress add directory %q: %w", entry.Name, err)
		}
		return nil
	}

	if strings.TrimSpace(password) == "" {
		header.Method = zip.Deflate
		file, err := os.Open(entry.SourcePath)
		if err != nil {
			return fmt.Errorf("zip_compress open %q: %w", entry.SourcePath, err)
		}
		defer file.Close()
		entryWriter, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("zip_compress add %q: %w", entry.Name, err)
		}
		if _, err := io.Copy(entryWriter, file); err != nil {
			return fmt.Errorf("zip_compress write %q: %w", entry.Name, err)
		}
		return nil
	}

	content, err := os.ReadFile(entry.SourcePath)
	if err != nil {
		return fmt.Errorf("zip_compress read %q: %w", entry.SourcePath, err)
	}
	crc := crc32.ChecksumIEEE(content)
	compressed, err := deflateRaw(content)
	if err != nil {
		return fmt.Errorf("zip_compress deflate %q: %w", entry.Name, err)
	}
	encrypted, err := zipCryptoEncrypt(password, compressed, crc)
	if err != nil {
		return fmt.Errorf("zip_compress encrypt %q: %w", entry.Name, err)
	}

	header.Method = zip.Deflate
	header.Flags |= zipEncryptedFlag
	header.CRC32 = crc
	header.CompressedSize64 = uint64(len(encrypted))
	header.UncompressedSize64 = uint64(len(content))
	header.CompressedSize = uint32(len(encrypted))
	header.UncompressedSize = uint32(len(content))
	entryWriter, err := writer.CreateRaw(header)
	if err != nil {
		return fmt.Errorf("zip_compress add encrypted %q: %w", entry.Name, err)
	}
	if _, err := entryWriter.Write(encrypted); err != nil {
		return fmt.Errorf("zip_compress write encrypted %q: %w", entry.Name, err)
	}
	return nil
}

func deflateRaw(content []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer, err := flate.NewWriter(&buffer, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(content); err != nil {
		_ = writer.Close()
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func executeZipExtract(config zipExtractConfig) (map[string]any, error) {
	reader, err := zip.OpenReader(config.FilePath)
	if err != nil {
		return nil, fmt.Errorf("zip_extract open %q: %w", config.FilePath, err)
	}
	defer reader.Close()

	rootAbs, err := filepath.Abs(config.SavePath)
	if err != nil {
		return nil, fmt.Errorf("zip_extract resolve save_path %q: %w", config.SavePath, err)
	}
	if err := os.MkdirAll(rootAbs, 0755); err != nil {
		return nil, fmt.Errorf("zip_extract create save_path %q: %w", config.SavePath, err)
	}
	rootReal, err := filepath.EvalSymlinks(rootAbs)
	if err != nil {
		return nil, fmt.Errorf("zip_extract resolve save_path %q: %w", config.SavePath, err)
	}

	fileCount := 0
	dirCount := 0
	encryptedCount := 0
	var bytesWritten int64
	for _, entry := range reader.File {
		if entry.Flags&zipEncryptedFlag != 0 {
			encryptedCount++
		}
		name, err := sanitizeZipArchiveName(entry.Name)
		if err != nil {
			return nil, fmt.Errorf("zip_extract entry %q: %w", entry.Name, err)
		}
		if name == "" {
			continue
		}
		targetPath, err := safeZipExtractTarget(rootReal, name)
		if err != nil {
			return nil, err
		}
		if entry.FileInfo().IsDir() || strings.HasSuffix(entry.Name, "/") {
			dirCount++
			if err := ensureZipExtractDirectory(rootReal, targetPath); err != nil {
				return nil, err
			}
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return nil, fmt.Errorf("zip_extract create directory %q: %w", targetPath, err)
			}
			continue
		}
		written, err := extractZipFileEntry(entry, rootReal, targetPath, config)
		if err != nil {
			return nil, err
		}
		fileCount++
		bytesWritten += written
	}
	return map[string]any{
		"file_path":         config.FilePath,
		"save_path":         config.SavePath,
		"entries":           len(reader.File),
		"files":             fileCount,
		"directories":       dirCount,
		"bytes":             bytesWritten,
		"encrypted_entries": encryptedCount,
	}, nil
}

func safeZipExtractTarget(rootReal string, name string) (string, error) {
	target := filepath.Join(rootReal, filepath.FromSlash(name))
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf("zip_extract resolve entry %q: %w", name, err)
	}
	if err := ensurePathInsideRoot(targetAbs, rootReal); err != nil {
		return "", fmt.Errorf("zip_extract entry %q is outside save_path %q", name, rootReal)
	}
	return targetAbs, nil
}

func ensureZipExtractDirectory(rootReal string, targetPath string) error {
	if info, err := os.Lstat(targetPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("zip_extract target directory %q is a symlink", targetPath)
	}
	return ensureZipExtractParent(rootReal, filepath.Dir(targetPath))
}

func extractZipFileEntry(entry *zip.File, rootReal string, targetPath string, config zipExtractConfig) (int64, error) {
	if !config.Overwrite {
		if _, err := os.Lstat(targetPath); err == nil {
			return 0, fmt.Errorf("zip_extract target %q already exists", targetPath)
		} else if !errors.Is(err, os.ErrNotExist) {
			return 0, fmt.Errorf("zip_extract check target %q: %w", targetPath, err)
		}
	}
	parent := filepath.Dir(targetPath)
	if err := ensureZipExtractParent(rootReal, parent); err != nil {
		return 0, err
	}
	if err := os.MkdirAll(parent, 0755); err != nil {
		return 0, fmt.Errorf("zip_extract create directory %q: %w", parent, err)
	}
	parentReal, err := filepath.EvalSymlinks(parent)
	if err != nil {
		return 0, fmt.Errorf("zip_extract resolve output directory %q: %w", parent, err)
	}
	if err := ensurePathInsideRoot(parentReal, rootReal); err != nil {
		return 0, fmt.Errorf("zip_extract output directory %q is outside save_path %q", parent, rootReal)
	}
	if info, err := os.Lstat(targetPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return 0, fmt.Errorf("zip_extract target %q is a symlink", targetPath)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("zip_extract check target %q: %w", targetPath, err)
	}
	mode := entry.Mode().Perm()
	if mode == 0 {
		mode = 0644
	}
	input, encrypted, err := openZipEntryReader(entry, config.Password)
	if err != nil {
		return 0, err
	}
	defer input.Close()

	output, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return 0, fmt.Errorf("zip_extract create %q: %w", targetPath, err)
	}

	var written int64
	if encrypted {
		hash := crc32.NewIEEE()
		written, err = io.Copy(io.MultiWriter(output, hash), input)
		if err == nil && uint64(written) != entry.UncompressedSize64 {
			err = fmt.Errorf("zip_extract entry %q size mismatch", entry.Name)
		}
		if err == nil && hash.Sum32() != entry.CRC32 {
			err = fmt.Errorf("zip_extract entry %q checksum mismatch; password may be wrong", entry.Name)
		}
	} else {
		written, err = io.Copy(output, input)
	}
	closeErr := output.Close()
	if err != nil {
		return written, fmt.Errorf("zip_extract write %q: %w", targetPath, err)
	}
	if closeErr != nil {
		return written, fmt.Errorf("zip_extract close %q: %w", targetPath, closeErr)
	}
	_ = os.Chtimes(targetPath, entry.Modified, entry.Modified)
	return written, nil
}

func ensureZipExtractParent(rootReal string, parent string) error {
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		return fmt.Errorf("zip_extract resolve output directory %q: %w", parent, err)
	}
	if err := ensurePathInsideRoot(parentAbs, rootReal); err != nil {
		return fmt.Errorf("zip_extract output directory %q is outside save_path %q", parent, rootReal)
	}
	current := parentAbs
	for {
		info, err := os.Lstat(current)
		if err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("zip_extract output directory %q uses symlink ancestor %q", parent, current)
			}
			if !info.IsDir() {
				return fmt.Errorf("zip_extract output directory ancestor %q is not a directory", current)
			}
			real, err := filepath.EvalSymlinks(current)
			if err != nil {
				return fmt.Errorf("zip_extract resolve output directory ancestor %q: %w", current, err)
			}
			if err := ensurePathInsideRoot(real, rootReal); err != nil {
				return fmt.Errorf("zip_extract output directory %q is outside save_path %q", parent, rootReal)
			}
			return nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("zip_extract check output directory ancestor %q: %w", current, err)
		}
		if current == rootReal || current == filepath.Dir(current) {
			return nil
		}
		current = filepath.Dir(current)
		if err := ensurePathInsideRoot(current, rootReal); err != nil {
			return fmt.Errorf("zip_extract output directory %q is outside save_path %q", parent, rootReal)
		}
	}
}

func openZipEntryReader(entry *zip.File, password string) (io.ReadCloser, bool, error) {
	if entry.Flags&zipEncryptedFlag == 0 {
		reader, err := entry.Open()
		if err != nil {
			return nil, false, fmt.Errorf("zip_extract open entry %q: %w", entry.Name, err)
		}
		return reader, false, nil
	}
	if strings.TrimSpace(password) == "" {
		return nil, true, fmt.Errorf("zip_extract entry %q is password protected", entry.Name)
	}
	raw, err := entry.OpenRaw()
	if err != nil {
		return nil, true, fmt.Errorf("zip_extract open raw entry %q: %w", entry.Name, err)
	}
	decrypted := newZipCryptoReader(raw, password)
	header := make([]byte, 12)
	if _, err := io.ReadFull(decrypted, header); err != nil {
		return nil, true, fmt.Errorf("zip_extract decrypt header %q: %w", entry.Name, err)
	}
	expected := byte(entry.CRC32 >> 24)
	if entry.Flags&zipDataDescFlag != 0 {
		expected = byte(entry.ModifiedTime >> 8)
	}
	if header[11] != expected {
		return nil, true, fmt.Errorf("zip_extract entry %q password check failed", entry.Name)
	}
	switch entry.Method {
	case zip.Store:
		return io.NopCloser(decrypted), true, nil
	case zip.Deflate:
		return flate.NewReader(decrypted), true, nil
	default:
		return nil, true, fmt.Errorf("zip_extract entry %q uses unsupported compression method %d", entry.Name, entry.Method)
	}
}

func sameFilePath(a string, b string) bool {
	aAbs, aErr := filepath.Abs(a)
	bAbs, bErr := filepath.Abs(b)
	if aErr != nil || bErr != nil {
		return false
	}
	return aAbs == bAbs
}

type zipCryptoKeys struct {
	key0 uint32
	key1 uint32
	key2 uint32
}

func newZipCryptoKeys(password string) *zipCryptoKeys {
	keys := &zipCryptoKeys{
		key0: 0x12345678,
		key1: 0x23456789,
		key2: 0x34567890,
	}
	for i := 0; i < len(password); i++ {
		keys.update(password[i])
	}
	return keys
}

func (keys *zipCryptoKeys) update(value byte) {
	keys.key0 = zipCryptoCRC32Update(keys.key0, value)
	keys.key1 += keys.key0 & 0xff
	keys.key1 = keys.key1*134775813 + 1
	keys.key2 = zipCryptoCRC32Update(keys.key2, byte(keys.key1>>24))
}

func (keys *zipCryptoKeys) streamByte() byte {
	temp := uint16(keys.key2) | 2
	return byte((temp * (temp ^ 1)) >> 8)
}

func zipCryptoCRC32Update(crc uint32, value byte) uint32 {
	return zipCryptoCRCTable[(crc^uint32(value))&0xff] ^ (crc >> 8)
}

func zipCryptoEncrypt(password string, plain []byte, crc uint32) ([]byte, error) {
	keys := newZipCryptoKeys(password)
	header := make([]byte, 12)
	if _, err := rand.Read(header[:11]); err != nil {
		return nil, err
	}
	header[11] = byte(crc >> 24)
	out := make([]byte, 0, len(header)+len(plain))
	out = appendZipCryptoEncrypted(out, keys, header)
	out = appendZipCryptoEncrypted(out, keys, plain)
	return out, nil
}

func appendZipCryptoEncrypted(out []byte, keys *zipCryptoKeys, plain []byte) []byte {
	for _, value := range plain {
		cipher := value ^ keys.streamByte()
		keys.update(value)
		out = append(out, cipher)
	}
	return out
}

type zipCryptoReader struct {
	reader io.Reader
	keys   *zipCryptoKeys
}

func newZipCryptoReader(reader io.Reader, password string) *zipCryptoReader {
	return &zipCryptoReader{
		reader: reader,
		keys:   newZipCryptoKeys(password),
	}
}

func (reader *zipCryptoReader) Read(p []byte) (int, error) {
	n, err := reader.reader.Read(p)
	for i := 0; i < n; i++ {
		cipher := p[i]
		plain := cipher ^ reader.keys.streamByte()
		reader.keys.update(plain)
		p[i] = plain
	}
	return n, err
}
