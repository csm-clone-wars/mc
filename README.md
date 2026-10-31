# mc
`mc` stands for `meta-c`:
    A pre-preprocessor for C and C++.
    
It generates `.h` header files for `.c` files and `.hpp` header files for `.cpp`, `.c++`, and `.cxx` files using custom macro syntax.

## Development Environment
- Operating System: Linux Mint
- Go Toolchain

## Example
The idea is:
In a single `.c` or `.cpp` file using custom comment macro syntax we write all our header content and give the file to `mc`, which generates the corresponding header file.
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

```c
// Even if empty each file should have the @pub macro at the top
//@pub
//@end

#include "hello.h"

int main(){
    hello();
    return 0;
}

```

- widget.cpp:

```cpp
/* Widget component */
//@pub
template <typename T>
class Widget {
public:
    Widget(T val);
};
//@end

template <typename T>
Widget<T>::Widget(T val) {}
```

## MC Rules
- Each `.c`, `.cpp`, `.c++`, or `.cxx` file should have `@pub` macro at the top.
- Before `@pub` macro whitespace characters and comments (both single-line `//` and block `/* ... */` comments) are allowed.
- `@pub` macro can be empty.
- Anything inside (between the macro comments `//@pub` and `//@end`) will go to the generated header file (`.h` for `.c`, `.hpp` for `.cpp`, `.c++`, `.cxx`).
- C++ constructs such as raw string literals (`R"(...)"`), block comments, string/char literals, templates, lambdas, and concepts won't break macro parsing.
- Every generated header file template includes `#pragma once` at the top, even if the `@pub` macro is empty.

## Commands
```
mc <no args>          -> Print mc version and help
mc -version,  -v      -> Print mc version
mc -help,     -h      -> Print mc commands list
mc -watch,    -w      -> Watch for any file changes and regenerate header files
mc hi.c widget.cpp    -> Generate .h/.hpp header files in include/ folder for given files
```

## License
MIT
