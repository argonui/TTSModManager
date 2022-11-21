# Description
This golang library is intended to be able to convert source-controllable config
and luascript into a functioning json file that can be loaded into TTS as a
workshop mod.

# Getting the binary

The binary is built by slsa-framework/slsa-github-generator and can be found attached to the latest release, for example: https://github.com/argonui/TTSModManager/releases/tag/v0.2.4/TTSModManager.exe for windows. In the examples i'll refer to the exe, but you can use the TTSModManager-Liunux with the same expected behavior.

# Example Usage
## Generate json from a directory
$moddir = directory to read from

```
TTSModManager.exe --moddir="C:\Users\USER\Documents\Projects\MyProject"
```

The finished json file is found in $moddir/output.json by default. if you'd like
to specify the output file you can use the `modfile` argument.

## Generate a directory from existing json file
$moddir = directory to write to
$modfile = existing tts mod file to read from

```
TTSModManager.exe --reverse --moddir="C:\Users\USER\Documents\Projects\MyProject" --modfile="C:\Users\USER\Documents\My Games\Tabletop Simulator\Mods\Workshop\existingMod.json"
```

If you'd like the bundled lua requirements to be written to the `src/` folder, pass `--writesrc`.

## Testing a TTS mod conversion
### reverse existing modfile into directory
$ttsmodfile = existing tts mod file to read from

$moddir = directory to write to
```
TTSModManager.exe --reverse --moddir="C:\Users\USER\Documents\Projects\MyProject" --modfile="C:\Users\USER\Documents\My Games\Tabletop Simulator\Mods\Workshop\existingMod.json"
```

### generate a modfile based on directory
$moddir = directory to read from
```
TTSModManager.exe --moddir="C:\Users\USER\Documents\Projects\MyProject"
```

## Running a local copy

If you are developing a feature and would like to run the tool, use this instead of `TTSModManager.exe`

```
go run main.go --moddir="..."
```
