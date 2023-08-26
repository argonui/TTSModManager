This golang library is intended to be able to convert source-controllable config
and luascript into a functioning json file that can be loaded into TTS as a
workshop mod.

# Main Features

* Handles Lua bundling!
* Handles XML bundling!
* Automatically builds on any incoming PR

# Getting the binary

The binary is built by slsa-framework/slsa-github-generator and can be found attached to the latest release, for example: https://github.com/argonui/TTSModManager/releases/tag/v0.2.4/TTSModManager.exe for windows. In the examples i'll refer to the exe, but you can use the TTSModManager-Liunux with the same expected behavior.

## Use automation in your repo

See https://github.com/argonui/TTSModManager.action for automated building of the mod on every PR / release.

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

## Working with downloadable content
TTS allows you to download content into a active game. This content must be a
json file in the form of a single object. In order to accomodate storing these
downloadable json files, TTSModManager can assemble and reverse these files.

### Generating a downloadable file

In this example foo.json the root of a partial mod object you'd like to
represent as a downloadable json.

```
TTSModManager.exe --moddir="C:\Users\USER\Documents\Projects\MyProject"
--objin="C:\Users\USER\Documents\Projects\MyProject\downloadable\content\foo.json"
--objout="C:\Users\USER\Documents\Projects\MyProject\to_be_downloaded.json"
```

### Reversing a downloadable file

This process assumes you already have your file you are used to downloading, and
want to decompose it into sub-objects and luascript etc.

Please note that **objout** is a directory and the trailing slash is needed.

```
TTSModManager.exe --moddir="C:\Users\USER\Documents\Projects\MyProject"
--reverse
--objin="C:\Users\USER\Documents\Projects\MyProject\ready_to_download.json"
--objout="C:\Users\USER\Documents\Projects\MyProject\downloadable\content\"
```


## Running a local copy

If you are developing a feature and would like to run the tool, use this instead of `TTSModManager.exe`

```
go run main.go --moddir="..."
```
