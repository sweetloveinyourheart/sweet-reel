package auth

type options struct {
	override ServiceAuthFuncOverride
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

// WithOverride customizes the function
func WithOverride(f ServiceAuthFuncOverride) Option {
	return func(o *options) {
		o.override = f
	}
}
