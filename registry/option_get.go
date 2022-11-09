package registry

import "context"

type GetOptions struct {
	Context context.Context
	Version string
}

func NewGetOption(opts ...GetOption) *GetOptions {
	o := &GetOptions{
		Context: context.TODO(),
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

type GetOption func(*GetOptions)

func GetContext(ctx context.Context) GetOption {
	return func(o *GetOptions) {
		o.Context = ctx
	}
}

func GetVersion(version string) GetOption {
	return func(o *GetOptions) {
		o.Version = version
	}
}
