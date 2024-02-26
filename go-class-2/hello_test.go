package hello

import "testing"

func TestSayHello(t *testing.T){
  subtests := []struct{
    items []string
    result string
  }{
    {
      result: "Hello, world!",
    },
    {
      items: []string{},
      result: "Hello, world!",
    },
    {
      items: []string{"Preeti"},
      result: "Hello, Preeti!",
    },
    {
      items: []string{"Preeti", "Aaishi", "Eesha"},
      result: "Hello, Preeti, Aaishi, Eesha!",
    },
  }

  for _, st := range subtests {
    if s := Say(st.items); s != st.result {
      t.Errorf("wanted %s (%v), got %s", st.result, st.items, s)
    }
  }
}