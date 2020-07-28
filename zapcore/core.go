package zapcore

// Core is a minimal, fast logger interface. It's designed for library authors
// to wrap in a more user-friendly API.
type Core interface {
	LevelEnabler

	// With adds structured context to the Core.
	With([]Field) Core
	// Check determines whether the supplied Entry should be logged (using the
	// embedded LevelEnabler and possibly some extra logic). If the entry
	// should be logged, the Core adds itself to the CheckedEntry and returns
	// the result.
	//
	// Callers must use Check before calling Write.
	Check(Entry, *CheckedEntry) *CheckedEntry
	// Write serializes the Entry and any Fields supplied at the log site and
	// writes them to their destination.
	//
	// If called, Write should always log the Entry and Fields; it should not
	// replicate the logic of Check.
	Write(Entry, []Field) error
	// Sync flushes buffered logs (if any).
	Sync() error
}

// NOpCore nop core
type NOpCore struct{}

// NewNopCore returns a no-op Core.
func NewNopCore() Core {
	return NOpCore{}
}

// Enabled enable func
func (NOpCore) Enabled(Level) bool {
	return false
}

// With with func
func (n NOpCore) With([]Field) Core {
	return n
}

// Check check func
func (NOpCore) Check(_ Entry, ce *CheckedEntry) *CheckedEntry {
	return ce
}

// Write write func
func (NOpCore) Write(Entry, []Field) error {
	return nil
}

// Sync sync func
func (NOpCore) Sync() error {
	return nil
}

// NewCore creates a Core that writes logs to a WriteSyncer.
func NewCore(enc Encoder, ws WriteSyncer, enab LevelEnabler) Core {
	return &IOCore{
		LevelEnabler: enab,
		enc:          enc,
		out:          ws,
	}
}

// IOCore io core
type IOCore struct {
	LevelEnabler
	enc Encoder
	out WriteSyncer
}

// With write func
func (c *IOCore) With(fields []Field) Core {
	clone := c.clone()
	addFields(clone.enc, fields)
	return clone
}

// Check check func
func (c *IOCore) Check(ent Entry, ce *CheckedEntry) *CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *IOCore) Write(ent Entry, fields []Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	_, err = c.out.Write(buf.Bytes())
	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		c.Sync()
	}
	return nil
}

// Sync sync func
func (c *IOCore) Sync() error {
	return c.out.Sync()
}

func (c *IOCore) clone() *IOCore {
	return &IOCore{
		LevelEnabler: c.LevelEnabler,
		enc:          c.enc.Clone(),
		out:          c.out,
	}
}
