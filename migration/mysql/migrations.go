package mysql

type Migrations []*Migration

func (m Migrations) Len() int      { return len(m) }
func (m Migrations) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m Migrations) Less(i, j int) bool {
	return m[i].Name < m[j].Name
}
