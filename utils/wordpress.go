package utils

type Wordpress struct {
	posts, categories, tags, pages, comments, media, users, adminUsers, themes, plugins, taxonomies int
	Name                                                                                            string
	MonitoredWordpress                                                                              string
	UserAgent                                                                                       string
	Auth                                                                                            WPAuth
}

type WPAuth struct {
	Use      bool
	Username string
	Password string
}

func NewWordpress(name string, monitor string, ua string, authuser string, authpass string, useAuth bool) *Wordpress {
	return &Wordpress{
		Name:               name,
		MonitoredWordpress: monitor,
		UserAgent:          ua,
		Auth: WPAuth{
			Use:      useAuth,
			Username: authuser,
			Password: authpass,
		},
	}
}
