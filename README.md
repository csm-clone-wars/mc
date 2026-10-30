# mc
`mc` stands for `meta-c`:
    A pre-preprocessor for C.
    
It generates .h header files for .c files using custom macro syntax.

## Development Environemnt
- Operating System: Linux Mint
- Go Toolchain

## Example
The idea is:
In a single .c file using custom comment macro syntax we write all our .h header content and give the file to `mc` and it will generate the header file.
- hello.c:

```c
// Hello module example!

//@pub
void hello();
//@end

#include<stdio.h>

void hello(){
    puts("Hello world!");
}

```
- main.c:

```
// Enven if empty each .c file should have the @pub macro at the top
//@pub
//@end

#include "hello.h"

int main(){
    hello();
    return 0;
}

```
## MC Rules
- Each .c file should have @pub macro at the top.
- Before @pub macro whitespace characters and single line comments are allowed.
- @pub macro can be empty.
- Anything inside (between the macro comments '//@pub' '//@end' will go to the .h header file).
- Every generated header file template includes `#pragma once` at the top, even if the @pub macro is empty.



## Commands
```
mc <no args>      -> Print mc version and help
mc -version,  -v  -> Print mc version
mc -help,     -h  -> Print mc commands list
mc -watch,    -w  -> Watch for any file changes and regenerate .h files
mc hi.c main.c    -> Generate .h header files in include/ folder for given .c files
```



## License
MIT
