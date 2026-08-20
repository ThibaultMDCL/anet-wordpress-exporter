package utils

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type WordpressCollector struct {
	posts      *prometheus.Desc
	categories *prometheus.Desc
	tags       *prometheus.Desc
	pages      *prometheus.Desc
	comments   *prometheus.Desc
	media      *prometheus.Desc
	users      *prometheus.Desc
	adminUsers *prometheus.Desc
	taxonomies *prometheus.Desc
	themes     *prometheus.Desc
	plugins    *prometheus.Desc
	Wp         *Wordpress
}

func (c *WordpressCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.posts
	ch <- c.categories
	ch <- c.tags
	ch <- c.pages
	ch <- c.comments
	ch <- c.media
	ch <- c.users
	ch <- c.adminUsers
	ch <- c.taxonomies
	ch <- c.themes
	ch <- c.plugins
}

func fetchErrCheck(err error, cat string) bool {
	if err != nil {
		log.Printf("unable to fetch  %s: %v, ", cat, err)
		return false
	}
	return true
}

func (c *WordpressCollector) Collect(ch chan<- prometheus.Metric) {
	var err error
	var data []byte

	// TODO : transformer ca par un ernorme foreach et une belle liste pour le scalinggg

	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/categories")
	if fetchErrCheck(err, "categories") {
		c.Wp.categories, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/posts")
	if fetchErrCheck(err, "posts") {
		c.Wp.posts, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/tags")
	if fetchErrCheck(err, "tags") {
		c.Wp.tags, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/pages")
	if fetchErrCheck(err, "pages") {
		c.Wp.pages, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/comments")
	if fetchErrCheck(err, "comments") {
		c.Wp.comments, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/media")
	if fetchErrCheck(err, "media") {
		c.Wp.media, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/users")
	if fetchErrCheck(err, "users") {
		c.Wp.users, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/taxonomies")
	if fetchErrCheck(err, "taxonomies") {
		c.Wp.taxonomies, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/themes")
	if fetchErrCheck(err, "themes") {
		c.Wp.themes, err = CountJSONItems(data)
		ErrCheck(err)
	}
	data, err = c.FetchJSONFromEndpoint("/wp-json/wp/v2/plugins")
	if fetchErrCheck(err, "plugins") {
		c.Wp.plugins, err = CountJSONItems(data)
		ErrCheck(err)
	}

	adminUsers, err := c.FetchTotalFromEndpoint("/wp-json/wp/v2/users?roles=administrator&per_page=1")
	if fetchErrCheck(err, "admin users") {
		c.Wp.adminUsers = adminUsers
	}

	ch <- prometheus.MustNewConstMetric(c.categories, prometheus.GaugeValue, float64(c.Wp.categories))
	ch <- prometheus.MustNewConstMetric(c.posts, prometheus.GaugeValue, float64(c.Wp.posts))
	ch <- prometheus.MustNewConstMetric(c.tags, prometheus.GaugeValue, float64(c.Wp.tags))
	ch <- prometheus.MustNewConstMetric(c.pages, prometheus.GaugeValue, float64(c.Wp.pages))
	ch <- prometheus.MustNewConstMetric(c.comments, prometheus.GaugeValue, float64(c.Wp.comments))
	ch <- prometheus.MustNewConstMetric(c.media, prometheus.GaugeValue, float64(c.Wp.media))
	ch <- prometheus.MustNewConstMetric(c.users, prometheus.GaugeValue, float64(c.Wp.users))
	ch <- prometheus.MustNewConstMetric(c.adminUsers, prometheus.GaugeValue, float64(c.Wp.adminUsers))
	ch <- prometheus.MustNewConstMetric(c.taxonomies, prometheus.GaugeValue, float64(c.Wp.taxonomies))
	ch <- prometheus.MustNewConstMetric(c.themes, prometheus.GaugeValue, float64(c.Wp.themes))
	ch <- prometheus.MustNewConstMetric(c.plugins, prometheus.GaugeValue, float64(c.Wp.plugins))
}

func NewWordpressCollector(w *Wordpress) *WordpressCollector {

	// debug
	fmt.Printf("NewWordpressCollector:\nSite: %v\nUse auth: %v\n", w.MonitoredWordpress, w.Auth.Use)

	return &WordpressCollector{
		Wp:         w,
		posts:      prometheus.NewDesc("wordpress_post_count", "WordPress posts count", nil, prometheus.Labels{"site": w.Name}),
		categories: prometheus.NewDesc("wordpress_category_count", "WordPress category count", nil, prometheus.Labels{"site": w.Name}),
		tags:       prometheus.NewDesc("wordpress_tag_count", "WordPress tags count", nil, prometheus.Labels{"site": w.Name}),
		pages:      prometheus.NewDesc("wordpress_page_count", "WordPress pages count", nil, prometheus.Labels{"site": w.Name}),
		comments:   prometheus.NewDesc("wordpress_comment_count", "WordPress comments count", nil, prometheus.Labels{"site": w.Name}),
		media:      prometheus.NewDesc("wordpress_media_count", "WordPress media files count", nil, prometheus.Labels{"site": w.Name}),
		users:      prometheus.NewDesc("wordpress_user_count", "WordPress users count", nil, prometheus.Labels{"site": w.Name}),
		adminUsers: prometheus.NewDesc("wordpress_admin_user_count", "WordPress administrator user count", nil, prometheus.Labels{"site": w.Name}),
		taxonomies: prometheus.NewDesc("wordpress_taxonomies_count", "WordPress taxonomy count", nil, prometheus.Labels{"site": w.Name}),
		themes:     prometheus.NewDesc("wordpress_theme_count", "WordPress theme count", nil, prometheus.Labels{"site": w.Name}),
		plugins:    prometheus.NewDesc("wordpress_plugin_count", "Wordpress plugin count", nil, prometheus.Labels{"site": w.Name}),
	}
}
