# vendordep
arrange go deps in vendor

```
currently only support:
git
```

# Usage

### 1. create vendordep.json file

```
{
  "Project": {
    "GroupId": "goodplayer",
    "Name": "vendordep",
    "ImportRootPath": "github.com/goodplayer/vendordep"
  },
  "Deps": [
    {
      "GroupId": "goodplayer",
      "Name": "vendordep",
      "ImportRootPath": "github.com/goodplayer/vendordep",
      "VcsType": "git",
      "VcsUrl": "https://github.com/goodplayer/vendordep.git",
      "Version": "81d6743ada34fcd511b5bf48281b44d8cbf4c7d6"
    }
  ]
}
```

### 2. run

```
vendordep get
```

# NOTE

best practice:

```
1. DO NOT place 'vendor' in project repository, just import using tools.
2. All imports ranges in the top same level, including inherited imports.
```
