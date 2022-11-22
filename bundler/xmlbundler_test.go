package bundler

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnbundleXML(t *testing.T) {
	input := `<!-- include Main -->
<Canvas raycastTarget="false">
  <Defaults>
    <Button color="a"/>
  </Defaults>

  <!-- include ui/CameraControl -->
  <Panel id="CameraControl"
         visibility="false"
  </Panel>

  <!-- include ui/CameraControl -->
  <!-- include ui/Shop -->
  <Panel id="shop.window" class="drag"
         offsetXY="480 -70"
         width="260" height="500">
         <!-- include deep -->
         <Button id="deepinclude">
         </Button>
         <!-- include deep -->
  </Panel>
  <!-- include ui/Shop -->

</Canvas>

<!-- include Main -->
`

	want := map[string]string{
		"__root": `<Include src="Main">
`,
		"Main": `<Canvas raycastTarget="false">
  <Defaults>
    <Button color="a"/>
  </Defaults>

  <Include src="ui/CameraControl">
  <Include src="ui/Shop">

</Canvas>
`,
		"deep": `<Button id="deepinclude">
</Button>`,
		"ui/Shop": `<Panel id="shop.window" class="drag"
       offsetXY="480 -70"
       width="260" height="500">
       <Include src="deep">
</Panel>`,
		"ui/CameraControl": `<Panel id="CameraControl"
       visibility="false"
</Panel>
`,
	}

	got, err := UnbundleAllXML(input)
	if err != nil {
		t.Fatalf("UnbundleAllXML(): %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}
