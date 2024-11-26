# TTS Mod Manager (TTSMM) - Savegame Handling Documentation

## Overview

The TTS Mod Manager (TTSMM) is a tool designed to manage savegames from **Tabletop Simulator (TTS)**. TTS savegames are large JSON files that encapsulate the entire state of a game session, including all objects, settings, and scripts. Due to their size and complexity, these files are not well-suited for source control systems like GitHub.  

TTSMM provides functionality to split savegames into separate files for individual objects and reconstruct them as needed. This makes it easier to manage and version control TTS savegames.

## Feature: "Saved Object" Output

In addition to handling regular savegames, TTSMM supports generating "Saved Object" files. Saved Objects are a special type of savegame in TTS with the same overall structure as regular savegames, but with many of the outer-layer fields left empty. These files are typically used to save and share individual objects or small groups of objects independently of a full savegame.

### Key Differences: Regular Savegames vs. Saved Objects

| Field            | Regular Savegame Value     | Saved Object Value    |
| ---------------- | -------------------------- | --------------------- |
| `SaveName`       | Populated with save name   | Empty (`""`)          |
| `Date`           | Timestamp of save creation | Empty (`""`)          |
| `VersionNumber`  | Current TTS version        | Empty (`""`)          |
| `GameMode`       | Game mode description      | Empty (`""`)          |
| `GameType`       | Type of game               | Empty (`""`)          |
| `GameComplexity` | Complexity descriptor      | Empty (`""`)          |
| `Tags`           | Array of tags              | Empty (`[]`)          |
| `Gravity`        | Physics gravity setting    | Default (`0.5`)       |
| `PlayArea`       | Play area scale            | Default (`0.5`)       |
| `Table`          | Table model used           | Empty (`""`)          |
| `Sky`            | Skybox setting             | Empty (`""`)          |
| `Note`           | Session note               | Empty (`""`)          |
| `TabStates`      | Tab states information     | Empty (`{}`)          |
| `LuaScript`      | Global Lua script content  | Empty (`""`)          |
| `LuaScriptState` | Lua script state data      | Empty (`""`)          |
| `XmlUI`          | XML UI data                | Empty (`""`)          |
| `ObjectStates`   | Array of game objects      | Array of game objects |
