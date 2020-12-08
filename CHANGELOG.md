0.3.0 (Unreleased)

NEW FEATURES:

- new builtin functions: `radians(angle)` and `degrees(angle)` for angles conversion
- new builting function: `fontsize([size])` to get or set the current font size
- new builtin functions: `text(str, x, y, ax, ay)` and `textWrapped(str, x, y, ax, ay, width, lineSpacing, align)`
- new examples: spirograph, hypocycloid showing the use of custom functions

BUG FIXES:

- Fix redundand rad->degrees conversion in `arc(...)` and `ellArc(...)` functions

CHANGES:

- `polygon(...)` builtin function now takes angles in radians
- `arc(...)` builtin function now takes angles in radians
- `ellArc(...)` builtin function now takes angles in radians
- `rotate(...)` builtin function now takes angles in radians

---

0.2.0 (Mon 7, 2020)

NEW FEATURES:

- better comparing between integers and floats 
- add new examples: dots, heart showing the use of custom function

BUG FIXES:

- ensure comparison between integers and floats in while loops works

---

0.1.0 (Dec 4, 2020)

First release