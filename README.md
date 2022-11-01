This golang library is intended to be able to convert source-controllable config
and luascript into a functioning json file that can be loaded into TTS as a
workshop mod.


## Example Usage

### Generate json from a directory

$moddir = directory to read from
$config/output.json = file to output to

go run main.go --moddir=C:\Users\USER\Documents\Projects\MyProject

### Generate a directory from existing json file

$config = directory to write to
$ttsmodfile = existing tts mod file to read from

go run main.go --reverse --moddir=C:\Users\USER\Documents\Projects\MyProject --modfile="C:\Users\USER\Documents\My Games\Tabletop Simulator\Mods\Workshop\existingMod.json"

### Testing a TTS mod conversion

// reverse existing modfile $ttsmodfile into directory $moddir
go run main.go --reverse --moddir=C:\Users\USER\Documents\Projects\MyProject --modfile="C:\Users\USER\Documents\My Games\Tabletop Simulator\Mods\Workshop\existingMod.json"

// generate a modfile based on $moddir directory
go run main.go --moddir=C:\Users\USER\Documents\Projects\MyProject

// compare the original modfile ($modfile) with new generated modfile ($altmodfile)
go test . --modfile="C:\Users\USER\Documents\My Games\Tabletop Simulator\Mods\Workshop\existingMod.json" --altmodfile=""C:\Users\USER\Documents\Projects\MyProject\output.json""
