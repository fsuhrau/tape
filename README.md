# tape
tape is a small tool to manage command line dependencies in your project dir

## what does it do?
- it will download different versions of command line tools in a "central" space
- it will provide symbolic links to specific versions in your project

## requirements
- Go

## works on
- MacOS (tested)
- Linux (untested)
- Windows (untested)

## Installation
### from source
``` bash
$ go get -u -v github.com/fsuhrau/tape
```

## usage
### init
initialize tape in your current project
will create a new .tape dir and provide a config file which contains the dependencies
```
tape init
```

### add
you can easily add a new dependency by using the add command
you can also add directories for example SDKs as dependency there for the download URL must provide a zip file for now
```
tape add binary_name https://exmaple.com/downloads/v1/url_to_binary_name
tape add your_sdk https://exmaple.com/downloads/v1/your_sdk.zip
```

### remove
remove a dependency by remove command
it will unlink it and remove it from the dependencies you central space will stay untouched
```
tape remove binary_name
```

### update
update works similar to add but will update the existing dependency with a new version
```
tape update binary_name https://exmaple.com/downloads/v2/url_to_binary_name
tape add your_sdk https://exmaple.com/downloads/v2/your_sdk.zip
```

### list
you can check all your current dependencies with the list command
```
tape list
```

### link
you can download and link your current dependencies with link
link will create a symbolic link in .tape/links/ which point to the correct version in your "central" space
```
tape link
```
