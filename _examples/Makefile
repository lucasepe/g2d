CC := go run ../main.go
CFLAGS := eval -d .

SOURCES ?= $(wildcard *.g2d)

.PHONY: all clean

# pattern rule to parse all .g2d files and generate .png
%.png: %.g2d
	$(CC) $(CFLAGS) $<

# rule to make a .png to each .g2d
all: clean $(SOURCES:%.g2d=%.png)

clean:
	rm -f *.png