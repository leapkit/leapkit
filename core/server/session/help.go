package session

import "github.com/gorilla/sessions"

// flashHelper is a helper function that can be used in templates
// to retrieve a flash message from the session. This function returns
// that helpers by receiving a pointer to the session.
func flashHelper(session *sessions.Session) func(string) string {
	return func(key string) string {
		val := session.Flashes(key)
		if len(val) == 0 {
			return ""
		}

		return val[0].(string)
	}
}
