<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [C++ Code Style Guide](#c-code-style-guide)
  - [Header Files](#header-files)
  - [Naming Conventions](#naming-conventions)
  - [Classes](#classes)
  - [Functions](#functions)
  - [Formatting](#formatting)
  - [Other C++ Features](#other-c-features)
  - [Comments](#comments)
  - [Memory Management](#memory-management)
  - [Error Handling](#error-handling)
  - [Tools](#tools)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# C++ Code Style Guide

This guide follows Google's C++ Style Guide.

## Header Files

- Every `.cpp` file should have an associated `.h` file
- Use `#pragma once` for header guards
- Include headers in the following order:
  1. Related header
  2. C system headers
  3. C++ standard library headers
  4. Other libraries' headers
  5. Your project's headers
- Forward declare when possible instead of including headers

## Naming Conventions

- File names: `snake_case.cc` and `snake_case.h`
- Type names (classes, structs, enums): `PascalCase`
- Variables and function names: `snake_case`
- Class member variables: `snake_case_` with trailing underscore
- Constants: `kConstantName` with leading 'k'
- Namespace names: `snake_case`
- Enumerators: `kEnumName` with leading 'k'
- Macro names: `MACRO_NAME` in all capitals with underscores

## Classes

- Organize class definitions as:
  1. Public section
  2. Protected section
  3. Private section
- Write a constructor if a class has member variables
- Use explicit for single-argument constructors
- Define destructors virtual if a class has virtual functions
- Prefer composition over inheritance

```cpp
class MyClass {
 public:
  explicit MyClass(int value);
  ~MyClass() override;

  void DoSomething();

 private:
  void HelperMethod();

  int value_;
  std::string name_;
};
```

## Functions

- Keep functions short and focused
- Prefer returning values over output parameters
- Use const for parameters that are not modified
- Write short inline functions in the header file
- Parameters ordering: inputs, then outputs

## Formatting

- Line length: 80 characters maximum
- Use 2 spaces for indentation (no tabs)
- Function declarations and definitions: return type on the same line
- Braces follow Stroustrup style:
  - Opening braces at the end of line
  - Closing braces on their own line
- If/else statements: braces required even for single-line blocks
- Switch statements: braces for each case optional

```cpp
if (condition) {
  DoOneThing();
} else {
  DoAnotherThing();
}

for (int i = 0; i < limit; ++i) {
  DoSomething();
}

void MyClass::MyFunction(int parameter) {
  DoSomething();
}
```

## Other C++ Features

- Use C++17 features when available
- Use smart pointers (`std::unique_ptr`, `std::shared_ptr`) instead of raw pointers
- Use `auto` only when it improves readability
- Use range-based for loops when possible
- Prefer `enum class` over `enum`
- Use `nullptr` instead of `NULL` or `0`
- Use `constexpr` for values known at compile time

## Comments

- Use `//` for comments, not `/* */`
- Every file, class, function should have a comment
- Use descriptive comments explaining why, not what

```cpp
// Returns the next available ID. Thread-safe.
int GetNextId();
```

## Memory Management

- Prefer automatic variables to heap allocation
- Use containers like `std::vector` instead of C arrays
- Follow the Rule of Three (or Rule of Five in C++11)
- Use RAII design (Resource Acquisition Is Initialization)

## Error Handling

- Use exceptions for exceptional cases only
- Use assertions (`CHECK()`) for invariants
- Use status objects for expected errors

## Tools

- Use clang-format with Google's style configuration
- Use clang-tidy for static analysis
- Configure your IDE to match Google C++ style
