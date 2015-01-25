package main

// Strings is a multiple-strings type for `flag`.
type Strings []string

func (s *Strings) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *Strings) Get() interface{} {
	return *s
}

func (s *Strings) String() string {
	return "..."
}
