package server

import "context"

// GracefulShutdown gracefully shuts down the server
func (s *Server) GracefulShutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}
