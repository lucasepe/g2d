0.5.0 (Unreleased)

NEW FEATURES:

- commandline tool can load `.g2d` scripts from your local filesystem or from HTTP (see [README](README.md))
- `viewport(...)` builtin function takes two more optional parameters `xOffset` and `yOffset` to eventually relocate the viewport
- new builtin function: `fillColor` to specify the fill color (takes `r`, `g`, `b`, `[a]` range [`0-255`] or hex color string `#rrggbbaa`)
- new builtin function: `strokeColor` to specify the stroke color (takes `r`, `g`, `b`, `[a]` range [`0-255`] or hex color string `#rrggbbaa`)
- new builtin function: `arcTo` like [WebKit Context2D arcTo](https://developer.mozilla.org/en-US/docs/Web/API/CanvasRenderingContext2D/arcTo)
- new builtin function: `fontSize(size)` to set the current font height
- mew builtin function: `star` to draw a star polygon
- mew builtin function: `fillAndStroke` to fill and stroke in one shot
- new builtin function: `imageGet` to load a PNG image form local file system
- new builtin function: `imageAt` to draw a PNG image at specified location on the main canvas
- new builtin function: `transform(x, y)` multiplies the point _x_, _y_ by the current matrix, returning a transformed position

BREAKING CHANGES:

- huge refactoring - GraphicContext is an interface, so, potentially we could use different drivers (PDF, SVG, etc.)
- REPL removed (now `g2d` works only as commandline tool)
- ~~`worldcoors(...)`~~ renamed to `viewport(...)`
- ~~`pensize(...)`~~ renamed to `strokeWeight(...)`

---

0.4.1 (Dec 15, 2020)

BUG FIXES:

- fix gorelaser build

---

0.4.0 (Dec 15, 2020)

NEW FEATURES:

- new object type: `Image`
- new builtin canvas functions: `width()`, `height()` to query canvas size at runtime
- new builtin canvas functions: `scale(sx, sy)`, `identity()` as matrix transformations
- new builtin canvas function: `loadPNG(filename)` to load an external PNG image
- new builtin canvas function: `image(im, x, y, ax, ay)` to draw a loaded PNG image
- new builtin canvas functions: `arcTo(x1, y1, x2, y2)` to draw a circular arc on path
- new builtin canvas functions: `rect(x, y, w, h, tl, tr, br, bl)` to draw a rectangle with different rounded angles
- new builtin math function: `map(value, start1, stop1, start2, stop2)` to re-map a number from one range to another
- new builtin math function: `lerp(start, stop, amt)` to calculate a number between two numbers at a specific increment
- new builtin math function: `min(n1, n2, ...n)` to calculate the minimum of a sequence
- new builtin math function: `max(n1, n2, ...n)` to calculate the maximum of a sequence


BREAKING CHANGES:

- massive code refactory
- renamed `screensize` function to `size`
- renamed `rectangle` function to `rect`
- renamed `saveState` function to `push`
- renamed `restoreState` function to `pop`

---

0.3.0 (Dec 9, 2020)

NEW FEATURES:

- canvas can be either square or rectangular; use `screensize(W,H)`
- new builtin math functions: `radians(angle)` and `degrees(angle)` for angles conversion
- new builtin canvas function: `fontsize([size])` to get or set the current font size
- new builtin canvas function: `text(str, x, y, ax, ay)` to write text on canvas
- new examples: spirograph, hypocycloid showing the use of custom functions
- new example: landscape to show howt to create not only squared images

BUG FIXES:

- Fix redundand rad->degrees conversion in `arc(...)` and `ellArc(...)` functions

CHANGES:

- `polygon(...)` builtin function now takes angles in radians
- `arc(...)` builtin function now takes angles in radians
- `ellArc(...)` builtin function now takes angles in radians
- `rotate(...)` builtin function now takes angles in radians

---

0.2.0 (Dec 7, 2020)

NEW FEATURES:

- better comparing between integers and floats 
- add new examples: dots, heart showing the use of custom function

BUG FIXES:

- ensure comparison between integers and floats in while loops works

---

0.1.0 (Dec 4, 2020)

First release