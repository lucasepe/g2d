package gg

import (
	"fmt"
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
)

// FontStyle defines bold and italic styles for the font
// It is possible to combine values for mixed styles, eg.
//     FontData.Style = FontStyleBold | FontStyleItalic
type FontStyle byte

const (
	// FontStyleNormal regular style
	FontStyleNormal FontStyle = iota
	// FontStyleBold bold style
	FontStyleBold
	// FontStyleItalic italic style
	FontStyleItalic
)

// FontData wraps font info
type FontData struct {
	Name  string
	Style FontStyle
}

// RegisterFont registers a font to the global cache
func RegisterFont(fontData FontData, font *truetype.Font) {
	fontCache.Store(fontData, font)
}

// GetFont retrieves a font from the cache
func GetFont(fontData FontData) *truetype.Font {
	font, err := fontCache.Load(fontData)
	if err != nil {
		log.Println(err)
		return nil
	}

	return font
}

// GetGlobalFontCache returns the global font cache
func GetGlobalFontCache() FontCache { return fontCache }

// FontCache interface can be passed to SetFontCache to change the
// way fonts are being stored and retrieved.
type FontCache interface {
	// Loads a truetype font represented by the FontData object passed as
	// argument.
	// The method returns an error if the font could not be loaded, either
	// because it didn't exist or the resource it was loaded from was corrupted.
	Load(FontData) (*truetype.Font, error)

	// Sets the truetype font that will be returned by Load when given the font
	// data passed as first argument.
	Store(FontData, *truetype.Font)
}

// SetFontCache changes the font cache backend used by the package.
// To restore the default font cache, call this function passing nil as argument.
func SetFontCache(cache FontCache) {
	if cache == nil {
		fontCache = newDefaultFontCache()
	} else {
		fontCache = cache
	}
}

// DefaultFontCache defines the in memory default fonts cache
type DefaultFontCache map[string]*truetype.Font

// Store implements the FontCache store function
func (fc DefaultFontCache) Store(fd FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

// Load implements the FontCache load font function
func (fc DefaultFontCache) Load(fd FontData) (*truetype.Font, error) {
	font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("font %s is not stored in font cache", fd.Name)
	}
	return font, nil
}

func newDefaultFontCache() FontCache {
	res := DefaultFontCache{}

	TTFs := map[string]([]byte){
		"goregular": goregular.TTF,
		"gobold":    gobold.TTF,
		"goitalic":  goitalic.TTF,
		"gomono":    gomono.TTF,
	}

	for fontName, TTF := range TTFs {
		font, err := truetype.Parse(TTF)
		if err != nil {
			panic(err)
		}
		res.Store(FontData{Name: fontName}, font)
	}

	return res
}

var (
	fontCache = newDefaultFontCache()
)
