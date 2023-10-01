// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cookiejar implements an in-memory RFC 6265-compliant http.CookieJar.
package cookie_entity

import (
	"fmt"
	"net"
	"net/http/cookiejar"
	"strings"
	"unicode"
	"unicode/utf8"
)

// CanonicalHost strips port from host if present and returns the canonicalized
// host name.
func CanonicalHost(host string) (string, error) {
	var err error
	if hasPort(host) {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return "", err
		}
	}
	// Strip trailing dot from fully qualified domain names.
	host = strings.TrimSuffix(host, ".")
	encoded, err := toASCII(host)
	if err != nil {
		return "", err
	}
	// We know this is ascii, no need to check.
	lower, _ := ToLower(encoded)
	return lower, nil
}

// hasPort reports whether host contains a port number. host may be a host
// name, an IPv4 or an IPv6 address.
func hasPort(host string) bool {
	colons := strings.Count(host, ":")
	if colons == 0 {
		return false
	}
	if colons == 1 {
		return true
	}
	return host[0] == '[' && strings.Contains(host, "]:")
}

// JarKey returns the key to use for a jar.
func JarKey(host string, psl cookiejar.PublicSuffixList) string {
	if isIP(host) {
		return host
	}

	var i int
	if psl == nil {
		i = strings.LastIndex(host, ".")
		if i <= 0 {
			return host
		}
	} else {
		suffix := psl.PublicSuffix(host)
		if suffix == host {
			return host
		}
		i = len(host) - len(suffix)
		if i <= 0 || host[i-1] != '.' {
			// The provided public suffix list psl is broken.
			// Storing cookies under host is a safe stopgap.
			return host
		}
		// Only len(suffix) is used to determine the jar key from
		// here on, so it is okay if psl.PublicSuffix("www.buggy.psl")
		// returns "com" as the jar key is generated from host.
	}
	prevDot := strings.LastIndex(host[:i-1], ".")
	return host[prevDot+1:]
}

// isIP reports whether host is an IP address.
func isIP(host string) bool {
	return net.ParseIP(host) != nil
}

// DefaultPath returns the directory part of a URL's path according to
// RFC 6265 section 5.1.4.
func DefaultPath(path string) string {
	if len(path) == 0 || path[0] != '/' {
		return "/" // Path is empty or malformed.
	}

	i := strings.LastIndex(path, "/") // Path starts with "/", so i != -1.
	if i == 0 {
		return "/" // Path has the form "/abc".
	}
	return path[:i] // Path is either of form "/abc/xyz" or "/abc/xyz/".
}

// toASCII converts a domain or domain label to its ASCII form. For example,
// toASCII("bücher.example.com") is "xn--bcher-kva.example.com", and
// toASCII("golang") is "golang".
func toASCII(s string) (string, error) {
	if Is(s) {
		return s, nil
	}
	labels := strings.Split(s, ".")
	for i, label := range labels {
		if !Is(label) {
			a, err := encode(acePrefix, label)
			if err != nil {
				return "", err
			}
			labels[i] = a
		}
	}
	return strings.Join(labels, "."), nil
}

// Is returns whether s is ASCII.
func Is(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// These parameter values are specified in section 5.
//
// All computation is done with int32s, so that overflow behavior is identical
// regardless of whether int is 32-bit or 64-bit.
const (
	base        int32 = 36
	damp        int32 = 700
	initialBias int32 = 72
	initialN    int32 = 128
	skew        int32 = 38
	tmax        int32 = 26
	tmin        int32 = 1
)

// encode encodes a string as specified in section 6.3 and prepends prefix to
// the result.
//
// The "while h < length(input)" line in the specification becomes "for
// remaining != 0" in the Go code, because len(s) in Go is in bytes, not runes.
func encode(prefix, s string) (string, error) {
	output := make([]byte, len(prefix), len(prefix)+1+2*len(s))
	copy(output, prefix)
	delta, n, bias := int32(0), initialN, initialBias
	b, remaining := int32(0), int32(0)
	for _, r := range s {
		if r < utf8.RuneSelf {
			b++
			output = append(output, byte(r))
		} else {
			remaining++
		}
	}
	h := b
	if b > 0 {
		output = append(output, '-')
	}
	for remaining != 0 {
		m := int32(0x7fffffff)
		for _, r := range s {
			if m > r && r >= n {
				m = r
			}
		}
		delta += (m - n) * (h + 1)
		if delta < 0 {
			return "", fmt.Errorf("cookiejar: invalid label %q", s)
		}
		n = m
		for _, r := range s {
			if r < n {
				delta++
				if delta < 0 {
					return "", fmt.Errorf("cookiejar: invalid label %q", s)
				}
				continue
			}
			if r > n {
				continue
			}
			q := delta
			for k := base; ; k += base {
				t := k - bias
				if t < tmin {
					t = tmin
				} else if t > tmax {
					t = tmax
				}
				if q < t {
					break
				}
				output = append(output, encodeDigit(t+(q-t)%(base-t)))
				q = (q - t) / (base - t)
			}
			output = append(output, encodeDigit(q))
			bias = adapt(delta, h+1, h == b)
			delta = 0
			h++
			remaining--
		}
		delta++
		n++
	}
	return string(output), nil
}

func encodeDigit(digit int32) byte {
	switch {
	case 0 <= digit && digit < 26:
		return byte(digit + 'a')
	case 26 <= digit && digit < 36:
		return byte(digit + ('0' - 26))
	}
	panic("cookiejar: internal error in punycode encoding")
}

// adapt is the bias adaptation function specified in section 6.1.
func adapt(delta, numPoints int32, firstTime bool) int32 {
	if firstTime {
		delta /= damp
	} else {
		delta /= 2
	}
	delta += delta / numPoints
	k := int32(0)
	for delta > ((base-tmin)*tmax)/2 {
		delta /= base - tmin
		k += base
	}
	return k + (base-tmin+1)*delta/(delta+skew)
}

// Strictly speaking, the remaining code below deals with IDNA (RFC 5890 and
// friends) and not Punycode (RFC 3492) per se.

// acePrefix is the ASCII Compatible Encoding prefix.
const acePrefix = "xn--"

// ToLower returns the lowercase version of s if s is ASCII and printable.
func ToLower(s string) (lower string, ok bool) {
	if !IsPrint(s) {
		return "", false
	}
	return strings.ToLower(s), true
}

// IsPrint returns whether s is ASCII and printable according to
// https://tools.ietf.org/html/rfc20#section-4.2.
func IsPrint(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < ' ' || s[i] > '~' {
			return false
		}
	}
	return true
}
