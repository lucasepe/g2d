0.4.0 (Unreleased)

NEW FEATURES:

- new object type: `Image`
- new builtin canvas functions: `width()`, `height()` to query canvas size at runtime
- new builtin canvas functions: `scale(sx, sy)`, `identity()` as matrix transformations
- new builtin canvas function: `loadPNG(filename)` to load an external PNG image
- new builtin canvas function: `drawImage(im, x, y, ax, ay)` to draw a loaded PNG image
- more examples

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