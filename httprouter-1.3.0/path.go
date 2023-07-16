// Copyright 2013 Julien Schmidt. All rights reserved.
// Based on the path package, Copyright 2009 The Go Authors.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package httprouter

// CleanPath is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//  1. Replace multiple slashes with a single slash.
//  2. Eliminate each . path name element (the current directory).
//  3. Eliminate each inner .. path name element (the parent directory)
//     along with the non-.. element that precedes it.
//  4. Eliminate .. elements that begin a rooted path:
//     that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned.
//
// 预处理 URL ，使其成为一个规范的路径
func CleanPath(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1 // r 是 p 的指针
	w := 1 // w 是 buf 的 指针

	// 如果不以 '/' 开头，此刻立即初始化 buf
	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	// trailing type: bool 判断最后面的 byte 是否 '/'
	trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)
	// 有点笨拙的方法，但是很快
	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++
		// 如果最后面的 byte 是 '.' 将其视作 '/' 并跳过
		case p[r] == '.' && r+1 == n:
			trailing = true
			r++
		// 如果遇到了 './' 则跳过
		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2
		// 如果遇到了 '..' 且满足以下任一条件：1. 这个出现在末尾 2. 下一个是 '/' (就是排除三个点或以上的情况)
		// trans: 如果遇到了 '../' 或者 末尾的 '..'
		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3 // 先移动读指针
			if w > 1 {
				// can backtrack
				w--
				if buf == nil {
					// 如果 buf 还没初始化说明 p 之前的都是无效的
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}
		// 如果遇到正常的 token，注意: '...' 也是正常的
		default:
			// real path element.
			// add slash if needed
			// 给之前写好的 token 加上 '/'
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}
			// copy element
			// 这里是使用 for 一次性将 token 写入进 buf
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	// 重新附加末尾斜杠
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}

// internal helper to lazily create a buffer if necessary
func bufApp(buf *[]byte, s string, w int, c byte) {
	// 如果 buf 还没初始化
	if *buf == nil {
		if s[w] == c {
			return
		}
		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}
