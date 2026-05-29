package scrapper

import (
	"net/url"
	"reflect"
	"testing"
)

func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := ExtractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %+v, got: %+v", expected, actual)
	}

}

func TestGetHeadingHTML(t *testing.T) {
	tcs := []struct {
		name      string
		inputHTML string

		expected string
	}{
		{
			name:      "with H1",
			inputHTML: "<html><body><nav><h1>Hello Scrapper</h1></nav></body></hmtl>",
			expected:  "Hello Scrapper",
		},
		{
			name:      "with h2",
			inputHTML: "<html><body><div><div><h2>a subtitle</h2><div></div><p>paraghrap here</p></div></div></body></hmtl>",
			expected:  "a subtitle",
		},
		{
			name:      "non heading",
			inputHTML: "<html><body><nav><a>some ling</a></nav><div>with text</div></body></hmtl>",
			expected:  "",
		},
		{
			name:      "with H1 and h2",
			inputHTML: "<html><body><h2>I'm first</h2><nav><h1>I'm more important</h1></nav></body></hmtl>",
			expected:  "I'm more important",
		},
		{
			name:      "invalid html",
			inputHTML: "",
			expected:  "",
		},
		{
			name:      "empty tags",
			inputHTML: "<html><body><h2></h2><nav><h1></h1></nav></body></hmtl>",
			expected:  "",
		},
		{
			name:      "empty h1 tag",
			inputHTML: "<html><body><h2>Hello I'm only here</h2><nav><h1></h1></nav></body></hmtl>",
			expected:  "Hello I'm only here",
		},
	}

	for i, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingHTML(tc.inputHTML)
			if tc.expected != actual {
				t.Fatalf("Test %v - '%s' FAIL: expected: %s, actual: %s", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tcs := []struct {
		name      string
		inputHTML string

		expected string
	}{
		{
			name: "with main and P",
			inputHTML: `
			<html>
			  <body>
					<p>outside</p>
					<main>
						<p>inside main</p>
					</main>
			  </body>
			</html>
		`,
			expected: "inside main",
		},
		{
			name: "no main P",
			inputHTML: `
			<html>
			  <body>
					<p>outside</p>
					<p>outside second</p>
			  </body>
			</html>
		`,
			expected: "outside",
		},
		{
			name: "with main and multi p",
			inputHTML: `
			<html>
			  <body>
					<p>outside</p>
					<main>
						<p>inside main</p>
						<p>inside main2</p>
						<p>inside main3</p>
					</main>
			  </body>
			</html>
		`,
			expected: "inside main",
		},
		{
			name: "none",
			inputHTML: `
			<html>
			  <body>
					<main>
					</main>
			  </body>
			</html>
		`,
			expected: "",
		},
	}

	for i, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			finded := getFirstParagraphFromHTML(tc.inputHTML)
			if tc.expected != finded {
				t.Fatalf("Test %v - '%s' FAIL: expected: %s, founded: %s ", i, tc.name, tc.expected, finded)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tcs := []struct {
		name      string
		inputHTML string
		inputURL  string

		expected []string
	}{
		{
			name:      "base get single url",
			inputURL:  "https://crawler-test.com",
			inputHTML: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,

			expected: []string{"https://crawler-test.com"},
		},
		{
			name:     "base get multi url",
			inputURL: "https://crawler-test.com",
			inputHTML: `<html><body>
			<a href="https://crawler-test.com"><span>Boot.dev</span></a>
			<a href="/data.html"><span>another</span></a>
			</body></html>`,

			expected: []string{"https://crawler-test.com", "https://crawler-test.com/data.html"},
		},
		{
			name:      "empty",
			inputURL:  "https://crawler-test.com",
			inputHTML: ``,

			expected: nil,
		},
		{
			name:     "non urls",
			inputURL: "https://crawler-test.com",
			inputHTML: `<html><body>
			<a><span>Boot.dev</span></a>
			<a><span>another</span></a>
			</body></html>`,

			expected: []string{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			url, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Fatalf("%s - FAIL with parsing err %s", tc.name, err)
			}
			result, err := getURLsFromHTML(tc.inputHTML, url)
			if err != nil {
				t.Fatalf("%s FAIL with unexpected error: %s", tc.name, err)
			}

			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("expected: %v, got: %v", tc.expected, result)
			}

		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	tcs := []struct {
		name     string
		inputHML string
		inputURL string

		expected []string
	}{
		{
			name:     "base case",
			inputURL: "https://crawler-test.com",
			inputHML: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,

			expected: []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:     "base case 1+ links",
			inputURL: "https://crawler-test.com",
			inputHML: `<html><body>
			<img src="/logo.png" alt="Logo">
			<img src="/logo2.png" alt="Logo2">
			<div>
				<img src="/logo-inside.png" alt="Logo-inside">
			</div>
			</body></html>`,

			expected: []string{"https://crawler-test.com/logo.png", "https://crawler-test.com/logo2.png", "https://crawler-test.com/logo-inside.png"},
		},
		{
			name:     "empty src",
			inputURL: "https://crawler-test.com",
			inputHML: `<html><body><img src="" alt="Logo"></body></html>`,

			expected: []string{},
		},
		{
			name:     "non img elements",
			inputURL: "https://crawler-test.com",
			inputHML: `<html><body>
			<div>
			<a hrf="#" >blah</a>
			</div>
			</body></html>`,

			expected: nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Fatalf("getImagesFromHTML:  error parsing url: %s", err)
			}
			result, err := getImagesFromHTML(tc.inputHML, parsed)
			if err != nil {
				t.Fatalf("getImagesFromHTML unexpected error when getting images: %s", err)
			}
			if !reflect.DeepEqual(tc.expected, result) {
				t.Fatalf("getImagesFromHTML FAIL: expect: %v, got: %v", tc.expected, result)
			}
		})
	}
}
