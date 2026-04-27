package rga

import "strings"

// AssetKind represents the type of asset (texture, sound, model, etc.).
type AssetKind string

// Asset kind constants.
const (
	KindUnknown AssetKind = ""
	KindModel   AssetKind = "model"
	KindTexture AssetKind = "texture"
	KindImage   AssetKind = "image"
	KindSound   AssetKind = "sound"
	KindMusic   AssetKind = "music"
	KindFont    AssetKind = "font"
	KindShader  AssetKind = "shader"
)

var knownKinds = map[AssetKind]struct{}{
	KindModel:   {},
	KindTexture: {},
	KindImage:   {},
	KindSound:   {},
	KindMusic:   {},
	KindFont:    {},
	KindShader:  {},
}

func (k AssetKind) String() string { return string(k) }

func (k AssetKind) DefaultDir() string {
	switch k {
	case KindModel:
		return "models"
	case KindTexture:
		return "textures"
	case KindImage:
		return "images"
	case KindSound, KindMusic:
		return "audio"
	case KindFont:
		return "fonts"
	case KindShader:
		return "shaders"
	default:
		return strings.ToLower(k.Plural())
	}
}

func (k AssetKind) Plural() string {
	s := string(k)
	if strings.HasSuffix(s, "h") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") && len(s) > 1 && !isVowel(s[len(s)-2]) {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

func (k AssetKind) IsKnown() bool {
	_, ok := knownKinds[k]
	return ok
}

func (k AssetKind) DefaultExtensions() []string {
	switch k {
	case KindModel:
		return []string{".obj", ".gltf", ".glb"}
	case KindTexture:
		return []string{".png", ".jpg", ".jpeg", ".bmp", ".gif"}
	case KindImage:
		return []string{".hdr", ".pic", ".ppm", ".pkm"}
	case KindSound:
		return []string{".wav", ".ogg"}
	case KindMusic:
		return []string{".mp3", ".flac", ".xm", ".mod"}
	case KindFont:
		return []string{".ttf", ".otf"}
	case KindShader:
		return []string{".glsl", ".frag", ".vert", ".vs", ".fs"}
	default:
		return nil
	}
}

func (k AssetKind) TypeName() string {
	switch k {
	case KindModel:
		return "ModelName"
	case KindTexture:
		return "TextureName"
	case KindImage:
		return "ImageName"
	case KindSound:
		return "SoundName"
	case KindMusic:
		return "MusicName"
	case KindFont:
		return "FontName"
	case KindShader:
		return "ShaderName"
	default:
		return capitalize(k.String()) + "Name"
	}
}