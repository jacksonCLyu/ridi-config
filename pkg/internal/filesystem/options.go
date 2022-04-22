package filesystem

type options struct {
	currentUser         string
	versioning          int
	proxyHost           string
	proxyPort           int
	maxHostConnections  int
	maxTotalConnections int
}

// Option option interface for config file system
type Option interface {
	apply(opts *options)
}

// WithCurrentUser with current user option
func WithCurrentUser(currentUser string) Option {
	return currentUserOption(currentUser)
}

// WithVersion with version option
func WithVersion(version int) Option {
	return versionOption(version)
}

// WithProxyHost with proxy host option
func WithProxyHost(proxyHost string) Option {
	return proxyHostOption(proxyHost)
}

// WithProxyPort with proxy port option
func WithProxyPort(proxyPort int) Option {
	return proxyPortOption(proxyPort)
}

// WithMaxHostConn with max host connections option
func WithMaxHostConn(maxHostConn int) Option {
	return maxHostConnOption(maxHostConn)
}

// WithMaxTotalConn with max total connections option
func WithMaxTotalConn(maxTotalConn int) Option {
	return maxTotalConnOption(maxTotalConn)
}

type currentUserOption string

func (o currentUserOption) apply(opts *options) {
	opts.currentUser = string(o)
}

type versionOption int

func (o versionOption) apply(opts *options) {
	opts.versioning = int(o)
}

type proxyHostOption string

func (o proxyHostOption) apply(opts *options) {
	opts.proxyHost = string(o)
}

type proxyPortOption int

func (o proxyPortOption) apply(opts *options) {
	opts.proxyPort = int(o)
}

type maxHostConnOption int

func (o maxHostConnOption) apply(opts *options) {
	opts.maxHostConnections = int(o)
}

type maxTotalConnOption int

func (o maxTotalConnOption) apply(opts *options) {
	opts.maxTotalConnections = int(o)
}
