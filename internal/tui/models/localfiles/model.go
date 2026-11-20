package localfiles

import (
    "os"
    "path/filepath"
    "sort"
)

type Model struct {
    Files       []FileItem
    SelectedIdx int
    Dir         string
    Width       int
}

type FileItem struct {
    Name    string
    Path    string
    Size    int64
    ModTime time.Time
}

func NewModel(dir string) Model {
    return Model{Dir: dir}
}

func (m *Model) Load() error {
    entries, err := os.ReadDir(m.Dir)
    if err != nil {
        return err
    }
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        ext := strings.ToLower(filepath.Ext(e.Name()))
        if ext == ".mp4" || ext == ".mkv" || ext == ".webm" || ext == ".mp3" || ext == ".m4a" {
            info, _ := e.Info()
            m.Files = append(m.Files, FileItem{
                Name:    e.Name(),
                Path:    filepath.Join(m.Dir, e.Name()),
                Size:    info.Size(),
                ModTime: info.ModTime(),
            })
        }
    }
    sort.Slice(m.Files, func(i, j int) bool {
        return m.Files[i].ModTime.After(m.Files[j].ModTime)
    })
    return nil
}   
