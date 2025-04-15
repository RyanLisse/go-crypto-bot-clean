package config

// SecureHeadersConfig contains configuration for secure HTTP headers
type SecureHeadersConfig struct {
	Enabled                      bool     `mapstructure:"enabled"`
	ContentSecurityPolicy        string   `mapstructure:"content_security_policy"`
	XContentTypeOptions          string   `mapstructure:"x_content_type_options"`
	XFrameOptions                string   `mapstructure:"x_frame_options"`
	XXSSProtection               string   `mapstructure:"x_xss_protection"`
	ReferrerPolicy               string   `mapstructure:"referrer_policy"`
	StrictTransportSecurity      string   `mapstructure:"strict_transport_security"`
	PermissionsPolicy            string   `mapstructure:"permissions_policy"`
	CrossOriginEmbedderPolicy    string   `mapstructure:"cross_origin_embedder_policy"`
	CrossOriginOpenerPolicy      string   `mapstructure:"cross_origin_opener_policy"`
	CrossOriginResourcePolicy    string   `mapstructure:"cross_origin_resource_policy"`
	CacheControl                 string   `mapstructure:"cache_control"`
	ExcludedPaths                []string `mapstructure:"excluded_paths"`
	CustomHeaders                map[string]string `mapstructure:"custom_headers"`
	RemoveServerHeader           bool     `mapstructure:"remove_server_header"`
	RemovePoweredByHeader        bool     `mapstructure:"remove_powered_by_header"`
	ContentSecurityPolicyReportOnly bool   `mapstructure:"content_security_policy_report_only"`
	ContentSecurityPolicyReportURI  string `mapstructure:"content_security_policy_report_uri"`
}

// GetDefaultSecureHeadersConfig returns the default secure headers configuration
func GetDefaultSecureHeadersConfig() SecureHeadersConfig {
	return SecureHeadersConfig{
		Enabled:                   true,
		ContentSecurityPolicy:     "default-src 'self'; script-src 'self'; object-src 'none'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; font-src 'self'; frame-src 'none'; connect-src 'self'",
		XContentTypeOptions:       "nosniff",
		XFrameOptions:             "DENY",
		XXSSProtection:            "1; mode=block",
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		StrictTransportSecurity:   "max-age=31536000; includeSubDomains",
		PermissionsPolicy:         "camera=(), microphone=(), geolocation=(), interest-cohort=()",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginOpenerPolicy:   "same-origin",
		CrossOriginResourcePolicy: "same-origin",
		CacheControl:              "no-store, max-age=0",
		ExcludedPaths:             []string{"/health", "/metrics", "/favicon.ico"},
		CustomHeaders:             map[string]string{},
		RemoveServerHeader:        true,
		RemovePoweredByHeader:     true,
		ContentSecurityPolicyReportOnly: false,
		ContentSecurityPolicyReportURI:  "",
	}
}
