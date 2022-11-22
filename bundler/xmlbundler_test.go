package bundler

import (
	"ModCreator/tests"
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
		"__root": `<Include src="Main"/>
`,
		"Main": `<Canvas raycastTarget="false">
  <Defaults>
    <Button color="a"/>
  </Defaults>

  <Include src="ui/CameraControl"/>
  <Include src="ui/Shop"/>

</Canvas>
`,
		"deep": `<Button id="deepinclude">
</Button>`,
		"ui/Shop": `<Panel id="shop.window" class="drag"
       offsetXY="480 -70"
       width="260" height="500">
       <Include src="deep"/>
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

func TestBundleXML(t *testing.T) {
	want := `<!-- include Main -->
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

	input := map[string]string{
		"__root": `<Include src="Main"/>
`,
		"Main.xml": `<Canvas raycastTarget="false">
  <Defaults>
    <Button color="a"/>
  </Defaults>

  <Include src="ui/CameraControl"/>
  <Include src="ui/Shop"/>

</Canvas>
`,
		"deep.xml": `<Button id="deepinclude">
</Button>`,
		"ui/Shop.xml": `<Panel id="shop.window" class="drag"
       offsetXY="480 -70"
       width="260" height="500">
       <Include src="deep"/>
</Panel>`,
		"ui/CameraControl.xml": `<Panel id="CameraControl"
       visibility="false"
</Panel>
`,
	}
	ff := tests.NewFF()
	ff.Fs = input

	got, err := BundleXML(input[Rootname], ff)
	if err != nil {
		t.Fatalf("UnbundleAllXML(): %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("want != got:\n%v\n", diff)
	}
}
