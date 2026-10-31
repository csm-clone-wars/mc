/*
 * Widget component example in C++
 * Demonstrates C++ constructs with @pub macro
 */

//@pub
#include <string>
#include <vector>

template <typename T>
class Widget {
private:
    T value;
    std::string marker = "//@end inside string literal";

public:
    explicit Widget(T val) : value(val) {}

    T get() const { return value; }

    const char* rawSample() const {
        return R"(
//@pub
void sample();
//@end
)";
    }
};
//@end

#include <iostream>

int main() {
    Widget<int> w(42);
    std::cout << "Widget value: " << w.get() << std::endl;
    return 0;
}
