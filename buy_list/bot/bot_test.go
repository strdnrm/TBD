package bot

import "testing"

func TestDig(t *testing.T) {
	p := New()
	expected := Creature{
		BurrowLength: p.BurrowLength + 5,
		Health:       p.Health - 30,
		Respect:      p.Respect,
		Weight:       p.Weight,
	}
	p.Dig("intensively")
	if p != expected {
		t.Error("Incorrect parameter change")
	}
	expected = Creature{
		BurrowLength: p.BurrowLength + 2,
		Health:       p.Health - 10,
		Respect:      p.Respect,
		Weight:       p.Weight,
	}
	p.Dig("lazily")
	if p != expected {
		t.Error("Incorrect parameter change")
	}
}
