package db

import "fmt"

// Service provides various functions to operate on Models. All implementations
// should use type assertions to guarantee prevention of runtime errors.
type Service interface {
	Bucket() string

	Clean(m Model, tx Tx) error
	Validate(m Model, tx Tx) error
	Initialize(m Model, tx Tx) error
	PersistOldProperties(n Model, o Model, tx Tx) error

	PersistHooks() *PersistHooks

	Marshal(m Model) ([]byte, error)
	Unmarshal(buf []byte) (Model, error)
}

// PersistHookFunc is a callback used as a hook function.
type PersistHookFunc = func(m Model, ser Service, tx Tx) error

// PersistHooks provides hook functions to be called before and after service
// operations.
type PersistHooks struct {
	PreCreateHooks  []PersistHookFunc
	PostCreateHooks []PersistHookFunc
	PreUpdateHooks  []PersistHookFunc
	PostUpdateHooks []PersistHookFunc
	PreDeleteHooks  []PersistHookFunc
	PostDeleteHooks []PersistHookFunc
}

// PreCreateHook executes all hook functions designated to be called before
// create operations.
func (hooks *PersistHooks) PreCreateHook(m Model, ser Service, tx Tx) error {
	return hooks.callHooks(hooks.PreCreateHooks, m, ser, tx)
}

// PostCreateHook executes all hook functions designated to be called after
// create operations.
func (hooks *PersistHooks) PostCreateHook(m Model, ser Service, tx Tx) error {
	return hooks.callHooks(hooks.PostCreateHooks, m, ser, tx)
}

// PreUpdateHook executes all hook functions designated to be called before
// update operations.
func (hooks *PersistHooks) PreUpdateHook(m Model, ser Service, tx Tx) error {
	return hooks.callHooks(hooks.PreUpdateHooks, m, ser, tx)
}

// PostUpdateHook executes all hook functions designated to be called after
// update operations.
func (hooks *PersistHooks) PostUpdateHook(m Model, ser Service, tx Tx) error {
	return hooks.callHooks(hooks.PostUpdateHooks, m, ser, tx)
}

// PreDeleteHook executes all hook functions designated to be called before
// delete operations.
func (hooks *PersistHooks) PreDeleteHook(m Model, ser Service, tx Tx) error {
	return hooks.callHooks(hooks.PreDeleteHooks, m, ser, tx)
}

// PostDeleteHook executes all hook functions designated to be called after
// delete operations.
func (hooks *PersistHooks) PostDeleteHook(m Model, ser Service, tx Tx) error {
	return hooks.callHooks(hooks.PostDeleteHooks, m, ser, tx)
}

func (hooks *PersistHooks) callHooks(list []PersistHookFunc, m Model, ser Service, tx Tx) error {
	for _, h := range list {
		if h == nil {
			continue
		}
		err := h(m, ser, tx)
		if err != nil {
			return fmt.Errorf("failed to execute persist hook: %w", err)
		}
	}

	return nil
}
