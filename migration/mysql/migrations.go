package mysql

const (
	MigrationExtension = ".sql"
)

type Migrations []*Migration

// Len implements the sort.Interface
func (m Migrations) Len() int { return len(m) }

// Swap implements the sort.Interface
func (m Migrations) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

// Less implements the sort.Interface
func (m Migrations) Less(i, j int) bool {
	return m[i].Name < m[j].Name
}

// func (m *Migrations) LoadFromDirectory(directoryPath, ext string) error {
// 	files, err := ioutil.ReadDir(directoryPath)
// 	if err != nil {
// 		return err
// 	}
// 	for i
// }
