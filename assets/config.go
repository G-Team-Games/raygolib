package assets

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type Kind string

const (
	KindModel   Kind = "model"
	KindTexture Kind = "texture"
	KindImage   Kind = "image"
	KindSound   Kind = "sound"
	KindMusic   Kind = "music"
	KindFont    Kind = "font"
	KindShader  Kind = "shader"
)

type ThreadPolicy int

const (
	ThreadPolicyStrict ThreadPolicy = iota
	ThreadPolicyQueueOnly
)

type AssetRef struct {
	Kind Kind
	Key  string
}

type Resolver interface {
	Resolve(kind Kind, key string) (string, error)
	KeysForPath(path string) []AssetRef
	WatchRoots() []string
}

type KindRule struct {
	Kind     Kind
	Include  []string
	Exclude  []string
	Priority int
}

type Config struct {
	Resolver      Resolver
	Rules         []KindRule
	HotReload     bool
	WatchDebounce time.Duration
	ThreadPolicy  ThreadPolicy
}

type Option func(*Config) error

func DefaultConfig() Config {
	return Config{
		Resolver:      NewFixedDirsResolver(assetsBasePath),
		Rules:         nil,
		HotReload:     false,
		WatchDebounce: 100 * time.Millisecond,
		ThreadPolicy:  ThreadPolicyStrict,
	}
}

type Manager struct {
	*AssetManager
	cfg Config
}

func NewManager(opts ...Option) (*Manager, error) {
	cfg := DefaultConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(&cfg); err != nil {
			return nil, err
		}
	}

	if cfg.Resolver == nil {
		return nil, fmt.Errorf("assets: resolver is required")
	}

	return &Manager{
		AssetManager: NewAssetManager(),
		cfg:          cfg,
	}, nil
}

func (m *Manager) Config() Config {
	return m.cfg
}

func WithConfig(cfg Config) Option {
	return func(dst *Config) error {
		*dst = cfg
		if dst.WatchDebounce <= 0 {
			dst.WatchDebounce = 100 * time.Millisecond
		}
		return nil
	}
}

func WithResolver(resolver Resolver) Option {
	return func(cfg *Config) error {
		if resolver == nil {
			return fmt.Errorf("assets: resolver is nil")
		}
		cfg.Resolver = resolver
		return nil
	}
}

func WithSingleRoot(root string) Option {
	return func(cfg *Config) error {
		cfg.Resolver = NewSingleRootResolver(root)
		return nil
	}
}

func WithFixedDirs(root string) Option {
	return func(cfg *Config) error {
		cfg.Resolver = NewFixedDirsResolver(root)
		return nil
	}
}

func WithRule(rule KindRule) Option {
	return func(cfg *Config) error {
		cfg.Rules = append(cfg.Rules, rule)
		return nil
	}
}

func WithRules(rules []KindRule) Option {
	return func(cfg *Config) error {
		cfg.Rules = append([]KindRule(nil), rules...)
		return nil
	}
}

func WithHotReload(enabled bool) Option {
	return func(cfg *Config) error {
		cfg.HotReload = enabled
		return nil
	}
}

func WithWatchDebounce(d time.Duration) Option {
	return func(cfg *Config) error {
		if d <= 0 {
			return fmt.Errorf("assets: watch debounce must be > 0")
		}
		cfg.WatchDebounce = d
		return nil
	}
}

func WithThreadPolicy(policy ThreadPolicy) Option {
	return func(cfg *Config) error {
		switch policy {
		case ThreadPolicyStrict, ThreadPolicyQueueOnly:
			cfg.ThreadPolicy = policy
			return nil
		default:
			return fmt.Errorf("assets: invalid thread policy: %d", policy)
		}
	}
}

type FixedDirsResolver struct {
	root string
	dirs map[Kind]string
}

func NewFixedDirsResolver(root string) *FixedDirsResolver {
	if root == "" {
		root = assetsBasePath
	}

	return &FixedDirsResolver{
		root: root,
		dirs: map[Kind]string{
			KindModel:   "models",
			KindTexture: "textures",
			KindImage:   "images",
			KindSound:   "audio",
			KindMusic:   "audio",
			KindFont:    "fonts",
			KindShader:  "shaders",
		},
	}
}

func (r *FixedDirsResolver) Resolve(kind Kind, key string) (string, error) {
	dir, ok := r.dirs[kind]
	if !ok {
		return "", fmt.Errorf("assets: unsupported kind %q", kind)
	}
	return filepath.Join(r.root, dir, key), nil
}

func (r *FixedDirsResolver) KeysForPath(path string) []AssetRef {
	norm := filepath.Clean(path)
	for kind, dir := range r.dirs {
		prefix := filepath.Join(r.root, dir) + string(filepath.Separator)
		if after, ok :=strings.CutPrefix(norm, prefix); ok  {
			return []AssetRef{{Kind: kind, Key: after}}
		}
	}
	return nil
}

func (r *FixedDirsResolver) WatchRoots() []string {
	roots := make([]string, 0, len(r.dirs))
	seen := map[string]struct{}{}
	for _, dir := range r.dirs {
		p := filepath.Join(r.root, dir)
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		roots = append(roots, p)
	}
	return roots
}

type SingleRootResolver struct {
	root string
}

func NewSingleRootResolver(root string) *SingleRootResolver {
	if root == "" {
		root = assetsBasePath
	}
	return &SingleRootResolver{root: root}
}

func (r *SingleRootResolver) Resolve(kind Kind, key string) (string, error) {
	return filepath.Join(r.root, key), nil
}

func (r *SingleRootResolver) KeysForPath(path string) []AssetRef {
	norm := filepath.Clean(path)
	prefix := filepath.Clean(r.root) + string(filepath.Separator)
	if strings.HasPrefix(norm, prefix) {
		key := strings.TrimPrefix(norm, prefix)
		return []AssetRef{{Kind: "", Key: key}}
	}
	return nil
}

func (r *SingleRootResolver) WatchRoots() []string {
	return []string{r.root}
}
