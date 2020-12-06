PRJ := $(shell basename "$(PWD)")
IN := $(PRJ)_%03d.png
OUT := $(PRJ).webm

webm:
	ffmpeg -framerate 16 -f image2 -i $(IN) -c:v libvpx-vp9 \
	-quality good -cpu-used 0 -b:v 500K -crf 12 \
	-pix_fmt yuva420p -movflags faststart $(OUT)

