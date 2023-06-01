dir-simulator acts a filesystem simulator, providing ability to run simple commands.
Usecase: simulate execution of basic filesystem commands, get output as you would in normal terminal (no actual changes are being made in the system).

Supported commands: dir, cd, up, mkdir, tree, mv
For examples of input and output please refer to resources directory.

to build the program run:
```
go build .
```
to run:
```
./dir-sumlator
```
program accepts optional flags:
```
./dir-simulator -input=path-to-input-file -output=path-to-output-file
```
default values are input.txt for input and output.txt for output
