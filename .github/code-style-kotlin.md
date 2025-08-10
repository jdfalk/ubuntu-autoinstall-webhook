<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Kotlin Code Style Guide](#kotlin-code-style-guide)
  - [Source File Structure](#source-file-structure)
  - [Naming Conventions](#naming-conventions)
  - [Formatting and Indentation](#formatting-and-indentation)
  - [Classes and Interfaces](#classes-and-interfaces)
  - [Functions and Lambdas](#functions-and-lambdas)
  - [Properties](#properties)
  - [Control Flow](#control-flow)
  - [Null Safety](#null-safety)
  - [Collections and Functional Programming](#collections-and-functional-programming)
  - [Coroutines](#coroutines)
  - [Tools](#tools)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Kotlin Code Style Guide

This guide follows the official Kotlin coding conventions and Google's Android
Kotlin style guide.

## Source File Structure

- Source files should be encoded in UTF-8
- Each source file contains a single top-level class or multiple related
  top-level declarations
- File names should be descriptive and match the top-level class name
  (PascalCase.kt)
- Source files are organized by feature or layer, not by type

## Naming Conventions

- Package names: lowercase, no underscores (e.g., `com.example.myapp`)
- Class/Interface names: `PascalCase`
- Functions/Properties/Local variables: `camelCase`
- Constants (compile-time constants): `UPPER_SNAKE_CASE`
- Backing properties: `_camelCase` (with underscore prefix)
- Object instances representing singletons: `PascalCase` (same as class names)

```kotlin
// Naming convention examples
package com.example.myapp

const val MAX_COUNT = 100

class DatabaseHelper {
    private val _connection: Connection? = null
    val connection: Connection?
        get() = _connection

    fun connectToDatabase(url: String): Boolean {
        // Implementation
    }
}

object Logger {
    fun log(message: String) {
        // Implementation
    }
}
```

## Formatting and Indentation

- Use 4 spaces for indentation (no tabs)
- Maximum line length of 100 characters
- Use a space after control flow keywords (`if`, `when`, `for`)
- No space after function name in function calls
- Curly braces begin on the same line as the declaration
- When a function contains only a single expression, it can be expressed as an
  expression body

```kotlin
// Formatting examples
fun calculateSum(a: Int, b: Int): Int {
    return a + b
}

// Expression body function
fun multiply(a: Int, b: Int) = a * b

if (condition) {
    doSomething()
}
```

## Classes and Interfaces

- Constructor parameters can be declared as properties using `val` or `var`
- Organize members in a logical order, typically:
  1. Properties
  2. Init blocks
  3. Constructors
  4. Methods
  5. Companion object
- Use data classes for classes that primarily hold data
- Prefer composition over inheritance
- Group and sort property accessors with their property declaration

```kotlin
data class User(
    val id: String,
    val name: String,
    val email: String
)

class ProfileManager(private val userRepository: UserRepository) {
    private var currentUser: User? = null

    init {
        loadLastUser()
    }

    fun loadUser(userId: String): User? {
        // Implementation
    }

    companion object {
        private const val USER_CACHE_SIZE = 100
    }
}
```

## Functions and Lambdas

- Keep functions small and focused
- Prefer named parameters for better readability when a function takes multiple
  parameters
- Use trailing lambda syntax when the lambda is the last parameter
- Use the `it` identifier for single-parameter lambdas
- Use parameter names for multi-parameter lambdas

```kotlin
// Function examples with named parameters
fun createUser(id: String, name: String, email: String) {
    // Implementation
}
// Called with named parameters
createUser(id = "123", name = "John", email = "john@example.com")

// Lambda examples
items.filter { it > 10 }

items.fold(0) { acc, item -> acc + item }
```

## Properties

- Prefer properties over getter/setter functions
- Use custom accessors when logic is required
- Place custom accessors with their property definitions
- Use backing fields only when necessary

```kotlin
class Rectangle(val width: Int, val height: Int) {
    val area: Int
        get() = width * height

    var margin: Int = 0
        set(value) {
            field = maxOf(0, value) // field refers to the backing field
        }
}
```

## Control Flow

- Prefer using `when` over chained `if-else` statements
- For conditionals with multiple branches, put each branch on a separate line
- Use expression form for simple conditions
- Use curly braces for multi-line blocks

```kotlin
// When statement example
when (x) {
    0 -> println("Zero")
    1 -> println("One")
    2, 3 -> println("Two or Three")
    else -> println("Other")
}

// If expression
val max = if (a > b) a else b
```

## Null Safety

- Avoid using the `!!` operator
- Use safe calls (`?.`) or the Elvis operator (`?:`) for nullable values
- Use `let` with safe call for executing code only when reference is non-null
- Declare variables as non-nullable when possible

```kotlin
// Null safety examples
val name: String? = user?.name

// Safe call with let
user?.let {
    println("User name is ${it.name}")
}

// Elvis operator
val userName = user?.name ?: "Unknown"
```

## Collections and Functional Programming

- Prefer collection operations like `map`, `filter`, etc. over loops
- Use sequence operations for large collections to avoid intermediate collection
  creation
- Use destructuring declarations where appropriate

```kotlin
// Collection processing
val doubled = numbers.map { it * 2 }
val sum = numbers.fold(0) { acc, number -> acc + number }

// Destructuring
for ((key, value) in map) {
    println("$key -> $value")
}
```

## Coroutines

- Use structured concurrency principles
- Properly manage coroutine scopes and contexts
- Use `suspend` functions for asynchronous operations
- Avoid using `runBlocking` in production code
- Handle exceptions in coroutines

```kotlin
// Coroutine example
suspend fun fetchUserData(userId: String): UserData {
    return withContext(Dispatchers.IO) {
        api.getUserData(userId)
    }
}

// Launching coroutines
viewModelScope.launch {
    try {
        val userData = fetchUserData(userId)
        processData(userData)
    } catch (e: Exception) {
        handleError(e)
    }
}
```

## Tools

- Use ktlint or detekt for static analysis
- Configure Android Studio/IntelliJ formatting according to this style guide
- Enable line length warnings in your IDE
