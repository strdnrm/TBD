package creature

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

func TestEat(t *testing.T) {
	p := New()
	expected := Creature{
		BurrowLength: p.BurrowLength,
		Health:       p.Health + 10,
		Respect:      p.Respect,
		Weight:       p.Weight + 15,
	}
	p.Eat("withered")
	if p != expected {
		t.Error("Incorrect parameter change")
	}
	expected = Creature{
		BurrowLength: p.BurrowLength,
		Health:       p.Health - 30,
		Respect:      p.Respect,
		Weight:       p.Weight,
	}
	p.Eat("green")
	if p != expected {
		t.Error("Incorrect parameter change")
	}
}

func TestFight(t *testing.T) {
	p := New()
	bfrfight := p
	p.Fight("normal")
	if p == bfrfight {
		t.Error("Incorrect parameter change")
	}
}

func TestSleep(t *testing.T) {
	p := New()
	expected := Creature{
		BurrowLength: p.BurrowLength - 2,
		Health:       p.Health + 20,
		Respect:      p.Respect - 2,
		Weight:       p.Weight - 5,
	}
	p.Sleep()
	if p != expected {
		t.Error("Incorrect parameter change")
	}
}
