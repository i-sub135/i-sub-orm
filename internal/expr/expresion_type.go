package expr

// Eq  => equality ("=")
// Neq => not equal ("!=")
// Gt  => greater than
// Lt  => less than
// In  => IN (...)
type Eq map[string]any
type Neq map[string]any
type Gt map[string]any
type Lt map[string]any
type In map[string][]any
