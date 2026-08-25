package circular

import "fmt"

type Store struct {
	Path  string
	Table *Table
}

func NewStore(path string) *Store {
	return &Store{Path: path}
}

func (s *Store) Load() (*Table, error) {
	if table, ok := s.cachedTable(); ok {
		return table, nil
	}
	return s.loadFromPath()
}

func (s *Store) cachedTable() (*Table, bool) {
	if s == nil || s.Table == nil {
		return nil, false
	}
	return s.Table, true
}

func (s *Store) loadFromPath() (*Table, error) {
	table, err := LoadCSV(s.Path)
	if err != nil {
		return nil, fmt.Errorf("load influence table %s: %w", s.Path, err)
	}
	bindStoreTable(s, &table)
	return s.Table, nil
}

func (s *Store) Influence(zRatio, rRatio float64) (float64, error) {
	table, err := s.Load()
	if err != nil {
		return 0, err
	}
	return lookupInfluence(table, zRatio, rRatio)
}

func (s *Store) Evaluate(l Load) (Result, error) {
	if err := l.Validate(); err != nil {
		return Result{}, err
	}
	influence, err := s.Influence(l.Z/l.A, l.R/l.A)
	if err != nil {
		return Result{}, err
	}
	return PressureResult(l, influence), nil
}

func (s *Store) Reload() (*Table, error) {
	if table, ok := s.cachedTable(); ok {
		return table, nil
	}
	return s.loadFromPath()
}

func (s *Store) Snapshot() (map[string]interface{}, error) {
	table, err := s.Load()
	if err != nil {
		return nil, err
	}
	return table.Snapshot(), nil
}

func (s *Store) TableLoaded() bool {
	return s.Table != nil
}

func (s *Store) Reset() {
	s.Table = nil
}

func (s *Store) PathString() string {
	return s.Path
}
