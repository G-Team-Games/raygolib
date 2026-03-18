package raygolib

type Game interface {
	Update(dt float32) error
	Draw()
}
