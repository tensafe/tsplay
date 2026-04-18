package tsplay_core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

type PageObservationOptions struct {
	URL          string
	Headless     bool
	ArtifactRoot string
	TimeoutMS    int
	MaxElements  int
	Context      context.Context
	RunID        string
	RunRoot      string
}

type PageObservation struct {
	URL             string                   `json:"url"`
	Title           string                   `json:"title,omitempty"`
	ArtifactRoot    string                   `json:"artifact_root,omitempty"`
	ScreenshotPath  string                   `json:"screenshot_path,omitempty"`
	DOMSnapshotPath string                   `json:"dom_snapshot_path,omitempty"`
	Elements        []PageObservationElement `json:"elements"`
	Errors          []string                 `json:"errors,omitempty"`
}

type PageObservationElement struct {
	Index              int                       `json:"index"`
	Tag                string                    `json:"tag,omitempty"`
	Type               string                    `json:"type,omitempty"`
	Role               string                    `json:"role,omitempty"`
	ID                 string                    `json:"id,omitempty"`
	Name               string                    `json:"name,omitempty"`
	Text               string                    `json:"text,omitempty"`
	Label              string                    `json:"label,omitempty"`
	Placeholder        string                    `json:"placeholder,omitempty"`
	AriaLabel          string                    `json:"aria_label,omitempty"`
	Href               string                    `json:"href,omitempty"`
	Value              string                    `json:"value,omitempty"`
	Visible            bool                      `json:"visible"`
	Enabled            bool                      `json:"enabled"`
	NearText           string                    `json:"near_text,omitempty"`
	PrimarySelector    string                    `json:"primary_selector,omitempty"`
	SelectorRationale  string                    `json:"selector_rationale,omitempty"`
	SelectorCandidates []string                  `json:"selector_candidates,omitempty"`
	SelectorDetails    []PageObservationSelector `json:"selector_details,omitempty"`
	BoundingBox        *PageObservationBox       `json:"bounding_box,omitempty"`
	Attributes         map[string]string         `json:"attributes,omitempty"`
}

type PageObservationBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

func ObservePage(options PageObservationOptions) (*PageObservation, error) {
	if strings.TrimSpace(options.URL) == "" {
		return nil, fmt.Errorf("url is required")
	}
	if options.Context != nil {
		if err := options.Context.Err(); err != nil {
			return nil, err
		}
	}

	pw, err := StartPlaywright()
	if err != nil {
		return nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(options.Headless),
	})
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %w", err)
	}
	closePlaywright := sync.OnceFunc(func() {
		if page != nil {
			_ = page.Close()
		}
		if browser != nil {
			_ = browser.Close()
		}
		_ = pw.Stop()
	})
	defer closePlaywright()
	stopWatcher := watchContextCancel(options.Context, closePlaywright)
	defer stopWatcher()

	timeout := options.TimeoutMS
	if timeout <= 0 {
		timeout = 30000
	}
	if _, err := page.Goto(options.URL, playwright.PageGotoOptions{
		Timeout:   playwright.Float(float64(timeout)),
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		return nil, fmt.Errorf("navigate to %q: %w", options.URL, err)
	}

	return ObserveLoadedPage(page, options)
}

func ObserveLoadedPage(page playwright.Page, options PageObservationOptions) (*PageObservation, error) {
	if page == nil {
		return nil, fmt.Errorf("page is nil")
	}

	artifactRoot := strings.TrimSpace(options.ArtifactRoot)
	if artifactRoot == "" {
		artifactRoot = DefaultFlowArtifactRoot
	}
	root, err := prepareRuntimeFileRoot(artifactRoot)
	if err != nil {
		return nil, fmt.Errorf("prepare artifact root %q: %w", artifactRoot, err)
	}

	dir := filepath.Join(root, "observe-"+time.Now().Format("20060102-150405.000000000"))
	if strings.TrimSpace(options.RunRoot) != "" {
		dir = filepath.Join(options.RunRoot, "observe")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create observation directory: %w", err)
	}

	title, err := page.Title()
	errors := []string{}
	if err != nil {
		errors = append(errors, fmt.Sprintf("title: %v", err))
	}

	observation := &PageObservation{
		URL:          page.URL(),
		Title:        title,
		ArtifactRoot: firstNonEmpty(strings.TrimSpace(options.RunRoot), root),
		Elements:     []PageObservationElement{},
		Errors:       errors,
	}

	screenshotPath := filepath.Join(dir, "observe.png")
	if _, err := page.Screenshot(playwright.PageScreenshotOptions{Path: playwright.String(screenshotPath), FullPage: playwright.Bool(true)}); err != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("screenshot: %v", err))
	} else {
		observation.ScreenshotPath = screenshotPath
	}

	domSnapshotPath := filepath.Join(dir, "dom_snapshot.json")
	if snapshot, err := ExtractSimplifiedElementWithXPathResult(page); err != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("dom snapshot: %v", err))
	} else if err := os.WriteFile(domSnapshotPath, []byte(snapshot), 0644); err != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("dom snapshot write: %v", err))
	} else {
		observation.DOMSnapshotPath = domSnapshotPath
	}

	elements, err := observeInteractiveElements(page)
	if err != nil {
		observation.Errors = append(observation.Errors, fmt.Sprintf("elements: %v", err))
		return observation, nil
	}
	maxElements := options.MaxElements
	if maxElements <= 0 {
		maxElements = 100
	}
	if len(elements) > maxElements {
		elements = elements[:maxElements]
	}
	for i := range elements {
		elements[i].Index = i + 1
	}
	observation.Elements = elements
	return observation, nil
}

func observeInteractiveElements(page playwright.Page) ([]PageObservationElement, error) {
	value, err := page.Evaluate(`() => {
		const candidates = [
			'a[href]',
			'button',
			'input',
			'textarea',
			'select',
			'[role="button"]',
			'[role="link"]',
			'[role="textbox"]',
			'[role="checkbox"]',
			'[role="radio"]',
			'[contenteditable="true"]',
			'[onclick]',
			'[tabindex]'
		].join(',');

		const clean = (value, max = 160) => {
			if (!value) return '';
			const text = String(value).replace(/\s+/g, ' ').trim();
			return text.length > max ? text.slice(0, max) + '...' : text;
		};
		const cssEscape = (value) => {
			if (window.CSS && window.CSS.escape) return window.CSS.escape(value);
			return String(value).replace(/["\\#.:>+~*^[\]$()=|/@]/g, '\\$&');
		};
		const quote = (value) => JSON.stringify(String(value));
		const addUnique = (items, value) => {
			if (value && !items.includes(value)) items.push(value);
		};
		const visible = (element) => {
			const style = window.getComputedStyle(element);
			const rect = element.getBoundingClientRect();
			return style.display !== 'none' &&
				style.visibility !== 'hidden' &&
				style.opacity !== '0' &&
				rect.width > 0 &&
				rect.height > 0;
		};
		const labelFor = (element) => {
			if (element.labels && element.labels.length > 0) {
				return clean(Array.from(element.labels).map(label => label.innerText || label.textContent).join(' '));
			}
			if (element.id) {
				const label = document.querySelector('label[for="' + cssEscape(element.id) + '"]');
				if (label) return clean(label.innerText || label.textContent);
			}
			const parentLabel = element.closest && element.closest('label');
			if (parentLabel) return clean(parentLabel.innerText || parentLabel.textContent);
			return '';
		};
		const generateXPath = (element) => {
			if (element.id) return '//*[@id="' + element.id.replace(/"/g, '\\"') + '"]';
			if (element === document.documentElement) return '/html';
			if (element === document.body) return '/html/body';
			let ix = 0;
			const siblings = element.parentNode ? Array.from(element.parentNode.childNodes) : [];
			for (let i = 0; i < siblings.length; i++) {
				const sibling = siblings[i];
				if (sibling === element) {
					return generateXPath(element.parentNode) + '/' + element.tagName.toLowerCase() + '[' + (ix + 1) + ']';
				}
				if (sibling.nodeType === 1 && sibling.tagName === element.tagName) ix++;
			}
			return '';
		};
		const inferredRole = (element, tag, inputType) => {
			const explicit = element.getAttribute('role');
			if (explicit) return explicit;
			if (tag === 'button' || inputType === 'button' || inputType === 'submit') return 'button';
			if (tag === 'a') return 'link';
			if (tag === 'select') return 'combobox';
			if (tag === 'textarea' || ['text', 'search', 'email', 'password', 'number', 'tel', 'url'].includes(inputType)) return 'textbox';
			if (inputType === 'checkbox') return 'checkbox';
			if (inputType === 'radio') return 'radio';
			return '';
		};
		const elementType = (tag, inputType) => {
			if (tag === 'input') return inputType || 'text';
			if (tag === 'a') return 'link';
			return tag;
		};
		const nearbyText = (element) => {
			const parent = element.closest('form, section, article, main, [role="dialog"], [role="region"], div');
			if (!parent || parent === element) return '';
			return clean(parent.innerText || parent.textContent, 220);
		};
		const attributesFor = (element) => {
			const names = ['data-testid', 'data-test', 'data-cy', 'id', 'name', 'type', 'placeholder', 'aria-label', 'href'];
			const attrs = {};
			for (const name of names) {
				const value = element.getAttribute(name);
				if (value) attrs[name] = value;
			}
			return attrs;
		};
		const selectorCandidates = (element, tag, inputType, text, label, placeholder, ariaLabel, href, role) => {
			const selectors = [];
			for (const attr of ['data-testid', 'data-test', 'data-cy']) {
				const value = element.getAttribute(attr);
				if (value) {
					addUnique(selectors, '[' + attr + '=' + quote(value) + ']');
					addUnique(selectors, tag + '[' + attr + '=' + quote(value) + ']');
				}
			}
			if (element.id) {
				addUnique(selectors, '#' + cssEscape(element.id));
				addUnique(selectors, tag + '#' + cssEscape(element.id));
			}
			if (element.name) {
				addUnique(selectors, tag + '[name=' + quote(element.name) + ']');
				if (inputType) addUnique(selectors, tag + '[name=' + quote(element.name) + '][type=' + quote(inputType) + ']');
			}
			if (placeholder) addUnique(selectors, tag + '[placeholder=' + quote(placeholder) + ']');
			if (ariaLabel) {
				addUnique(selectors, '[aria-label=' + quote(ariaLabel) + ']');
				addUnique(selectors, tag + '[aria-label=' + quote(ariaLabel) + ']');
			}
			if (href && tag === 'a') addUnique(selectors, 'a[href=' + quote(href) + ']');
			if (label && role) addUnique(selectors, 'role=' + role + '[name=' + quote(label) + ']');
			if (text && ['button', 'a'].includes(tag)) addUnique(selectors, 'text=' + quote(text));
			if (text && role) addUnique(selectors, 'role=' + role + '[name=' + quote(text) + ']');
			const xpath = generateXPath(element);
			if (xpath) addUnique(selectors, 'xpath=' + xpath);
			return selectors;
		};

		return Array.from(document.querySelectorAll(candidates))
			.filter((element) => visible(element))
			.map((element) => {
				const tag = element.tagName.toLowerCase();
				const inputType = tag === 'input' ? (element.getAttribute('type') || 'text').toLowerCase() : '';
				const text = clean(element.innerText || element.textContent || element.value);
				const label = labelFor(element);
				const placeholder = clean(element.getAttribute('placeholder'));
				const ariaLabel = clean(element.getAttribute('aria-label'));
				const role = inferredRole(element, tag, inputType);
				const rect = element.getBoundingClientRect();
				const enabled = !element.disabled && element.getAttribute('aria-disabled') !== 'true';
				return {
					tag,
					type: elementType(tag, inputType),
					role,
					id: element.id || '',
					name: element.getAttribute('name') || '',
					text,
					label,
					placeholder,
					aria_label: ariaLabel,
					href: element.getAttribute('href') || '',
					value: clean(element.value),
					visible: true,
					enabled,
					near_text: nearbyText(element),
					selector_candidates: selectorCandidates(element, tag, inputType, text, label, placeholder, ariaLabel, element.getAttribute('href') || '', role),
					bounding_box: {
						x: rect.x,
						y: rect.y,
						width: rect.width,
						height: rect.height
					},
					attributes: attributesFor(element)
				};
			})
			.sort((a, b) => {
				if (a.bounding_box.y !== b.bounding_box.y) return a.bounding_box.y - b.bounding_box.y;
				return a.bounding_box.x - b.bounding_box.x;
			});
	}`)
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal observed elements: %w", err)
	}
	var elements []PageObservationElement
	if err := json.Unmarshal(encoded, &elements); err != nil {
		return nil, fmt.Errorf("decode observed elements: %w", err)
	}
	for i := range elements {
		normalizeObservedSelectorDiagnostics(&elements[i])
	}
	sort.SliceStable(elements, func(i, j int) bool {
		if elements[i].BoundingBox == nil || elements[j].BoundingBox == nil {
			return i < j
		}
		if elements[i].BoundingBox.Y != elements[j].BoundingBox.Y {
			return elements[i].BoundingBox.Y < elements[j].BoundingBox.Y
		}
		return elements[i].BoundingBox.X < elements[j].BoundingBox.X
	})
	return elements, nil
}
