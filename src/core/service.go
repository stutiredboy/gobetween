package core

/**
 * Service interface
 */
type Service interface {
	Apply(Server) error
	Forget(Server) error
}
