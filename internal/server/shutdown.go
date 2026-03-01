package server

import "context"

// GracefulShutdown gracefully shuts down the server
func (s *Server) GracefulShutdown(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}
