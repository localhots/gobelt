package ptrto

// String returns a pointer to a given string value.
func String(s string) *string { return &s }

// Bool returns a pointer to a given boolean value.
func Bool(b bool) *bool { return &b }

// Int returns a pointer to a given int value.
func Int(i int) *int { return &i }

// Int8 returns a pointer to a given int8 value.
func Int8(i int8) *int8 { return &i }

// Int16 returns a pointer to a given int16 value.
func Int16(i int16) *int16 { return &i }

// Int32 returns a pointer to a given int32 value.
func Int32(i int32) *int32 { return &i }

// Int64 returns a pointer to a given int64 value.
func Int64(i int64) *int64 { return &i }

// Uint returns a pointer to a given uint value.
func Uint(i uint) *uint { return &i }

// Uint8 returns a pointer to a given uint8 value.
func Uint8(i uint8) *uint8 { return &i }

// Uint16 returns a pointer to a given uint16 value.
func Uint16(i uint16) *uint16 { return &i }

// Uint32 returns a pointer to a given uint32 value.
func Uint32(i uint32) *uint32 { return &i }

// Uint64 returns a pointer to a given uint64 value.
func Uint64(i uint64) *uint64 { return &i }

// Float32 returns a pointer to a given float32 value.
func Float32(i float32) *float32 { return &i }

// Float64 returns a pointer to a given float64 value.
func Float64(i float64) *float64 { return &i }
