package config

import "testing"

func TestProxyEnvEmpty(t *testing.T) {
	if got := proxyEnv(""); got != nil {
		t.Fatalf("empty proxy should return nil, got %v", got)
	}
	if got := proxyEnv("   "); got != nil {
		t.Fatalf("whitespace proxy should return nil, got %v", got)
	}
}

func TestProxyEnvSetsAllVariants(t *testing.T) {
	env := proxyEnv("socks5://127.0.0.1:9050")
	for _, key := range []string{"HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY", "http_proxy", "https_proxy", "all_proxy"} {
		if env[key] != "socks5://127.0.0.1:9050" {
			t.Fatalf("%s = %q, want proxy", key, env[key])
		}
	}
	if env["NO_PROXY"] != "localhost,127.0.0.1" {
		t.Fatalf("NO_PROXY = %q", env["NO_PROXY"])
	}
}

func TestProxyEnvTrimsWhitespace(t *testing.T) {
	env := proxyEnv("  http://p:8080  ")
	if env["HTTP_PROXY"] != "http://p:8080" {
		t.Fatalf("expected trimmed proxy, got %q", env["HTTP_PROXY"])
	}
}
